package ws

import (
	"encoding/json"
	"time"

	inventoryEvents "eshop-monolith/internal/inventory/events"
	orderEvents "eshop-monolith/internal/order/events"
	paymentEvents "eshop-monolith/internal/payment/events"
	flashEvents "eshop-monolith/internal/flashsale/events"
)

// PushMessage WebSocket 推送消息
type PushMessage struct {
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// NewPushMessage 从事件创建推送消息
func NewPushMessage(event interface{}) *PushMessage {
	now := time.Now().UnixMilli()

	switch e := event.(type) {
	case orderEvents.OrderPaidEvent:
		return &PushMessage{Type: "order.paid", Timestamp: now, Data: e}
	case orderEvents.OrderShippedEvent:
		return &PushMessage{Type: "order.shipped", Timestamp: now, Data: e}
	case orderEvents.OrderDeliveredEvent:
		return &PushMessage{Type: "order.delivered", Timestamp: now, Data: e}
	case orderEvents.OrderCancelledEvent:
		return &PushMessage{Type: "order.cancelled", Timestamp: now, Data: e}

	case paymentEvents.PaymentSuccessEvent:
		return &PushMessage{Type: "payment.success", Timestamp: now, Data: e}
	case paymentEvents.PaymentFailedEvent:
		return &PushMessage{Type: "payment.failed", Timestamp: now, Data: e}
	case paymentEvents.RefundCreatedEvent:
		return &PushMessage{Type: "refund.created", Timestamp: now, Data: e}

	case flashEvents.FlashOrderCreatedEvent:
		return &PushMessage{Type: "flash.created", Timestamp: now, Data: e}
	case flashEvents.FlashOrderPaidEvent:
		return &PushMessage{Type: "flash.paid", Timestamp: now, Data: e}
	case flashEvents.FlashOrderCancelledEvent:
		return &PushMessage{Type: "flash.cancelled", Timestamp: now, Data: e}

	case inventoryEvents.InventoryLowEvent:
		return &PushMessage{Type: "inventory.low", Timestamp: now, Data: e}
	}

	return nil
}

// Marshal 序列化为 JSON
func (m *PushMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// extractUserID 从事件中提取 user_id
func extractUserID(event interface{}) int64 {
	switch e := event.(type) {
	case orderEvents.OrderPaidEvent:
		return parseCustomerID(e.CustomerID)
	case orderEvents.OrderShippedEvent:
		return parseCustomerID(e.CustomerID)
	case orderEvents.OrderDeliveredEvent:
		return parseCustomerID(e.CustomerID)
	case orderEvents.OrderCancelledEvent:
		return parseCustomerID(e.CustomerID)

	case flashEvents.FlashOrderCreatedEvent:
		return e.UserID
	case flashEvents.FlashOrderPaidEvent:
		return e.UserID
	case flashEvents.FlashOrderCancelledEvent:
		return e.UserID

	case paymentEvents.PaymentSuccessEvent:
		// PaymentSuccessEvent 没有用户 ID，通过 OrderType 区分
		// 实际推送由 OrderPaidEvent/FailedEvent 完成
		return 0
	case paymentEvents.PaymentFailedEvent:
		return 0
	case paymentEvents.RefundCreatedEvent:
		return 0
	}

	return 0
}

// parseCustomerID 将 CustomerID (string) 解析为 int64
func parseCustomerID(customerID string) int64 {
	if customerID == "" {
		return 0
	}
	var id int64
	for _, c := range customerID {
		if c < '0' || c > '9' {
			return 0
		}
		id = id*10 + int64(c-'0')
		if id < 0 {
			return 0
		}
	}
	return id
}
