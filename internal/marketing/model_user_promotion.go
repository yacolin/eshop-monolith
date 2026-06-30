package marketing

import (
	"time"

	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type UserPromotion struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64          `gorm:"not null;uniqueIndex:uk_user_promo;index:idx_user" json:"user_id"`
	PromotionID int64          `gorm:"not null;uniqueIndex:uk_user_promo;index:idx_promotion" json:"promotion_id"`
	AcquireTime time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"acquire_time"`
	ExpireTime  *time.Time     `gorm:"type:datetime" json:"expire_time"`
	Status      int8           `gorm:"not null;default:1;index:idx_status_expire" json:"status"`
	UsedTime    *time.Time     `gorm:"type:datetime" json:"used_time"`
	OrderID     *int64         `gorm:"index:idx_order" json:"order_id"`
	QueueToken  string         `gorm:"type:varchar(64);default:''" json:"queue_token"`
	CreatedBy   *int64         `json:"created_by"`
	UpdatedBy   *int64         `json:"updated_by"`
	CreatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserPromotion) TableName() string { return "mkt_user_promotions" }
