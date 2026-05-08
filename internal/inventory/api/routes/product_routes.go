package routes

import (
	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProductRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, bus *eventbus.Bus) {
	productRepo := repositories.NewProductRepository(db)
	productService := service.NewProductService(productRepo, bus, db)
	productHandler := handlers.NewProductHandler(productService)

	products := v1.Group("/products")
	{
		products.GET("", productHandler.ListProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.GET("/category/:category_id", productHandler.ListProductsByCategory)
	}

	// 需要认证的产品写操作
	auth := v1.Group("/products")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", productHandler.CreateProduct)
		auth.PUT("/:id", productHandler.UpdateProduct)
		auth.DELETE("/:id", productHandler.DeleteProduct)
	}
}
