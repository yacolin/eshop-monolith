package models

import (
	"time"

	cartDomain "eshop-monolith/internal/cart/domain/models"
	"eshop-monolith/pkg/utils"
)

// CartPO 购物车持久化对象
type CartPO struct {
	ID        int64         `gorm:"primaryKey;autoIncrement"`
	UserID    int64         `gorm:"index;not null"`
	SessionID string        `gorm:"type:varchar(255);index"`
	Items     []CartItemPO  `gorm:"foreignKey:CartID"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
}

func (CartPO) TableName() string { return "carts" }

func (po *CartPO) ToDomain() *cartDomain.Cart {
	items := make([]cartDomain.CartItem, len(po.Items))
	for i, item := range po.Items {
		items[i] = *item.ToDomain()
	}
	return &cartDomain.Cart{
		ID:        po.ID,
		UserID:    po.UserID,
		SessionID: po.SessionID,
		Items:     items,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func CartFromDomain(c *cartDomain.Cart) *CartPO {
	items := make([]CartItemPO, len(c.Items))
	for i, item := range c.Items {
		items[i] = *CartItemFromDomain(&item)
	}
	return &CartPO{
		ID:        c.ID,
		UserID:    c.UserID,
		SessionID: c.SessionID,
		Items:     items,
		CreatedAt: time.Time(c.CreatedAt),
		UpdatedAt: time.Time(c.UpdatedAt),
	}
}
