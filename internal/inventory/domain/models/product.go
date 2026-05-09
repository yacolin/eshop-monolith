package models

import "eshop-monolith/pkg/utils"

// Product 产品领域模型
type Product struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       int64           `json:"price"` // 价格，单位：分
	SKU         string          `json:"sku"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 产品表名
func (Product) TableName() string {
	return "products"
}
