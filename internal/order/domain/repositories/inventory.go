package repositories

import "context"

// InventoryForOrder 订单模块对库存模块的需求接口
// 避免订单服务直接依赖 inventory 模块的具体 repo
type InventoryForOrder interface {
	ReserveInventory(ctx context.Context, productID int64, quantity int) error
	ReleaseInventory(ctx context.Context, productID int64, quantity int) error
	DeductInventory(ctx context.Context, productID int64, quantity int) error
}
