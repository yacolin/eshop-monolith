package product

import (
	"eshop-monolith/internal/pkg/utils"

	"gorm.io/gorm"
)

// Product 产品领域模型
type Product struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Price       int64  `gorm:"type:bigint;not null" json:"price"` // 价格，单位：分
	SKU         string `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 产品表名
func (Product) TableName() string {
	return "products"
}
