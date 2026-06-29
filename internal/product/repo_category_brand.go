package product

import (
	"context"

	"gorm.io/gorm"
)

type IcategoryBrandRepository interface {
	AddBatch(tx *gorm.DB, items []CategoryBrand) error
	DeleteByCategory(tx *gorm.DB, categoryID int64) error
	FindBrandsByCategory(ctx context.Context, categoryID int64) ([]CategoryBrand, error)
	FindCategoriesByBrand(ctx context.Context, brandID int64) ([]CategoryBrand, error)
}

type CategoryBrandRepository struct {
	db *gorm.DB
}

func NewCategoryBrandRepository(db *gorm.DB) IcategoryBrandRepository {
	return &CategoryBrandRepository{db: db}
}

func (r *CategoryBrandRepository) AddBatch(tx *gorm.DB, items []CategoryBrand) error {
	if len(items) == 0 {
		return nil
	}
	return tx.Create(&items).Error
}

func (r *CategoryBrandRepository) DeleteByCategory(tx *gorm.DB, categoryID int64) error {
	return tx.Where("category_id = ?", categoryID).Delete(&CategoryBrand{}).Error
}

func (r *CategoryBrandRepository) FindBrandsByCategory(ctx context.Context, categoryID int64) ([]CategoryBrand, error) {
	var list []CategoryBrand
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryBrandRepository) FindCategoriesByBrand(ctx context.Context, brandID int64) ([]CategoryBrand, error) {
	var list []CategoryBrand
	err := r.db.WithContext(ctx).Where("brand_id = ?", brandID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}
