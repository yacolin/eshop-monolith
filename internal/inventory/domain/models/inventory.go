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
//
// 核心业务关系:  可用库存 = quantity - reserved
//   - quantity:  实际物理库存总量
//   - reserved:  已预占（下单未支付）的库存量
//
// 库存生命周期:
//   下单(预留):    reserved += N               — 占坑, 防止超卖
//   支付(扣减):    reserved -= N, quantity -= N — 释放预留 + 实际出库
//   取消(释放):    reserved -= N               — 退坑, 库存归还
type Inventory struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`   // 实际物理库存。卖出一件减一
	Status    string `json:"status"`     // 库存状态: instock / lowstock / outofstock
	Reserved  int    `json:"reserved"`   // 已预订(下单未支付)数量。下单+1, 支付/取消-1
	Threshold int    `json:"threshold"`  // 低库存预警阈值。quantity<=threshold 时自动标为 lowstock

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
