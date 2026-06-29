package trade

// ── 订单事件 ─────────────────────────────────────

type OrderPaidEvent struct {
	CustomerID  string `json:"customer_id"`
	OrderID     int64  `json:"order_id"`
	TotalAmount int64  `json:"total_amount"`
}

func (OrderPaidEvent) RoutingKey() string { return "order.paid" }

type OrderShippedEvent struct {
	CustomerID string `json:"customer_id"`
	OrderID    int64  `json:"order_id"`
}

func (OrderShippedEvent) RoutingKey() string { return "order.shipped" }

type OrderDeliveredEvent struct {
	CustomerID string `json:"customer_id"`
	OrderID    int64  `json:"order_id"`
}

func (OrderDeliveredEvent) RoutingKey() string { return "order.delivered" }

type OrderCancelledEvent struct {
	CustomerID string `json:"customer_id"`
	OrderID    int64  `json:"order_id"`
	UserID     int64  `json:"user_id"`
}

func (OrderCancelledEvent) RoutingKey() string { return "order.cancelled" }

// ── 支付事件 ─────────────────────────────────────

type PaymentSuccessEvent struct {
	OrderType string `json:"order_type"`
	OrderID   int64  `json:"order_id"`
	Amount    int64  `json:"amount"`
}

func (PaymentSuccessEvent) RoutingKey() string { return "payment.success" }

type PaymentFailedEvent struct{}

func (PaymentFailedEvent) RoutingKey() string { return "payment.failed" }

type RefundCreatedEvent struct{}

func (RefundCreatedEvent) RoutingKey() string { return "payment.refund.created" }

type RefundFailedEvent struct {
	RefundID      int64  `json:"refund_id"`
	OrderID       int64  `json:"order_id"`
	FailureReason string `json:"failure_reason"`
}

func (RefundFailedEvent) RoutingKey() string { return "payment.refund.failed" }
