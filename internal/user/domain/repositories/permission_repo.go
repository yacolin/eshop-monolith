package repositories

import (
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

type IpermissionRepository interface {
	Create(permission *userModels.Permission) error
	GetByID(id int64) (*userModels.Permission, error)
	GetByName(name string) (*userModels.Permission, error)
	Update(permission *userModels.Permission) error
	Delete(id int64) error

	List(limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error)
	ByCategory(category string, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error)
	ByResource(resource string, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error)
	ByRoleID(roleID int64, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error)
	ByStatus(status int, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error)

	ExistsByName(name string) (bool, error)
	GetPermissionsByRoleIDs(roleIDs []int64) ([]*userModels.Permission, error)
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

func (r *permissionRepository) Create(permission *userModels.Permission) error {
	po := models.PermissionFromDomain(permission)
	if err := r.db.Create(po).Error; err != nil {
		return err
	}
	permission.ID = po.ID
	return nil
}

func (r *permissionRepository) GetByID(id int64) (*userModels.Permission, error) {
	var po models.PermissionPO
	err := r.db.Where("id = ?", id).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *permissionRepository) GetByName(name string) (*userModels.Permission, error) {
	var po models.PermissionPO
	err := r.db.Where("name = ?", name).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *permissionRepository) Update(permission *userModels.Permission) error {
	po := models.PermissionFromDomain(permission)
	return r.db.Save(po).Error
}

func (r *permissionRepository) Delete(id int64) error {
	return r.db.Delete(&models.PermissionPO{}, "id = ?", id).Error
}

// List 返回权限列表，默认按 sort ASC, created_at DESC 排序
func (r *permissionRepository) List(limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error) {
	var pos []*models.PermissionPO
	var total int64

	db := r.db.Model(&models.PermissionPO{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.ApplyOrder(db, sortBy, order, "sort ASC, created_at DESC")
	err := db.Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]*userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = po.ToDomain()
	}
	return permissions, total, err
}

// ByCategory 按分类查询，默认按 sort ASC, created_at DESC 排序
func (r *permissionRepository) ByCategory(category string, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error) {
	var pos []*models.PermissionPO
	var total int64

	db := r.db.Model(&models.PermissionPO{}).Where("category = ?", category)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.ApplyOrder(db, sortBy, order, "sort ASC, created_at DESC")
	err := db.Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]*userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = po.ToDomain()
	}
	return permissions, total, err
}

// ByResource 按资源查询，默认按 sort ASC, created_at DESC 排序
func (r *permissionRepository) ByResource(resource string, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error) {
	var pos []*models.PermissionPO
	var total int64

	db := r.db.Model(&models.PermissionPO{}).Where("resource = ?", resource)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.ApplyOrder(db, sortBy, order, "sort ASC, created_at DESC")
	err := db.Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]*userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = po.ToDomain()
	}
	return permissions, total, err
}

func (r *permissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.PermissionPO{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ByStatus 按状态查询，默认按 sort ASC, created_at DESC 排序
func (r *permissionRepository) ByStatus(status int, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error) {
	var pos []*models.PermissionPO
	var total int64

	db := r.db.Model(&models.PermissionPO{}).Where("status = ?", status)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.ApplyOrder(db, sortBy, order, "sort ASC, created_at DESC")
	err := db.Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]*userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = po.ToDomain()
	}
	return permissions, total, err
}

// ByRoleID 按角色查询，默认按 permissions.sort ASC, permissions.created_at DESC 排序
func (r *permissionRepository) ByRoleID(roleID int64, limit, offset int, sortBy, order string) ([]*userModels.Permission, int64, error) {
	var pos []*models.PermissionPO
	var total int64

	db := r.db.Model(&models.PermissionPO{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Where("role_permissions.deleted_at IS NULL")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.ApplyOrder(db, sortBy, order, "permissions.sort ASC, permissions.created_at DESC")
	err := db.Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]*userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = po.ToDomain()
	}
	return permissions, total, err
}

func (r *permissionRepository) GetPermissionsByRoleIDs(roleIDs []int64) ([]*userModels.Permission, error) {
	var pos []*models.PermissionPO

	err := r.db.Model(&models.PermissionPO{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.status = ?", 1).
		Distinct("permissions.*").
		Find(&pos).Error

	if err != nil {
		return nil, err
	}

	permissions := make([]*userModels.Permission, len(pos))
	for i, po := range pos {
		permissions[i] = po.ToDomain()
	}
	return permissions, err
}

func (r *permissionRepository) HasPermissionByRoleIDs(roleIDs []int64, permissionName string) (bool, error) {
	var count int64

	err := r.db.Model(&models.PermissionPO{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.name = ?", permissionName).
		Where("permissions.status = ?", 1).
		Count(&count).Error

	return count > 0, err
}

func (r *permissionRepository) AssignPermissionToRoleByID(roleID, permissionID int64) error {
	rolePermission := &models.RolePermissionPO{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	return r.db.Create(rolePermission).Error
}

func (r *permissionRepository) RemovePermissionFromRoleByID(roleID int64, permissionID int64) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermissionPO{}).Error
}
