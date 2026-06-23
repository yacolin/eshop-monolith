package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	domain "eshop-monolith/internal/inventory/domain/models"

	"gorm.io/gorm"
)

type IattributeRepository interface {
	// Attribute CRUD
	Create(ctx context.Context, attr *domain.Attribute) error
	FindByID(ctx context.Context, id int64) (*domain.Attribute, error)
	FindAll(ctx context.Context, offset, limit int) ([]domain.Attribute, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, attr *domain.Attribute) error
	Delete(ctx context.Context, id int64) error

	// AttributeValue CRUD
	CreateValue(ctx context.Context, val *domain.AttributeValue) error
	FindValueByID(ctx context.Context, id int64) (*domain.AttributeValue, error)
	FindValuesByAttributeID(ctx context.Context, attributeID int64) ([]domain.AttributeValue, error)
	UpdateValue(ctx context.Context, val *domain.AttributeValue) error
	DeleteValue(ctx context.Context, id int64) error
}

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) IattributeRepository {
	return &AttributeRepository{db: db}
}

// ── Attribute ──────────────────────────────────────────────────

func (r *AttributeRepository) Create(ctx context.Context, attr *domain.Attribute) error {
	po := models.AttributeFromDomain(attr)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	attr.ID = po.ID
	return nil
}

func (r *AttributeRepository) FindByID(ctx context.Context, id int64) (*domain.Attribute, error) {
	var po models.AttributePO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *AttributeRepository) FindAll(ctx context.Context, offset, limit int) ([]domain.Attribute, error) {
	var pos []models.AttributePO
	if err := r.db.WithContext(ctx).Order("sort_order asc, id asc").Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, err
	}
	attrs := make([]domain.Attribute, len(pos))
	for i, po := range pos {
		attrs[i] = *po.ToDomain()
	}
	return attrs, nil
}

func (r *AttributeRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.AttributePO{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *AttributeRepository) Update(ctx context.Context, attr *domain.Attribute) error {
	po := models.AttributeFromDomain(attr)
	return r.db.WithContext(ctx).Select("name", "sort_order").Where("id = ?", attr.ID).Updates(po).Error
}

func (r *AttributeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.AttributePO{}, "id = ?", id).Error
}

// ── AttributeValue ─────────────────────────────────────────────

func (r *AttributeRepository) CreateValue(ctx context.Context, val *domain.AttributeValue) error {
	po := models.AttributeValueFromDomain(val)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	val.ID = po.ID
	return nil
}

func (r *AttributeRepository) FindValueByID(ctx context.Context, id int64) (*domain.AttributeValue, error) {
	var po models.AttributeValuePO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *AttributeRepository) FindValuesByAttributeID(ctx context.Context, attributeID int64) ([]domain.AttributeValue, error) {
	var pos []models.AttributeValuePO
	if err := r.db.WithContext(ctx).Where("attribute_id = ?", attributeID).Order("sort_order asc, id asc").Find(&pos).Error; err != nil {
		return nil, err
	}
	vals := make([]domain.AttributeValue, len(pos))
	for i, po := range pos {
		vals[i] = *po.ToDomain()
	}
	return vals, nil
}

func (r *AttributeRepository) UpdateValue(ctx context.Context, val *domain.AttributeValue) error {
	po := models.AttributeValueFromDomain(val)
	return r.db.WithContext(ctx).Select("value", "sort_order").Where("id = ?", val.ID).Updates(po).Error
}

func (r *AttributeRepository) DeleteValue(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.AttributeValuePO{}, "id = ?", id).Error
}
