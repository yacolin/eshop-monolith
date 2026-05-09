package models

import "eshop-monolith/pkg/utils"

// InventoryStatus 枚举类型
type InventoryStatus string

// 定义枚举值
const (
	InventoryStatusInStock    = "instock"
	InventoryStatusOutOfStock = "outofstock"
	InventoryStatusLowStock   = "lowstock"
)

// Inventory 库存领域模型
type Inventory struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
	Reserved  int    `json:"reserved"`   // 已预订数量
	Threshold int    `json:"threshold"` // 低库存阈值

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 库存表名
func (Inventory) TableName() string {
	return "inventories"
}

// UpdateStatus 更新库存状态
func (i *Inventory) UpdateStatus() {
	if i.Quantity <= 0 {
		i.Status = InventoryStatusOutOfStock
	} else if i.Quantity <= i.Threshold {
		i.Status = InventoryStatusLowStock
	} else {
		i.Status = InventoryStatusInStock
	}
}
