package models

import "eshop-monolith/pkg/utils"

type FlashStatus string

const (
	FlashStatusPending  FlashStatus = "pending"
	FlashStatusActive   FlashStatus = "active"
	FlashStatusFinished FlashStatus = "finished"
)

type FlashActivity struct {
	ID         int64       `json:"id"`
	ProductID  int64       `json:"product_id"`
	FlashPrice int64       `json:"flash_price"`
	TotalStock int         `json:"total_stock"`
	SoldStock  int         `json:"sold_stock"`
	StartTime  utils.Timestamp `json:"start_time"`
	EndTime    utils.Timestamp `json:"end_time"`
	Status     string      `json:"status"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (FlashActivity) TableName() string {
	return "flash_activities"
}