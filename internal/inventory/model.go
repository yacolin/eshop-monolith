package inventory

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	SkuID        int64          `gorm:"not null;uniqueIndex:uk_sku_warehouse" json:"sku_id"`
	WarehouseID  int64          `gorm:"not null;default:0;uniqueIndex:uk_sku_warehouse" json:"warehouse_id"`
	Quantity     int64          `gorm:"not null;default:0" json:"quantity"`
	Reserved     int64          `gorm:"not null;default:0" json:"reserved"`
	InTransit    int64          `gorm:"not null;default:0" json:"in_transit"`
	Threshold    int64          `gorm:"not null;default:10" json:"threshold"`
	MaxThreshold int64          `gorm:"not null;default:999999" json:"max_threshold"`
	Status       string         `gorm:"type:varchar(20);not null;default:'instock'" json:"status"`
	LastCountedAt *time.Time    `gorm:"type:datetime" json:"last_counted_at"`
	LastCountedBy string        `gorm:"type:varchar(50);default:''" json:"last_counted_by"`
	CreatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (Inventory) TableName() string { return "sp_inventories" }

// Available 可售库存（计算字段，对应 VIRTUAL 生成列）
func (i *Inventory) Available() int64 { return i.Quantity - i.Reserved }
