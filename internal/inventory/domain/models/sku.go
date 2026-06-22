package models

import "eshop-monolith/pkg/utils"

// Sku 库存量单元（规格变体）
type Sku struct {
	ID        int64             `json:"id"`
	ProductID int64             `json:"product_id"`
	Name      string            `json:"name"`
	Price     int64             `json:"price"`
	SKUCode   string            `json:"sku_code"`
	Image     string            `json:"image,omitempty"`
	Spec      map[string]string `json:"spec,omitempty"`
	CreatedAt utils.Timestamp   `json:"created_at"`
	UpdatedAt utils.Timestamp   `json:"updated_at"`
}

func (Sku) TableName() string { return "skus" }

// SkuAttribute SKU-属性值关联
type SkuAttribute struct {
	SkuID           int64 `json:"sku_id"`
	AttributeID     int64 `json:"attribute_id"`
	AttributeValueID int64 `json:"attribute_value_id"`
}

func (SkuAttribute) TableName() string { return "sku_attributes" }
