package user

import (
	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type LoginHistory struct {
	ID         int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64           `gorm:"not null;index:idx_user_id" json:"user_id"`
	Provider   string          `gorm:"type:varchar(20);not null;default:''" json:"provider"`
	IP         string          `gorm:"type:varchar(50);not null;default:''" json:"ip"`
	UserAgent  string          `gorm:"type:varchar(500);not null;default:''" json:"user_agent"`
	DeviceID   string          `gorm:"type:varchar(100);not null;default:''" json:"device_id"`
	Event      string          `gorm:"type:varchar(20);not null" json:"event"`
	Status     string          `gorm:"type:varchar(20);not null" json:"status"`
	FailReason string          `gorm:"type:varchar(255);not null;default:''" json:"fail_reason"`
	CreatedAt  utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);index:idx_created_at" json:"created_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (LoginHistory) TableName() string { return "usr_login_histories" }

// LoginEvent 常量
const (
	LoginEventLogin    = "login"
	LoginEventLogout   = "logout"
	LoginEventRefresh  = "refresh"
	LoginEventRevoke   = "revoke"
	LoginEventPassword = "password"
)

// LoginStatus 常量
const (
	LoginStatusSuccess = "success"
	LoginStatusFailed  = "failed"
)
