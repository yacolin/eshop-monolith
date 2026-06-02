package models

import (
	"time"

	domain "eshop-monolith/internal/flashsale/domain/models"
	"eshop-monolith/pkg/utils"
)

type FlashOrderPO struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	ActivityID  int64     `gorm:"type:bigint;not null;index"`
	UserID      int64     `gorm:"type:bigint;not null;index"`
	ProductID   int64     `gorm:"type:bigint;not null"`
	Quantity    int       `gorm:"type:int;not null;default:1"`
	FlashPrice  int64     `gorm:"type:bigint;not null"`
	TotalAmount int64     `gorm:"type:bigint;not null"`
	Status      string    `gorm:"type:varchar(20);not null;index"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
}

func (FlashOrderPO) TableName() string { return "flash_orders" }

func (po *FlashOrderPO) ToDomain() *domain.FlashOrder {
	return &domain.FlashOrder{
		ID:          po.ID,
		ActivityID:  po.ActivityID,
		UserID:      po.UserID,
		ProductID:   po.ProductID,
		Quantity:    po.Quantity,
		FlashPrice:  po.FlashPrice,
		TotalAmount: po.TotalAmount,
		Status:      po.Status,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func FlashOrderFromDomain(o *domain.FlashOrder) *FlashOrderPO {
	return &FlashOrderPO{
		ID:          o.ID,
		ActivityID:  o.ActivityID,
		UserID:      o.UserID,
		ProductID:   o.ProductID,
		Quantity:    o.Quantity,
		FlashPrice:  o.FlashPrice,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   time.Time(o.CreatedAt),
		UpdatedAt:   time.Time(o.UpdatedAt),
	}
}