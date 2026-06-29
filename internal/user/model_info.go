package user

import (
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64          `gorm:"not null;uniqueIndex:uk_user_id" json:"user_id"`
	Gender    int8           `gorm:"type:tinyint;not null;default:0" json:"gender"`
	Birthday  *time.Time     `gorm:"type:date" json:"birthday"`
	Bio       string         `gorm:"type:varchar(200);not null;default:''" json:"bio"`
	Country   string         `gorm:"type:varchar(32);not null;default:''" json:"country"`
	Province  string         `gorm:"type:varchar(32);not null;default:''" json:"province"`
	City      string         `gorm:"type:varchar(32);not null;default:''" json:"city"`
	ZipCode   string         `gorm:"type:varchar(10);not null;default:''" json:"zip_code"`
	Language  string         `gorm:"type:varchar(10);not null;default:'zh-CN'" json:"language"`
	Timezone  string         `gorm:"type:varchar(32);not null;default:'Asia/Shanghai'" json:"timezone"`
	CreatedAt time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (UserInfo) TableName() string { return "usr_infos" }
