package routes

import (
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProductRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, bus *eventbus.Bus) *service.ProductService {
	productRepo := repositories.NewProductRepository(db)
	productService := service.NewProductService(productRepo, repos.Inventory, repos.Sku, bus, db, repos.Redis)
	productHandler := handlers.NewProductHandler(productService)

	products := v1.Group("/products")
	{
		products.GET("", productHandler.ListProducts)
		products.GET("/cursor", productHandler.ListProductsByCursor)
		products.GET("/cache", productHandler.ListCachedProducts)
		products.GET("/cache-cursor", productHandler.ListCachedProductsByCursor)
		products.GET("/cache/:id", productHandler.GetCachedProduct)
		products.POST("/cache/warmup", productHandler.WarmupCache)
		products.GET("/:id", productHandler.GetProduct)
		products.GET("/:id/detail", productHandler.GetProductWithSkus)
		products.GET("/:id/enriched", productHandler.GetProductWithCategory)
		products.GET("/enriched", productHandler.ListProductsWithCategory)
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

	return productService
}
