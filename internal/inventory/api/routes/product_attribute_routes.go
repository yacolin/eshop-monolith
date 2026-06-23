package routes

import (
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProductAttributeRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	attrRepo := repositories.NewProductAttributeRepository(db)
	svc := service.NewProductAttributeService(attrRepo, repos.Sku, repos.Product, db)
	h := handlers.NewProductAttributeHandler(svc)

	products := v1.Group("/products")
	{
		products.GET("/:id/attributes", h.GetProductAttributes)
		products.PUT("/:id/attributes", h.UpdateProductAttributes)
		products.POST("/:id/skus/batch", h.BatchCreateSkus)
	}
}
