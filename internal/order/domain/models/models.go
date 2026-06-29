package models

type Order struct {
	ID         int64
	CustomerID string
	OrderNo    string
	Items      []OrderItem
}

type OrderItem struct {
	ID      int64
	OrderID int64
	ProductID string
}
