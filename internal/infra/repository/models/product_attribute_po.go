package models

import domain "eshop-monolith/internal/inventory/domain/models"

// ProductAttributeValuePO 产品-属性-属性值 持久化对象
type ProductAttributeValuePO struct {
	ProductID       int64 `gorm:"primaryKey"`
	AttributeID     int64 `gorm:"primaryKey"`
	AttributeValueID int64 `gorm:"primaryKey"`
}

func (ProductAttributeValuePO) TableName() string { return "product_attribute_values" }

func (po *ProductAttributeValuePO) ToDomain() *domain.ProductAttributeValue {
	return &domain.ProductAttributeValue{
		ProductID:       po.ProductID,
		AttributeID:     po.AttributeID,
		AttributeValueID: po.AttributeValueID,
	}
}

func ProductAttributeValueFromDomain(v *domain.ProductAttributeValue) *ProductAttributeValuePO {
	return &ProductAttributeValuePO{
		ProductID:       v.ProductID,
		AttributeID:     v.AttributeID,
		AttributeValueID: v.AttributeValueID,
	}
}
