package repositories

import (
	"context"

	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/pkg/query"

	"gorm.io/gorm"
)

// Repository 分类仓储接口
type IcategoryRepository interface {
	// Create 创建分类
	Create(ctx context.Context, category *models.Category) error
	// FindByID 根据ID查询分类
	FindByID(ctx context.Context, id int64) (*models.Category, error)
	// ListAll 列出所有分类
	ListAll(ctx context.Context) ([]models.Category, error)
	// ListRoot 列出根分类
	ListRoot(ctx context.Context) ([]models.Category, error)
	// ListByParent 列出子分类
	ListByParent(ctx context.Context, parentID int64) ([]models.Category, error)
	// Update 更新分类
	Update(ctx context.Context, category *models.Category) error
	// Delete 删除分类
	Delete(ctx context.Context, id int64) error

	// ListCategories 列出分类
	ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error)
	// CountCategories 统计分类数量
	CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error)
}

// CategoryRepository 分类仓储实现
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) IcategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建分类
func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// FindByID 根据ID查询分类
func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*models.Category, error) {
	var c models.Category
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// FindByPath 根据路径查询分类
func (r *CategoryRepository) FindByPath(ctx context.Context, path string) (*models.Category, error) {
	var c models.Category
	err := r.db.WithContext(ctx).First(&c, "path = ?", path).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListRoot 列出根分类
func (r *CategoryRepository) ListRoot(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Where("parent_id IS NULL").Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ListByParent 列出子分类
func (r *CategoryRepository) ListByParent(ctx context.Context, parentID int64) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ListAll 列出所有分类
func (r *CategoryRepository) ListAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 更新分类
func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

// ListCategories 列出分类（支持查询条件）
func (r *CategoryRepository) ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error) {
	var categories []models.Category
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
func (r *CategoryRepository) CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error) {
	var total int64
	db := r.applyQueryConditions(ctx, q)

	// 执行统计（不需要排序）
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r *CategoryRepository) applyQueryConditions(ctx context.Context, q dto.CategoryListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Category{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.ParentID != nil {
		db = db.Where("parent_id = ?", *q.ParentID)
	}
	return db
}

// applyOrder 应用排序
func (r *CategoryRepository) applyOrder(db *gorm.DB, q dto.CategoryListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}
