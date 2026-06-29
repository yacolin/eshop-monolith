package inventory

import "time"

type InventoryLog struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SkuID          int64     `gorm:"not null;index:idx_sku_id" json:"sku_id"`
	WarehouseID    int64     `gorm:"not null;default:0" json:"warehouse_id"`
	ChangeType     string    `gorm:"type:varchar(30);not null;index:idx_change_type" json:"change_type"`
	BeforeQuantity int64     `gorm:"not null;default:0" json:"before_quantity"`
	AfterQuantity  int64     `gorm:"not null;default:0" json:"after_quantity"`
	BeforeReserved int64     `gorm:"not null;default:0" json:"before_reserved"`
	AfterReserved  int64     `gorm:"not null;default:0" json:"after_reserved"`
	ChangeAmount   int64     `gorm:"not null;default:0" json:"change_amount"`
	ReferenceID    string    `gorm:"type:varchar(64);default:'';index:idx_reference_id" json:"reference_id"`
	Operator       string    `gorm:"type:varchar(50);default:''" json:"operator"`
	Note           string    `gorm:"type:varchar(500);default:''" json:"note"`
	CreatedAt      time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
}

func (InventoryLog) TableName() string { return "sp_inventory_logs" }
