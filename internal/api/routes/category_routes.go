package routes

import (
	"eshop-monolith/internal/api/handlers"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
)

func registerCategoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories) {
	categoryService := service.NewCategoryService(repos.Category, nil)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	categories := v1.Group("/categories")
	{
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/root", categoryHandler.ListRootCategories)
		categories.GET("/:id/children", categoryHandler.ListSubCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}

