package models

import "eshop-monolith/pkg/utils"

type User struct {
	ID     int64 `json:"id"`
	Status int   `json:"status"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`

	UserInfo *UserInfo `json:"user_info,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsActive() bool {
	return u.Status == 1
}
