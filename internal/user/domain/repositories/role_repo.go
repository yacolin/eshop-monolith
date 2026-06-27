package repositories

import (
	"context"
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

// UserRolePermissionRow 用户角色+权限 enriched 查询结果行（单次 LEFT JOIN 产出）
type UserRolePermissionRow struct {
	// 角色字段
	RoleID          int64           `gorm:"column:role_id"`
	RoleName        string          `gorm:"column:role_name"`
	RoleDisplayName string          `gorm:"column:role_display_name"`
	RoleDescription string          `gorm:"column:role_description"`
	RoleStatus      int             `gorm:"column:role_status"`
	RoleSort        int             `gorm:"column:role_sort"`
	RoleIsSystem    bool            `gorm:"column:role_is_system"`
	RoleCreatedAt   utils.Timestamp `gorm:"column:role_created_at"`
	RoleUpdatedAt   utils.Timestamp `gorm:"column:role_updated_at"`
	// 权限字段（LEFT JOIN 无匹配时为 NULL）
	PermID          *int64           `gorm:"column:perm_id"`
	PermName        *string          `gorm:"column:perm_name"`
	PermDisplayName *string          `gorm:"column:perm_display_name"`
	PermDescription *string          `gorm:"column:perm_description"`
	PermResource    *string          `gorm:"column:perm_resource"`
	PermAction      *string          `gorm:"column:perm_action"`
	PermCategory    *string          `gorm:"column:perm_category"`
	PermSort        *int             `gorm:"column:perm_sort"`
	PermStatus      *int             `gorm:"column:perm_status"`
	PermCreatedAt   *utils.Timestamp `gorm:"column:perm_created_at"`
	PermUpdatedAt   *utils.Timestamp `gorm:"column:perm_updated_at"`
}

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
	// GetUserRolesEnriched 获取用户角色+权限（单次 LEFT JOIN 替代 Preload）
	GetUserRolesEnriched(ctx context.Context, userID int64) ([]UserRolePermissionRow, error)
	// GetUserRoleNames 获取用户角色名列表（仅 name，不加 Preload，给 JWT 生成用）
	GetUserRoleNames(ctx context.Context, userID int64) ([]string, error)
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

// GetUserRolesEnriched 获取用户角色+权限（单次 LEFT JOIN 替代 Preload("Permissions")）
func (r *roleRepository) GetUserRolesEnriched(ctx context.Context, userID int64) ([]UserRolePermissionRow, error) {
	var rows []UserRolePermissionRow
	err := r.db.WithContext(ctx).
		Table("roles").
		Select(`roles.id AS role_id, roles.name AS role_name, roles.display_name AS role_display_name,
				roles.description AS role_description, roles.status AS role_status,
				roles.sort AS role_sort, roles.is_system AS role_is_system,
				roles.created_at AS role_created_at, roles.updated_at AS role_updated_at,
				p.id AS perm_id, p.name AS perm_name, p.display_name AS perm_display_name,
				p.description AS perm_description, p.resource AS perm_resource,
				p.action AS perm_action, p.category AS perm_category,
				p.sort AS perm_sort, p.status AS perm_status,
				p.created_at AS perm_created_at, p.updated_at AS perm_updated_at`).
		Joins("JOIN user_roles ur ON ur.role_id = roles.id AND ur.user_id = ? AND ur.deleted_at IS NULL", userID).
		Joins("LEFT JOIN role_permissions rp ON rp.role_id = roles.id AND rp.deleted_at IS NULL").
		Joins("LEFT JOIN permissions p ON p.id = rp.permission_id AND p.status = 1").
		Where("roles.status = ?", 1).
		Order("roles.sort ASC, roles.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetUserRoleNames 获取用户角色名列表（仅 name，不加载 Permissions，用于 JWT 生成）
func (r *roleRepository) GetUserRoleNames(ctx context.Context, userID int64) ([]string, error) {
	var names []string
	err := r.db.WithContext(ctx).
		Table("roles").
		Select("roles.name").
		Joins("JOIN user_roles ur ON ur.role_id = roles.id AND ur.user_id = ? AND ur.deleted_at IS NULL", userID).
		Where("roles.status = ?", 1).
		Order("roles.sort ASC").
		Pluck("roles.name", &names).Error
	if err != nil {
		return nil, err
	}
	return names, nil
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
