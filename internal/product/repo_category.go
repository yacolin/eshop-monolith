package product

import (
	"context"

	"gorm.io/gorm"
)

type IcategoryRepository interface {
	Create(ctx context.Context, cat *Category) error
	FindByID(ctx context.Context, id int64) (*Category, error)
	FindByName(ctx context.Context, name string, parentID int64) (*Category, error)
	List(ctx context.Context, name string, status, level *int8, page, size int) ([]Category, int64, error)
	ListRoot(ctx context.Context) ([]Category, error)
	ListByParent(ctx context.Context, parentID int64) ([]Category, error)
	ListByLevel(ctx context.Context, level int8) ([]Category, error)
	ListAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, cat *Category) error
	Delete(ctx context.Context, id int64) error
	CountByParentID(ctx context.Context, parentID int64) (int64, error)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) IcategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, cat *Category) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*Category, error) {
	var cat Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error
	return &cat, err
}

func (r *CategoryRepository) FindByName(ctx context.Context, name string, parentID int64) (*Category, error) {
	var cat Category
	err := r.db.WithContext(ctx).Where("name = ? AND parent_id = ?", name, parentID).First(&cat).Error
	return &cat, err
}

func (r *CategoryRepository) ListRoot(ctx context.Context) ([]Category, error) {
	var list []Category
	err := r.db.WithContext(ctx).Where("parent_id = 0").Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) ListByParent(ctx context.Context, parentID int64) ([]Category, error) {
	var list []Category
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) ListByLevel(ctx context.Context, level int8) ([]Category, error) {
	var list []Category
	err := r.db.WithContext(ctx).Where("level = ?", level).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) ListAll(ctx context.Context) ([]Category, error) {
	var list []Category
	err := r.db.WithContext(ctx).Order("level ASC, sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) List(ctx context.Context, name string, status, level *int8, page, size int) ([]Category, int64, error) {
	db := r.db.WithContext(ctx).Model(&Category{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if level != nil {
		db = db.Where("level = ?", *level)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Category
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("level ASC, sort_order ASC, id ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *CategoryRepository) Update(ctx context.Context, cat *Category) error {
	return r.db.WithContext(ctx).Model(cat).Updates(cat).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Category{}).Error
}

func (r *CategoryRepository) CountByParentID(ctx context.Context, parentID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Category{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}
