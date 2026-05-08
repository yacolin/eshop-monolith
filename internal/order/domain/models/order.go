package models

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

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
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID  string `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	TotalAmount int64  `gorm:"type:bigint;not null" json:"total_amount"` // 订单总金额，单位：分
	Currency    string `gorm:"type:varchar(10);default:CNY" json:"currency"`
	Status      string `gorm:"type:varchar(20);not null;index" json:"status"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// TableName 表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项
type OrderItem struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int64  `gorm:"type:bigint;not null;index" json:"order_id"`
	ProductID string `gorm:"type:varchar(36);not null" json:"product_id"`
	Quantity  int    `gorm:"not null" json:"quantity"`
	UnitPrice int64  `gorm:"type:bigint;not null" json:"unit_price"` // 单价，单位：分
	Amount    int64  `gorm:"type:bigint;not null" json:"amount"`     // 单项小计，单位：分 = UnitPrice * Quantity
}

// TableName 表名
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate 创建前钩子
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.Status == "" {
		o.Status = OrderStatusPending
	}
	return nil
}

// BeforeCreate 创建前钩子
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	return nil
}
