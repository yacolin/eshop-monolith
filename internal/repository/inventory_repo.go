package repository

import (
	"context"

	"gorm.io/gorm"

	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/domain/shared"
)

// InventoryRepository 库存仓储实现
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return InventoryRepository{db: db}
}

// CreateInventory 创建库存
func (r InventoryRepository) CreateInventory(ctx context.Context, inventory *inventory.Inventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

// FindInventoryByProductID 根据产品ID查询库存
func (r InventoryRepository) FindInventoryByProductID(ctx context.Context, productID int64) (*inventory.Inventory, error) {
	var inv inventory.Inventory
	err := r.db.WithContext(ctx).First(&inv, "product_id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// ReserveInventory 预占库存
func (r InventoryRepository) ReserveInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv inventory.Inventory
		if err := tx.First(&inv, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 检查库存是否足够
		if inv.Quantity-inv.ReservedQuantity < quantity {
			return shared.ErrInsufficientInventory
		}

		// 预占库存
		inv.ReservedQuantity += quantity
		return tx.Save(&inv).Error
	})
}

// ReleaseInventory 释放库存
func (r InventoryRepository) ReleaseInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv inventory.Inventory
		if err := tx.First(&inv, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 检查预占库存是否足够
		if inv.ReservedQuantity < quantity {
			return shared.ErrInsufficientInventory
		}

		// 释放库存
		inv.ReservedQuantity -= quantity
		return tx.Save(&inv).Error
	})
}

// UpdateInventory 更新库存
func (r InventoryRepository) UpdateInventory(ctx context.Context, inventory *inventory.Inventory) error {
	return r.db.WithContext(ctx).Save(inventory).Error
}

// ListInventories 列出所有库存
func (r InventoryRepository) ListInventories(ctx context.Context, page, pageSize int) ([]inventory.Inventory, int64, error) {
	var inventories []inventory.Inventory
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&inventory.Inventory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Order("id ASC").Offset(offset).Limit(pageSize).Find(&inventories).Error
	if err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}