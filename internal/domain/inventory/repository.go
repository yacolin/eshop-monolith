package inventory

import (
	"context"
)

// Repository 库存仓储接口
type Repository interface {
	// CreateInventory 创建库存
	CreateInventory(ctx context.Context, inventory *Inventory) error
	// FindInventoryByProductID 根据产品ID查询库存
	FindInventoryByProductID(ctx context.Context, productID int64) (*Inventory, error)
	// ReserveInventory 预占库存
	ReserveInventory(ctx context.Context, productID int64, quantity int) error
	// ReleaseInventory 释放库存
	ReleaseInventory(ctx context.Context, productID int64, quantity int) error
	// UpdateInventory 更新库存
	UpdateInventory(ctx context.Context, inventory *Inventory) error
	// ListInventories 列出所有库存
	ListInventories(ctx context.Context, query InventoryListQuery, offset, limit int) ([]Inventory, error)
	CountInventories(ctx context.Context, query InventoryListQuery) (int64, error)
}
