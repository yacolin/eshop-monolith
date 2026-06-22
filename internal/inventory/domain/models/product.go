package models

import "eshop-monolith/pkg/utils"

// Product 产品领域模型 (SPU)
type Product struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	MinPrice    int64           `json:"min_price"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 产品表名
func (Product) TableName() string {
	return "products"
}
