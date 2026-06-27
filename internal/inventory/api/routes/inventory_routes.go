package routes

import (
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, rabbit *rabbitmq.Client) {
	inventoryService := service.NewInventoryService(repos.Inventory, repos.Sku, repos.Product, rabbit)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	inventories := v1.Group("/inventories")
	{
		inventories.GET("", inventoryHandler.ListInventories)
		inventories.GET("/enriched", inventoryHandler.ListInventoriesEnriched)
		inventories.GET("/sku/:skuId", inventoryHandler.GetInventoryBySkuID)
	}

	// 需要认证的库存写操作
	auth := v1.Group("/inventories")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", inventoryHandler.CreateInventory)
		auth.POST("/batch", inventoryHandler.BatchCreateInventory)
		auth.PUT("/:id", inventoryHandler.UpdateInventory)
		auth.POST("/reserve", inventoryHandler.ReserveInventory)
		auth.POST("/release", inventoryHandler.ReleaseInventory)
	}
}
