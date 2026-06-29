package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	OrderID       int64
	OrderType     string
	Amount        int64
	Currency      string
	PaymentMethod string
	Status        string
	Metadata      string `gorm:"type:json"`
}
