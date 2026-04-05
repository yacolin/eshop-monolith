package repositories

import (
	"context"
	"eshop-monolith/internal/user/domain/models"

	"gorm.io/gorm"
)

type IroleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id int64) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]models.Role, int64, error)
	ByStatus(ctx context.Context, status int, limit, offset int) ([]models.Role, int64, error)
	ByUserID(ctx context.Context, userID int64) ([]models.Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID int64) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]models.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IroleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id int64) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error
}

func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Role{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) ByStatus(ctx context.Context, status int, limit, offset int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Role{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) ByUserID(ctx context.Context, userID int64) ([]models.Role, error) {
	var roles []models.Role

	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.deleted_at IS NULL").
		Find(&roles).Error

	return roles, err
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID int64) error {
	userRole := &models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error) {
	var roles []models.Role

	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.deleted_at IS NULL").
		Where("roles.status = ?", 1).
		Find(&roles).Error

	return roles, err
}

func (r *roleRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	rolePermission := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	return r.db.WithContext(ctx).Create(rolePermission).Error
}

func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]models.Permission, error) {
	var permissions []models.Permission

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.status = ?", 1).
		Find(&permissions).Error

	return permissions, err
}
