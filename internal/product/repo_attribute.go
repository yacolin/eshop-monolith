package product

import (
	"context"

	"gorm.io/gorm"
)

type IattributeRepository interface {
	Create(ctx context.Context, attr *Attribute) error
	FindByID(ctx context.Context, id int64) (*Attribute, error)
	ListByCategory(ctx context.Context, categoryID int64) ([]Attribute, error)
	ListSearchable(ctx context.Context, categoryID int64) ([]Attribute, error)
	ListSkuSpec(ctx context.Context, categoryID int64) ([]Attribute, error)
	ListAll(ctx context.Context) ([]Attribute, error)
	Update(ctx context.Context, attr *Attribute) error
	Delete(ctx context.Context, id int64) error
}

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) IattributeRepository {
	return &AttributeRepository{db: db}
}

func (r *AttributeRepository) Create(ctx context.Context, attr *Attribute) error {
	return r.db.WithContext(ctx).Create(attr).Error
}

func (r *AttributeRepository) FindByID(ctx context.Context, id int64) (*Attribute, error) {
	var attr Attribute
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&attr).Error
	return &attr, err
}

func (r *AttributeRepository) ListByCategory(ctx context.Context, categoryID int64) ([]Attribute, error) {
	var list []Attribute
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *AttributeRepository) ListSearchable(ctx context.Context, categoryID int64) ([]Attribute, error) {
	var list []Attribute
	err := r.db.WithContext(ctx).
		Where("category_id = ? AND searchable = 1 AND status = 1", categoryID).
		Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *AttributeRepository) ListSkuSpec(ctx context.Context, categoryID int64) ([]Attribute, error) {
	var list []Attribute
	err := r.db.WithContext(ctx).
		Where("category_id = ? AND is_sku_spec = 1 AND status = 1", categoryID).
		Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *AttributeRepository) ListAll(ctx context.Context) ([]Attribute, error) {
	var list []Attribute
	err := r.db.WithContext(ctx).Order("category_id ASC, sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *AttributeRepository) Update(ctx context.Context, attr *Attribute) error {
	return r.db.WithContext(ctx).Model(attr).Updates(attr).Error
}

func (r *AttributeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Attribute{}).Error
}
