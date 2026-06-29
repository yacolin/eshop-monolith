package user

import "eshop-monolith/pkg/query"

type CreatePermissionReq struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Category    string `json:"category"`
	SortOrder   *int   `json:"sort_order"`
}

type UpdatePermissionReq struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

type PermissionListReq struct {
	query.Pagination
	Category string `form:"category"`
	Resource string `form:"resource"`
	RoleID   int64  `form:"role_id"`
}

type PermissionListResult struct {
	Total int64         `json:"total"`
	List  []*Permission `json:"list"`
}

type CheckPermissionsReq struct {
	PermissionNames []string `json:"permission_names" binding:"required"`
}

type CheckPermissionsResult struct {
	Permissions map[string]bool `json:"permissions"`
}

type AssignPermissionsReq struct {
	PermissionIDs []int64 `json:"permission_ids" binding:"required"`
}
