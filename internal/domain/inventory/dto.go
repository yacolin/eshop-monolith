package inventory

import "eshop-monolith/internal/pkg/query"

// InventoryListQuery 库存列表查询参数
type InventoryListQuery struct {
	query.Pagination
	ProductName string `form:"product_name"`      // 产品名称模糊搜索
	SKU         string `form:"sku"`               // SKU精确搜索
	Status      string `form:"status"`            // 库存状态
	LowStock    *bool  `form:"low_stock"`         // 是否低库存
	SortBy      string `form:"sort_by"`           // 排序字段，例如 quantity, reserved, created_at
	Order       string `form:"order,default=asc"` // asc or desc
}

// InventoryListResult 库存列表结果（使用泛型）
type InventoryListResult = query.ListResult[Inventory]

// CreateInventoryDTO 创建库存请求
type CreateInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=0"`
	Threshold int   `json:"threshold" binding:"min=0"`
}

// UpdateInventoryDTO 更新库存请求
type UpdateInventoryDTO struct {
	Quantity  *int `json:"quantity"`
	Threshold *int `json:"threshold"`
	Reserved  *int `json:"reserved"`
}

// ReserveInventoryDTO 预订库存请求
type ReserveInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

// ReleaseInventoryDTO 释放库存请求
type ReleaseInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

// AdjustInventoryDTO 调整库存请求
type AdjustInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required"` // 正数增加，负数减少
}
