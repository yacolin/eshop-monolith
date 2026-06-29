package product

import (
	"context"

	"gorm.io/gorm"
)

type IspuRepository interface {
	CreateSPUWithTx(tx *gorm.DB, spu *SPU) error
	CreateSKUWithTx(tx *gorm.DB, sku *SKU) error
	CreateDescriptionWithTx(tx *gorm.DB, desc *Description) error
	CreateProductAttrWithTx(tx *gorm.DB, pa *ProductAttribute) error
	UpdateSPUPriceWithTx(tx *gorm.DB, productID int64) error

	FindByID(ctx context.Context, id int64) (*SPU, error)
	FindSKUsByProductID(ctx context.Context, productID int64) ([]SKU, error)
	FindDescriptionByProductID(ctx context.Context, productID int64) (*Description, error)
	FindProductAttrsByProductID(ctx context.Context, productID int64) ([]ProductAttribute, error)
	FindSKUByCode(ctx context.Context, code string) (*SKU, error)
	FindSKUByID(ctx context.Context, id int64) (*SKU, error)

	List(ctx context.Context, name string, categoryID, brandID *int64, status *int8, priceMin, priceMax int64, page, size int) ([]SPU, int64, error)

	Update(ctx context.Context, spu *SPU) error
	UpdateSKU(ctx context.Context, sku *SKU) error
	Delete(ctx context.Context, id int64) error
	DeleteSKU(ctx context.Context, id int64) error
}

type SpuRepository struct {
	db *gorm.DB
}

func NewSpuRepository(db *gorm.DB) IspuRepository {
	return &SpuRepository{db: db}
}

func (r *SpuRepository) CreateSPUWithTx(tx *gorm.DB, spu *SPU) error {
	return tx.Create(spu).Error
}

func (r *SpuRepository) CreateSKUWithTx(tx *gorm.DB, sku *SKU) error {
	return tx.Create(sku).Error
}

func (r *SpuRepository) CreateDescriptionWithTx(tx *gorm.DB, desc *Description) error {
	return tx.Create(desc).Error
}

func (r *SpuRepository) CreateProductAttrWithTx(tx *gorm.DB, pa *ProductAttribute) error {
	return tx.Create(pa).Error
}

func (r *SpuRepository) UpdateSPUPriceWithTx(tx *gorm.DB, productID int64) error {
	return tx.Exec(`
		UPDATE sp_products
		SET min_price = COALESCE((SELECT MIN(price) FROM sp_skus WHERE product_id = ? AND status = 1 AND deleted_at IS NULL), 0),
		    max_price = COALESCE((SELECT MAX(price) FROM sp_skus WHERE product_id = ? AND status = 1 AND deleted_at IS NULL), 0)
		WHERE id = ?`, productID, productID, productID).Error
}

func (r *SpuRepository) FindByID(ctx context.Context, id int64) (*SPU, error) {
	var spu SPU
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&spu).Error
	return &spu, err
}

func (r *SpuRepository) FindSKUsByProductID(ctx context.Context, productID int64) ([]SKU, error) {
	var list []SKU
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("id ASC").Find(&list).Error
	return list, err
}

func (r *SpuRepository) FindDescriptionByProductID(ctx context.Context, productID int64) (*Description, error) {
	var desc Description
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&desc).Error
	return &desc, err
}

func (r *SpuRepository) FindProductAttrsByProductID(ctx context.Context, productID int64) ([]ProductAttribute, error) {
	var list []ProductAttribute
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *SpuRepository) FindSKUByCode(ctx context.Context, code string) (*SKU, error) {
	var sku SKU
	err := r.db.WithContext(ctx).Where("sku_code = ?", code).First(&sku).Error
	return &sku, err
}

func (r *SpuRepository) FindSKUByID(ctx context.Context, id int64) (*SKU, error) {
	var sku SKU
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sku).Error
	return &sku, err
}

func (r *SpuRepository) List(ctx context.Context, name string, categoryID, brandID *int64, status *int8, priceMin, priceMax int64, page, size int) ([]SPU, int64, error) {
	db := r.db.WithContext(ctx).Model(&SPU{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if categoryID != nil {
		db = db.Where("category_id = ?", *categoryID)
	}
	if brandID != nil {
		db = db.Where("brand_id = ?", *brandID)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if priceMin > 0 {
		db = db.Where("min_price >= ?", priceMin)
	}
	if priceMax > 0 {
		db = db.Where("max_price <= ?", priceMax)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []SPU
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("sort_order DESC, id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *SpuRepository) Update(ctx context.Context, spu *SPU) error {
	return r.db.WithContext(ctx).Model(spu).Updates(spu).Error
}

func (r *SpuRepository) UpdateSKU(ctx context.Context, sku *SKU) error {
	return r.db.WithContext(ctx).Model(sku).Updates(sku).Error
}

func (r *SpuRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&SPU{}).Error
}

func (r *SpuRepository) DeleteSKU(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&SKU{}).Error
}
