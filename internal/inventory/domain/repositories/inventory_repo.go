package repositories

import (
	"context"

	"eshop-monolith/internal/infra/domain/shared"
	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/inventory/api/dto"
	invModels "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

// IinventoryRepository 库存仓储接口
type IinventoryRepository interface {
	// CreateInventory 创建库存
	CreateInventory(ctx context.Context, inv *invModels.Inventory) error
	// FindInventoryByProductID 根据产品ID查询库存
	FindInventoryByProductID(ctx context.Context, productID int64) (*invModels.Inventory, error)
	// ReserveInventory 预占库存
	ReserveInventory(ctx context.Context, productID int64, quantity int) error
	// ReleaseInventory 释放库存
	ReleaseInventory(ctx context.Context, productID int64, quantity int) error
	// DeductInventory 扣减库存
	DeductInventory(ctx context.Context, productID int64, quantity int) error
	// UpdateInventory 更新库存
	UpdateInventory(ctx context.Context, inv *invModels.Inventory) error
	// ListInventories 列出所有库存
	ListInventories(ctx context.Context, query dto.InventoryListQuery, offset, limit int) ([]invModels.Inventory, error)
	CountInventories(ctx context.Context, query dto.InventoryListQuery) (int64, error)
}

// InventoryRepository 库存仓储实现
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *gorm.DB) IinventoryRepository {
	return &InventoryRepository{db: db}
}

// CreateInventory 创建库存
func (r *InventoryRepository) CreateInventory(ctx context.Context, inv *invModels.Inventory) error {
	po := models.InventoryFromDomain(inv)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	inv.ID = po.ID
	return nil
}

// GetInventoryByID 根据ID查询库存
func (r *InventoryRepository) GetInventoryByID(ctx context.Context, id string) (*invModels.Inventory, error) {
	var po models.InventoryPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindInventoryByProductID 根据产品ID查询库存
func (r *InventoryRepository) FindInventoryByProductID(ctx context.Context, productID int64) (*invModels.Inventory, error) {
	var po models.InventoryPO
	err := r.db.WithContext(ctx).First(&po, "product_id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ReserveInventory 预占库存
func (r *InventoryRepository) ReserveInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 检查库存是否足够
		if po.Quantity-po.Reserved < quantity {
			return shared.ErrInsufficientInventory
		}

		// 预占库存
		po.Reserved += quantity
		return tx.Save(&po).Error
	})
}

// ReleaseInventory 释放库存
func (r *InventoryRepository) ReleaseInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 检查预占库存是否足够
		if po.Reserved < quantity {
			return shared.ErrInsufficientInventory
		}

		// 释放库存
		po.Reserved -= quantity
		return tx.Save(&po).Error
	})
}

// ReserveWithTx 在已有事务内预占库存（不自行提交/回滚）
func (r *InventoryRepository) ReserveWithTx(tx *gorm.DB, productID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
		return err
	}
	if po.Quantity-po.Reserved < quantity {
		return shared.ErrInsufficientInventory
	}
	po.Reserved += quantity
	return tx.Save(&po).Error
}

// DeductWithTx 在已有事务内扣减库存（不自行提交/回滚）
func (r *InventoryRepository) DeductWithTx(tx *gorm.DB, productID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
		return err
	}
	if po.Reserved < quantity {
		return shared.ErrInsufficientInventory
	}
	po.Quantity -= quantity
	po.Reserved -= quantity
	if po.Quantity <= 0 {
		po.Status = string(invModels.InventoryStatusOutOfStock)
	} else if po.Quantity <= po.Threshold {
		po.Status = string(invModels.InventoryStatusLowStock)
	} else {
		po.Status = string(invModels.InventoryStatusInStock)
	}
	return tx.Save(&po).Error
}

// ReleaseWithTx 在已有事务内释放库存（不自行提交/回滚）
func (r *InventoryRepository) ReleaseWithTx(tx *gorm.DB, productID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
		return err
	}
	if po.Reserved < quantity {
		return shared.ErrInsufficientInventory
	}
	po.Reserved -= quantity
	return tx.Save(&po).Error
}

// UpdateInventory 更新库存
func (r *InventoryRepository) UpdateInventory(ctx context.Context, inv *invModels.Inventory) error {
	// 更新库存状态
	inv.UpdateStatus()
	po := models.InventoryFromDomain(inv)
	return r.db.WithContext(ctx).Save(po).Error
}

// DeductInventory 扣减库存（确认订单时使用）
func (r *InventoryRepository) DeductInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 检查预占库存是否足够
		if po.Reserved < quantity {
			return shared.ErrInsufficientInventory
		}

		// 扣减库存：减少实际库存和预占库存
		po.Quantity -= quantity
		po.Reserved -= quantity
		// UpdateStatus handled by BeforeCreate hook or manually
		if po.Quantity <= 0 {
			po.Status = string(invModels.InventoryStatusOutOfStock)
		} else if po.Quantity <= po.Threshold {
			po.Status = string(invModels.InventoryStatusLowStock)
		} else {
			po.Status = string(invModels.InventoryStatusInStock)
		}
		return tx.Save(&po).Error
	})
}

// IncreaseInventory 增加库存（取消订单或入库时使用）
func (r *InventoryRepository) IncreaseInventory(ctx context.Context, productID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "product_id = ?", productID).Error; err != nil {
			return err
		}

		// 增加实际库存
		po.Quantity += quantity
		if po.Quantity <= 0 {
			po.Status = string(invModels.InventoryStatusOutOfStock)
		} else if po.Quantity <= po.Threshold {
			po.Status = string(invModels.InventoryStatusLowStock)
		} else {
			po.Status = string(invModels.InventoryStatusInStock)
		}
		return tx.Save(&po).Error
	})
}

// FindLowStockInventories 查询低库存商品
func (r *InventoryRepository) FindLowStockInventories(ctx context.Context) ([]invModels.Inventory, error) {
	var pos []models.InventoryPO
	err := r.db.WithContext(ctx).
		Where("status = ? OR status = ?", invModels.InventoryStatusLowStock, invModels.InventoryStatusOutOfStock).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	inventories := make([]invModels.Inventory, len(pos))
	for i, po := range pos {
		inventories[i] = *po.ToDomain()
	}
	return inventories, nil
}

// FindInventoryByID 根据ID查询库存
func (r *InventoryRepository) FindInventoryByID(ctx context.Context, id int64) (*invModels.Inventory, error) {
	var po models.InventoryPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// DeleteInventory 删除库存
func (r *InventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.InventoryPO{}, id).Error
}

// ListInventories 列出所有库存
func (r *InventoryRepository) ListInventories(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]invModels.Inventory, error) {
	var pos []models.InventoryPO
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	// 执行查询
	err := db.Offset(offset).Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	inventories := make([]invModels.Inventory, len(pos))
	for i, po := range pos {
		inventories[i] = *po.ToDomain()
	}
	return inventories, nil
}

// CountInventories 查询库存总数
func (r *InventoryRepository) CountInventories(ctx context.Context, q dto.InventoryListQuery) (int64, error) {
	var total int64
	db := r.applyQueryConditions(ctx, q)

	// 执行统计（不需要排序）
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r *InventoryRepository) applyQueryConditions(ctx context.Context, q dto.InventoryListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.InventoryPO{})
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
func (r *InventoryRepository) applyOrder(db *gorm.DB, q dto.InventoryListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}
