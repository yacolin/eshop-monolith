package models

import (
	"eshop-monolith/pkg/utils"
)

// Cart 购物车模型
type Cart struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"index;not null"`
	SessionID string    `json:"session_id" gorm:"type:varchar(255);index"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt utils.Timestamp `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Cart) TableName() string {
	return "carts"
}

// CartItem 购物车项模型
type CartItem struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CartID    int64     `json:"cart_id" gorm:"index;not null"`
	ProductID int64     `json:"product_id" gorm:"index;not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Price     int64     `json:"price" gorm:"not null"` // 商品单价，单位：分
	SKU       string    `json:"sku" gorm:"type:varchar(100)"` // 商品SKU
	CreatedAt utils.Timestamp `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
