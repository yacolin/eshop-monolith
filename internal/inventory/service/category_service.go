package service

import (
	"context"

	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"
)

// CategoryService 分类服务
type CategoryService struct {
	repo    repositories.IcategoryRepository
	catAttr repositories.IcategoryAttributeRepository
	bus     *eventbus.Bus
}

// NewCategoryService 创建分类服务
func NewCategoryService(
	repo repositories.IcategoryRepository,
	catAttr repositories.IcategoryAttributeRepository,
	bus *eventbus.Bus,
) *CategoryService {
	return &CategoryService{
		repo:    repo,
		catAttr: catAttr,
		bus:     bus,
	}
}

type CategoryListResult struct {
	Total int64             `json:"total"`
	List  []models.Category `json:"list"`
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryDTO) (*models.Category, error) {
	// 创建分类
	newCategory := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	// 保存分类
	if err := s.repo.Create(ctx, newCategory); err != nil {
		return nil, err
	}

	// 发布分类创建事件
	s.bus.Publish(events.CategoryCreatedEvent{
		CategoryID: newCategory.ID,
		Name:       newCategory.Name,
		ParentID:   newCategory.ParentID,
	})

	return newCategory, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return category, nil
}

// ListAllCategories 列出所有分类
func (s *CategoryService) ListAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListAll(ctx)
}

// ListRootCategories 列出根分类
func (s *CategoryService) ListRootCategories(ctx context.Context) (*dto.CategoryListResult, error) {
	list, err := s.repo.ListRoot(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.CategoryListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

// ListSubCategories 列出子分类
func (s *CategoryService) ListSubCategories(ctx context.Context, parentID int64) ([]models.Category, error) {
	return s.repo.ListByParent(ctx, parentID)
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(ctx context.Context, id int64, req *dto.UpdateCategoryDTO) (*models.Category, error) {
	// 获取分类
	existingCategory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新分类信息
	if req.Name != nil {
		existingCategory.Name = *req.Name
	}
	if req.Description != nil {
		existingCategory.Description = *req.Description
	}
	if req.ParentID != nil {
		existingCategory.ParentID = req.ParentID
	}

	// 保存分类
	if err := s.repo.Update(ctx, existingCategory); err != nil {
		return nil, err
	}

	// 发布分类更新事件
	s.bus.Publish(events.CategoryUpdatedEvent{
		CategoryID: existingCategory.ID,
		Name:       existingCategory.Name,
		ParentID:   existingCategory.ParentID,
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
	s.bus.Publish(events.CategoryDeletedEvent{
		CategoryID: existingCategory.ID,
		Name:       existingCategory.Name,
		ParentID:   existingCategory.ParentID,
	})

	return nil
}

func (s *CategoryService) ListCategories(ctx context.Context, q dto.CategoryListQuery) (*CategoryListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.repo.ListCategories(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountCategories(ctx, q)
	if err != nil {
		return nil, err
	}

	return &CategoryListResult{
		List:  list,
		Total: total,
	}, nil
}


// ── 品类-属性关联 ──────────────────────────────────────────────

// GetCategoryAttributes 获取品类关联的属性列表
func (s *CategoryService) GetCategoryAttributes(ctx context.Context, categoryID int64) ([]dto.CategoryAttributeResponse, error) {
	attrs, err := s.catAttr.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.CategoryAttributeResponse, len(attrs))
	for i, a := range attrs {
		items[i] = dto.CategoryAttributeResponse{
			AttributeID:   a.ID,
			AttributeName: a.Name,
		}
	}
	return items, nil
}

// SetCategoryAttributes 全量替换品类关联的属性
func (s *CategoryService) SetCategoryAttributes(ctx context.Context, categoryID int64, req *dto.SetCategoryAttributesDTO) error {
	return s.catAttr.SetByCategoryID(ctx, categoryID, req.AttributeIDs)
}
