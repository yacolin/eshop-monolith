package eventbus

import (
	"eshop-monolith/internal/flashsale/events"
	"eshop-monolith/pkg/logger"
)

// RegisterFlashHandlers 注册闪购事件处理器
func RegisterFlashHandlers(bus *Bus) {
	bus.Subscribe("flashsale.FlashOrderCreatedEvent", handleFlashOrderCreated)
	bus.Subscribe("flashsale.FlashOrderPaidEvent", handleFlashOrderPaid)
	bus.Subscribe("flashsale.FlashOrderCancelledEvent", handleFlashOrderCancelled)
}

// handleFlashOrderCreated 处理闪购订单创建事件
func handleFlashOrderCreated(event interface{}) {
	e, ok := event.(events.FlashOrderCreatedEvent)
	if !ok {
		return
	}
	logger.Info("Flash order created",
		"order_id", e.OrderID,
		"user_id", e.UserID,
		"activity_id", e.ActivityID,
		"product_id", e.ProductID,
		"amount", e.Amount,
	)
	// TODO: 发送通知、记录风控日志等
}

// handleFlashOrderPaid 处理闪购订单支付成功事件
func handleFlashOrderPaid(event interface{}) {
	e, ok := event.(events.FlashOrderPaidEvent)
	if !ok {
		return
	}
	logger.Info("Flash order paid",
		"order_id", e.OrderID,
		"user_id", e.UserID,
		"activity_id", e.ActivityID,
		"product_id", e.ProductID,
		"amount", e.Amount,
	)
	// TODO: 发送发货通知等
}

// handleFlashOrderCancelled 处理闪购订单取消事件
func handleFlashOrderCancelled(event interface{}) {
	e, ok := event.(events.FlashOrderCancelledEvent)
	if !ok {
		return
	}
	logger.Info("Flash order cancelled",
		"order_id", e.OrderID,
		"user_id", e.UserID,
		"activity_id", e.ActivityID,
	)
	// TODO: 退款处理等
}
