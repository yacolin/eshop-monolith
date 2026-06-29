package promotion

import (
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PromotionHandler struct {
	svc *PromotionService
}

func NewPromotionHandler(svc *PromotionService) *PromotionHandler {
	return &PromotionHandler{svc: svc}
}

func (h *PromotionHandler) Create(c *gin.Context) {
	var req CreatePromotionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Create(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *PromotionHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.GetByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *PromotionHandler) List(c *gin.Context) {
	var req PromotionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.List(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *PromotionHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdatePromotionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Update(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *PromotionHandler) Delete(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Delete(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func RegisterPromotionRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewPromotionRepository(db)
	svc := NewPromotionService(repo, db)
	couponSvc := NewCouponService(repo, db)
	h := NewPromotionHandler(svc)
	ch := NewCouponHandler(svc, couponSvc)

	promo := v1.Group("/promotions")
	{
		promo.GET("", h.List)
		promo.GET("/:id", h.GetByID)
	}
	auth := v1.Group("/promotions")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.Create)
		auth.PUT("/:id", h.Update)
		auth.DELETE("/:id", h.Delete)
	}

	// 优惠券
	coupon := v1.Group("/coupons")
	coupon.Use(middleware.JWTAuth())
	{
		coupon.POST("/claim", ch.Claim)
		coupon.POST("/use", ch.Use)
		coupon.GET("/mine", ch.ListUserCoupons)
	}

	// 秒杀
	flashSvc := NewFlashService(repo, db)
	fh := NewFlashHandler(flashSvc)
	flash := v1.Group("/flash")
	flash.Use(middleware.JWTAuth())
	{
		flash.POST("/buy", fh.Buy)
		flash.POST("/confirm", fh.Confirm)
	}
}
