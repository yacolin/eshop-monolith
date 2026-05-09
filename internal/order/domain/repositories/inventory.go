package repositories

import (
	"context"

	"gorm.io/gorm"
)

// InventoryForOrder 订单模块对库存模块的需求接口
// 避免订单服务直接依赖 inventory 模块的具体 repo
type InventoryForOrder interface {
	ReserveInventory(ctx context.Context, productID int64, quantity int) error
	ReleaseInventory(ctx context.Context, productID int64, quantity int) error
	DeductInventory(ctx context.Context, productID int64, quantity int) error

	// ReserveWithTx 在已有事务内预占库存
	ReserveWithTx(tx *gorm.DB, productID int64, quantity int) error
	// DeductWithTx 在已有事务内扣减库存
	DeductWithTx(tx *gorm.DB, productID int64, quantity int) error
	// ReleaseWithTx 在已有事务内释放库存
	ReleaseWithTx(tx *gorm.DB, productID int64, quantity int) error
}
