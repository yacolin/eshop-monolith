package user

import (
	"context"
)

// Repository 用户仓储接口
type Repository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error
	// FindByID 根据ID查询用户
	FindByID(ctx context.Context, id int64) (*User, error)
	// FindByUsername 根据用户名查询用户
	FindByUsername(ctx context.Context, username string) (*User, error)
	// FindByEmail 根据邮箱查询用户
	FindByEmail(ctx context.Context, email string) (*User, error)
	// Update 更新用户
	Update(ctx context.Context, user *User) error
	// Delete 删除用户
	Delete(ctx context.Context, id int64) error

	// FindRoleByName 根据名称查询角色
	FindRoleByName(ctx context.Context, name string) (*Role, error)
	// ListRoles 列出所有角色
	ListRoles(ctx context.Context) ([]Role, error)

	// FindPermissionByName 根据名称查询权限
	FindPermissionByName(ctx context.Context, name string) (*Permission, error)
	// ListPermissions 列出所有权限
	ListPermissions(ctx context.Context) ([]Permission, error)

	// CheckPermission 检查用户是否有指定权限
	CheckPermission(ctx context.Context, userID int64, permissionName string) (bool, error)
}
