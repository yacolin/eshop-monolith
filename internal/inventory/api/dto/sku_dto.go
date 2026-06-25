package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

type CreateSkuDTO struct {
	ProductID int64             `json:"product_id" binding:"required"`
	Name      string            `json:"name" binding:"required,max=255"`
	Price     int64             `json:"price" binding:"required,gt=0"`
	SKUCode   string            `json:"sku_code" binding:"required,max=100"`
	Image     string            `json:"image" binding:"max=500"`
	Spec      map[string]string `json:"spec"`
}

type UpdateSkuDTO struct {
	Name     *string            `json:"name" binding:"omitempty,max=255"`
	Price    *int64             `json:"price" binding:"omitempty,gt=0"`
	SKUCode  *string            `json:"sku_code" binding:"omitempty,max=100"`
	Image    *string            `json:"image" binding:"omitempty,max=500"`
	Spec     map[string]string  `json:"spec"`
}

type SkuResponse struct {
	ID        int64             `json:"id"`
	ProductID int64             `json:"product_id"`
	Name      string            `json:"name"`
	Price     int64             `json:"price"`
	SKUCode   string            `json:"sku_code"`
	Image     string            `json:"image,omitempty"`
	Spec      map[string]string `json:"spec,omitempty"`
	CreatedAt utils.Timestamp   `json:"created_at"`
	UpdatedAt utils.Timestamp   `json:"updated_at"`
}

// SkuListQuery SKU 列表查询参数
type SkuListQuery struct {
	query.Pagination
	ProductID *int64 `form:"product_id"`       // 可选，按产品ID筛选
	Name      string `form:"name"`             // 可选，SKU名称模糊搜索
	SKUCode   string `form:"sku_code"`         // 可选，SKU编码精确搜索
	PriceMin  *int64 `form:"price_min"`        // 可选，最低价格（分）
	PriceMax  *int64 `form:"price_max"`        // 可选，最高价格（分）
	SortBy    string `form:"sort_by"`          // 排序字段 (id, name, price, created_at)
	Order     string `form:"order,default=asc"` // 排序方向 (asc, desc)
}

type SkuListResult struct {
	List  []SkuResponse `json:"list"`
	Total int64         `json:"total"`
}

// SkuDetailResponse SKU 详情响应（含库存信息）
type SkuDetailResponse struct {
	ID                int64             `json:"id"`
	ProductID         int64             `json:"product_id"`
	Name              string            `json:"name"`
	Price             int64             `json:"price"`
	SKUCode           string            `json:"sku_code"`
	Image             string            `json:"image,omitempty"`
	Spec              map[string]string `json:"spec,omitempty"`
	AvailableQuantity int               `json:"available_quantity"`
	InventoryStatus   string            `json:"inventory_status"`
	CreatedAt         utils.Timestamp   `json:"created_at"`
	UpdatedAt         utils.Timestamp   `json:"updated_at"`
}

// SkuDetailToResponse 将 SKU + 库存信息映射为 SkuDetailResponse
func SkuDetailToResponse(s *models.Sku, availableQuantity int, inventoryStatus string) SkuDetailResponse {
	return SkuDetailResponse{
		ID:                s.ID,
		ProductID:         s.ProductID,
		Name:              s.Name,
		Price:             s.Price,
		SKUCode:           s.SKUCode,
		Image:             s.Image,
		Spec:              s.Spec,
		AvailableQuantity: availableQuantity,
		InventoryStatus:   inventoryStatus,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func SkuToResponse(s *models.Sku) SkuResponse {
	return SkuResponse{
		ID: s.ID, ProductID: s.ProductID, Name: s.Name, Price: s.Price,
		SKUCode: s.SKUCode, Image: s.Image, Spec: s.Spec,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}
