package user

import (
	"time"
)

// User 用户领域模型
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex"`
	Email     string    `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	Password  string    `json:"password" gorm:"type:varchar(100)"`
	Role      string    `json:"role" gorm:"type:varchar(20)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role 角色模型
type Role struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(50);uniqueIndex"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Permission 权限模型
type Permission struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(50);uniqueIndex"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       int64 `json:"role_id" gorm:"primaryKey"`
	PermissionID int64 `json:"permission_id" gorm:"primaryKey"`
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate() error {
	if u.Role == "" {
		u.Role = "user"
	}
	return nil
}

// BeforeCreate 创建前钩子
func (r *Role) BeforeCreate() error {
	return nil
}

// BeforeCreate 创建前钩子
func (p *Permission) BeforeCreate() error {
	return nil
}
