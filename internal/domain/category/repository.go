package category

import (
	"context"
)

// Repository 分类仓储接口
type Repository interface {
	// Create 创建分类
	Create(ctx context.Context, category *Category) error
	// FindByID 根据ID查询分类
	FindByID(ctx context.Context, id int64) (*Category, error)
	// ListAll 列出所有分类
	ListAll(ctx context.Context) ([]Category, error)
	// ListRoot 列出根分类
	ListRoot(ctx context.Context) ([]Category, error)
	// ListByParent 列出子分类
	ListByParent(ctx context.Context, parentID int64) ([]Category, error)
	// Update 更新分类
	Update(ctx context.Context, category *Category) error
	// Delete 删除分类
	Delete(ctx context.Context, id int64) error

	// ListCategories 列出分类
	ListCategories(ctx context.Context, q CategoryListQuery, offset, limit int) ([]Category, error)
	// CountCategories 统计分类数量
	CountCategories(ctx context.Context, q CategoryListQuery) (int64, error)
}
