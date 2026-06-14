package routes

import (
	"log"

	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/order/api/handlers"
	orderRepos "eshop-monolith/internal/order/domain/repositories"
	"eshop-monolith/internal/order/service"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, bus *eventbus.Bus) *service.OrderService {
	invForOrder, ok := repos.Inventory.(orderRepos.InventoryForOrder)
	if !ok {
		log.Fatal("Inventory repository does not implement InventoryForOrder interface")
	}
	orderService := service.NewOrderService(db, repos.Order, invForOrder, bus)
	orderHandler := handlers.NewOrderHandler(orderService)

	orders := v1.Group("/orders")
	{
		orders.GET("", orderHandler.ListOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.POST("", orderHandler.CreateOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
		orders.GET("/items", orderHandler.ListAllOrderItems)
		orders.GET("/:id/items", orderHandler.GetOrderItems)
		orders.POST("/:id/cancel", orderHandler.CancelOrder)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	}

	// 根据用户ID获取订单列表
	users := v1.Group("/users")
	{
		users.GET("/:user_id/orders", orderHandler.GetOrdersByUserID)
	}

	return orderService
}
