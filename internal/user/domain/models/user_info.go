package models

import (
	"eshop-monolith/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// UserInfo 用户详细信息模型（对应 User.md 中的 user_profile）
// 保存可变个人信息：nickname, avatar, gender 等
type UserInfo struct {
	ID       int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID   int64      `gorm:"not null;uniqueIndex" json:"user_id"`
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`
	Nickname string     `gorm:"type:varchar(50)" json:"nickname"`
	Gender   int        `gorm:"type:tinyint;default:0" json:"gender"` // 0:未知 1:男 2:女
	Birthday *time.Time `json:"birthday"`
	Address  string     `gorm:"type:varchar(255)" json:"address"`
	Bio      string     `gorm:"type:text" json:"bio"`
	Country  string     `gorm:"type:varchar(50)" json:"country"`
	Province string     `gorm:"type:varchar(50)" json:"province"`
	City     string     `gorm:"type:varchar(50)" json:"city"`
	ZipCode  string     `gorm:"type:varchar(20)" json:"zip_code"`
	Language string     `gorm:"type:varchar(20);default:zh-CN" json:"language"`
	Timezone string     `gorm:"type:varchar(50);default:Asia/Shanghai" json:"timezone"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (UserInfo) TableName() string {
	return "user_infos"
}
