package models

import (
	"time"

	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// InventoryPO 库存持久化对象
type InventoryPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	SkuID int64          `gorm:"uniqueIndex"`
	Quantity  int            `gorm:"not null;default:0"`
	Status    string         `gorm:"type:varchar(20);not null;default:'instock'"`
	Reserved  int            `gorm:"not null;default:0"`
	Threshold int            `gorm:"not null;default:10"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (InventoryPO) TableName() string { return "inventories" }

// BeforeCreate GORM 钩子
func (po *InventoryPO) BeforeCreate(tx *gorm.DB) error {
	if po.Status == "" {
		if po.Quantity <= po.Threshold {
			if po.Quantity <= 0 {
				po.Status = string(domain.InventoryStatusOutOfStock)
			} else {
				po.Status = string(domain.InventoryStatusLowStock)
			}
		} else {
			po.Status = string(domain.InventoryStatusInStock)
		}
	}
	return nil
}

func (po *InventoryPO) ToDomain() *domain.Inventory {
	return &domain.Inventory{
		ID:        po.ID,
		SkuID: po.SkuID,
		Quantity:  po.Quantity,
		Status:    po.Status,
		Reserved:  po.Reserved,
		Threshold: po.Threshold,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func InventoryFromDomain(inv *domain.Inventory) *InventoryPO {
	return &InventoryPO{
		ID:        inv.ID,
		SkuID: inv.SkuID,
		Quantity:  inv.Quantity,
		Status:    inv.Status,
		Reserved:  inv.Reserved,
		Threshold: inv.Threshold,
		CreatedAt: time.Time(inv.CreatedAt),
		UpdatedAt: time.Time(inv.UpdatedAt),
	}
}
