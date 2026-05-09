package models

import (
	"eshop-monolith/pkg/utils"
	"time"
)

// UserInfo 用户详细信息模型（对应 User.md 中的 user_profile）
// 保存可变个人信息：nickname, avatar, gender 等
type UserInfo struct {
	ID       int64      `json:"id"`
	UserID   int64      `json:"user_id"`
	Avatar   string     `json:"avatar"`
	Nickname string     `json:"nickname"`
	Gender   int        `json:"gender"` // 0:未知 1:男 2:女
	Birthday *time.Time `json:"birthday"`
	Address  string     `json:"address"`
	Bio      string     `json:"bio"`
	Country  string     `json:"country"`
	Province string     `json:"province"`
	City     string     `json:"city"`
	ZipCode  string     `json:"zip_code"`
	Language string     `json:"language"`
	Timezone string     `json:"timezone"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (UserInfo) TableName() string {
	return "user_infos"
}
