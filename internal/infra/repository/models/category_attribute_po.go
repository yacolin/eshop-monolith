package models

import domain "eshop-monolith/internal/inventory/domain/models"

// CategoryAttributePO 品类与规格属性多对多关联持久化对象
type CategoryAttributePO struct {
	CategoryID  int64 `gorm:"primaryKey"`
	AttributeID int64 `gorm:"primaryKey"`
}

func (CategoryAttributePO) TableName() string { return "category_attributes" }

func (po *CategoryAttributePO) ToDomain() *domain.CategoryAttribute {
	return &domain.CategoryAttribute{
		CategoryID:  po.CategoryID,
		AttributeID: po.AttributeID,
	}
}

func CategoryAttributeFromDomain(ca *domain.CategoryAttribute) *CategoryAttributePO {
	return &CategoryAttributePO{
		CategoryID:  ca.CategoryID,
		AttributeID: ca.AttributeID,
	}
}
