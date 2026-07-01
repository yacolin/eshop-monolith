package inventory

import (
	"context"

	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

type IinventoryRepository interface {
	FindBySkuForUpdate(tx *gorm.DB, skuID, warehouseID int64) (*Inventory, error)
	UpdateWithTx(tx *gorm.DB, inv *Inventory) error
	FindBySku(ctx context.Context, skuID, warehouseID int64) (*Inventory, error)
	UpsertWithTx(tx *gorm.DB, inv *Inventory) error
	CreateLogWithTx(tx *gorm.DB, log *InventoryLog) error
	ListLogs(ctx context.Context, skuID int64, changeType string, page, size int) ([]InventoryLog, int64, error)
	FindOrCreateWithTx(tx *gorm.DB, skuID, warehouseID int64) (*Inventory, error)
}

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) IinventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) FindBySkuForUpdate(tx *gorm.DB, skuID, warehouseID int64) (*Inventory, error) {
	var inv Inventory
	err := tx.Where("sku_id = ? AND warehouse_id = ?", skuID, warehouseID).
		Session(&gorm.Session{}).First(&inv).Error
	return &inv, err
}

func (r *InventoryRepository) UpdateWithTx(tx *gorm.DB, inv *Inventory) error {
	return tx.Model(inv).Select("Quantity", "Reserved", "Status", "LastCountedAt", "LastCountedBy").Updates(inv).Error
}

func (r *InventoryRepository) FindBySku(ctx context.Context, skuID, warehouseID int64) (*Inventory, error) {
	var inv Inventory
	err := r.db.WithContext(ctx).Where("sku_id = ? AND warehouse_id = ?", skuID, warehouseID).First(&inv).Error
	return &inv, err
}

func (r *InventoryRepository) UpsertWithTx(tx *gorm.DB, inv *Inventory) error {
	return tx.Where("sku_id = ? AND warehouse_id = ?", inv.SkuID, inv.WarehouseID).
		Assign(map[string]interface{}{
			"quantity": gorm.Expr("quantity + ?", inv.Quantity),
		}).
		FirstOrCreate(inv).Error
}

func (r *InventoryRepository) CreateLogWithTx(tx *gorm.DB, log *InventoryLog) error {
	return tx.Create(log).Error
}

func (r *InventoryRepository) ListLogs(ctx context.Context, skuID int64, changeType string, page, size int) ([]InventoryLog, int64, error) {
	db := r.db.WithContext(ctx).Model(&InventoryLog{}).Where("sku_id = ?", skuID)
	if changeType != "" {
		db = db.Where("change_type = ?", changeType)
	}
	return query.ConcurrentCountList[InventoryLog](db.Order("id DESC"), page, size)
}

func (r *InventoryRepository) FindOrCreateWithTx(tx *gorm.DB, skuID, warehouseID int64) (*Inventory, error) {
	var inv Inventory
	err := tx.Where("sku_id = ? AND warehouse_id = ?", skuID, warehouseID).First(&inv).Error
	if err == nil {
		return &inv, nil
	}
	if err == gorm.ErrRecordNotFound {
		inv = Inventory{
			SkuID:       skuID,
			WarehouseID: warehouseID,
			Status:      "instock",
		}
		return &inv, tx.Create(&inv).Error
	}
	return nil, err
}
