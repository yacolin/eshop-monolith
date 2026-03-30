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

// ListCategories 列出分类（支持查询条件）
func (r CategoryRepository) ListCategories(ctx context.Context, q category.CategoryListQuery, offset, limit int) ([]category.Category, error) {
	var categories []category.Category
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	// 执行查询
	err := db.Offset(offset).Limit(limit).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// CountCategories 统计分类数量
func (r CategoryRepository) CountCategories(ctx context.Context, q category.CategoryListQuery) (int64, error) {
	var total int64
	db := r.applyQueryConditions(ctx, q)

	// 执行统计（不需要排序）
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r CategoryRepository) applyQueryConditions(ctx context.Context, q category.CategoryListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&category.Category{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.ParentID != nil {
		db = db.Where("parent_id = ?", *q.ParentID)
	}
	return db
}

// applyOrder 应用排序
func (r CategoryRepository) applyOrder(db *gorm.DB, q category.CategoryListQuery) *gorm.DB {
	order := "id asc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "asc"
		}
		order = q.SortBy + " " + ord
	}
	return db.Order(order)
}
