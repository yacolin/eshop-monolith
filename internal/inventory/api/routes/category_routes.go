package routes

import (
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, rabbit *rabbitmq.Client) *service.CategoryService {
	categoryService := service.NewCategoryService(repos.Category, repos.CategoryAttribute, rabbit, repos.Redis)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	categories := v1.Group("/categories")
	{
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/root", categoryHandler.ListRootCategories)
		categories.GET("/non-root", categoryHandler.ListNonRootCategories)
		categories.GET("/cache", categoryHandler.ListCachedCategories)
		categories.GET("/cache/:id", categoryHandler.GetCachedCategory)
		categories.POST("/cache/warmup", categoryHandler.WarmupCategoryCache)
		categories.GET("/:id/children", categoryHandler.ListSubCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		// 品类-属性关联
		categories.GET("/:id/attributes", categoryHandler.GetCategoryAttributes)
	}

	// 需要认证的分类写操作
	auth := v1.Group("/categories")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", categoryHandler.CreateCategory)
		auth.PUT("/:id", categoryHandler.UpdateCategory)
		auth.DELETE("/:id", categoryHandler.DeleteCategory)
		// 品类-属性关联写操作
		auth.PUT("/:id/attributes", categoryHandler.SetCategoryAttributes)
	}

	return categoryService
}
