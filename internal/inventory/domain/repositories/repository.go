package repositories

import (
	"context"

	"gorm.io/gorm"
)

// IinventoryRepository 旧库存仓储桩（保留兼容旧模块引用）
// 实现 order 模块所需的 InventoryForOrder 接口
type IinventoryRepository interface {
	ReserveInventory(ctx context.Context, skuID int64, quantity int) error
	ReleaseInventory(ctx context.Context, skuID int64, quantity int) error
	DeductInventory(ctx context.Context, skuID int64, quantity int) error
	ReserveWithTx(tx *gorm.DB, skuID int64, quantity int) error
	DeductWithTx(tx *gorm.DB, skuID int64, quantity int) error
	ReleaseWithTx(tx *gorm.DB, skuID int64, quantity int) error
	RestoreWithTx(tx *gorm.DB, skuID int64, quantity int) error
}

func NewInventoryRepository(db interface{}) IinventoryRepository {
	return &inventoryRepository{}
}

type inventoryRepository struct{}

func (r *inventoryRepository) ReserveInventory(ctx context.Context, skuID int64, quantity int) error {
	return nil
}
func (r *inventoryRepository) ReleaseInventory(ctx context.Context, skuID int64, quantity int) error {
	return nil
}
func (r *inventoryRepository) DeductInventory(ctx context.Context, skuID int64, quantity int) error {
	return nil
}
func (r *inventoryRepository) ReserveWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	return nil
}
func (r *inventoryRepository) DeductWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	return nil
}
func (r *inventoryRepository) ReleaseWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	return nil
}
func (r *inventoryRepository) RestoreWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	return nil
}
