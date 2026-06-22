package models

import "eshop-monolith/pkg/utils"

// Attribute 规格属性维度
type Attribute struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	SortOrder int             `json:"sort_order"`
	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (Attribute) TableName() string { return "attribute_attributes" }

// AttributeValue 规格属性可选值
type AttributeValue struct {
	ID          int64           `json:"id"`
	AttributeID int64           `json:"attribute_id"`
	Value       string          `json:"value"`
	SortOrder   int             `json:"sort_order"`
	CreatedAt   utils.Timestamp `json:"created_at"`
	UpdatedAt   utils.Timestamp `json:"updated_at"`
}

func (AttributeValue) TableName() string { return "attribute_values" }
