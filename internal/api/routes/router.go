package routes

import (
	cartRoutes "eshop-monolith/internal/cart/api/routes"
	invRoutes "eshop-monolith/internal/inventory/api/routes"
	orderRoutes "eshop-monolith/internal/order/api/routes"
	payRoutes "eshop-monolith/internal/payment/api/routes"
	userRoutes "eshop-monolith/internal/user/api/routes"

	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/pkg/config"
	"eshop-monolith/internal/pkg/middleware"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, repos *repository.Repositories, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 创建事件总线实例
	bus := eventbus.NewBus()

	// 添加全局错误处理中间件
	router.Use(middleware.ErrorHandler())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{
				"status":  "ok",
				"version": "v1",
			})
		})

		// 公开路由（按领域拆分注册）
		invRoutes.RegisterCategoryRoutes(v1, repos)
		invRoutes.RegisterProductRoutes(v1, repos, db, bus)
		invRoutes.RegisterInventoryRoutes(v1, repos)

		orderRoutes.RegisterOrderRoutes(v1, repos)
		userRoutes.RegisterUserRoutes(v1, repos)
		payRoutes.RegisterPaymentRoutes(v1, repos, bus, db)
		cartRoutes.RegisterCartRoutes(v1, repos)
		userRoutes.RegisterAuthRoutes(v1, repos, db)
		userRoutes.RegisterPermissionRoutes(v1, repos, db)
		userRoutes.RegisterRoleRoutes(v1, repos, db)

		// 需要认证的路由组
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// 这里将添加需要认证的路由
			// 例如：订单、用户管理等路由
		}
	}

	return router
}
