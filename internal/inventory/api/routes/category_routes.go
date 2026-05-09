package routes

import (
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, bus *eventbus.Bus) {
	categoryService := service.NewCategoryService(repos.Category, bus)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	categories := v1.Group("/categories")
	{
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/root", categoryHandler.ListRootCategories)
		categories.GET("/:id/children", categoryHandler.ListSubCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
	}

	// 需要认证的分类写操作
	auth := v1.Group("/categories")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", categoryHandler.CreateCategory)
		auth.PUT("/:id", categoryHandler.UpdateCategory)
		auth.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
