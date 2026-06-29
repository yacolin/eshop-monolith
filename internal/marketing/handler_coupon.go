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
