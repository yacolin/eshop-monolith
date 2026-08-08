package staff

import (
	"context"

	"gorm.io/gorm"
)

type IroleRepository interface {
	List(ctx context.Context) ([]SysRole, error)
	FindByID(ctx context.Context, id int64) (*SysRole, error)
	Create(ctx context.Context, role *SysRole) error
	Update(ctx context.Context, role *SysRole) error
	Delete(ctx context.Context, id int64) error // 软删角色并级联清理关联
	GetPermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error)
	ReplaceRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IroleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) List(ctx context.Context) ([]SysRole, error) {
	var roles []SysRole
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*SysRole, error) {
	var role SysRole
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *RoleRepository) Create(ctx context.Context, role *SysRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) Update(ctx context.Context, role *SysRole) error {
	return r.db.WithContext(ctx).Model(role).Updates(map[string]interface{}{
		"display_name": role.DisplayName,
		"description":  role.Description,
		"sort_order":   role.SortOrder,
		"status":       role.Status,
	}).Error
}

// Delete 事务内:软删角色 → 清理角色权限绑定与员工角色绑定
func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&SysRolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&SysStaffRole{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&SysRole{}).Error
	})
}

func (r *RoleRepository) GetPermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Model(&SysRolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &ids).Error
	return ids, err
}

// ReplaceRolePermissions 全量替换角色权限绑定(scope_type=platform, scope_id=0)
func (r *RoleRepository) ReplaceRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&SysRolePermission{}).Error; err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			rp := &SysRolePermission{RoleID: roleID, PermissionID: pid, ScopeType: "platform", ScopeID: 0}
			if err := tx.Create(rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
