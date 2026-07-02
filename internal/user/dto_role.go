package user

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Status      *int8  `json:"status" binding:"omitempty"`
	SortOrder   *int   `json:"sort_order"`
	RoleType    string `json:"role_type" binding:"omitempty,oneof=custom builtin"`
}

type UpdateRoleReq struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
	Status      *int8   `json:"status" binding:"omitempty"`
	SortOrder   *int    `json:"sort_order"`
}

type RoleListResult struct {
	Total int64   `json:"total"`
	List  []*Role `json:"list"`
}

type AssignRoleReq struct {
	RoleID int64 `json:"role_id" binding:"required"`
}
