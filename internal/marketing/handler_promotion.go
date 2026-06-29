package marketing

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

// Create 创建促销
// @Summary 创建促销
// @Tags promotions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreatePromotionReq true "促销信息"
// @Success 200 {object} response.Response{data=Promotion}
// @Router /api/v1/promotions [post]
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

// GetByID 获取促销详情
// @Summary 获取促销详情
// @Tags promotions
// @Produce json
// @Param id path int true "促销ID"
// @Success 200 {object} response.Response{data=Promotion}
// @Router /api/v1/promotions/{id} [get]
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

// List 促销列表
// @Summary 促销列表
// @Tags promotions
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=PromotionListResult}
// @Router /api/v1/promotions [get]
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

// Update 更新促销
// @Summary 更新促销
// @Tags promotions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "促销ID"
// @Param request body UpdatePromotionReq true "促销信息"
// @Success 200 {object} response.Response{data=Promotion}
// @Router /api/v1/promotions/{id} [put]
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

// Delete 删除促销
// @Summary 删除促销
// @Tags promotions
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "促销ID"
// @Success 200 {object} response.Response
// @Router /api/v1/promotions/{id} [delete]
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
		coupon.GET("", ch.ListUserCoupons)
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
