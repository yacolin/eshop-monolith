package models

import (
	"time"

	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// CategoryPO 分类持久化对象
type CategoryPO struct {
	ID          int64           `gorm:"primaryKey;autoIncrement"`
	Name        string          `gorm:"type:varchar(100);not null"`
	Description string          `gorm:"type:text"`
	ParentID    *int64          `gorm:"index"`
	Parent      *CategoryPO     `gorm:"foreignKey:ParentID"`
	Children    []CategoryPO    `gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
}

func (CategoryPO) TableName() string { return "categories" }

func (po *CategoryPO) ToDomain() *domain.Category {
	var parent *domain.Category
	if po.Parent != nil {
		parent = po.Parent.ToDomain()
	}
	var children []domain.Category
	for i := range po.Children {
		children = append(children, *po.Children[i].ToDomain())
	}
	return &domain.Category{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		ParentID:    po.ParentID,
		Parent:      parent,
		Children:    children,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func CategoryFromDomain(c *domain.Category) *CategoryPO {
	var parent *CategoryPO
	if c.Parent != nil {
		parent = CategoryFromDomain(c.Parent)
	}
	children := make([]CategoryPO, len(c.Children))
	for i, child := range c.Children {
		children[i] = *CategoryFromDomain(&child)
	}
	return &CategoryPO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ParentID:    c.ParentID,
		Parent:      parent,
		Children:    children,
		CreatedAt:   time.Time(c.CreatedAt),
		UpdatedAt:   time.Time(c.UpdatedAt),
	}
}
