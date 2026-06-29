package product

import "eshop-monolith/pkg/query"

type CreateCategoryReq struct {
	Name      string `json:"name" binding:"required,max=100"`
	ParentID  int64  `json:"parent_id"`
	IconURL   string `json:"icon_url" binding:"max=512"`
	SortOrder int    `json:"sort_order"`
}

type UpdateCategoryReq struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	IconURL   *string `json:"icon_url" binding:"omitempty,max=512"`
	SortOrder *int    `json:"sort_order"`
	Status    *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

type CategoryListReq struct {
	query.Pagination
	Name   string `form:"name"`
	Status *int8  `form:"status"`
	Level  *int8  `form:"level"`
}
