package models

import (
	"time"

	domain "eshop-monolith/internal/order/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// OrderPO 订单持久化对象
type OrderPO struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	OrderNo        string         `gorm:"type:varchar(32);not null"`
	CustomerID     string         `gorm:"type:varchar(36);not null;index"`
	TotalAmount    int64          `gorm:"type:bigint;not null"`
	DiscountAmount int64          `gorm:"type:bigint;not null;default:0;comment:优惠金额(分)"`
	CouponID       *int64         `gorm:"type:bigint;index;comment:使用的优惠券模板ID"`
	Currency       string         `gorm:"type:varchar(10);default:CNY"`
	Status         string         `gorm:"type:varchar(20);not null;index"`
	CreatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Items          []OrderItemPO  `gorm:"foreignKey:OrderID"`
}

func (OrderPO) TableName() string { return "orders" }

func (po *OrderPO) ToDomain() *domain.Order {
	items := make([]domain.OrderItem, len(po.Items))
	for i, item := range po.Items {
		items[i] = *item.ToDomain()
	}
	return &domain.Order{
		ID:             po.ID,
		OrderNo:        po.OrderNo,
		CustomerID:     po.CustomerID,
		TotalAmount:    po.TotalAmount,
		DiscountAmount: po.DiscountAmount,
		CouponID:       po.CouponID,
		Currency:       po.Currency,
		Status:         po.Status,
		CreatedAt:      utils.Timestamp(po.CreatedAt),
		UpdatedAt:      utils.Timestamp(po.UpdatedAt),
		Items:          items,
	}
}

func OrderFromDomain(o *domain.Order) *OrderPO {
	items := make([]OrderItemPO, len(o.Items))
	for i, item := range o.Items {
		items[i] = *OrderItemFromDomain(&item)
	}
	return &OrderPO{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		CustomerID:     o.CustomerID,
		TotalAmount:    o.TotalAmount,
		DiscountAmount: o.DiscountAmount,
		CouponID:       o.CouponID,
		Currency:       o.Currency,
		Status:         o.Status,
		CreatedAt:      time.Time(o.CreatedAt),
		UpdatedAt:      time.Time(o.UpdatedAt),
		Items:          items,
	}
}
