package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
)

// ProductListQuery 商品列表查询参数
// ProductListQuery 产品列表查询参数
type ProductListQuery struct {
	query.Pagination
	Name   string `form:"name"`              // 产品名称模糊搜索
	SKU    string `form:"sku"`               // SKU精确搜索
	SortBy string `form:"sort_by"`           // 排序字段，例如 name, price, created_at
	Order  string `form:"order,default=asc"` // asc or desc
}

// ProductListResult 商品列表结果（使用泛型）
type ProductListResult = query.ListResult[models.Product]

// CreateProductDTO 创建商品请求
type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description string  `json:"description" binding:"max=65535"`
	Price       int64   `json:"price" binding:"required,gt=0"`
	SKU         string  `json:"sku" binding:"required,max=100"`
	CategoryIDs []int64 `json:"category_ids" binding:"omitempty,dive,gt=0"`
}

// UpdateProductDTO 更新商品请求
type UpdateProductDTO struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Description *string `json:"description" binding:"omitempty,max=65535"`
	Price       *int64  `json:"price" binding:"omitempty,gt=0"`
	CategoryIDs []int64 `json:"category_ids" binding:"omitempty,dive,gt=0"`
}
