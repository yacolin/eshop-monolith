package repository

import (
	"context"

	"gorm.io/gorm"

	"eshop-monolith/internal/domain/category"
)

// CategoryRepository 分类仓储实现
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{db: db}
}

// Create 创建分类
func (r CategoryRepository) Create(ctx context.Context, category *category.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// FindByID 根据ID查询分类
func (r CategoryRepository) FindByID(ctx context.Context, id int64) (*category.Category, error) {
	var c category.Category
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// FindByPath 根据路径查询分类
func (r CategoryRepository) FindByPath(ctx context.Context, path string) (*category.Category, error) {
	var c category.Category
	err := r.db.WithContext(ctx).First(&c, "path = ?", path).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListRoot 列出根分类
func (r CategoryRepository) ListRoot(ctx context.Context) ([]category.Category, error) {
	var categories []category.Category
	err := r.db.WithContext(ctx).Where("parent_id IS NULL").Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ListByParent 列出子分类
func (r CategoryRepository) ListByParent(ctx context.Context, parentID int64) ([]category.Category, error) {
	var categories []category.Category
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ListAll 列出所有分类
func (r CategoryRepository) ListAll(ctx context.Context) ([]category.Category, error) {
	var categories []category.Category
	err := r.db.WithContext(ctx).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 更新分类
func (r CategoryRepository) Update(ctx context.Context, category *category.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete 删除分类
func (r CategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&category.Category{}, "id = ?", id).Error
}
