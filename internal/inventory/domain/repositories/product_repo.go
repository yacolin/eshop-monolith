package repositories

import (
	"context"

	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

// Repository 产品仓储接口
type IproductRepository interface {
	// Create 创建产品
	Create(ctx context.Context, product *models.Product) error
	// FindByID 根据ID查询产品
	FindByID(ctx context.Context, id int64) (*models.Product, error)
	// FindBySKU 根据SKU查询产品
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	// ListByCategory 根据分类查询产品
	ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]models.Product, int64, error)
	// ListAll 列出所有产品
	ListAll(ctx context.Context, page, pageSize int) ([]models.Product, int64, error)
	// Update 更新产品
	Update(ctx context.Context, product *models.Product) error
	// Delete 删除产品
	Delete(ctx context.Context, id int64) error

	// ListProducts 列出产品
	ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error)
	// CountProducts 统计产品数量
	CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error)
}

// ProductRepository 产品仓储实现
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓储
func NewProductRepository(db *gorm.DB) IproductRepository {
	return &ProductRepository{db: db}
}

// Create 创建产品
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// FindByID 根据ID查询产品
func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBySKU 根据SKU查询产品
func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).First(&p, "sku = ?", sku).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListByCategory 根据分类查询产品
func (r *ProductRepository) ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
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
func (r *ProductRepository) ListAll(ctx context.Context, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
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
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 删除产品
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

// ListProducts 列出产品（支持查询条件）
func (r *ProductRepository) ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error) {
	var products []models.Product
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	// 执行查询
	err := db.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

// CountProducts 统计产品数量
func (r *ProductRepository) CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error) {
	var total int64
	db := r.applyQueryConditions(ctx, q)

	// 执行统计（不需要排序）
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r *ProductRepository) applyQueryConditions(ctx context.Context, q dto.ProductListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Product{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.SKU != "" {
		db = db.Where("sku = ?", q.SKU)
	}
	return db
}

// applyOrder 应用排序
func (r *ProductRepository) applyOrder(db *gorm.DB, q dto.ProductListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}
