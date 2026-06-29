package review

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/user"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"eshop-monolith/pkg/errcode"
)

type ReviewHandler struct {
	svc *ReviewService
}

func NewReviewHandler(svc *ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// CreateReview 创建评论
// @Summary 创建评论
// @Tags reviews
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateReviewReq true "评论信息"
// @Success 200 {object} response.Response{data=ReviewResp}
// @Router /api/v1/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	v, _ := c.Get("user_id")
	userID, _ := v.(int64)

	var req CreateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	rv, err := h.svc.CreateReview(c, CreateReviewInput{
		ProductID:   req.ProductID,
		UserID:      userID,
		OrderItemID: req.OrderItemID,
		Rating:      req.Rating,
		Content:     req.Content,
		Media:       req.Media,
	})
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, toReviewResp(rv))
}

// ListProductReviews 产品评论列表
// @Summary 产品评论列表
// @Tags reviews
// @Produce json
// @Param id path int true "产品ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=ReviewListResult}
// @Router /api/v1/products/{id}/reviews [get]
func (h *ReviewHandler) ListProductReviews(c *gin.Context) {
	productID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	page, size := parsePageSize(c)
	list, total, err := h.svc.ListProductReviews(c, productID, []string{string(ReviewStatusApproved)}, page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &ReviewListResult{
		Total: total,
		List:  toReviewRespList(list),
	})
}

// ListMyReviews 我的评论
// @Summary 我的评论
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=ReviewListResult}
// @Router /api/v1/reviews/me [get]
func (h *ReviewHandler) ListMyReviews(c *gin.Context) {
	v, _ := c.Get("user_id")
	userID, _ := v.(int64)

	page, size := parsePageSize(c)
	list, total, err := h.svc.ListUserReviews(c, userID, page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &ReviewListResult{
		Total: total,
		List:  toReviewRespList(list),
	})
}

// GetProductRating 产品评分汇总
// @Summary 产品评分汇总
// @Tags reviews
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=ProductRatingResp}
// @Router /api/v1/products/{id}/rating [get]
func (h *ReviewHandler) GetProductRating(c *gin.Context) {
	productID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	summary, err := h.svc.GetProductRating(c, productID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, toProductRatingResp(summary))
}

// DeleteMyReview 删除我的评论
// @Summary 删除我的评论
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews/{id} [delete]
func (h *ReviewHandler) DeleteMyReview(c *gin.Context) {
	v, _ := c.Get("user_id")
	userID, _ := v.(int64)

	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	rv, err := h.svc.GetReview(c, reviewID)
	if err != nil {
		c.Error(err)
		return
	}
	if rv.UserID != userID {
		c.Error(errcode.ErrReviewNotOwner)
		return
	}
	if err := h.svc.DeleteReview(c, reviewID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ListPendingReviews 待审核评论列表
// @Summary 待审核评论列表
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=ReviewListResult}
// @Router /api/v1/admin/reviews/pending [get]
func (h *ReviewHandler) ListPendingReviews(c *gin.Context) {
	page, size := parsePageSize(c)
	list, total, err := h.svc.ListPendingReviews(c, page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &ReviewListResult{
		Total: total,
		List:  toReviewRespList(list),
	})
}

// ModerateReview 审核评论
// @Summary 审核评论
// @Tags reviews
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param request body ModerateReviewReq true "审核信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/reviews/{id}/moderate [patch]
func (h *ReviewHandler) ModerateReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req ModerateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.ModerateReview(c, reviewID, req.Status); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ReplyReview 商家回复
// @Summary 商家回复
// @Tags reviews
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param request body ReplyReviewReq true "回复内容"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/reviews/{id}/reply [post]
func (h *ReviewHandler) ReplyReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req ReplyReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.ReplyReview(c, reviewID, req.Reply); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// AdminDeleteReview 管理员删除评论
// @Summary 管理员删除评论
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/reviews/{id} [delete]
func (h *ReviewHandler) AdminDeleteReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.DeleteReview(c, reviewID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ── helpers ────────────────────────────────────────

func parsePageSize(c *gin.Context) (int, int) {
	page, size := 1, 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(c.Query("size")); err == nil && s > 0 && s <= 100 {
		size = s
	}
	return page, size
}

func toReviewResp(r *Review) ReviewResp {
	resp := ReviewResp{
		ID:          r.ID,
		ProductID:   r.ProductID,
		UserID:      r.UserID,
		UserName:    r.UserName,
		UserAvatar:  r.UserAvatar,
		OrderItemID: r.OrderItemID,
		OrderNo:     r.OrderNo,
		Rating:      r.Rating,
		Content:     r.Content,
		Media:       r.GetMedia(),
		Status:      r.Status,
		Reply:       r.Reply,
		CreatedAt:   r.CreatedAt.UnixMilli(),
		UpdatedAt:   r.UpdatedAt.UnixMilli(),
	}
	if r.ReplyAt != nil {
		resp.ReplyAt = r.ReplyAt.UnixMilli()
	}
	return resp
}

func toReviewRespList(list []*Review) []*ReviewResp {
	resp := make([]*ReviewResp, len(list))
	for i, r := range list {
		v := toReviewResp(r)
		resp[i] = &v
	}
	return resp
}

func toProductRatingResp(s *ProductRatingSummary) ProductRatingResp {
	return ProductRatingResp{
		ProductID:     s.ProductID,
		AverageRating: s.AverageRating,
		ReviewCount:   s.ReviewCount,
		Rating1Count:  s.Rating1Count,
		Rating2Count:  s.Rating2Count,
		Rating3Count:  s.Rating3Count,
		Rating4Count:  s.Rating4Count,
		Rating5Count:  s.Rating5Count,
	}
}

// ── Routes ────────────────────────────────────────

func RegisterReviewRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	repo := NewReviewRepository(db)
	findOrderByItem := buildOrderLookup(repos, db)
	findUser := buildUserLookup(repos)
	svc := NewReviewService(repo, findOrderByItem, findUser)
	h := NewReviewHandler(svc)

	// 公开路由
	products := v1.Group("/products")
	{
		products.GET("/:id/reviews", h.ListProductReviews)
		products.GET("/:id/rating", h.GetProductRating)
	}

	// 用户路由（需要登录）
	reviews := v1.Group("/reviews")
	reviews.Use(middleware.JWTAuth())
	{
		reviews.POST("", h.CreateReview)
		reviews.GET("/me", h.ListMyReviews)
		reviews.DELETE("/:id", h.DeleteMyReview)
	}

	// 管理端路由（需要管理员权限）
	roleCfg := user.NewRequireRoleConfig(repos.Role)
	admin := v1.Group("/admin/reviews")
	admin.Use(middleware.JWTAuth(), user.RequireAdmin(roleCfg))
	{
		admin.GET("/pending", h.ListPendingReviews)
		admin.PATCH("/:id/moderate", h.ModerateReview)
		admin.POST("/:id/reply", h.ReplyReview)
		admin.DELETE("/:id", h.AdminDeleteReview)
	}
}

func buildOrderLookup(repos *repository.Repositories, db *gorm.DB) OrderByItemLookup {
	return func(ctx context.Context, orderItemID int64) (*OrderSnapshot, error) {
		var orderID int64
		if err := db.WithContext(ctx).Table("order_items").
			Select("order_id").Where("id = ?", orderItemID).
			Scan(&orderID).Error; err != nil {
			return nil, err
		}
		if orderID == 0 {
			return nil, gorm.ErrRecordNotFound
		}
		order, err := repos.Order.FindByID(ctx, orderID)
		if err != nil {
			return nil, err
		}
		items := make([]OrderItemSnapshot, 0, len(order.Items))
		for _, it := range order.Items {
			items = append(items, OrderItemSnapshot{
				ID:        it.ID,
				ProductID: it.ProductID,
			})
		}
		return &OrderSnapshot{
			ID:         order.ID,
			CustomerID: order.CustomerID,
			OrderNo:    order.OrderNo,
			Items:      items,
		}, nil
	}
}

func buildUserLookup(repos *repository.Repositories) UserInfoLookup {
	return func(ctx context.Context, userID int64) (*UserInfoSnapshot, error) {
		user, err := repos.User.FindByID(ctx, userID)
		if err != nil {
			return &UserInfoSnapshot{}, nil
		}
		return &UserInfoSnapshot{
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		}, nil
	}
}
