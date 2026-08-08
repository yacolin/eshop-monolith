package staff

import (
	"time"
)

// StaffLoginHistory B 端员工登录历史(sys_login_histories,无软删列)
type StaffLoginHistory struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StaffID       int64     `gorm:"not null;index:idx_staff_id" json:"staff_id"`
	LoginIP       string    `gorm:"type:varchar(50);not null;default:''" json:"login_ip"`
	LoginDevice   string    `gorm:"type:varchar(100);not null;default:''" json:"login_device"`
	LoginLocation string    `gorm:"type:varchar(100);not null;default:''" json:"login_location"`
	LoginMethod   string    `gorm:"type:varchar(20);not null;default:''" json:"login_method"`
	LoginStatus   int8      `gorm:"type:tinyint(1);not null;default:1" json:"login_status"` // 1-成功 0-失败
	FailureReason string    `gorm:"type:varchar(100);not null;default:''" json:"failure_reason"`
	CreatedAt     time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);index:idx_created_at" json:"created_at"`
}

func (StaffLoginHistory) TableName() string { return "sys_login_histories" }
