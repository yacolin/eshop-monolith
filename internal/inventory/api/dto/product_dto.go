package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// CachedProductItem 精简缓存 DTO，仅包含列表展示所需字段
type CachedProductItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	MinPrice int64  `json:"min_price"`
}

// ProductListQuery 商品列表查询参数
type ProductListQuery struct {
	query.Pagination
	Name       string `form:"name"`              // 产品名称模糊搜索
	SKU        string `form:"sku"`               // SKU精确搜索
	CategoryID *int64 `form:"category_id"`       // 分类ID筛选
	SortBy     string `form:"sort_by"`           // 排序字段，例如 name, price, created_at
	Order      string `form:"order,default=asc"` // asc or desc
}

// ProductListResult 商品列表结果
type ProductListResult struct {
	Total int64           `json:"total"`
	List  []models.Product `json:"list"`
}

// CreateProductDTO 创建商品请求
type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description string  `json:"description" binding:"max=65535"`
	CategoryIDs []int64 `json:"category_ids" binding:"omitempty,dive,gt=0"`
}

// UpdateProductDTO 更新商品请求
type UpdateProductDTO struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Description *string `json:"description" binding:"omitempty,max=65535"`
	CategoryIDs []int64 `json:"category_ids" binding:"omitempty,dive,gt=0"`
}

// ProductCategoryBrief 产品分类简要信息
type ProductCategoryBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProductResponse 产品简要响应
type ProductResponse struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	MinPrice    int64           `json:"min_price"`
	CreatedAt   utils.Timestamp `json:"created_at"`
	UpdatedAt   utils.Timestamp `json:"updated_at"`
}

// ProductToResponse 将领域模型转换为响应 DTO
func ProductToResponse(p *models.Product) ProductResponse {
	return ProductResponse{
		ID: p.ID, Name: p.Name, Description: p.Description,
		MinPrice:  p.MinPrice,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

// ProductDetailDTO 产品详情（聚合产品和库存信息，字段平摊）
type ProductDetailDTO struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	MinPrice    int64                  `json:"min_price"`
	Categories  []ProductCategoryBrief `json:"categories"`
	Skus        []SkuResponse          `json:"skus,omitempty"`
	Quantity    int                    `json:"quantity"`
	Status      string                 `json:"status"`
	Reserved    int                    `json:"reserved"`
	Threshold   int                    `json:"threshold"`
	CreatedAt   utils.Timestamp        `json:"created_at"`
	UpdatedAt   utils.Timestamp        `json:"updated_at"`
}

// ProductWithSkusResponse 产品详情（含 SKU 列表）
type ProductWithSkusResponse struct {
	Product ProductResponse `json:"product"`
	Skus    []SkuDetailResponse `json:"skus"`
}

// ProductWithCategoryDTO 产品列表项（含分类信息）
type ProductWithCategoryDTO struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	MinPrice    int64                  `json:"min_price"`
	Categories  []ProductCategoryBrief `json:"categories"`
	CreatedAt   utils.Timestamp        `json:"created_at"`
	UpdatedAt   utils.Timestamp        `json:"updated_at"`
}

// ProductWithCategoryListResult swaggo 兼容的具体类型
type ProductWithCategoryListResult struct {
	List  []ProductWithCategoryDTO `json:"list"`
	Total int64                    `json:"total"`
}

// ProductCursorQuery 游标分页查询参数（深分页优化，基于主键 ID 游标）
type ProductCursorQuery struct {
	Cursor     int64  `form:"cursor"`               // 游标（上一页最后一条的 ID，首次查询传 0）
	Size       int    `form:"size,default=20"`      // 每页条数
	Name       string `form:"name"`                 // 产品名称模糊搜索
	SKU        string `form:"sku"`                  // SKU精确搜索
	CategoryID *int64 `form:"category_id"`          // 分类ID筛选
}

// Normalize 校验并规范化游标查询参数
func (q *ProductCursorQuery) Normalize() {
	if q.Size <= 0 {
		q.Size = 20
	}
	if q.Size > 100 {
		q.Size = 100
	}
	if q.Cursor < 0 {
		q.Cursor = 0
	}
}

// ProductCursorResult 游标分页结果
type ProductCursorResult struct {
	List       []models.Product `json:"list"`
	NextCursor int64            `json:"next_cursor"`  // 下一页游标值（无更多数据时为 0）
	HasMore    bool             `json:"has_more"`     // 是否还有更多数据
}

// ProductCacheCursorQuery 缓存游标查询参数
type ProductCacheCursorQuery struct {
	Cursor     int64  `form:"cursor"`               // 游标（上一页最后一条的 ID，首次查询传 0）
	Size       int    `form:"size,default=20"`      // 每页条数
	CategoryID *int64 `form:"category_id"`          // 分类ID筛选
}

func (q *ProductCacheCursorQuery) Normalize() {
	if q.Size <= 0 {
		q.Size = 20
	}
	if q.Size > 100 {
		q.Size = 100
	}
	if q.Cursor < 0 {
		q.Cursor = 0
	}
}

// ProductCacheCursorResult 缓存游标结果
type ProductCacheCursorResult struct {
	List       []CachedProductItem `json:"list"`
	NextCursor int64               `json:"next_cursor"`
	HasMore    bool                `json:"has_more"`
}
