package models

import "eshop-monolith/pkg/utils"

type FlashOrderStatus string

const (
	FlashOrderStatusPending   FlashOrderStatus = "pending"
	FlashOrderStatusPaid      FlashOrderStatus = "paid"
	FlashOrderStatusCancelled FlashOrderStatus = "cancelled"
)

type FlashOrder struct {
	ID         int64       `json:"id"`
	ActivityID int64       `json:"activity_id"`
	UserID     int64       `json:"user_id"`
	ProductID  int64       `json:"product_id"`
	Quantity   int         `json:"quantity"`
	FlashPrice int64       `json:"flash_price"`
	TotalAmount int64      `json:"total_amount"`
	Status     string      `json:"status"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (FlashOrder) TableName() string {
	return "flash_orders"
}