package models

import (
	"time"

	cartDomain "eshop-monolith/internal/cart/domain/models"
	"eshop-monolith/pkg/utils"
)

// CartItemPO 购物车项持久化对象
type CartItemPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	CartID    int64     `gorm:"index;not null"`
	ProductID int64     `gorm:"index;not null"`
	Quantity  int       `gorm:"not null;default:1"`
	Price     int64     `gorm:"not null"`
	SKU       string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (CartItemPO) TableName() string { return "cart_items" }

func (po *CartItemPO) ToDomain() *cartDomain.CartItem {
	return &cartDomain.CartItem{
		ID:        po.ID,
		CartID:    po.CartID,
		ProductID: po.ProductID,
		Quantity:  po.Quantity,
		Price:     po.Price,
		SKU:       po.SKU,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func CartItemFromDomain(item *cartDomain.CartItem) *CartItemPO {
	return &CartItemPO{
		ID:        item.ID,
		CartID:    item.CartID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
		SKU:       item.SKU,
		CreatedAt: time.Time(item.CreatedAt),
		UpdatedAt: time.Time(item.UpdatedAt),
	}
}
