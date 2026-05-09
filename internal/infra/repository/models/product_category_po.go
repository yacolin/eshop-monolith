package models

import (
	shared "eshop-monolith/internal/infra/domain/shared"
)

// ProductCategoryPO 产品与分类多对多关联持久化对象
type ProductCategoryPO struct {
	ProductID  int64 `gorm:"primaryKey"`
	CategoryID int64 `gorm:"primaryKey"`
}

func (ProductCategoryPO) TableName() string { return "product_categories" }

func (po *ProductCategoryPO) ToDomain() *shared.ProductCategory {
	return &shared.ProductCategory{
		ProductID:  po.ProductID,
		CategoryID: po.CategoryID,
	}
}

func ProductCategoryFromDomain(pc *shared.ProductCategory) *ProductCategoryPO {
	return &ProductCategoryPO{
		ProductID:  pc.ProductID,
		CategoryID: pc.CategoryID,
	}
}
