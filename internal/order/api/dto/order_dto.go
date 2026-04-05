package dto

import (
	"eshop-monolith/internal/order/domain/models"
	"eshop-monolith/internal/pkg/query"
)

// OrderListQuery 支持通过 query string 进行过滤、排序和分页
type OrderListQuery struct {
	query.Pagination
	CustomerID *int64   `form:"customer_id"`       // 用户ID过滤
	Status     string   `form:"status"`            // 订单状态过滤
	MinPrice   *float64 `form:"min_price"`         // 价格区间下限
	MaxPrice   *float64 `form:"max_price"`         // 价格区间上限
	SortBy     string   `form:"sort_by"`           // 排序字段，例如 total_price, created_at
	Order      string   `form:"order,default=asc"` // asc or desc
}

// OrderListResult 订单列表结果（使用泛型）
type OrderListResult = query.ListResult[models.Order]

type CreateOrderDTO struct {
	CustomerID string               `json:"customer_id" binding:"required"`
	Currency   string               `json:"currency"` // 可选，默认 CNY
	Items      []CreateOrderItemDTO `json:"items" binding:"required,min=1,dive"`
}

type UpdateOrderDTO struct {
	UserID     *int64   `json:"user_id"`
	ProductID  *int64   `json:"product_id"`
	Quantity   *int64   `json:"quantity"`
	TotalPrice *float64 `json:"total_price"`
	Status     string   `form:"status"`
}

// CreateOrderItemReq 订单项
type CreateOrderItemDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	UnitPrice int64  `json:"unit_price" binding:"required,min=0"` // 单价，单位：分
}
