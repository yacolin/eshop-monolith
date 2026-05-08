package routes

import (
	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, bus *eventbus.Bus) {
	inventoryService := service.NewInventoryService(repos.Inventory, bus)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	inventories := v1.Group("/inventories")
	{
		inventories.GET("", inventoryHandler.ListInventories)
		inventories.GET("/product/:productId", inventoryHandler.GetInventoryByProductID)
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
