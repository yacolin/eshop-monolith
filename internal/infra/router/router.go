package router

import (
	"context"
	"sync/atomic"

	"eshop-monolith/internal/base"
	"eshop-monolith/internal/dashboard"
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/marketing"
	"eshop-monolith/internal/product"
	"eshop-monolith/internal/review"
	"eshop-monolith/internal/trade"
	"eshop-monolith/internal/user"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/rabbitmq/consumers"
	"eshop-monolith/internal/infra/repository"
	ws "eshop-monolith/internal/infra/ws"
	"eshop-monolith/pkg/config"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"

	_ "eshop-monolith/docs" // swagger docs, 由 swag CLI 生成

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, repos *repository.Repositories, db *gorm.DB, mqClient *rabbitmq.Client) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 创建 WebSocket Hub 并启动（传入Redis客户端支持断线重连和增量同步）
	wsHub := ws.NewHub(repos.Redis)
	go wsHub.Run()

	// 设置用户信息查询回调（用于实时在线事件广播）
	wsHub.SetUserInfoProvider(func(uid int64) (string, string, error) {
		ctx := context.Background()

		user, err := repos.User.FindByID(ctx, uid)
		if err != nil {
			return "", "", err
		}

		return user.Username, user.Nickname, nil
	})

	// 添加全局错误处理中间件
	router.Use(middleware.ErrorHandler())

	// 声明 service 变量（在 v1 block 内赋值, 在 block 外用于事件处理器注册和启动预热）
	var dashboardSvc *dashboard.DashboardService
	var notifSvc *base.NotificationService
	var warmupDone atomic.Bool
	warmupDone.Store(true)

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		if !warmupDone.Load() {
			c.JSON(503, gin.H{"status": "warming_up"})
			return
		}
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// Prometheus 监控指标
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 声明 service 变量（在 v1 block 内赋值, 在 block 外用于事件处理器注册和启动预热）
	v1 := router.Group("/api/v1")
	{

		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{
				"status":  "ok",
				"version": "v1",
			})
		})

		product.RegisterBrandRoutes(v1, db, repos.Redis)
		product.RegisterCategoryRoutes(v1, db, repos.Redis)
		product.RegisterAttributeRoutes(v1, db)
		product.RegisterProductRoutes(v1, db, repos.Redis)
		product.RegisterCategoryBrandRoutes(v1, db)
		product.RegisterSKURoutes(v1, db)
		marketing.RegisterPromotionRoutes(v1, db, repos.Redis)
		trade.RegisterTradeRoutes(v1, db, mqClient, func() { dashboardSvc.InvalidateCache(context.Background()) })
		inventory.RegisterInventoryRoutes(v1, db)
		// categorySvc = invRoutes.RegisterCategoryRoutes(v1, repos, mqClient)
		// productSvc = invRoutes.RegisterProductRoutes(v1, repos, db, mqClient)
		// invRoutes.RegisterInventoryRoutes(v1, repos, mqClient)
		// invRoutes.RegisterSkuRoutes(v1, repos, db, mqClient)
		// invRoutes.RegisterProductAttributeRoutes(v1, repos, db)
		// invRoutes.RegisterAttributeRoutes(v1, repos, db)

		// 优惠券系统（需先于订单初始化，用于结算时优惠校验）
		// [DEPRECATED] couponRoutes.RegisterCouponRoutes(v1, repos, db, mqClient)
		// [DEPRECATED] couponRoutes.RegisterPromotionRoutes(v1, repos, db, mqClient)

		// [DEPRECATED] orderSvc = orderRoutes.RegisterOrderRoutes(v1, repos, db, mqClient, couponSvc)
		tokenService := user.NewTokenService("your-secret-key-change-in-production", repos.Role)
		user.RegisterUserRoutes(v1, db, repos.User, repos.UserInfo, repos.Role)
		user.RegisterAuthRoutes(v1, db, tokenService, repos.Role, repos.User, repos.UserInfo, repos.LoginHistory)
		user.RegisterPermissionRoutes(v1, db, repos.Permission, repos.Role)
		user.RegisterRoleRoutes(v1, db, repos.Role)
		notifSvc = base.RegisterNotificationRoutes(v1, repos, db)
		review.RegisterReviewRoutes(v1, repos, db)
		user.RegisterAddressRoutes(v1, db)
		dashboardSvc = dashboard.RegisterDashboardRoutes(v1, repos, db, mqClient)

		// WebSocket 路由
		ws.RegisterWSRoutes(v1, wsHub)

		// 需要认证的路由组
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// 这里将添加需要认证的路由
			// 例如：订单、用户管理等路由
		}
	}

	// [DEPRECATED] 旧缓存预热与新商品中心不兼容，暂时禁用
	// if productSvc != nil {
	// 	go func() {
	// 		logger.Info("Starting product cache warmup...")
	// 		total, err := productSvc.WarmupProductCache(context.Background())
	// 		if err != nil {
	// 			logger.Error("Product cache warmup failed", "error", err)
	// 		} else {
	// 			logger.Info("Product cache warmup completed", "total", total)
	// 		}
	// 		warmupDone.Store(true)
	// 	}()
	// }
	// // 启动时异步预热分类缓存
	// if categorySvc != nil {
	// 	go func() {
	// 		logger.Info("Starting category cache warmup...")
	// 		total, err := categorySvc.WarmupCategoryCache(context.Background())
	// 		if err != nil {
	// 			logger.Error("Category cache warmup failed", "error", err)
	// 		} else {
	// 			logger.Info("Category cache warmup completed", "total", total)
	// 		}
	// 	}()
	// }

	// Start RabbitMQ consumers（WS 独立，业务合并）
	go func() {
		if err := consumers.StartWSConsumer(context.Background(), mqClient, wsHub); err != nil {
			logger.Error("启动 WS 消费者失败", "error", err)
		}
	}()
	go func() {
		if err := consumers.StartBusinessConsumer(context.Background(), mqClient, notifSvc); err != nil {
			logger.Error("启动业务消费者失败", "error", err)
		}
	}()

	return router
}
