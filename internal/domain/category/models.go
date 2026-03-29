package category

import (
	"eshop-monolith/internal/domain/product"
	"eshop-monolith/internal/pkg/utils"

	"gorm.io/gorm"
)

// Category 分类领域模型
type Category struct {
	ID          int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string            `gorm:"type:varchar(100);not null" json:"name"`
	Description string            `gorm:"type:text" json:"description"`
	ParentID    *int64            `gorm:"index" json:"parent_id"` // 父分类ID，支持层级结构
	Parent      *Category         `gorm:"foreignKey:ParentID" json:"parent"`
	Children    []Category        `gorm:"foreignKey:ParentID" json:"children"`
	Products    []product.Product `gorm:"many2many:product_categories;" json:"products"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 分类表名
func (Category) TableName() string {
	return "categories"
}
