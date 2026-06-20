package handlers

import (
	"strconv"
	"time"

	"eshop-monolith/internal/coupon/api/dto"
	couponModels "eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/internal/coupon/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	couponService *service.CouponService
}

func NewCouponHandler(couponService *service.CouponService) *CouponHandler {
	return &CouponHandler{couponService: couponService}
}

// getCurrentUserID 从上下文中获取当前用户ID
func getCurrentUserID(c *gin.Context) int64 {
	v, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	switch id := v.(type) {
	case float64:
		return int64(id)
	case uint:
		return int64(id)
	case int64:
		return id
	case string:
		if parsed, err := strconv.ParseInt(id, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}

// CreateCoupon 创建优惠券模板
// @Summary 创建优惠券模板
// @Description 创建新的优惠券模板，支持满减券、折扣券、代金券
// @Tags coupons
// @Accept json
// @Produce json
// @Param coupon body dto.CreateCouponReq true "优惠券信息"
// @Success 200 {object} response.Response{data=dto.CouponResponse}
// @Router /api/v1/admin/coupons [post]
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req dto.CreateCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		c.Error(err)
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		c.Error(err)
		return
	}

	coupon := &couponModels.Coupon{
		Name:        req.Name,
		Description: req.Description,
		CouponType:  couponModels.CouponType(req.CouponType),
		Scope:       couponModels.CouponScope(req.Scope),
		ScopeValue:  req.ScopeValue,
		Value:       req.Value,
		MinAmount:   req.MinAmount,
		MaxDiscount: req.MaxDiscount,
		TotalStock:  req.TotalStock,
		UserLimit:   req.UserLimit,
		StartTime:   utils.Timestamp(startTime),
		EndTime:     utils.Timestamp(endTime),
		ValidDays:   req.ValidDays,
	}

	if err := h.couponService.CreateCoupon(c, coupon); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toCouponResponse(coupon))
}

// UpdateCoupon 更新优惠券模板
// @Summary 更新优惠券模板
// @Description 更新指定优惠券模板的信息
// @Tags coupons
// @Accept json
// @Produce json
// @Param id path int true "优惠券模板ID"
// @Param coupon body dto.UpdateCouponReq true "更新信息"
// @Success 200 {object} response.Response{data=dto.CouponResponse}
// @Router /api/v1/admin/coupons/{id} [put]
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	existing, err := h.couponService.GetCoupon(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Scope != "" {
		existing.Scope = couponModels.CouponScope(req.Scope)
	}
	if req.ScopeValue != "" {
		existing.ScopeValue = req.ScopeValue
	}
	if req.Value > 0 {
		existing.Value = req.Value
	}
	existing.MinAmount = req.MinAmount
	existing.MaxDiscount = req.MaxDiscount
	if req.UserLimit > 0 {
		existing.UserLimit = req.UserLimit
	}
	if req.StartTime != "" {
		t, parseErr := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if parseErr != nil {
			c.Error(parseErr)
			return
		}
		existing.StartTime = utils.Timestamp(t)
	}
	if req.EndTime != "" {
		t, parseErr := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if parseErr != nil {
			c.Error(parseErr)
			return
		}
		existing.EndTime = utils.Timestamp(t)
	}
	if req.ValidDays > 0 {
		existing.ValidDays = req.ValidDays
	}
	if req.Status != "" {
		existing.Status = couponModels.CouponStatus(req.Status)
	}

	if err := h.couponService.UpdateCoupon(c, existing); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toCouponResponse(existing))
}

// GetCoupon 获取优惠券详情
// @Summary 获取优惠券详情
// @Description 根据ID获取优惠券模板详情
// @Tags coupons
// @Accept json
// @Produce json
// @Param id path int true "优惠券模板ID"
// @Success 200 {object} response.Response{data=dto.CouponResponse}
// @Router /api/v1/coupons/{id} [get]
func (h *CouponHandler) GetCoupon(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	coupon, err := h.couponService.GetCoupon(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toCouponResponse(coupon))
}

// ListCoupons 优惠券模板列表
// @Summary 优惠券模板列表
// @Description 分页查询优惠券模板列表
// @Tags coupons
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=dto.CouponListResult}
// @Router /api/v1/coupons [get]
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	list, total, err := h.couponService.ListCoupons(c, int(page), int(pageSize))
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.CouponResponse, len(list))
	for i, cp := range list {
		items[i] = toCouponResponse(&cp)
	}

	response.Success(c, dto.CouponListResult{
		Total: total,
		List:  items,
	})
}

// ClaimCoupon 领取优惠券
// @Summary 领取优惠券
// @Description 当前登录用户领取指定优惠券
// @Tags coupons
// @Accept json
// @Produce json
// @Param claim body dto.ClaimCouponReq true "领取请求"
// @Success 200 {object} response.Response{data=dto.UserCouponResponse}
// @Router /api/v1/coupons/claim [post]
func (h *CouponHandler) ClaimCoupon(c *gin.Context) {
	var req dto.ClaimCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userID := getCurrentUserID(c)

	uc, err := h.couponService.ClaimCoupon(c, userID, req.CouponID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toUserCouponResponse(uc))
}

// UseCoupon 使用优惠券
// @Summary 使用优惠券
// @Description 结算时使用优惠券，标记为已使用并返回抵扣金额
// @Tags coupons
// @Accept json
// @Produce json
// @Param use body dto.UseCouponReq true "使用请求"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/v1/coupons/use [post]
func (h *CouponHandler) UseCoupon(c *gin.Context) {
	var req dto.UseCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userID := getCurrentUserID(c)

	discount, err := h.couponService.UseCoupon(c, req.UserCouponID, userID, req.OrderNo, req.OrderAmount)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"discount": discount,
		"order_no": req.OrderNo,
	})
}

// GetUserCoupons 获取用户优惠券列表
// @Summary 获取用户优惠券列表
// @Description 当前登录用户的优惠券列表，可按状态筛选
// @Tags coupons
// @Accept json
// @Produce json
// @Param status query string false "状态过滤：unused/used/expired"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=dto.UserCouponListResult}
// @Router /api/v1/coupons/mine [get]
func (h *CouponHandler) GetUserCoupons(c *gin.Context) {
	userID := getCurrentUserID(c)
	if userID <= 0 {
		// 如果没登录，尝试从路径参数获取
		if parsed, err := utils.ParseIntParam(c, "user_id"); err == nil {
			userID = parsed
		}
	}

	status := couponModels.UserCouponStatus(c.DefaultQuery("status", ""))
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	list, total, err := h.couponService.GetUserCoupons(c, userID, status, int(page), int(pageSize))
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.UserCouponResponse, len(list))
	for i, uc := range list {
		items[i] = toUserCouponResponse(&uc)
	}

	response.Success(c, dto.UserCouponListResult{
		Total: total,
		List:  items,
	})
}

// GetUsableCoupons 获取用户可用优惠券
// @Summary 获取用户可用优惠券
// @Description 当前登录用户所有未使用且未过期的优惠券
// @Tags coupons
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.UserCouponResponse}
// @Router /api/v1/coupons/usable [get]
func (h *CouponHandler) GetUsableCoupons(c *gin.Context) {
	userID := getCurrentUserID(c)

	list, err := h.couponService.GetUsableCoupons(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.UserCouponResponse, len(list))
	for i, uc := range list {
		items[i] = toUserCouponResponse(&uc)
	}

	response.Success(c, items)
}

// ValidateCoupon 预校验优惠券
// @Summary 预校验优惠券
// @Description 结算前校验优惠券是否可用，返回预估抵扣金额
// @Tags coupons
// @Accept json
// @Produce json
// @Param validate body object{user_coupon_id=int,order_amount=int} true "校验请求"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/v1/coupons/validate [post]
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	userID := getCurrentUserID(c)

	var req struct {
		UserCouponID int64 `json:"user_coupon_id" binding:"required"`
		OrderAmount  int64 `json:"order_amount" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	discount, err := h.couponService.ValidateCoupon(c, req.UserCouponID, userID, req.OrderAmount)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"valid":    true,
		"discount": discount,
	})
}

// toCouponResponse 转换优惠券模板到响应
func toCouponResponse(cp *couponModels.Coupon) dto.CouponResponse {
	return dto.CouponResponse{
		ID:          cp.ID,
		Name:        cp.Name,
		Description: cp.Description,
		CouponType:  string(cp.CouponType),
		Scope:       string(cp.Scope),
		ScopeValue:  cp.ScopeValue,
		Value:       cp.Value,
		MinAmount:   cp.MinAmount,
		MaxDiscount: cp.MaxDiscount,
		TotalStock:  cp.TotalStock,
		RemainStock: cp.RemainStock,
		UserLimit:   cp.UserLimit,
		StartTime:   cp.StartTime,
		EndTime:     cp.EndTime,
		ValidDays:   cp.ValidDays,
		Status:      string(cp.Status),
		CreatedAt:   cp.CreatedAt,
		UpdatedAt:   cp.UpdatedAt,
	}
}

// toUserCouponResponse 转换用户优惠券到响应
func toUserCouponResponse(uc *couponModels.UserCoupon) dto.UserCouponResponse {
	resp := dto.UserCouponResponse{
		ID:         uc.ID,
		UserID:     uc.UserID,
		CouponID:   uc.CouponID,
		CouponCode: uc.CouponCode,
		OrderNo:    uc.OrderNo,
		Status:     string(uc.Status),
		ExpireAt:   utils.Timestamp(uc.ExpireAt),
		CreatedAt:  uc.CreatedAt,
	}
	if uc.UsedAt != nil {
		t := utils.Timestamp(*uc.UsedAt)
		resp.UsedAt = &t
	}
	return resp
}
