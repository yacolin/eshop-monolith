package eventbus

import (
	"eshop-monolith/internal/infra/ws"
)

// RegisterWSHandlers 注册 WebSocket 事件处理器
// 将业务事件通过 WebSocket 推送给在线用户
func RegisterWSHandlers(bus *Bus, hub *ws.Hub) {
	// 订单事件
	bus.Subscribe("order.OrderPaidEvent", hub.HandleEvent)
	bus.Subscribe("order.OrderShippedEvent", hub.HandleEvent)
	bus.Subscribe("order.OrderDeliveredEvent", hub.HandleEvent)
	bus.Subscribe("order.OrderCancelledEvent", hub.HandleEvent)

	// 支付事件
	bus.Subscribe("payment.PaymentSuccessEvent", hub.HandleEvent)
	bus.Subscribe("payment.PaymentFailedEvent", hub.HandleEvent)
	bus.Subscribe("payment.RefundCreatedEvent", hub.HandleEvent)

	// 闪购事件
	bus.Subscribe("flashsale.FlashOrderCreatedEvent", hub.HandleEvent)
	bus.Subscribe("flashsale.FlashOrderPaidEvent", hub.HandleEvent)
	bus.Subscribe("flashsale.FlashOrderCancelledEvent", hub.HandleEvent)

	// 库存预警
	bus.Subscribe("inventory.InventoryLowEvent", hub.HandleEvent)
}
