package routes

import (
	"eshop-monolith/internal/coupon/api/handlers"
	couponRepos "eshop-monolith/internal/coupon/domain/repositories"
	couponSvc "eshop-monolith/internal/coupon/service"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterCouponRoutes 注册优惠券相关路由
func RegisterCouponRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, bus *eventbus.Bus) *couponSvc.CouponService {
	// 初始化仓储
	couponRepo := couponRepos.NewCouponRepository(db)
	userCouponRepo := couponRepos.NewUserCouponRepository(db)

	// 初始化服务
	couponService := couponSvc.NewCouponService(db, couponRepo, userCouponRepo, bus)

	// 初始化 Handler
	couponHandler := handlers.NewCouponHandler(couponService)

	// 公共路由（无需登录）
	coupon := v1.Group("/coupons")
	{
		coupon.GET("", couponHandler.ListCoupons)
		coupon.GET("/:id", couponHandler.GetCoupon)
	}

	// 需要认证的路由
	auth := v1.Group("/coupons")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/claim", couponHandler.ClaimCoupon)
		auth.GET("/mine", couponHandler.GetUserCoupons)
		auth.GET("/usable", couponHandler.GetUsableCoupons)
		auth.POST("/use", couponHandler.UseCoupon)
		auth.POST("/validate", couponHandler.ValidateCoupon)
	}

	// 管理员路由（优惠券管理）
	admin := v1.Group("/admin/coupons")
	admin.Use(middleware.JWTAuth())
	{
		admin.POST("", couponHandler.CreateCoupon)
		admin.PUT("/:id", couponHandler.UpdateCoupon)
		admin.GET("", couponHandler.ListCoupons)
		admin.GET("/:id", couponHandler.GetCoupon)
	}

	return couponService
}

// RegisterPromotionRoutes 注册促销活动相关路由
func RegisterPromotionRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, bus *eventbus.Bus) *couponSvc.PromotionService {
	// 初始化仓储
	promotionRepo := couponRepos.NewPromotionRepository(db)
	promotionProdRepo := couponRepos.NewPromotionProductRepository(db)

	// 初始化服务
	promotionService := couponSvc.NewPromotionService(db, promotionRepo, promotionProdRepo, bus)

	// 初始化 Handler
	promotionHandler := handlers.NewPromotionHandler(promotionService)

	// 公共路由（无需登录）
	promo := v1.Group("/promotions")
	{
		promo.GET("/active", promotionHandler.GetActivePromotions)
		promo.GET("", promotionHandler.ListPromotions)
		promo.GET("/:id", promotionHandler.GetPromotion)
		promo.GET("/:id/products", promotionHandler.GetPromotionProducts)
	}

	// 管理员路由（促销管理）
	admin := v1.Group("/admin/promotions")
	admin.Use(middleware.JWTAuth())
	{
		admin.POST("", promotionHandler.CreatePromotion)
		admin.PUT("/:id", promotionHandler.UpdatePromotion)
		admin.PUT("/:id/status", promotionHandler.UpdatePromotionStatus)
		admin.POST("/:id/products", promotionHandler.LinkProducts)
		admin.GET("/:id/products", promotionHandler.GetPromotionProducts)
		admin.GET("", promotionHandler.ListPromotions)
		admin.GET("/:id", promotionHandler.GetPromotion)
	}

	return promotionService
}
