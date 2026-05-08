package models

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

// Permission 权限模型
// RBAC 中的最小授权单元，格式为: 资源:操作，如 order:create, product:read
type Permission struct {
	ID          int64  `gorm:"type:int;primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"` // 权限名称，如：order:create
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`     // 显示名称，如：创建订单
	Description string `gorm:"type:text" json:"description"`                       // 描述
	Resource    string `gorm:"type:varchar(50);not null;index" json:"resource"`    // 资源：order, product, user 等
	Action      string `gorm:"type:varchar(50);not null;index" json:"action"`      // 操作：create, read, update, delete 等
	Category    string `gorm:"type:varchar(50)" json:"category"`                   // 分类：business, system, admin 等
	Sort        int    `gorm:"type:int;default:0" json:"sort"`                     // 排序
	Status      int    `gorm:"type:tinyint;default:1" json:"status"`               // 1:启用 2:禁用

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}

// PermissionCategory 权限分类常量
// const (
// 	PermissionCategoryBusiness = "business" // 业务权限
// 	PermissionCategorySystem   = "system"   // 系统权限
// 	PermissionCategoryAdmin    = "admin"    // 管理权限
// )

// PermissionAction 权限操作常量
// const (
// 	ActionCreate  = "create"
// 	ActionRead    = "read"
// 	ActionUpdate  = "update"
// 	ActionDelete  = "delete"
// 	ActionList    = "list"
// 	ActionExport  = "export"
// 	ActionImport  = "import"
// 	ActionApprove = "approve"
// 	ActionReject  = "reject"
// )

// RolePermission 角色与权限的关联表
type RolePermission struct {
	ID           int64           `gorm:"type:bigint;primaryKey" json:"id"`
	RoleID       int64           `gorm:"type:bigint;not null;index:idx_role_id" json:"role_id"`
	PermissionID int64           `gorm:"type:bigint;not null;index:idx_permission_id" json:"permission_id"`
	CreatedAt    utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`

	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	Role       *Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
