package service

import (
	"context"

	"eshop-monolith/internal/domain/category"
	"eshop-monolith/internal/eventbus"
)

// CategoryService 分类服务
type CategoryService struct {
	repo category.Repository
	bus  *eventbus.Bus
}

// NewCategoryService 创建分类服务
func NewCategoryService(repo category.Repository, bus *eventbus.Bus) *CategoryService {
	return &CategoryService{
		repo: repo,
		bus:  bus,
	}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id"`
	Level       int    `json:"level"`
	Path        string `json:"path"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id"`
	Level       int    `json:"level"`
	Path        string `json:"path"`
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*category.Category, error) {
	// 创建分类
	newCategory := &category.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Level:       req.Level,
		Path:        req.Path,
	}

	// 保存分类
	if err := s.repo.Create(ctx, newCategory); err != nil {
		return nil, err
	}

	// 发布分类创建事件
	s.bus.Publish(category.CategoryCreatedEvent{
		CategoryID: newCategory.ID,
		Name:       newCategory.Name,
		ParentID:   newCategory.ParentID,
		Level:      newCategory.Level,
		Path:       newCategory.Path,
	})

	return newCategory, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int64) (*category.Category, error) {
	return s.repo.FindByID(ctx, id)
}

// GetCategoryByPath 根据路径获取分类
func (s *CategoryService) GetCategoryByPath(ctx context.Context, path string) (*category.Category, error) {
	return s.repo.FindByPath(ctx, path)
}

// ListRootCategories 列出根分类
func (s *CategoryService) ListRootCategories(ctx context.Context) ([]category.Category, error) {
	return s.repo.ListRoot(ctx)
}

// ListSubCategories 列出子分类
func (s *CategoryService) ListSubCategories(ctx context.Context, parentID int64) ([]category.Category, error) {
	return s.repo.ListByParent(ctx, parentID)
}

// ListAllCategories 列出所有分类
func (s *CategoryService) ListAllCategories(ctx context.Context) ([]category.Category, error) {
	return s.repo.ListAll(ctx)
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) (*category.Category, error) {
	// 获取分类
	existingCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新分类信息
	existingCategory.Name = req.Name
	existingCategory.Description = req.Description
	existingCategory.ParentID = req.ParentID
	existingCategory.Level = req.Level
	existingCategory.Path = req.Path

	// 保存分类
	if err := s.repo.Update(ctx, existingCategory); err != nil {
		return nil, err
	}

	// 发布分类更新事件
	s.bus.Publish(category.CategoryUpdatedEvent{
		CategoryID: existingCategory.ID,
		Name:       existingCategory.Name,
		ParentID:   existingCategory.ParentID,
		Level:      existingCategory.Level,
		Path:       existingCategory.Path,
	})

	return existingCategory, nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	// 获取分类
	existingCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除分类
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 发布分类删除事件
	s.bus.Publish(category.CategoryDeletedEvent{
		CategoryID: existingCategory.ID,
		Name:       existingCategory.Name,
		ParentID:   existingCategory.ParentID,
		Level:      existingCategory.Level,
		Path:       existingCategory.Path,
	})

	return nil
}