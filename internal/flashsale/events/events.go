package events

type FlashOrderCreatedEvent struct {
	UserID  int64
	OrderID int64
	Amount  int64
}

type FlashOrderPaidEvent struct {
	OrderID int64
	UserID  int64
}

type FlashOrderCancelledEvent struct {
	OrderID int64
	UserID  int64
}
