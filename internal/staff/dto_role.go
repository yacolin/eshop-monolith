package staff

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	DisplayName string `json:"display_name" binding:"max=100"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateRoleReq struct {
	DisplayName string `json:"display_name" binding:"max=100"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order"`
	Status      *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}

type AssignPermissionsReq struct {
	PermissionIDs []int64 `json:"permission_ids" binding:"required"`
}

type RoleListItem struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	RoleType        string `json:"role_type"`
	Status          int8   `json:"status"`
	SortOrder       int    `json:"sort_order"`
	PermissionCount int64  `json:"permission_count"`
}
