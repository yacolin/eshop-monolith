package routes

import (
	"eshop-monolith/internal/api/handlers"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
)

func registerInventoryRoutes(v1 *gin.RouterGroup, repos *repository.Repositories) {
	inventoryService := service.NewInventoryService(repos.Inventory, nil)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	inventories := v1.Group("/inventories")
	{
		inventories.GET("", inventoryHandler.ListInventories)
		inventories.GET("/product/:productId", inventoryHandler.GetInventoryByProductID)
	}
}

