package dto

import (
	"eshop-monolith/pkg/query"
)

// AddToCartDTO 添加商品到购物车请求
type AddToCartDTO struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	SKU       string `json:"sku"`
}

// UpdateCartItemDTO 更新购物车项请求
type UpdateCartItemDTO struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartItemResponse 购物车项响应
type CartItemResponse struct {
	ID          int64  `json:"id"`
	ProductID   int64  `json:"product_id"`
	Quantity    int    `json:"quantity"`
	Price       int64  `json:"price"` // 商品单价，单位：分
	SKU         string `json:"sku"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"` // 库存状态
}

// CartResponse 购物车响应
type CartResponse struct {
	ID         int64              `json:"id"`
	UserID     int64              `json:"user_id"`
	Items      []CartItemResponse `json:"items"`
	TotalItems int                `json:"total_items"`
	TotalPrice int64              `json:"total_price"` // 总价，单位：分
}

// CartListQuery 购物车列表查询参数
type CartListQuery struct {
	query.Pagination
	UserID    int64  `form:"user_id"`
	SessionID string `form:"session_id"`
}

// CartListResult 购物车列表结果
type CartListResult struct {
	Total int64           `json:"total"`
	List  []CartResponse  `json:"list"`
}
