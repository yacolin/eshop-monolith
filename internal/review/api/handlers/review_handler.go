package handlers

import (
	"strconv"

	"eshop-monolith/internal/review/api/dto"
	"eshop-monolith/internal/review/domain/models"
	"eshop-monolith/internal/review/service"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 评论处理器
type ReviewHandler struct {
	svc *service.ReviewService
}

// NewReviewHandler 创建评论处理器
func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// CreateReview 创建评论（已登录用户）
// @Summary 创建评论
// @Description 已登录用户创建评论或评分，需提供订单项ID以校验已购买
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body dto.CreateReviewReq true "评论信息"
// @Success 200 {object} response.Response{data=dto.ReviewResp}
// @Router /api/v1/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.CreateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	rv, err := h.svc.CreateReview(c, service.CreateReviewInput{
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

	response.Success(c, dto.ToReviewResp(rv))
}

// ListProductReviews 查询产品评论列表（公开，仅展示已审核通过的评论）
// @Summary 查询产品评论列表
// @Description 根据产品ID查询已审核通过的评论列表，支持分页
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=dto.ReviewListResp}
// @Router /api/v1/products/{id}/reviews [get]
func (h *ReviewHandler) ListProductReviews(c *gin.Context) {
	productID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.ListReviewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	req.Pagination.Normalize()

	// 公开端点仅返回审核通过的评论
	result, err := h.svc.ListProductReviews(c, productID, []models.ReviewStatus{models.ReviewStatusApproved}, req.Page, req.Size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, &dto.ReviewListResp{
		Total: result.Total,
		List:  dto.ToReviewRespList(result.List),
	})
}

// ListMyReviews 查询当前用户的评论（含所有状态）
// @Summary 查询我的评论
// @Description 查询当前登录用户的所有评论（含所有审核状态），支持分页
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=dto.ReviewListResp}
// @Router /api/v1/reviews/me [get]
func (h *ReviewHandler) ListMyReviews(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.ListReviewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	req.Pagination.Normalize()

	result, err := h.svc.ListUserReviews(c, userID, req.Page, req.Size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, &dto.ReviewListResp{
		Total: result.Total,
		List:  dto.ToReviewRespList(result.List),
	})
}

// GetProductRating 获取产品评分汇总（公开）
// @Summary 获取产品评分汇总
// @Description 根据产品ID获取评分汇总，包含平均分和各星级数量
// @Tags reviews
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=dto.ProductRatingResp}
// @Router /api/v1/products/{id}/rating [get]
func (h *ReviewHandler) GetProductRating(c *gin.Context) {
	productID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	summary, err := h.svc.GetProductRating(c, productID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto.ToProductRatingResp(summary))
}

// DeleteMyReview 删除评论（仅评论者本人）
// @Summary 删除我的评论
// @Description 删除本人发表的评论，需校验评论归属
// @Tags reviews
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/reviews/{id} [delete]
func (h *ReviewHandler) DeleteMyReview(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	// 归属校验：仅本人可删除
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

	response.Success(c, gin.H{"message": "已删除"})
}

// ---------------------------------------------------------------------------
// 管理端接口
// ---------------------------------------------------------------------------

// ListPendingReviews 查询待审核评论列表（管理端）
// @Summary 查询待审核评论列表
// @Description 管理员查询所有待审核评论列表，支持分页
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=dto.ReviewListResp}
// @Router /api/v1/admin/reviews/pending [get]
func (h *ReviewHandler) ListPendingReviews(c *gin.Context) {
	var req dto.ListReviewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	req.Pagination.Normalize()

	result, err := h.svc.ListPendingReviews(c, req.Page, req.Size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, &dto.ReviewListResp{
		Total: result.Total,
		List:  dto.ToReviewRespList(result.List),
	})
}

// ModerateReview 审核评论（通过/拒绝/隐藏）
// @Summary 审核评论
// @Description 管理员审核评论，可设置为通过(approved)、拒绝(rejected)或隐藏(hidden)
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Param moderate body dto.ModerateReviewReq true "审核信息"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/admin/reviews/{id}/moderate [patch]
func (h *ReviewHandler) ModerateReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.ModerateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.ModerateReview(c, reviewID, req.Status); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "已更新审核状态", "status": req.Status})
}

// ReplyReview 商家回复评论
// @Summary 商家回复评论
// @Description 管理员/商家回复指定评论
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Param reply body dto.ReplyReviewReq true "回复内容"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/admin/reviews/{id}/reply [post]
func (h *ReviewHandler) ReplyReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.ReplyReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.ReplyReview(c, reviewID, req.Reply); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "回复成功"})
}

// AdminDeleteReview 管理员删除评论
// @Summary 管理员删除评论
// @Description 管理员强制删除指定评论
// @Tags reviews
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/admin/reviews/{id} [delete]
func (h *ReviewHandler) AdminDeleteReview(c *gin.Context) {
	reviewID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	if err := h.svc.DeleteReview(c, reviewID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "已删除"})
}

// getCurrentUserID 从 Gin Context 获取当前用户 ID
func getCurrentUserID(c *gin.Context) (int64, error) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, errcode.ErrUnauthorized
	}
	switch id := v.(type) {
	case float64:
		return int64(id), nil
	case uint:
		return int64(id), nil
	case int64:
		return id, nil
	case int:
		return int64(id), nil
	case string:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return 0, errcode.ErrUnauthorized
		}
		return n, nil
	default:
		return 0, errcode.ErrUnauthorized
	}
}
