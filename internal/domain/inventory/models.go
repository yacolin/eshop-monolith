package inventory

import (
	"time"
)

// Inventory 库存领域模型
type Inventory struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID        int64     `json:"product_id" gorm:"uniqueIndex"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BeforeCreate 创建前钩子
func (i *Inventory) BeforeCreate() error {
	return nil
}
