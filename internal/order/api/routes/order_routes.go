package routes

import (
	"eshop-monolith/internal/order/api/handlers"
	"eshop-monolith/internal/order/service"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(v1 *gin.RouterGroup, repos *repository.Repositories) {
	orderService := service.NewOrderService(repos.Order, repos.Inventory, nil)
	orderHandler := handlers.NewOrderHandler(orderService)

	orders := v1.Group("/orders")
	{
		orders.GET("", orderHandler.ListOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.POST("", orderHandler.CreateOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
		orders.POST("/:id/cancel", orderHandler.CancelOrder)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	}

	// 根据用户ID获取订单列表
	users := v1.Group("/users")
	{
		users.GET("/:user_id/orders", orderHandler.GetOrdersByUserID)
	}
}
