package models

import "eshop-monolith/pkg/utils"

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	IsSystem    bool   `json:"is_system"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`

	Permissions []Permission `json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

type UserRole struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`

	CreatedAt utils.Timestamp `json:"created_at"`

	Role *Role `json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
