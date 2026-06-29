package trade

// ── 订单事件 ─────────────────────────────────────

type OrderPaidEvent struct {
	CustomerID  string
	OrderID     int64
	TotalAmount int64
}

type OrderShippedEvent struct {
	CustomerID string
	OrderID    int64
}

type OrderDeliveredEvent struct {
	CustomerID string
	OrderID    int64
}

type OrderCancelledEvent struct {
	CustomerID string
	OrderID    int64
	UserID     int64
}

// ── 支付事件 ─────────────────────────────────────

type PaymentSuccessEvent struct {
	OrderType string
	OrderID   int64
	Amount    int64
}

type PaymentFailedEvent struct{}

type RefundCreatedEvent struct{}

type RefundFailedEvent struct {
	RefundID      int64
	OrderID       int64
	FailureReason string
}
