package product

import (
	"context"
)

// Repository 产品仓储接口
type Repository interface {
	// Create 创建产品
	Create(ctx context.Context, product *Product) error
	// FindByID 根据ID查询产品
	FindByID(ctx context.Context, id int64) (*Product, error)
	// FindBySKU 根据SKU查询产品
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	// ListByCategory 根据分类查询产品
	ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]Product, int64, error)
	// ListAll 列出所有产品
	ListAll(ctx context.Context, page, pageSize int) ([]Product, int64, error)
	// Update 更新产品
	Update(ctx context.Context, product *Product) error
	// Delete 删除产品
	Delete(ctx context.Context, id int64) error

	// ListProducts 列出产品
	ListProducts(ctx context.Context, q ProductListQuery, offset, limit int) ([]Product, error)
	// CountProducts 统计产品数量
	CountProducts(ctx context.Context, q ProductListQuery) (int64, error)
}
