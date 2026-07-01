package product

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

type Attribute struct {
	ID         int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string          `gorm:"type:varchar(100);not null" json:"name"`
	CategoryID int64           `gorm:"not null;index:idx_category" json:"category_id"`
	InputType  int8            `gorm:"not null;default:1" json:"input_type"` // 1-文本 2-单选 3-多选 4-数字
	Values     string          `gorm:"type:json" json:"values,omitempty"`    // 可选值 JSON 数组
	Unit       string          `gorm:"type:varchar(20);default:''" json:"unit"`
	Required   int8            `gorm:"not null;default:0" json:"required"`
	Searchable int8            `gorm:"not null;default:0;index:idx_searchable" json:"searchable"`
	IsSkuSpec  int8            `gorm:"not null;default:0;index:idx_is_sku_spec" json:"is_sku_spec"`
	SortOrder  int             `gorm:"not null;default:0" json:"sort_order"`
	Status     int8            `gorm:"not null;default:1" json:"status"`
	CreatedAt  utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index:idx_deleted_at" json:"-"`
}

func (Attribute) TableName() string { return "sp_attributes" }
