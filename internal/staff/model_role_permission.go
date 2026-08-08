package staff

import (
	"time"

	"gorm.io/gorm"
)

// SysRolePermission 角色-权限关联(sys_role_permissions)
type SysRolePermission struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       int64          `gorm:"not null;uniqueIndex:uk_role_permission" json:"role_id"`
	PermissionID int64          `gorm:"not null;uniqueIndex:uk_role_permission" json:"permission_id"`
	ScopeType    string         `gorm:"type:varchar(20);not null;default:'platform'" json:"scope_type"`
	ScopeID      int64          `gorm:"not null;default:0" json:"scope_id"`
	CreatedAt    time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (SysRolePermission) TableName() string { return "sys_role_permissions" }
