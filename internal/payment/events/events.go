package events

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
