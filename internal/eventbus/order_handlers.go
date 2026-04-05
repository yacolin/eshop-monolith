package eventbus

import (
	"eshop-monolith/internal/order/events"
	"eshop-monolith/internal/pkg/logger"
)

// RegisterOrderHandlers 注册订单事件处理器
func RegisterOrderHandlers(bus *Bus) {
	bus.Subscribe("order.OrderCreatedEvent", handleOrderCreated)
	bus.Subscribe("order.OrderStatusChangedEvent", handleOrderStatusChanged)
	bus.Subscribe("order.OrderCancelledEvent", handleOrderCancelled)
}

// handleOrderCreated 处理订单创建事件
func handleOrderCreated(event interface{}) {
	e, ok := event.(events.OrderCreatedEvent)
	if !ok {
		return
	}
	logger.Info("Order created", "order_id", e.OrderID, "customer_id", e.CustomerID, "total_amount", e.TotalAmount)
	// 这里可以添加发送通知、记录日志等逻辑
}

// handleOrderStatusChanged 处理订单状态变更事件
func handleOrderStatusChanged(event interface{}) {
	e, ok := event.(events.OrderPaidEvent)
	if !ok {
		return
	}
	logger.Info("Order status changed", "order_id", e.OrderID, "customer_id", e.CustomerID, "total_amount", e.TotalAmount)
	// 这里可以添加状态变更通知等逻辑
}

// handleOrderCancelled 处理订单取消事件
func handleOrderCancelled(event interface{}) {
	e, ok := event.(events.OrderCancelledEvent)
	if !ok {
		return
	}
	logger.Info("Order cancelled", "order_id", e.OrderID, "customer_id", e.CustomerID)
	// 这里可以添加退款处理等逻辑
}
