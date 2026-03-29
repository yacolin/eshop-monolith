package repository

import (
	"context"

	"gorm.io/gorm"

	"eshop-monolith/internal/domain/user"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

// Create 创建用户
func (r UserRepository) Create(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据ID查询用户
func (r UserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var foundUser user.User
	err := r.db.WithContext(ctx).First(&foundUser, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// FindByUsername 根据用户名查询用户
func (r UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var foundUser user.User
	err := r.db.WithContext(ctx).First(&foundUser, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// FindByEmail 根据邮箱查询用户
func (r UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var foundUser user.User
	err := r.db.WithContext(ctx).First(&foundUser, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// Update 更新用户
func (r UserRepository) Update(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r UserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id).Error
}

// FindRoleByName 根据名称查询角色
func (r UserRepository) FindRoleByName(ctx context.Context, name string) (*user.Role, error) {
	var role user.Role
	err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles 列出所有角色
func (r UserRepository) ListRoles(ctx context.Context) ([]user.Role, error) {
	var roles []user.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// FindPermissionByName 根据名称查询权限
func (r UserRepository) FindPermissionByName(ctx context.Context, name string) (*user.Permission, error) {
	var permission user.Permission
	err := r.db.WithContext(ctx).First(&permission, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// ListPermissions 列出所有权限
func (r UserRepository) ListPermissions(ctx context.Context) ([]user.Permission, error) {
	var permissions []user.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// CheckPermission 检查用户是否有指定权限
func (r UserRepository) CheckPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("users").Joins("JOIN roles ON users.role = roles.name").Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").Where("users.id = ? AND permissions.name = ?", userID, permissionName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
