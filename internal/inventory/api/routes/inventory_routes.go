package routes

import (
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, bus *eventbus.Bus) {
	inventoryService := service.NewInventoryService(repos.Inventory, bus)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	inventories := v1.Group("/inventories")
	{
		inventories.GET("", inventoryHandler.ListInventories)
		inventories.GET("/sku/:skuId", inventoryHandler.GetInventoryBySkuID)
	}

	// 需要认证的库存写操作
	auth := v1.Group("/inventories")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", inventoryHandler.CreateInventory)
		auth.PUT("/:id", inventoryHandler.UpdateInventory)
		auth.POST("/reserve", inventoryHandler.ReserveInventory)
		auth.POST("/release", inventoryHandler.ReleaseInventory)
	}
}
