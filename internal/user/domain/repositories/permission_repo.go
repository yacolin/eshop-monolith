package repositories

import (
	"eshop-monolith/internal/user/domain/models"

	"gorm.io/gorm"
)

type IpermissionRepository interface {
	Create(permission *models.Permission) error
	GetByID(id int64) (*models.Permission, error)
	GetByName(name string) (*models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id int64) error

	List(limit, offset int) ([]*models.Permission, int64, error)
	ByCategory(category string, limit, offset int) ([]*models.Permission, int64, error)
	ByResource(resource string, limit, offset int) ([]*models.Permission, int64, error)
	ByRoleID(roleID int64, limit, offset int) ([]*models.Permission, int64, error)
	ByStatus(status int, limit, offset int) ([]*models.Permission, int64, error)

	ExistsByName(name string) (bool, error)
	GetPermissionsByRoleIDs(roleIDs []int64) ([]*models.Permission, error)
	HasPermissionByRoleIDs(roleIDs []int64, permissionName string) (bool, error)

	AssignPermissionToRoleByID(roleID, permissionID int64) error
	RemovePermissionFromRoleByID(roleID, permissionID int64) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) IpermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetByID(id int64) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id int64) error {
	return r.db.Delete(&models.Permission{}, "id = ?", id).Error
}

func (r *permissionRepository) List(limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByCategory(category string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).Where("category = ?", category)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByResource(resource string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).Where("resource = ?", resource)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) ByStatus(status int, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByRoleID(roleID int64, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Where("role_permissions.deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("permissions.sort ASC, permissions.created_at DESC").
		Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) GetPermissionsByRoleIDs(roleIDs []int64) ([]*models.Permission, error) {
	var permissions []*models.Permission

	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.status = ?", 1).
		Distinct("permissions.*").
		Find(&permissions).Error

	return permissions, err
}

func (r *permissionRepository) HasPermissionByRoleIDs(roleIDs []int64, permissionName string) (bool, error) {
	var count int64

	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.name = ?", permissionName).
		Where("permissions.status = ?", 1).
		Count(&count).Error

	return count > 0, err
}

func (r *permissionRepository) AssignPermissionToRoleByID(roleID, permissionID int64) error {
	rolePermission := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	return r.db.Create(rolePermission).Error
}

func (r *permissionRepository) RemovePermissionFromRoleByID(roleID int64, permissionID int64) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}
