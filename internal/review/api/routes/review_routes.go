package routes

import (
	"context"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/review/api/handlers"
	"eshop-monolith/internal/review/domain/repositories"
	"eshop-monolith/internal/review/service"
	usermw "eshop-monolith/internal/user/middleware"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterReviewRoutes 注册评论与评分相关路由
func RegisterReviewRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, rabbit *rabbitmq.Client) {
	reviewRepo := repositories.NewReviewRepository(db)
	if err := reviewRepo.AutoMigrate(); err != nil {
		panic("failed to auto migrate review tables: " + err.Error())
	}

	// 桥接 order / user 模块（端口适配器）：
	// 将其他模块的领域模型转换为本模块服务所需的快照，避免跨模块类型耦合。
	findOrderByItem := buildOrderLookup(repos, db)
	findUser := buildUserLookup(repos)

	reviewSvc := service.NewReviewService(reviewRepo, rabbit, findOrderByItem, findUser)
	reviewHandler := handlers.NewReviewHandler(reviewSvc)

	// ---- 公开路由（无需登录）：浏览评论与评分 ----
	products := v1.Group("/products")
	{
		products.GET("/:id/reviews", reviewHandler.ListProductReviews)
		products.GET("/:id/rating", reviewHandler.GetProductRating)
	}

	// ---- 用户路由（需要登录）：发评论、看自己评论、删自己评论 ----
	reviews := v1.Group("/reviews")
	reviews.Use(middleware.JWTAuth())
	{
		reviews.POST("", reviewHandler.CreateReview)
		reviews.GET("/me", reviewHandler.ListMyReviews)
		// 删除评论：仅能删除本人评论（归属校验在 handler 内完成）
		reviews.DELETE("/:id", reviewHandler.DeleteMyReview)
	}

	// ---- 管理端路由（需要管理员权限）：审核、回复、删除 ----
	roleConfig := usermw.NewRequireRoleConfig(repos.Role)
	admin := v1.Group("/admin/reviews")
	admin.Use(middleware.JWTAuth(), usermw.RequireAdmin(roleConfig))
	{
		admin.GET("/pending", reviewHandler.ListPendingReviews)
		admin.PATCH("/:id/moderate", reviewHandler.ModerateReview)
		admin.POST("/:id/reply", reviewHandler.ReplyReview)
		admin.DELETE("/:id", reviewHandler.AdminDeleteReview)
	}
}

// buildOrderLookup 构造「订单项 ID → 订单快照」适配器
//
// 实现路径：order_items.order_id → orders(含 items) → 快照。
// 由于 order 模块仓储接口（IorderRepository）在本包外不可直接命名，
// 这里通过 repos.Order 调用其方法（鸭子类型），并借助 db 解析 order_item → order_id。
func buildOrderLookup(repos *repository.Repositories, db *gorm.DB) service.OrderByItemLookup {
	return func(ctx context.Context, orderItemID int64) (*service.OrderSnapshot, error) {
		// 1. 由 order_item_id 解析 order_id
		var orderID int64
		if err := db.WithContext(ctx).Table("order_items").
			Select("order_id").Where("id = ?", orderItemID).
			Scan(&orderID).Error; err != nil {
			return nil, err
		}
		if orderID == 0 {
			return nil, gorm.ErrRecordNotFound
		}

		// 2. 复用 order 模块仓储查订单（含 items）
		order, err := repos.Order.FindByID(ctx, orderID)
		if err != nil {
			return nil, err
		}

		// 3. 转换为快照（与 order 领域模型解耦）
		items := make([]service.OrderItemSnapshot, 0, len(order.Items))
		for _, it := range order.Items {
			items = append(items, service.OrderItemSnapshot{
				ID:        it.ID,
				ProductID: it.ProductID, // order 模块中 ProductID 为 string
			})
		}
		return &service.OrderSnapshot{
			ID:         order.ID,
			CustomerID: order.CustomerID,
			OrderNo:    order.OrderNo,
			Items:      items,
		}, nil
	}
}

// buildUserLookup 构造「用户 ID → 用户信息快照」适配器
func buildUserLookup(repos *repository.Repositories) service.UserInfoLookup {
	return func(ctx context.Context, userID int64) (*service.UserInfoSnapshot, error) {
		info, err := repos.UserInfo.GetUserInfoByUserID(ctx, userID)
		if err != nil {
			// 用户信息缺失时返回空快照，使评论仍可展示（昵称为空）
			if err == gorm.ErrRecordNotFound {
				return &service.UserInfoSnapshot{}, nil
			}
			return nil, err
		}
		return &service.UserInfoSnapshot{
			Nickname: info.Nickname,
			Avatar:   info.Avatar,
		}, nil
	}
}
