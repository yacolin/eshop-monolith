package marketing

import "eshop-monolith/pkg/utils"

type UsageLog struct {
	ID            int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	PromotionID   int64           `gorm:"not null;index:idx_promotion" json:"promotion_id"`
	UserID        int64           `gorm:"not null;index:idx_user" json:"user_id"`
	OrderID       *int64          `gorm:"index:idx_order" json:"order_id"`
	ActionType    int8            `gorm:"not null" json:"action_type"`
	FailReason    string          `gorm:"type:varchar(200);default:''" json:"fail_reason"`
	RequestParams string          `gorm:"type:json" json:"request_params"`
	CreatedAt     utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_created" json:"created_at"`
}

func (UsageLog) TableName() string { return "mkt_promotion_usage_logs" }
