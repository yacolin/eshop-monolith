package events

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
}
