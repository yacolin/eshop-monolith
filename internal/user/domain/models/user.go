package models

import (
	"eshop-monolith/internal/pkg/utils"

	"gorm.io/gorm"
)

type User struct {
	ID     int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Status int   `gorm:"type:tinyint;default:1" json:"status"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	UserInfo *UserInfo `gorm:"foreignKey:UserID" json:"user_info,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsActive() bool {
	return u.Status == 1
}
