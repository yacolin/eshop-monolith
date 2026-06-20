package repositories

import (
	"context"
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

type IroleRepository interface {
	Create(ctx context.Context, role *userModels.Role) error
	GetByID(ctx context.Context, id int64) (*userModels.Role, error)
	GetByName(ctx context.Context, name string) (*userModels.Role, error)
	Update(ctx context.Context, role *userModels.Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]userModels.Role, int64, error)
	ByStatus(ctx context.Context, status int, limit, offset int) ([]userModels.Role, int64, error)
	ByUserID(ctx context.Context, userID int64) ([]userModels.Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID int64) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]userModels.Role, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]userModels.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IroleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *userModels.Role) error {
	po := models.RoleFromDomain(role)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	role.ID = po.ID
	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id int64) (*userModels.Role, error) {
	var po models.RolePO
	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", id).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*userModels.Role, error) {
	var po models.RolePO
	err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *roleRepository) Update(ctx context.Context, role *userModels.Role) error {
	po := models.RoleFromDomain(role)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.RolePO{}, "id = ?", id).Error
}

func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]userModels.Role, int64, error) {
	var pos []models.RolePO
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RolePO{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	roles := make([]userModels.Role, len(pos))
	for i, po := range pos {
		roles[i] = *po.ToDomain()
	}
	return roles, total, err
}

func (r *roleRepository) ByStatus(ctx context.Context, status int, limit, offset int) ([]userModels.Role, int64, error) {
	var pos []models.RolePO
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RolePO{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	roles := make([]userModels.Role, len(pos))
	for i, po := range pos {
		roles[i] = *po.ToDomain()
	}
	return roles, total, err
}

func (r *roleRepository) ByUserID(ctx context.Context, userID int64) ([]userModels.Role, error) {
	var pos []models.RolePO

	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.deleted_at IS NULL").
		Find(&pos).Error

	if err != nil {
		return nil, err
	}

	roles := make([]userModels.Role, len(pos))
	for i, po := range pos {
		roles[i] = *po.ToDomain()
	}
	return roles, err
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID int64) error {
	userRole := &models.UserRolePO{
		UserID: userID,
		RoleID: roleID,
	}

	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRolePO{}).Error
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID int64) ([]userModels.Role, error) {
	var pos []models.RolePO

	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.deleted_at IS NULL").
		Where("roles.status = ?", 1).
		Find(&pos).Error

	if err != nil {
		return nil, err
	}

	roles := make([]userModels.Role, len(pos))
	for i, po := range pos {
		roles[i] = *po.ToDomain()
	}
	return roles, err
}

func (r *roleRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	rolePermission := &models.RolePermissionPO{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	return r.db.WithContext(ctx).Create(rolePermission).Error
}

func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermissionPO{}).Error
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]userModels.Permission, error) {
	var pos []models.PermissionPO

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.status = ?", 1).
		Find(&pos).Error

	if err != nil {
		return nil, err
	}

	permissions := make([]userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = *po.ToDomain()
	}
	return permissions, err
}
