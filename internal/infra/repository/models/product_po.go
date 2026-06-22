package models

import (
	"time"

	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// ProductPO 产品持久化对象 (SPU)
type ProductPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	MinPrice    int64          `gorm:"type:bigint;default:0"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ProductPO) TableName() string { return "products" }

func (po *ProductPO) ToDomain() *domain.Product {
	return &domain.Product{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		MinPrice:    po.MinPrice,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func ProductFromDomain(p *domain.Product) *ProductPO {
	return &ProductPO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		MinPrice:    p.MinPrice,
		CreatedAt:   time.Time(p.CreatedAt),
		UpdatedAt:   time.Time(p.UpdatedAt),
	}
}
