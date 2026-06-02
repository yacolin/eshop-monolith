package router

import (
	"context"

	cartRoutes "eshop-monolith/internal/cart/api/routes"
	invRoutes "eshop-monolith/internal/inventory/api/routes"
	orderRoutes "eshop-monolith/internal/order/api/routes"
	orderSvcPkg "eshop-monolith/internal/order/service"
	payRoutes "eshop-monolith/internal/payment/api/routes"
	userRoutes "eshop-monolith/internal/user/api/routes"
	flashRoutes "eshop-monolith/internal/flashsale/api/routes"
	flashSvcPkg "eshop-monolith/internal/flashsale/service"

	"eshop-monolith/internal/infra/eventbus"
	paymentEvents "eshop-monolith/internal/payment/events"
	"eshop-monolith/pkg/config"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, repos *repository.Repositories, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 创建事件总线实例
	bus := eventbus.NewBus()

	// 注册基础事件处理器（日志记录等）
	eventbus.RegisterHandlers(bus)

	// 添加全局错误处理中间件
	router.Use(middleware.ErrorHandler())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// Prometheus 监控指标
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 声明 service 变量（在 v1 block 内赋值, 在 block 外用于事件处理器注册）
	var orderSvc *orderSvcPkg.OrderService
	var flashSvc *flashSvcPkg.FlashService

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
		invRoutes.RegisterCategoryRoutes(v1, repos, bus)
		invRoutes.RegisterProductRoutes(v1, repos, db, bus)
		invRoutes.RegisterInventoryRoutes(v1, repos, bus)

		orderSvc = orderRoutes.RegisterOrderRoutes(v1, repos, db, bus)
		userRoutes.RegisterUserRoutes(v1, repos)
		payRoutes.RegisterPaymentRoutes(v1, repos, bus, db)
		cartRoutes.RegisterCartRoutes(v1, repos)
		flashSvc = flashRoutes.RegisterFlashRoutes(v1, repos, db, bus)
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

	// 注册支付成功事件业务处理器（在 service 创建后, 通过闭包注入依赖）
	bus.Subscribe("payment.PaymentSuccessEvent", func(event interface{}) {
		e, ok := event.(paymentEvents.PaymentSuccessEvent)
		if !ok {
			return
		}
		switch e.OrderType {
		case "flash":
			if err := flashSvc.HandlePaidSuccess(context.Background(), e.OrderID); err != nil {
				logger.Error("flash order paid handler failed",
					"order_id", e.OrderID, "error", err)
			}
		default:
			if err := orderSvc.HandlePaidSuccess(context.Background(), e.OrderID); err != nil {
				logger.Error("order paid handler failed",
					"order_id", e.OrderID, "error", err)
			}
		}
	})

	return router
}
