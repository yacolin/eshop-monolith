package models

// ProductAttributeValue 产品-属性-属性值 关联
type ProductAttributeValue struct {
	ProductID       int64 `json:"product_id"`
	AttributeID     int64 `json:"attribute_id"`
	AttributeValueID int64 `json:"attribute_value_id"`
}

func (ProductAttributeValue) TableName() string { return "product_attribute_values" }
