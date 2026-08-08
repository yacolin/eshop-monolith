package staff

import (
	"context"

	"gorm.io/gorm"
)

type IstaffRoleRepository interface {
	GetRoleNamesByStaffID(ctx context.Context, staffID int64) ([]string, error)
	ReplaceStaffRoles(ctx context.Context, staffID int64, roleIDs []int64) error
	// HasPermission 实时校验员工是否拥有指定权限(join 四表,角色/权限均需启用)
	HasPermission(ctx context.Context, staffID int64, permissionName string) (bool, error)
}

type StaffRoleRepository struct {
	db *gorm.DB
}

func NewStaffRoleRepository(db *gorm.DB) IstaffRoleRepository {
	return &StaffRoleRepository{db: db}
}

func (r *StaffRoleRepository) GetRoleNamesByStaffID(ctx context.Context, staffID int64) ([]string, error) {
	var names []string
	err := r.db.WithContext(ctx).Model(&SysRole{}).
		Joins("JOIN sys_staff_roles sr ON sr.role_id = sys_roles.id AND sr.deleted_at IS NULL").
		Where("sr.staff_id = ? AND sys_roles.status = 1", staffID).
		Pluck("sys_roles.name", &names).Error
	return names, err
}

// ReplaceStaffRoles 全量替换员工角色绑定
func (r *StaffRoleRepository) ReplaceStaffRoles(ctx context.Context, staffID int64, roleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("staff_id = ?", staffID).Delete(&SysStaffRole{}).Error; err != nil {
			return err
		}
		for _, rid := range roleIDs {
			if err := tx.Create(&SysStaffRole{StaffID: staffID, RoleID: rid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *StaffRoleRepository) HasPermission(ctx context.Context, staffID int64, permissionName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysStaffRole{}).
		Joins("JOIN sys_roles r ON r.id = sys_staff_roles.role_id AND r.status = 1 AND r.deleted_at IS NULL").
		Joins("JOIN sys_role_permissions rp ON rp.role_id = r.id AND rp.deleted_at IS NULL").
		Joins("JOIN sys_permissions p ON p.id = rp.permission_id AND p.status = 1 AND p.deleted_at IS NULL").
		Where("sys_staff_roles.staff_id = ? AND sys_staff_roles.deleted_at IS NULL AND p.name = ?", staffID, permissionName).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
