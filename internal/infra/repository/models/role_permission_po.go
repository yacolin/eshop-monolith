package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// RolePermissionPO 角色权限关联持久化对象
type RolePermissionPO struct {
	ID           int64          `gorm:"type:bigint;primaryKey"`
	RoleID       int64          `gorm:"type:bigint;not null;index:idx_role_id"`
	PermissionID int64          `gorm:"type:bigint;not null;index:idx_permission_id"`
	CreatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Permission   *PermissionPO  `gorm:"foreignKey:PermissionID"`
	Role         *RolePO        `gorm:"foreignKey:RoleID"`
}

func (RolePermissionPO) TableName() string { return "role_permissions" }

func (po *RolePermissionPO) ToDomain() *userDomain.RolePermission {
	var permission *userDomain.Permission
	if po.Permission != nil {
		permission = po.Permission.ToDomain()
	}
	var role *userDomain.Role
	if po.Role != nil {
		role = po.Role.ToDomain()
	}
	return &userDomain.RolePermission{
		ID:           po.ID,
		RoleID:       po.RoleID,
		PermissionID: po.PermissionID,
		CreatedAt:    utils.Timestamp(po.CreatedAt),
		Permission:   permission,
		Role:         role,
	}
}

func RolePermissionFromDomain(rp *userDomain.RolePermission) *RolePermissionPO {
	return &RolePermissionPO{
		ID:           rp.ID,
		RoleID:       rp.RoleID,
		PermissionID: rp.PermissionID,
		CreatedAt:    time.Time(rp.CreatedAt),
	}
}
