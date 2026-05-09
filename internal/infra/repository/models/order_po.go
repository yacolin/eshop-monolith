package models

import (
	"time"

	domain "eshop-monolith/internal/order/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// OrderPO 订单持久化对象
type OrderPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	CustomerID  string         `gorm:"type:varchar(36);not null;index"`
	TotalAmount int64          `gorm:"type:bigint;not null"`
	Currency    string         `gorm:"type:varchar(10);default:CNY"`
	Status      string         `gorm:"type:varchar(20);not null;index"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Items       []OrderItemPO  `gorm:"foreignKey:OrderID"`
}

func (OrderPO) TableName() string { return "orders" }

func (po *OrderPO) ToDomain() *domain.Order {
	items := make([]domain.OrderItem, len(po.Items))
	for i, item := range po.Items {
		items[i] = *item.ToDomain()
	}
	return &domain.Order{
		ID:          po.ID,
		CustomerID:  po.CustomerID,
		TotalAmount: po.TotalAmount,
		Currency:    po.Currency,
		Status:      po.Status,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
		Items:       items,
	}
}

func OrderFromDomain(o *domain.Order) *OrderPO {
	items := make([]OrderItemPO, len(o.Items))
	for i, item := range o.Items {
		items[i] = *OrderItemFromDomain(&item)
	}
	return &OrderPO{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		TotalAmount: o.TotalAmount,
		Currency:    o.Currency,
		Status:      o.Status,
		CreatedAt:   time.Time(o.CreatedAt),
		UpdatedAt:   time.Time(o.UpdatedAt),
		Items:       items,
	}
}
