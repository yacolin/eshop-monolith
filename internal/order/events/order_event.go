package events

import "time"

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount int64     `json:"total_amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount int64     `json:"total_amount"`
	Currency    string    `json:"currency"`
	PaidAt      time.Time `json:"paid_at"`
}

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
	OrderID    int64     `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	ShippedAt  time.Time `json:"shipped_at"`
}

// OrderDeliveredEvent 订单送达事件
type OrderDeliveredEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	DeliveredAt time.Time `json:"delivered_at"`
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount int64     `json:"total_amount"`
	CancelledAt time.Time `json:"cancelled_at"`
}
