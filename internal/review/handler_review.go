package review

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type ReviewHandler struct {
	svc *ReviewService
}

func NewReviewHandler(svc *ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

func currentUserID(c *gin.Context) int64 {
	v, _ := c.Get("user_id")
	switch id := v.(type) {
	case int64:
		return id
	case uint:
		return int64(id)
	case float64:
		return int64(id)
	case int:
		return int64(id)
	}
	return 0
}

// CreateReview 创建评价
// @Summary 创建评价
// @Tags reviews
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateReviewReq true "评价信息"
// @Success 200 {object} response.Response{data=ReviewResp}
// @Router /api/v1/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req CreateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	rv, err := h.svc.CreateReview(c, CreateReviewInput{
		UserID:          currentUserID(c),
		OrderID:         req.OrderID,
		OrderItemID:     req.OrderItemID,
		SpuID:           req.SpuID,
		SkuID:           req.SkuID,
		OverallRating:   req.OverallRating,
		QualityRating:   req.QualityRating,
		LogisticsRating: req.LogisticsRating,
		ServiceRating:   req.ServiceRating,
		Content:         req.Content,
		IsAnonymous:     req.IsAnonymous,
	})
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, toReviewResp(rv))
}

// ListProductReviews 商品评价列表
// @Summary 商品评价列表
// @Tags reviews
// @Produce json
// @Param id path int true "商品SPU ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=ReviewListResult}
// @Router /api/v1/products/{id}/reviews [get]
func (h *ReviewHandler) ListProductReviews(c *gin.Context) {
	spuID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	page, size := parsePageSize(c)
	list, total, err := h.svc.ListBySpu(c, spuID, []int8{ReviewStatusApproved}, page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &ReviewListResult{Total: total, List: toReviewRespList(list)})
}

// ListMyReviews 我的评价
// @Summary 我的评价
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=ReviewListResult}
// @Router /api/v1/reviews/me [get]
func (h *ReviewHandler) ListMyReviews(c *gin.Context) {
	page, size := parsePageSize(c)
	list, total, err := h.svc.ListByUser(c, currentUserID(c), page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &ReviewListResult{Total: total, List: toReviewRespList(list)})
}

// GetProductRating 商品评分汇总
// @Summary 商品评分汇总
// @Tags reviews
// @Produce json
// @Param id path int true "商品SPU ID"
// @Success 200 {object} response.Response{data=ReviewRatingResp}
// @Router /api/v1/products/{id}/rating [get]
func (h *ReviewHandler) GetProductRating(c *gin.Context) {
	spuID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	rating, err := h.svc.GetRatingSummary(c, spuID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, toRatingResp(rating))
}

// DeleteMyReview 删除我的评价
// @Summary 删除我的评价
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "评价ID"
// @Success 200 {object} response.Response
// @Router /api/v1/reviews/{id} [delete]
func (h *ReviewHandler) DeleteMyReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.DeleteReview(c, currentUserID(c), reviewID, false); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ListPendingReviews 待审核评价列表
// @Summary 待审核评价列表
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=ReviewListResult}
// @Router /api/v1/admin/reviews/pending [get]
func (h *ReviewHandler) ListPendingReviews(c *gin.Context) {
	page, size := parsePageSize(c)
	list, total, err := h.svc.ListPending(c, page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, &ReviewListResult{Total: total, List: toReviewRespList(list)})
}

// ModerateReview 审核评价
// @Summary 审核评价
// @Tags reviews
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "评价ID"
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
	if err := h.svc.ModerateReview(c, reviewID, currentUserID(c), req.Status, req.Reason); err != nil {
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
// @Param id path int true "评价ID"
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
	if _, err := h.svc.ReplyReview(c, reviewID, currentUserID(c), req.Content); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// AdminDeleteReview 管理员删除评价
// @Summary 管理员删除评价
// @Tags reviews
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "评价ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/reviews/{id} [delete]
func (h *ReviewHandler) AdminDeleteReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.DeleteReview(c, currentUserID(c), reviewID, true); err != nil {
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

func toReviewResp(r *Review) *ReviewResp {
	return &ReviewResp{
		ID:              r.ID,
		UserID:          r.UserID,
		OrderID:         r.OrderID,
		OrderItemID:     r.OrderItemID,
		SpuID:           r.SpuID,
		SkuID:           r.SkuID,
		OverallRating:   r.OverallRating,
		QualityRating:   r.QualityRating,
		LogisticsRating: r.LogisticsRating,
		ServiceRating:   r.ServiceRating,
		Content:         r.Content,
		IsAnonymous:     r.IsAnonymous,
		Status:          r.Status,
		RejectReason:    r.RejectReason,
		ReplyCount:      r.ReplyCount,
		LikeCount:       r.LikeCount,
		HelpfulCount:    r.HelpfulCount,
		CreatedAt:       r.CreatedAt.UnixMilli(),
		UpdatedAt:       r.UpdatedAt.UnixMilli(),
	}
}

func toReviewRespList(list []*Review) []*ReviewResp {
	resp := make([]*ReviewResp, len(list))
	for i, r := range list {
		resp[i] = toReviewResp(r)
	}
	return resp
}

func toRatingResp(r *ReviewRating) ReviewRatingResp {
	return ReviewRatingResp{
		SpuID:              r.SpuID,
		AvgOverallRating:   r.AvgOverallRating,
		AvgQualityRating:   r.AvgQualityRating,
		AvgLogisticsRating: r.AvgLogisticsRating,
		AvgServiceRating:   r.AvgServiceRating,
		Rating5Count:       r.Rating5Count,
		Rating4Count:       r.Rating4Count,
		Rating3Count:       r.Rating3Count,
		Rating2Count:       r.Rating2Count,
		Rating1Count:       r.Rating1Count,
		TotalReviews:       r.TotalReviews,
		WithMediaCount:     r.WithMediaCount,
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
	admin := v1.Group("/admin/reviews")
	admin.Use(middleware.JWTAuth(), middleware.RequireAdmin())
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
		if err := db.WithContext(ctx).Table("tx_order_items").
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
		items, err := repos.Order.ListItems(ctx, orderID)
		if err != nil {
			return nil, err
		}
		itemsSnapshot := make([]OrderItemSnapshot, 0, len(items))
		for _, it := range items {
			itemsSnapshot = append(itemsSnapshot, OrderItemSnapshot{ID: it.ID, ProductID: strconv.FormatInt(it.ProductID, 10)})
		}
		return &OrderSnapshot{
			ID: order.ID, CustomerID: strconv.FormatInt(order.UserID, 10), OrderNo: order.OrderNo, Items: itemsSnapshot,
		}, nil
	}
}

func buildUserLookup(repos *repository.Repositories) UserInfoLookup {
	return func(ctx context.Context, userID int64) (*UserInfoSnapshot, error) {
		user, err := repos.User.FindByID(ctx, userID)
		if err != nil {
			return &UserInfoSnapshot{}, nil
		}
		return &UserInfoSnapshot{Nickname: user.Nickname, Avatar: user.Avatar}, nil
	}
}
