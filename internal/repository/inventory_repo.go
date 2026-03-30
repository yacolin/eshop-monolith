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
		if inv.Quantity-inv.Reserved < quantity {
			return shared.ErrInsufficientInventory
		}

		// 预占库存
		inv.Reserved += quantity
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
		if inv.Reserved < quantity {
			return shared.ErrInsufficientInventory
		}

		// 释放库存
		inv.Reserved -= quantity
		return tx.Save(&inv).Error
	})
}

// UpdateInventory 更新库存
func (r InventoryRepository) UpdateInventory(ctx context.Context, inventory *inventory.Inventory) error {
	// 更新库存状态
	inventory.UpdateStatus()
	return r.db.WithContext(ctx).Save(inventory).Error
}

// DeductInventory 扣减库存（确认订单时使用）
func (r InventoryRepository) DeductInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv inventory.Inventory
		if err := tx.First(&inv, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 检查预占库存是否足够
		if inv.Reserved < quantity {
			return shared.ErrInsufficientInventory
		}

		// 扣减库存：减少实际库存和预占库存
		inv.Quantity -= quantity
		inv.Reserved -= quantity
		inv.UpdateStatus()
		return tx.Save(&inv).Error
	})
}

// IncreaseInventory 增加库存（取消订单或入库时使用）
func (r InventoryRepository) IncreaseInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv inventory.Inventory
		if err := tx.First(&inv, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 增加实际库存
		inv.Quantity += quantity
		inv.UpdateStatus()
		return tx.Save(&inv).Error
	})
}

// FindLowStockInventories 查询低库存商品
func (r InventoryRepository) FindLowStockInventories(ctx context.Context) ([]inventory.Inventory, error) {
	var inventories []inventory.Inventory
	err := r.db.WithContext(ctx).
		Where("status = ? OR status = ?", inventory.InventoryStatusLowStock, inventory.InventoryStatusOutOfStock).
		Find(&inventories).Error
	if err != nil {
		return nil, err
	}
	return inventories, nil
}

// FindInventoryByID 根据ID查询库存
func (r InventoryRepository) FindInventoryByID(ctx context.Context, id int64) (*inventory.Inventory, error) {
	var inv inventory.Inventory
	err := r.db.WithContext(ctx).First(&inv, id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// DeleteInventory 删除库存
func (r InventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&inventory.Inventory{}, id).Error
}

// ListInventories 列出所有库存
func (r InventoryRepository) ListInventories(ctx context.Context, q inventory.InventoryListQuery, offset, limit int) ([]inventory.Inventory, error) {
	var inventories []inventory.Inventory
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	// 执行查询
	err := db.Offset(offset).Limit(limit).Find(&inventories).Error
	if err != nil {
		return nil, err
	}

	return inventories, nil
}

// CountInventories 查询库存总数
func (r InventoryRepository) CountInventories(ctx context.Context, q inventory.InventoryListQuery) (int64, error) {
	var total int64
	db := r.applyQueryConditions(ctx, q)

	// 执行统计（不需要排序）
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r InventoryRepository) applyQueryConditions(ctx context.Context, q inventory.InventoryListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&inventory.Inventory{})
	db = db.Joins("JOIN products ON inventories.product_id = products.id")

	if q.ProductName != "" {
		db = db.Where("products.name LIKE ?", "%"+q.ProductName+"%")
	}
	if q.SKU != "" {
		db = db.Where("products.sku = ?", q.SKU)
	}
	if q.Status != "" {
		db = db.Where("inventories.status = ?", q.Status)
	}
	if q.LowStock != nil && *q.LowStock {
		db = db.Where("inventories.quantity <= inventories.threshold AND inventories.quantity > 0")
	}
	return db
}

// applyOrder 应用排序
func (r InventoryRepository) applyOrder(db *gorm.DB, q inventory.InventoryListQuery) *gorm.DB {
	order := "id asc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "asc"
		}
		order = q.SortBy + " " + ord
	}
	return db.Order(order)
}
