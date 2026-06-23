package models

// CategoryAttribute 品类与规格属性多对多关联
type CategoryAttribute struct {
	CategoryID  int64 `json:"category_id"`
	AttributeID int64 `json:"attribute_id"`
}

func (CategoryAttribute) TableName() string { return "category_attributes" }
