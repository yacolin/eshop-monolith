package routes

import (
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories) {
	inventoryService := service.NewInventoryService(repos.Inventory, nil)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	inventories := v1.Group("/inventories")
	{
		inventories.GET("", inventoryHandler.ListInventories)
		inventories.GET("/product/:productId", inventoryHandler.GetInventoryByProductID)
	}
}
