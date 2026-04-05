package dto

import (
	"eshop-monolith/internal/pkg/query"
	"eshop-monolith/internal/user/domain/models"
)

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Category    string `json:"category"`
	Sort        int    `json:"sort"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Sort        *int    `json:"sort"`
	Status      *int    `json:"status"`
}

// ListPermissionsResponse 权限列表响应
type ListPermissionsResponse struct {
	Permissions []*models.Permission `json:"permissions"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
}

// ListRolePermissionsResponse 角色权限列表响应
type ListRolePermissionsResponse struct {
	RolePermissions []*models.RolePermission `json:"role_permissions"`
	Total           int64                    `json:"total"`
	Page            int                      `json:"page"`
	PageSize        int                      `json:"page_size"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	IsSystem    bool   `json:"is_system"`
}

type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
	Sort        *int    `json:"sort"`
}

type ListRolesResponse struct {
	Roles    []models.Role `json:"roles"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// PermissionListQuery 权限列表查询
type PermissionListQuery struct {
	query.Pagination
	Category string `form:"category"`
	Resource string `form:"resource"`
	RoleID   int64  `form:"role_id"`
}

type CheckPermissionsRequest struct {
	PermissionNames []string `json:"permission_names" binding:"required" example:"order:create,product:read"`
}

// CheckPermissionsResponse 检查权限响应
type CheckPermissionsResponse struct {
	Permissions map[string]bool `json:"permissions"`
}
