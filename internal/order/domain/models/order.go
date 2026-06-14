package models

import "eshop-monolith/pkg/utils"

// OrderStatus 枚举类型
type OrderStatus string

// 定义枚举值
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"
)

// 金额统一以「分」为单位存储，避免浮点精度问题（1 元 = 100 分）

// Order 订单
type Order struct {
	ID          int64          `json:"id"`
	OrderNo     string         `json:"order_no"`     // 订单号，全局唯一
	CustomerID  string         `json:"customer_id"`
	TotalAmount int64          `json:"total_amount"` // 订单总金额，单位：分
	Currency    string         `json:"currency"`
	Status      string         `json:"status"`
	CreatedAt   utils.Timestamp `json:"created_at"`
	UpdatedAt   utils.Timestamp `json:"updated_at"`

	Items []OrderItem `json:"items,omitempty"`
}

// TableName 表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项
// @Description 订单内的商品项
type OrderItem struct {
	ID        int64  `json:"id"`
	OrderID   int64  `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"` // 单价，单位：分
	Amount    int64  `json:"amount"`     // 单项小计，单位：分 = UnitPrice * Quantity
}

// TableName 表名
func (OrderItem) TableName() string {
	return "order_items"
}
