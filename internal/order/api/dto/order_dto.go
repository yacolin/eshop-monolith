package dto

import (
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// OrderResponse 订单响应
type OrderResponse struct {
	ID             int64              `json:"id"`
	OrderNo        string             `json:"order_no"`
	CustomerID     string             `json:"customer_id"`
	TotalAmount    int64              `json:"total_amount"`    // 订单总金额（已扣优惠），单位：分
	DiscountAmount int64              `json:"discount_amount"` // 优惠金额，单位：分
	CouponID       *int64             `json:"coupon_id"`       // 使用的优惠券模板ID
	Currency       string             `json:"currency"`
	Status         string             `json:"status"`
	CreatedAt      utils.Timestamp    `json:"created_at"`
	UpdatedAt      utils.Timestamp    `json:"updated_at"`
	Items          []OrderItemResponse `json:"items,omitempty"`
}

// OrderItemResponse 订单项响应
type OrderItemResponse struct {
	ID        int64  `json:"id"`
	OrderID   int64  `json:"order_id"`
	OrderNo   string `json:"order_no"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"` // 单价，单位：分
	Amount    int64  `json:"amount"`     // 单项小计，单位：分
}

// OrderListQuery 支持通过 query string 进行过滤、排序和分页
type OrderListQuery struct {
	query.Pagination
	CustomerID *int64   `form:"customer_id"`       // 用户ID过滤
	OrderNo    string   `form:"order_no"`          // 订单号搜索
	Status     string   `form:"status"`            // 订单状态过滤
	MinPrice   *float64 `form:"min_price"`         // 价格区间下限
	MaxPrice   *float64 `form:"max_price"`         // 价格区间上限
	SortBy     string   `form:"sort_by"`           // 排序字段，例如 total_price, created_at
	Order      string   `form:"order,default=asc"` // asc or desc
}

// OrderListResult 订单列表结果
type OrderListResult struct {
	Total int64          `json:"total"`
	List  []OrderResponse `json:"list"`
}

// OrderItemListResult 订单项列表结果
type OrderItemListResult struct {
	Total int64              `json:"total"`
	List  []OrderItemResponse `json:"list"`
}

// OrderItemListQuery 全量订单项查询参数
type OrderItemListQuery struct {
	query.Pagination
	OrderNo string `form:"order_no"`           // 按订单号筛选
	SortBy  string `form:"sort_by"`            // 排序字段，例如 id, order_id, amount
	Order   string `form:"order,default=asc"`  // asc or desc
}

// UserOrderListQuery 用户订单列表查询参数
type UserOrderListQuery struct {
	query.Pagination
}

// CreateOrderDTO 创建订单请求
// @Description 创建订单的请求体
type CreateOrderDTO struct {
	CustomerID   string               `json:"customer_id" binding:"required"`
	Currency     string               `json:"currency"`      // 可选，默认 CNY
	UserCouponID *int64               `json:"user_coupon_id"` // 可选，使用的用户优惠券ID
	Items        []CreateOrderItemDTO `json:"items" binding:"required,min=1,dive"`
}

// UpdateOrderDTO 更新订单请求
type UpdateOrderDTO struct {
	UserID     *int64   `json:"user_id"`
	ProductID  *int64   `json:"product_id"`
	Quantity   *int64   `json:"quantity"`
	TotalPrice *float64 `json:"total_price"`
	Status     string   `form:"status"`
}

// UpdateOrderStatusDTO 更新订单状态请求
type UpdateOrderStatusDTO struct {
	Status string `json:"status" binding:"required"`
}

// CreateOrderItemReq 订单项
type CreateOrderItemDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	UnitPrice int64  `json:"unit_price" binding:"required,min=0"` // 单价，单位：分
}
