package trade

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      int64          `gorm:"not null;index:idx_order_id" json:"order_id"`
	OrderNo      string         `gorm:"type:varchar(32);not null;index:idx_order_no" json:"order_no"`
	SkuID        int64          `gorm:"not null;index:idx_sku_id" json:"sku_id"`
	ProductID    int64          `gorm:"not null;default:0" json:"product_id"`
	SkuCode      string         `gorm:"type:varchar(100);default:''" json:"sku_code"`
	ProductName  string         `gorm:"type:varchar(200);not null;default:''" json:"product_name"`
	SkuSpec      string         `gorm:"type:json" json:"sku_spec,omitempty"`
	Image        string         `gorm:"type:varchar(512);default:''" json:"image"`
	Price        int64          `gorm:"not null" json:"price"`
	Quantity     int            `gorm:"not null;default:1" json:"quantity"`
	Subtotal     int64          `gorm:"not null;default:0" json:"subtotal"`
	RefundStatus string         `gorm:"type:varchar(20);not null;default:'none'" json:"refund_status"`
	RefundAmount int64          `gorm:"not null;default:0" json:"refund_amount"`
	CreatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (OrderItem) TableName() string { return "tx_order_items" }
