package marketing

import (
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	svc      *PromotionService
	couponSvc *CouponService
}

func NewCouponHandler(svc *PromotionService, couponSvc *CouponService) *CouponHandler {
	return &CouponHandler{svc: svc, couponSvc: couponSvc}
}

// Claim 领取优惠券
// @Summary 领取优惠券
// @Tags coupons
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body ClaimCouponReq true "领取信息"
// @Success 200 {object} response.Response
// @Router /api/v1/coupons/claim [post]
func (h *CouponHandler) Claim(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req ClaimCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.couponSvc.Claim(c, userID.(int64), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Use 使用优惠券
// @Summary 使用优惠券
// @Tags coupons
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body UseCouponReq true "使用信息"
// @Success 200 {object} response.Response
// @Router /api/v1/coupons/use [post]
func (h *CouponHandler) Use(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req UseCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.couponSvc.Use(c, userID.(int64), &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ListUserCoupons 用户优惠券列表
// @Summary 用户优惠券列表
// @Tags coupons
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=UserPromotionListResult}
// @Router /api/v1/coupons [get]
func (h *CouponHandler) ListUserCoupons(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req UserPromotionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.couponSvc.ListUserCoupons(c, userID.(int64), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}
