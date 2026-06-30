package product

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

type ProductAttribute struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   int64          `gorm:"not null;uniqueIndex:uk_product_attribute" json:"product_id"`
	AttributeID int64          `gorm:"not null;uniqueIndex:uk_product_attribute" json:"attribute_id"`
	Value       string         `gorm:"type:varchar(500);not null" json:"value"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (ProductAttribute) TableName() string { return "sp_product_attributes" }
