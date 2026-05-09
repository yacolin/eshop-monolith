package models

import "eshop-monolith/pkg/utils"

// Category 分类领域模型
type Category struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ParentID    *int64      `json:"parent_id"` // 父分类ID，支持层级结构
	Parent      *Category   `json:"parent"`
	Children    []Category  `json:"children"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 分类表名
func (Category) TableName() string {
	return "categories"
}
