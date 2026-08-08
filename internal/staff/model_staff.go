package staff

import (
	"time"

	"gorm.io/gorm"
)

// Staff B 端员工(sys_staff)
type Staff struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex:uk_username;not null" json:"username"`
	PasswordHash string         `gorm:"type:varchar(255);not null;default:''" json:"-"`
	RealName     string         `gorm:"type:varchar(50);not null;default:''" json:"real_name"`
	Email        string         `gorm:"type:varchar(100)" json:"email"`
	Phone        string         `gorm:"type:varchar(20)" json:"phone"`
	Avatar       string         `gorm:"type:varchar(512);not null;default:''" json:"avatar"`
	Status       int8           `gorm:"type:tinyint(1);not null;default:1;index:idx_status" json:"status"`
	LastLoginIP  string         `gorm:"type:varchar(50);not null;default:''" json:"last_login_ip"`
	LastLoginAt  *time.Time     `gorm:"type:datetime(3)" json:"last_login_at"`
	CreatedBy    int64          `gorm:"not null;default:0" json:"created_by"`
	CreatedAt    time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (Staff) TableName() string { return "sys_staff" }
