package models

type Payment struct {
	OrderID       int64
	OrderType     string
	Amount        int64
	Currency      string
	PaymentMethod string
	Status        string
	Metadata      string
}
