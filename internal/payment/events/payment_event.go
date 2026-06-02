package events

// PaymentCreatedEvent 支付创建事件
type PaymentCreatedEvent struct {
	PaymentID      int64  `json:"payment_id"`
	OrderID        int64  `json:"order_id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	PaymentMethod  string `json:"payment_method"`
}

// PaymentStatusUpdatedEvent 支付状态更新事件
type PaymentStatusUpdatedEvent struct {
	PaymentID      int64  `json:"payment_id"`
	OrderID        int64  `json:"order_id"`
	Status         string `json:"status"`
	PreviousStatus string `json:"previous_status"`
	TransactionID  string `json:"transaction_id,omitempty"`
}

// PaymentFailedEvent 支付失败事件
type PaymentFailedEvent struct {
	PaymentID     int64  `json:"payment_id"`
	OrderID       int64  `json:"order_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	FailureReason string `json:"failure_reason"`
}

// PaymentSuccessEvent 支付成功事件（下游 handler 据此进行订单状态更新+库存扣减）
type PaymentSuccessEvent struct {
	PaymentID int64  `json:"payment_id"`
	OrderID   int64  `json:"order_id"`
	Amount    int64  `json:"amount"`
	OrderType string `json:"order_type"` // "order" 常规订单, "flash" 闪购订单
}

// RefundCreatedEvent 退款创建事件
type RefundCreatedEvent struct {
	RefundID     int64  `json:"refund_id"`
	PaymentID    int64  `json:"payment_id"`
	OrderID      int64  `json:"order_id"`
	RefundAmount int64  `json:"refund_amount"`
	RefundReason string `json:"refund_reason"`
}

// RefundStatusUpdatedEvent 退款状态更新事件
type RefundStatusUpdatedEvent struct {
	RefundID       int64  `json:"refund_id"`
	PaymentID      int64  `json:"payment_id"`
	OrderID        int64  `json:"order_id"`
	Status         string `json:"status"`
	PreviousStatus string `json:"previous_status"`
	TransactionID  string `json:"transaction_id,omitempty"`
}

// RefundFailedEvent 退款失败事件
type RefundFailedEvent struct {
	RefundID      int64  `json:"refund_id"`
	PaymentID     int64  `json:"payment_id"`
	OrderID       int64  `json:"order_id"`
	RefundAmount  int64  `json:"refund_amount"`
	FailureReason string `json:"failure_reason"`
}
