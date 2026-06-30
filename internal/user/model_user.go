package user

import (
	"time"

	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type User struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string         `gorm:"type:varchar(50);uniqueIndex:uk_username;not null;default:''" json:"username"`
	PasswordHash  string         `gorm:"type:varchar(255);not null;default:''" json:"-"`
	Email         string         `gorm:"type:varchar(100);uniqueIndex:uk_email;not null;default:''" json:"email"`
	EmailVerified bool           `gorm:"not null;default:false" json:"email_verified"`
	Phone         string         `gorm:"type:varchar(20);uniqueIndex:uk_phone;not null;default:''" json:"phone"`
	PhoneVerified bool           `gorm:"not null;default:false" json:"phone_verified"`
	Avatar        string         `gorm:"type:varchar(512);not null;default:''" json:"avatar"`
	Nickname      string         `gorm:"type:varchar(50);not null;default:''" json:"nickname"`
	Status        int8           `gorm:"type:tinyint;not null;default:1;index:idx_status" json:"status"`
	RegisterIP    string         `gorm:"type:varchar(50);not null;default:''" json:"register_ip"`
	RegisterSource string        `gorm:"type:varchar(20);not null;default:''" json:"register_source"`
	LastLoginIP   string         `gorm:"type:varchar(50);not null;default:''" json:"last_login_ip"`
	LastLoginAt   *time.Time     `gorm:"type:datetime(3)" json:"last_login_at"`
	CreatedAt utils.Timestamp      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (User) TableName() string { return "usr_users" }
