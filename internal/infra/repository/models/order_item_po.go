package models

import (
	"time"

	domain "eshop-monolith/internal/order/domain/models"
	"gorm.io/gorm"
)

// OrderItemPO 订单项持久化对象
type OrderItemPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	OrderID   int64          `gorm:"type:bigint;not null;index"`
	SkuID     int64          `gorm:"index;default:0"`
	ProductID string         `gorm:"type:varchar(36);not null"`
	Quantity  int            `gorm:"not null"`
	UnitPrice int64          `gorm:"type:bigint;not null"`
	Amount    int64          `gorm:"type:bigint;not null"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (OrderItemPO) TableName() string { return "order_items" }

func (po *OrderItemPO) ToDomain() *domain.OrderItem {
	return &domain.OrderItem{
		ID:        po.ID,
		OrderID:   po.OrderID,
		SkuID:     po.SkuID,
		ProductID: po.ProductID,
		Quantity:  po.Quantity,
		UnitPrice: po.UnitPrice,
		Amount:    po.Amount,
	}
}

func OrderItemFromDomain(item *domain.OrderItem) *OrderItemPO {
	return &OrderItemPO{
		ID:        item.ID,
		OrderID:   item.OrderID,
		SkuID:     item.SkuID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		Amount:    item.Amount,
	}
}
