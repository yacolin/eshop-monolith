package events

import "time"

// FlashOrderCreatedEvent 闪购订单创建事件
type FlashOrderCreatedEvent struct {
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	ActivityID int64     `json:"activity_id"`
	ProductID  int64     `json:"product_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

func (e FlashOrderCreatedEvent) RoutingKey() string { return "flash-order.created" }

// FlashOrderPaidEvent 闪购订单支付成功事件
type FlashOrderPaidEvent struct {
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	ActivityID int64     `json:"activity_id"`
	ProductID  int64     `json:"product_id"`
	Amount     int64     `json:"amount"`
	PaidAt     time.Time `json:"paid_at"`
}

func (e FlashOrderPaidEvent) RoutingKey() string { return "flash-order.paid" }

// FlashOrderCancelledEvent 闪购订单取消事件
type FlashOrderCancelledEvent struct {
	OrderID     int64     `json:"order_id"`
	UserID      int64     `json:"user_id"`
	ActivityID  int64     `json:"activity_id"`
	CancelledAt time.Time `json:"cancelled_at"`
}

func (e FlashOrderCancelledEvent) RoutingKey() string { return "flash-order.cancelled" }
