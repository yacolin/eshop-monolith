package repository

import (
	"context"

	"gorm.io/gorm"

	"eshop-monolith/internal/domain/product"
)

// ProductRepository 产品仓储实现
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓储
func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db: db}
}

// Create 创建产品
func (r ProductRepository) Create(ctx context.Context, product *product.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// FindByID 根据ID查询产品
func (r ProductRepository) FindByID(ctx context.Context, id int64) (*product.Product, error) {
	var p product.Product
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBySKU 根据SKU查询产品
func (r ProductRepository) FindBySKU(ctx context.Context, sku string) (*product.Product, error) {
	var p product.Product
	err := r.db.WithContext(ctx).First(&p, "sku = ?", sku).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListByCategory 根据分类查询产品
func (r ProductRepository) ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&product.Product{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListAll 列出所有产品
func (r ProductRepository) ListAll(ctx context.Context, page, pageSize int) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&product.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update 更新产品
func (r ProductRepository) Update(ctx context.Context, product *product.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 删除产品
func (r ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&product.Product{}, "id = ?", id).Error
}