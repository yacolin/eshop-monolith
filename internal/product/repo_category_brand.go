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
	FindBrandDetailsByCategory(ctx context.Context, categoryID int64) ([]CategoryBrandDetail, error)
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

func (r *CategoryBrandRepository) FindBrandDetailsByCategory(ctx context.Context, categoryID int64) ([]CategoryBrandDetail, error) {
	var list []CategoryBrandDetail
	err := r.db.WithContext(ctx).Table("sp_category_brands scb").
		Select("scb.brand_id, b.name AS brand_name, b.english_name, b.logo_url, b.first_letter, MIN(scb.sort_order) AS sort_order").
		Joins("JOIN sp_brands b ON b.id = scb.brand_id").
		Where("scb.category_id IN (SELECT id FROM sp_categories WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM sp_categories WHERE id = ?), ?, '/%'))", categoryID, categoryID, categoryID).
		Where("b.status = 1").
		Group("scb.brand_id").
		Order("sort_order ASC, scb.brand_id ASC").
		Scan(&list).Error
	return list, err
}
