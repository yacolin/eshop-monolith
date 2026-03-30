package routes

import (
	"eshop-monolith/internal/api/handlers"
	"eshop-monolith/internal/pkg/config"
	"eshop-monolith/internal/pkg/middleware"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupCategoryRoutes 设置分类相关路由
func setupCategoryRoutes(router *gin.RouterGroup, repos *repository.Repositories) {
	// 初始化分类服务
	categoryService := service.NewCategoryService(repos.Category, nil)

	// 初始化分类处理器
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// 分类路由
	categories := router.Group("/categories")
	{
		// 列出所有分类
		categories.GET("", categoryHandler.ListCategories)

		// 列出根分类
		categories.GET("/root", categoryHandler.ListRootCategories)

		// 列出子分类
		categories.GET("/:id/children", categoryHandler.ListSubCategories)

		// 根据ID获取分类
		categories.GET("/:id", categoryHandler.GetCategoryByID)

		// 创建分类
		categories.POST("", categoryHandler.CreateCategory)

		// 更新分类
		categories.PUT("/:id", categoryHandler.UpdateCategory)

		// 删除分类
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}

// setupProductRoutes 设置产品相关路由
func setupProductRoutes(router *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	// 初始化产品服务
	productService := service.NewProductService(repos.Product, nil, db)

	// 初始化产品处理器
	productHandler := handlers.NewProductHandler(productService)

	// 产品路由
	products := router.Group("/products")
	{
		// 列出所有产品
		products.GET("", productHandler.ListProducts)

		// 根据ID获取产品
		products.GET("/:id", productHandler.GetProduct)

		// 根据分类获取产品
		products.GET("/category/:category_id", productHandler.ListProductsByCategory)
	}
}

// setupInventoryRoutes 设置库存相关路由
func setupInventoryRoutes(router *gin.RouterGroup, repos *repository.Repositories) {
	// 初始化库存服务
	inventoryService := service.NewInventoryService(repos.Inventory, nil)

	// 初始化库存处理器
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// 库存路由
	inventories := router.Group("/inventories")
	{
		// 列出所有库存
		inventories.GET("", inventoryHandler.ListInventories)
	}
}

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, repos *repository.Repositories, db *gorm.DB) *gin.Engine {
	router := gin.Default()

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

		// 公开路由
		// 分类路由（公开）
		setupCategoryRoutes(v1, repos)
		// 产品路由（公开）
		setupProductRoutes(v1, repos, db)

		// 库存路由（公开）
		setupInventoryRoutes(v1, repos)

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
