package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/pkg/errcode"
	invModels "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// recalcStatus 根据 可用库存=Quantity-Reserved 重新计算状态
func recalcStatus(po *models.InventoryPO) {
	available := po.Quantity - po.Reserved
	switch {
	case available <= 0:
		po.Status = string(invModels.InventoryStatusOutOfStock)
	case available <= po.Threshold:
		po.Status = string(invModels.InventoryStatusLowStock)
	default:
		po.Status = string(invModels.InventoryStatusInStock)
	}
}

// InventoryEnrichedRow 库存 enriched 查询结果行，含 SKU 和产品信息（单次 LEFT JOIN 产出）
type InventoryEnrichedRow struct {
	ID          int64           `gorm:"column:id"`
	SkuID       int64           `gorm:"column:sku_id"`
	SkuName     string          `gorm:"column:sku_name"`
	SkuCode     string          `gorm:"column:sku_code"`
	ProductID   int64           `gorm:"column:product_id"`
	ProductName string          `gorm:"column:product_name"`
	Quantity    int             `gorm:"column:quantity"`
	Status      string          `gorm:"column:status"`
	Reserved    int             `gorm:"column:reserved"`
	Threshold   int             `gorm:"column:threshold"`
	CreatedAt   utils.Timestamp `gorm:"column:created_at"`
	UpdatedAt   utils.Timestamp `gorm:"column:updated_at"`
}

// IinventoryRepository 库存仓储接口
type IinventoryRepository interface {
	// CreateInventory 创建库存
	CreateInventory(ctx context.Context, inv *invModels.Inventory) error
	// BatchCreateInventory 批量创建库存
	BatchCreateInventory(ctx context.Context, invs []*invModels.Inventory) error
	// FindInventoryBySkuID 根据SKU ID查询库存
	FindInventoryBySkuID(ctx context.Context, skuID int64) (*invModels.Inventory, error)
	// ReserveInventory 预占库存
	ReserveInventory(ctx context.Context, skuID int64, quantity int) error
	// ReleaseInventory 释放库存
	ReleaseInventory(ctx context.Context, skuID int64, quantity int) error
	// DeductInventory 扣减库存
	DeductInventory(ctx context.Context, skuID int64, quantity int) error
	// UpdateInventory 更新库存
	UpdateInventory(ctx context.Context, inv *invModels.Inventory) error
	// ListInventories 列出所有库存
	ListInventories(ctx context.Context, query dto.InventoryListQuery, offset, limit int) ([]invModels.Inventory, error)
	CountInventories(ctx context.Context, query dto.InventoryListQuery) (int64, error)
	// FindInventoriesBySkuIDs 根据 SKU ID 列表批量查询库存，返回 map[skuID]Inventory
	FindInventoriesBySkuIDs(ctx context.Context, skuIDs []int64) (map[int64]*invModels.Inventory, error)
	// ListInventoriesEnriched 列出库存（含 SKU 名称和产品名称，单次 LEFT JOIN 查询）
	ListInventoriesEnriched(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]InventoryEnrichedRow, error)
	// CountInventoriesEnriched 统计 enriched 库存数量（使用同样的 JOIN 条件）
	CountInventoriesEnriched(ctx context.Context, q dto.InventoryListQuery) (int64, error)
}

// InventoryRepository 库存仓储实现
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *gorm.DB) IinventoryRepository {
	return &InventoryRepository{db: db}
}

// BatchCreateInventory 批量创建库存
func (r *InventoryRepository) BatchCreateInventory(ctx context.Context, invs []*invModels.Inventory) error {
	pos := make([]models.InventoryPO, len(invs))
	for i, inv := range invs {
		inv.UpdateStatus()
		pos[i] = *models.InventoryFromDomain(inv)
	}
	return r.db.WithContext(ctx).CreateInBatches(pos, 100).Error
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

// FindInventoryBySkuID 根据SKU ID查询库存
func (r *InventoryRepository) FindInventoryBySkuID(ctx context.Context, skuID int64) (*invModels.Inventory, error) {
	var po models.InventoryPO
	err := r.db.WithContext(ctx).First(&po, "sku_id = ?", skuID).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ReserveInventory 预占库存
func (r *InventoryRepository) ReserveInventory(ctx context.Context, skuID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "sku_id = ?", skuID).Error; err != nil {
			return err
		}

		// 检查库存是否足够
		if po.Quantity-po.Reserved < quantity {
			return errcode.ErrInsufficientInventory
		}

		// 预占库存
		po.Reserved += quantity
		recalcStatus(&po)
		return tx.Save(&po).Error
	})
}

// ReleaseInventory 释放库存
func (r *InventoryRepository) ReleaseInventory(ctx context.Context, skuID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "sku_id = ?", skuID).Error; err != nil {
			return err
		}

		// 检查预占库存是否足够
		if po.Reserved < quantity {
			return errcode.ErrInsufficientInventory
		}

		// 释放库存
		po.Reserved -= quantity
		recalcStatus(&po)
		return tx.Save(&po).Error
	})
}

// ReserveWithTx 在已有事务内预占库存（不自行提交/回滚）
func (r *InventoryRepository) ReserveWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "sku_id = ?", skuID).Error; err != nil {
		return err
	}
	if po.Quantity-po.Reserved < quantity {
		return errcode.ErrInsufficientInventory
	}
	po.Reserved += quantity
	recalcStatus(&po)
	return tx.Save(&po).Error
}

// DeductWithTx 在已有事务内扣减库存（不自行提交/回滚）
func (r *InventoryRepository) DeductWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "sku_id = ?", skuID).Error; err != nil {
		return err
	}
	if po.Reserved < quantity {
		return errcode.ErrInsufficientInventory
	}
	po.Quantity -= quantity
	po.Reserved -= quantity
	recalcStatus(&po)
	return tx.Save(&po).Error
}

// ReleaseWithTx 在已有事务内释放库存（不自行提交/回滚）
func (r *InventoryRepository) ReleaseWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "sku_id = ?", skuID).Error; err != nil {
		return err
	}
	if po.Reserved < quantity {
		return errcode.ErrInsufficientInventory
	}
	po.Reserved -= quantity
	recalcStatus(&po)
	return tx.Save(&po).Error
}

// RestoreWithTx 在已有事务内恢复已扣减库存（用于支付后退款场景）
func (r *InventoryRepository) RestoreWithTx(tx *gorm.DB, skuID int64, quantity int) error {
	var po models.InventoryPO
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "sku_id = ?", skuID).Error; err != nil {
		return err
	}
	po.Quantity += quantity
	recalcStatus(&po)
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
func (r *InventoryRepository) DeductInventory(ctx context.Context, skuID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "sku_id = ?", skuID).Error; err != nil {
			return err
		}

		// 检查预占库存是否足够
		if po.Reserved < quantity {
			return errcode.ErrInsufficientInventory
		}

		// 扣减库存：减少实际库存和预占库存
		po.Quantity -= quantity
		po.Reserved -= quantity
		recalcStatus(&po)
		return tx.Save(&po).Error
	})
}

// IncreaseInventory 增加库存（取消订单或入库时使用）
func (r *InventoryRepository) IncreaseInventory(ctx context.Context, skuID int64, quantity int) error {
	// 使用事务保证原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var po models.InventoryPO
		if err := tx.First(&po, "sku_id = ?", skuID).Error; err != nil {
			return err
		}

		// 增加实际库存
		po.Quantity += quantity
		recalcStatus(&po)
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

// FindInventoriesBySkuIDs 根据 SKU ID 列表批量查询库存
func (r *InventoryRepository) FindInventoriesBySkuIDs(ctx context.Context, skuIDs []int64) (map[int64]*invModels.Inventory, error) {
	if len(skuIDs) == 0 {
		return map[int64]*invModels.Inventory{}, nil
	}
	var pos []models.InventoryPO
	if err := r.db.WithContext(ctx).Where("sku_id IN ?", skuIDs).Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make(map[int64]*invModels.Inventory, len(pos))
	for i := range pos {
		inv := pos[i].ToDomain()
		result[inv.SkuID] = inv
	}
	return result, nil
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

// ListInventoriesEnriched 列出库存（含 SKU 名称和产品名称，单次 LEFT JOIN 查询替代 N+1 补全）
func (r *InventoryRepository) ListInventoriesEnriched(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]InventoryEnrichedRow, error) {
	db := r.db.WithContext(ctx).Table("inventories").
		Select(`inventories.id, inventories.sku_id, skus.name AS sku_name, skus.sku_code,
				skus.product_id, products.name AS product_name,
				inventories.quantity, inventories.status, inventories.reserved,
				inventories.threshold, inventories.created_at, inventories.updated_at`).
		Joins("LEFT JOIN skus ON inventories.sku_id = skus.id").
		Joins("LEFT JOIN products ON skus.product_id = products.id")

	db = r.applyEnrichedConditions(db, q)
	db = query.ApplyOrder(db, q.SortBy, q.Order, "inventories.id ASC")

	var rows []InventoryEnrichedRow
	if err := db.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// CountInventoriesEnriched 统计 enriched 库存数量
func (r *InventoryRepository) CountInventoriesEnriched(ctx context.Context, q dto.InventoryListQuery) (int64, error) {
	var total int64
	db := r.db.WithContext(ctx).Table("inventories").
		Joins("LEFT JOIN skus ON inventories.sku_id = skus.id").
		Joins("LEFT JOIN products ON skus.product_id = products.id")

	db = r.applyEnrichedConditions(db, q)
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// applyEnrichedConditions 应用 enriched 查询条件（使用正确的 JOIN 链路而非 applyQueryConditions 的 buggy JOIN）
func (r *InventoryRepository) applyEnrichedConditions(db *gorm.DB, q dto.InventoryListQuery) *gorm.DB {
	if q.ProductID > 0 {
		db = db.Where("skus.product_id = ?", q.ProductID)
	}
	if q.ProductName != "" {
		db = db.Where("products.name LIKE ?", "%"+q.ProductName+"%")
	}
	if q.SkuName != "" {
		db = db.Where("skus.name LIKE ?", "%"+q.SkuName+"%")
	}
	if q.SKUCode != "" {
		db = db.Where("skus.sku_code = ?", q.SKUCode)
	}
	if q.Status != "" {
		db = db.Where("inventories.status = ?", q.Status)
	}
	if q.LowStock != nil && *q.LowStock {
		db = db.Where("inventories.quantity <= inventories.threshold AND inventories.quantity > 0")
	}
	return db
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r *InventoryRepository) applyQueryConditions(ctx context.Context, q dto.InventoryListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.InventoryPO{})
	db = db.Joins("JOIN products ON inventories.sku_id = products.id")

	if q.ProductID > 0 {
		db = db.Where("inventories.sku_id IN (SELECT id FROM skus WHERE product_id = ?)", q.ProductID)
	}
	if q.ProductName != "" {
		db = db.Where("products.name LIKE ?", "%"+q.ProductName+"%")
	}
	if q.SkuName != "" || q.SKUCode != "" {
		db = db.Joins("JOIN skus ON inventories.sku_id = skus.id")
	}
	if q.SkuName != "" {
		db = db.Where("skus.name LIKE ?", "%"+q.SkuName+"%")
	}
	if q.SKUCode != "" {
		db = db.Where("skus.sku_code = ?", q.SKUCode)
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
