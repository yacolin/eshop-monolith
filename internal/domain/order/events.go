package order

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
	TotalAmount int64  `json:"total_amount"`
}

// OrderStatusChangedEvent 订单状态变更事件
type OrderStatusChangedEvent struct {
	OrderID    string `json:"order_id"`
	OldStatus  string `json:"old_status"`
	NewStatus  string `json:"new_status"`
	UserID     string `json:"user_id"`
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
	TotalAmount int64  `json:"total_amount"`
}