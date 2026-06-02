package dto

import "time"

type CreateActivityReq struct {
	ProductID  int64 `json:"product_id" binding:"required"`
	FlashPrice int64 `json:"flash_price" binding:"required,min=1"`
	TotalStock int   `json:"total_stock" binding:"required,min=1"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
}

type FlashBuyReq struct {
	ActivityID int64 `json:"activity_id" binding:"required"`
	UserID     int64 `json:"user_id" binding:"required"`
}

type FlashBuyResp struct {
	Success    bool   `json:"success"`
	OrderID    int64  `json:"order_id,omitempty"`
	Message    string `json:"message"`
}

type ActivityResp struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	FlashPrice int64     `json:"flash_price"`
	TotalStock int       `json:"total_stock"`
	SoldStock  int       `json:"sold_stock"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
}