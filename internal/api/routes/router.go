package routes

import (
	"eshop-monolith/internal/api/handlers"
	"eshop-monolith/internal/pkg/config"
	"eshop-monolith/internal/pkg/middleware"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
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
		categories.GET("/:parent_id/children", categoryHandler.ListSubCategories)
	}
}

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, repos *repository.Repositories) *gin.Engine {
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
