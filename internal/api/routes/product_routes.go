package routes

import (
	"eshop-monolith/internal/api/handlers"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerProductRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	productService := service.NewProductService(repos.Product, nil, db)
	productHandler := handlers.NewProductHandler(productService)

	products := v1.Group("/products")
	{
		products.GET("", productHandler.ListProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.GET("/category/:category_id", productHandler.ListProductsByCategory)
	}
}

