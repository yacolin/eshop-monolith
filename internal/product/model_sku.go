package product

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

type SKU struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID     int64          `gorm:"not null;index:idx_product_id" json:"product_id"`
	SkuCode       string         `gorm:"type:varchar(100);not null;uniqueIndex:uk_sku_code" json:"sku_code"`
	Barcode       string         `gorm:"type:varchar(50);default:'';uniqueIndex:uk_barcode" json:"barcode"`
	Spec          string         `gorm:"type:json;not null" json:"spec"`
	Price         int64          `gorm:"not null" json:"price"`
	MarketPrice   int64          `gorm:"not null;default:0" json:"market_price"`
	CostPrice     int64          `gorm:"not null;default:0" json:"cost_price"`
	Weight        float64        `gorm:"type:decimal(10,2);not null;default:0" json:"weight"`
	Volume        float64        `gorm:"type:decimal(10,2);not null;default:0" json:"volume"`
	Length        float64        `gorm:"type:decimal(10,2);not null;default:0" json:"length"`
	Width         float64        `gorm:"type:decimal(10,2);not null;default:0" json:"width"`
	Height        float64        `gorm:"type:decimal(10,2);not null;default:0" json:"height"`
	MinPurchaseQty int           `gorm:"not null;default:1" json:"min_purchase_qty"`
	MaxPurchaseQty int           `gorm:"not null;default:0" json:"max_purchase_qty"`
	Image         string         `gorm:"type:varchar(512);default:''" json:"image"`
	Status        int8           `gorm:"not null;default:1" json:"status"`
	CreatedAt     utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (SKU) TableName() string { return "sp_skus" }
