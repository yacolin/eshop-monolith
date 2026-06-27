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

func (e OrderCreatedEvent) RoutingKey() string { return "order.created" }

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount int64     `json:"total_amount"`
	Currency    string    `json:"currency"`
	PaidAt      time.Time `json:"paid_at"`
}

func (e OrderPaidEvent) RoutingKey() string { return "order.paid" }

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
	OrderID    int64     `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	ShippedAt  time.Time `json:"shipped_at"`
}

func (e OrderShippedEvent) RoutingKey() string { return "order.shipped" }

// OrderDeliveredEvent 订单送达事件
type OrderDeliveredEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	DeliveredAt time.Time `json:"delivered_at"`
}

func (e OrderDeliveredEvent) RoutingKey() string { return "order.delivered" }

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	OrderID     int64     `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount int64     `json:"total_amount"`
	CancelledAt time.Time `json:"cancelled_at"`
}

func (e OrderCancelledEvent) RoutingKey() string { return "order.cancelled" }
