package notification

import (
	"time"

	"gorm.io/gorm"
)

const (
	NotifTypeSystem  = "system"
	NotifTypeOrder   = "order"
	NotifTypePayment = "payment"
	NotifTypeFlash   = "flash"
	NotifTypeAdmin   = "admin"
)

type Notification struct {
	ID        int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64            `gorm:"not null;index" json:"user_id"`
	Title     string           `gorm:"type:varchar(200);not null" json:"title"`
	Content   string           `gorm:"type:text;not null" json:"content"`
	Type      string           `gorm:"type:varchar(20);not null;index" json:"type"`
	IsRead    bool             `gorm:"not null;default:false;index" json:"is_read"`
	ReadAt    *time.Time       `gorm:"type:datetime(3)" json:"read_at,omitempty"`
	CreatedAt time.Time        `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt time.Time        `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"type:datetime(3);index" json:"-"`
}

func (Notification) TableName() string { return "bse_notifications" }
