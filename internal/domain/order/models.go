package order

import (
	"time"
)

// Order 订单领域模型
type Order struct {
	ID              int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int64           `json:"user_id" gorm:"index"`
	TotalAmount     int64           `json:"total_amount"`
	Status          string          `json:"status"`
	ShippingAddress ShippingAddress `json:"shipping_address" gorm:"type:json"`
	Items           []OrderItem     `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// OrderItem 订单商品
type OrderItem struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int64     `json:"order_id" gorm:"index"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}

// ShippingAddress 收货地址
type ShippingAddress struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Address  string `json:"address"`
}

// OrderStatus 订单状态
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"
)

// BeforeCreate 创建前钩子
func (o *Order) BeforeCreate() error {
	if o.Status == "" {
		o.Status = OrderStatusPending
	}
	return nil
}

// BeforeCreate 创建前钩子
func (oi *OrderItem) BeforeCreate() error {
	return nil
}
