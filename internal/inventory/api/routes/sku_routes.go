package routes

import (
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSkuRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, rabbit *rabbitmq.Client) {
	skuSvc := service.NewSkuService(repos.Sku, repos.Product, rabbit, db)
	skuH := handlers.NewSkuHandler(skuSvc)

	skus := v1.Group("/skus")
	{
		skus.POST("", skuH.CreateSku)
		skus.GET("", skuH.ListSkus)
		skus.GET("/:id", skuH.GetSku)
		skus.PUT("/:id", skuH.UpdateSku)
		skus.DELETE("/:id", skuH.DeleteSku)
	}
}
