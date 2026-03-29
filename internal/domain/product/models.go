package product

import (
	"time"
)

// Product 产品领域模型
type Product struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	SKU         string    `json:"sku" gorm:"uniqueIndex"`
	CategoryID  int64     `json:"category_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}