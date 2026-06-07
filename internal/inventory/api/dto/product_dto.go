package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// CachedProductItem 精简缓存 DTO，仅包含列表展示所需字段
type CachedProductItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	SKU   string `json:"sku"`
}

// ProductListQuery 商品列表查询参数
// ProductListQuery 产品列表查询参数
type ProductListQuery struct {
	query.Pagination
	Name       string `form:"name"`                // 产品名称模糊搜索
	SKU        string `form:"sku"`                 // SKU精确搜索
	CategoryID *int64 `form:"category_id"`         // 分类ID筛选
	SortBy     string `form:"sort_by"`             // 排序字段，例如 name, price, created_at
	Order      string `form:"order,default=asc"`   // asc or desc
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

// ProductCategoryBrief 产品分类简要信息
type ProductCategoryBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProductDetailDTO 产品详情（聚合产品和库存信息，字段平摊）
type ProductDetailDTO struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Price        int64               `json:"price"`
	SKU          string              `json:"sku"`
	CategoryID   int64               `json:"category_id"`
	CategoryName string              `json:"category_name"`
	Quantity     int                 `json:"quantity"`
	Status       string              `json:"status"`
	Reserved     int                 `json:"reserved"`
	Threshold    int                 `json:"threshold"`
	CreatedAt    utils.Timestamp     `json:"created_at"`
	UpdatedAt    utils.Timestamp     `json:"updated_at"`
}

// ProductWithCategoryDTO 产品列表项（含分类信息）
type ProductWithCategoryDTO struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        int64           `json:"price"`
	SKU          string          `json:"sku"`
	CategoryID   int64           `json:"category_id"`
	CategoryName string          `json:"category_name"`
	CreatedAt    utils.Timestamp `json:"created_at"`
	UpdatedAt    utils.Timestamp `json:"updated_at"`
}

// ProductWithCategoryListResult swaggo 兼容的具体类型
type ProductWithCategoryListResult struct {
	List  []ProductWithCategoryDTO `json:"list"`
	Total int64                    `json:"total"`
}
