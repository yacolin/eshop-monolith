package models

import (
	"time"

	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// ProductPO 产品持久化对象
type ProductPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Price       int64          `gorm:"type:bigint;not null"`
	SKU         string         `gorm:"type:varchar(100);uniqueIndex;not null"`
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
		Price:       po.Price,
		SKU:         po.SKU,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func ProductFromDomain(p *domain.Product) *ProductPO {
	return &ProductPO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		SKU:         p.SKU,
		CreatedAt:   time.Time(p.CreatedAt),
		UpdatedAt:   time.Time(p.UpdatedAt),
	}
}
