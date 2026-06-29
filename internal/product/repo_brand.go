package product

import (
	"context"

	"gorm.io/gorm"
)

type IbrandRepository interface {
	Create(ctx context.Context, brand *Brand) error
	FindByID(ctx context.Context, id int64) (*Brand, error)
	FindByName(ctx context.Context, name string) (*Brand, error)
	List(ctx context.Context, name, firstLetter string, status *int8, page, size int) ([]Brand, int64, error)
	Update(ctx context.Context, brand *Brand) error
	Delete(ctx context.Context, id int64) error
}

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) IbrandRepository {
	return &BrandRepository{db: db}
}

func (r *BrandRepository) Create(ctx context.Context, brand *Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

func (r *BrandRepository) FindByID(ctx context.Context, id int64) (*Brand, error) {
	var brand Brand
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&brand).Error
	return &brand, err
}

func (r *BrandRepository) FindByName(ctx context.Context, name string) (*Brand, error) {
	var brand Brand
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&brand).Error
	return &brand, err
}

func (r *BrandRepository) List(ctx context.Context, name, firstLetter string, status *int8, page, size int) ([]Brand, int64, error) {
	db := r.db.WithContext(ctx).Model(&Brand{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if firstLetter != "" {
		db = db.Where("first_letter = ?", firstLetter)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Brand
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("sort_order DESC, id ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *BrandRepository) Update(ctx context.Context, brand *Brand) error {
	return r.db.WithContext(ctx).Model(brand).Updates(brand).Error
}

func (r *BrandRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Brand{}).Error
}
