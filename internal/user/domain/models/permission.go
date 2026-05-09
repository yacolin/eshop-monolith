package models

import "eshop-monolith/pkg/utils"

// Permission 权限模型
// RBAC 中的最小授权单元，格式为: 资源:操作，如 order:create, product:read
type Permission struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`         // 权限名称，如：order:create
	DisplayName string `json:"display_name"` // 显示名称，如：创建订单
	Description string `json:"description"`  // 描述
	Resource    string `json:"resource"`     // 资源：order, product, user 等
	Action      string `json:"action"`       // 操作：create, read, update, delete 等
	Category    string `json:"category"`     // 分类：business, system, admin 等
	Sort        int    `json:"sort"`         // 排序
	Status      int    `json:"status"`       // 1:启用 2:禁用

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色与权限的关联表
type RolePermission struct {
	ID           int64           `json:"id"`
	RoleID       int64           `json:"role_id"`
	PermissionID int64           `json:"permission_id"`
	CreatedAt    utils.Timestamp `json:"created_at"`

	Permission *Permission `json:"permission,omitempty"`
	Role       *Role       `json:"role,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
