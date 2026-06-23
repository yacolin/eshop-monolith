package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/inventory/api/dto"
	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

type IskuRepository interface {
	Create(ctx context.Context, sku *domain.Sku) error
	FindByID(ctx context.Context, id int64) (*domain.Sku, error)
	FindByIDs(ctx context.Context, ids []int64) ([]domain.Sku, error)
	FindByProductID(ctx context.Context, productID int64) ([]domain.Sku, error)
	Update(ctx context.Context, sku *domain.Sku) error
	Delete(ctx context.Context, id int64) error
	CreateWithTx(tx *gorm.DB, sku *domain.Sku) error
	DeleteByProductIDWithTx(tx *gorm.DB, productID int64) error
	FindAll(ctx context.Context, q dto.SkuListQuery, offset, limit int) ([]domain.Sku, error)
	Count(ctx context.Context, q dto.SkuListQuery) (int64, error)
}

type SkuRepository struct{ db *gorm.DB }

func NewSkuRepository(db *gorm.DB) IskuRepository { return &SkuRepository{db: db} }

func (r *SkuRepository) Create(ctx context.Context, sku *domain.Sku) error {
	po := models.SkuFromDomain(sku)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	sku.ID = po.ID
	return nil
}

func (r *SkuRepository) FindByID(ctx context.Context, id int64) (*domain.Sku, error) {
	var po models.SkuPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *SkuRepository) FindByIDs(ctx context.Context, ids []int64) ([]domain.Sku, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var pos []models.SkuPO
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&pos).Error; err != nil {
		return nil, err
	}
	skus := make([]domain.Sku, len(pos))
	for i, po := range pos {
		skus[i] = *po.ToDomain()
	}
	return skus, nil
}

func (r *SkuRepository) FindByProductID(ctx context.Context, productID int64) ([]domain.Sku, error) {
	var pos []models.SkuPO
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("id asc").Find(&pos).Error; err != nil {
		return nil, err
	}
	skus := make([]domain.Sku, len(pos))
	for i, po := range pos {
		skus[i] = *po.ToDomain()
	}
	return skus, nil
}

func (r *SkuRepository) Update(ctx context.Context, sku *domain.Sku) error {
	po := models.SkuFromDomain(sku)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *SkuRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.SkuPO{}, "id = ?", id).Error
}

func (r *SkuRepository) CreateWithTx(tx *gorm.DB, sku *domain.Sku) error {
	po := models.SkuFromDomain(sku)
	if err := tx.Create(po).Error; err != nil {
		return err
	}
	sku.ID = po.ID
	return nil
}

func (r *SkuRepository) DeleteByProductIDWithTx(tx *gorm.DB, productID int64) error {
	return tx.Where("product_id = ?", productID).Delete(&models.SkuPO{}).Error
}

func (r *SkuRepository) FindAll(ctx context.Context, q dto.SkuListQuery, offset, limit int) ([]domain.Sku, error) {
	var pos []models.SkuPO
	db := r.applySkuConditions(ctx, q)
	db = query.ApplyOrder(db, q.SortBy, q.Order, "id asc")

	if err := db.Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, err
	}
	skus := make([]domain.Sku, len(pos))
	for i, po := range pos {
		skus[i] = *po.ToDomain()
	}
	return skus, nil
}

func (r *SkuRepository) Count(ctx context.Context, q dto.SkuListQuery) (int64, error) {
	var total int64
	db := r.applySkuConditions(ctx, q)
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *SkuRepository) applySkuConditions(ctx context.Context, q dto.SkuListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.SkuPO{})
	if q.ProductID != nil {
		db = db.Where("product_id = ?", *q.ProductID)
	}
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.SKUCode != "" {
		db = db.Where("sku_code = ?", q.SKUCode)
	}
	if q.PriceMin != nil {
		db = db.Where("price >= ?", *q.PriceMin)
	}
	if q.PriceMax != nil {
		db = db.Where("price <= ?", *q.PriceMax)
	}
	return db
}
