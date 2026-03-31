package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/pkg/query"
)

// CategoryListQuery 分类列表查询参数
type CategoryListQuery struct {
	query.Pagination
	Name     string  `form:"name"`              // 分类名称模糊搜索
	ParentID *string `form:"parent_id"`         // 父分类ID
	SortBy   string  `form:"sort_by"`           // 排序字段，例如 name, created_at
	Order    string  `form:"order,default=asc"` // asc or desc
}

// CategoryListResult 分类列表结果（使用泛型）
type CategoryListResult = query.ListResult[models.Category]

// CreateCategoryDTO 创建分类请求
type CreateCategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id"`
}

// UpdateCategoryDTO 更新分类请求
type UpdateCategoryDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *int64  `json:"parent_id"`
}
