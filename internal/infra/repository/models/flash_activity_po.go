package models

import (
	"time"

	domain "eshop-monolith/internal/flashsale/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

type FlashActivityPO struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	ProductID  int64          `gorm:"type:bigint;not null;index"`
	FlashPrice int64          `gorm:"type:bigint;not null"`
	TotalStock int            `gorm:"type:int;not null"`
	SoldStock  int            `gorm:"type:int;not null;default:0"`
	StartTime  time.Time      `gorm:"type:timestamp;not null"`
	EndTime    time.Time      `gorm:"type:timestamp;not null"`
	Status     string         `gorm:"type:varchar(20);not null;index"`
	CreatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (FlashActivityPO) TableName() string { return "flash_activities" }

func (po *FlashActivityPO) ToDomain() *domain.FlashActivity {
	return &domain.FlashActivity{
		ID:         po.ID,
		ProductID:  po.ProductID,
		FlashPrice: po.FlashPrice,
		TotalStock: po.TotalStock,
		SoldStock:  po.SoldStock,
		StartTime:  utils.Timestamp(po.StartTime),
		EndTime:    utils.Timestamp(po.EndTime),
		Status:     po.Status,
		CreatedAt:  utils.Timestamp(po.CreatedAt),
		UpdatedAt:  utils.Timestamp(po.UpdatedAt),
	}
}

func FlashActivityFromDomain(a *domain.FlashActivity) *FlashActivityPO {
	return &FlashActivityPO{
		ID:         a.ID,
		ProductID:  a.ProductID,
		FlashPrice: a.FlashPrice,
		TotalStock: a.TotalStock,
		SoldStock:  a.SoldStock,
		StartTime:  time.Time(a.StartTime),
		EndTime:    time.Time(a.EndTime),
		Status:     a.Status,
		CreatedAt:  time.Time(a.CreatedAt),
		UpdatedAt:  time.Time(a.UpdatedAt),
	}
}