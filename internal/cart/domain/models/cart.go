package models

import "eshop-monolith/pkg/utils"

// Cart 购物车模型
type Cart struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	SessionID string     `json:"session_id"`
	Items     []CartItem `json:"items"`
	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (Cart) TableName() string {
	return "carts"
}

// CartItem 购物车项模型
type CartItem struct {
	ID        int64     `json:"id"`
	CartID    int64     `json:"cart_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int64     `json:"price"` // 商品单价，单位：分
	SKU       string    `json:"sku"`   // 商品SKU
	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
