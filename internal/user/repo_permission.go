package user

import (
	"context"

	"gorm.io/gorm"
)

type IpermissionRepository interface {
	Create(ctx context.Context, perm *Permission) error
	FindByID(ctx context.Context, id int64) (*Permission, error)
	FindByName(ctx context.Context, name string) (*Permission, error)
	Update(ctx context.Context, perm *Permission) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, category, resource string, roleID int64, page, size int) ([]Permission, int64, error)
	GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]Permission, error)
	HasPermissionByRoleIDs(ctx context.Context, roleIDs []int64, permissionName string) (bool, error)
	AssignToRole(ctx context.Context, roleID, permissionID int64) error
	RemoveFromRole(ctx context.Context, roleID, permissionID int64) error
	GetByRoleID(ctx context.Context, roleID int64) ([]Permission, error)
}

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) IpermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, perm *Permission) error {
	return r.db.WithContext(ctx).Create(perm).Error
}

func (r *PermissionRepository) FindByID(ctx context.Context, id int64) (*Permission, error) {
	var p Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	return &p, err
}

func (r *PermissionRepository) FindByName(ctx context.Context, name string) (*Permission, error) {
	var p Permission
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&p).Error
	return &p, err
}

func (r *PermissionRepository) Update(ctx context.Context, perm *Permission) error {
	return r.db.WithContext(ctx).Model(perm).Updates(perm).Error
}

func (r *PermissionRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Permission{}).Error
}

func (r *PermissionRepository) List(ctx context.Context, category, resource string, roleID int64, page, size int) ([]Permission, int64, error) {
	var list []Permission
	var total int64

	db := r.db.WithContext(ctx).Model(&Permission{})
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if resource != "" {
		db = db.Where("resource = ?", resource)
	}
	if roleID > 0 {
		db = db.Joins("JOIN usr_role_permissions ON usr_role_permissions.permission_id = usr_permissions.id").
			Where("usr_role_permissions.role_id = ?", roleID).
			Where("usr_role_permissions.deleted_at IS NULL")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, total, err
}

func (r *PermissionRepository) GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]Permission, error) {
	var perms []Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN usr_role_permissions ON usr_role_permissions.permission_id = usr_permissions.id").
		Where("usr_role_permissions.role_id IN ?", roleIDs).
		Where("usr_role_permissions.deleted_at IS NULL").
		Where("usr_permissions.status = ?", 1).
		Distinct("usr_permissions.*").
		Order("usr_permissions.sort_order ASC").
		Find(&perms).Error
	return perms, err
}

func (r *PermissionRepository) HasPermissionByRoleIDs(ctx context.Context, roleIDs []int64, permissionName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Permission{}).
		Joins("JOIN usr_role_permissions ON usr_role_permissions.permission_id = usr_permissions.id").
		Where("usr_role_permissions.role_id IN ?", roleIDs).
		Where("usr_role_permissions.deleted_at IS NULL").
		Where("usr_permissions.name = ?", permissionName).
		Where("usr_permissions.status = ?", 1).
		Count(&count).Error
	return count > 0, err
}

func (r *PermissionRepository) AssignToRole(ctx context.Context, roleID, permissionID int64) error {
	rp := RolePermission{RoleID: roleID, PermissionID: permissionID}
	return r.db.WithContext(ctx).Create(&rp).Error
}

func (r *PermissionRepository) RemoveFromRole(ctx context.Context, roleID, permissionID int64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&RolePermission{}).Error
}

func (r *PermissionRepository) GetByRoleID(ctx context.Context, roleID int64) ([]Permission, error) {
	var perms []Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN usr_role_permissions ON usr_role_permissions.permission_id = usr_permissions.id").
		Where("usr_role_permissions.role_id = ?", roleID).
		Where("usr_role_permissions.deleted_at IS NULL").
		Where("usr_permissions.status = ?", 1).
		Order("usr_permissions.sort_order ASC").
		Find(&perms).Error
	return perms, err
}
