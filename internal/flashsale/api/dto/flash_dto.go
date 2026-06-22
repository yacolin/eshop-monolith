package dto

import (
	"time"

	"eshop-monolith/internal/flashsale/domain/models"
)

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

// ActivityCursorQuery 游标分页查询参数（深分页优化）
type ActivityCursorQuery struct {
	Cursor int64  `form:"cursor"`          // 游标（上一页最后一条的 ID，首次查询传 0）
	Size   int    `form:"size"`                 // 每页条数，默认 20
	Status string `form:"status"`          // 筛选状态：pending/active/finished
}

func (q *ActivityCursorQuery) Normalize() {
	if q.Size <= 0 {
		q.Size = 20
	}
	if q.Size > 100 {
		q.Size = 100
	}
	if q.Cursor < 0 {
		q.Cursor = 0
	}
}

// ActivityCursorResult 游标分页结果
type ActivityCursorResult struct {
	List       []models.FlashActivity `json:"list"`
	NextCursor int64                  `json:"next_cursor"`
	HasMore    bool                   `json:"has_more"`
}