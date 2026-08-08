package staff

import (
	"time"

	"eshop-monolith/pkg/query"
)

// ── Request ──

type StaffListReq struct {
	query.Pagination
	Keyword string `form:"keyword"`
	Status  *int8  `form:"status"`
}

type CreateStaffReq struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name" binding:"max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"max=20"`
	Status   *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}

type UpdateStaffReq struct {
	RealName string `json:"real_name" binding:"max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"max=20"`
	Avatar   string `json:"avatar" binding:"max=512"`
	Status   *int8  `json:"status" binding:"omitempty,oneof=0 1"`
	Password string `json:"password" binding:"omitempty,min=6"` // 重置密码,空表示不修改
}

type AssignRolesReq struct {
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}

// ── Response ──

type StaffListItem struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	RealName    string     `json:"real_name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Status      int8       `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	Roles       []string   `json:"roles"`
}

type StaffListResult struct {
	Total int64            `json:"total"`
	List  []*StaffListItem `json:"list"`
}
