package product

import (
	"context"
	"errors"
	"fmt"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type CategoryService struct {
	repo IcategoryRepository
}

func NewCategoryService(repo IcategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

const maxCategoryLevel = 3

// Create 创建类目，自动计算 level 和 path
func (s *CategoryService) Create(ctx context.Context, req *CreateCategoryReq) (*Category, error) {
	// 同级名称唯一
	existing, err := s.repo.FindByName(ctx, req.Name, req.ParentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errcode.ErrCategoryNameExists
	}

	cat := &Category{
		Name:      req.Name,
		ParentID:  req.ParentID,
		IconURL:   req.IconURL,
		SortOrder: req.SortOrder,
		Status:    1,
	}

	// 计算层级
	if req.ParentID == 0 {
		cat.Level = 1
	} else {
		parent, err := s.repo.FindByID(ctx, req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errcode.ErrCategoryParentNotFound
			}
			return nil, err
		}
		if parent.Level >= maxCategoryLevel {
			return nil, errcode.ErrCategoryLevelExceed
		}
		cat.Level = parent.Level + 1
	}

	if req.SortOrder > 0 {
		cat.SortOrder = req.SortOrder
	}

	// 先创建（获取 ID），再更新 path
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	cat.Path = fmt.Sprintf("%s%d/", cat.Path, cat.ID)
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*Category, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryNotFound
		}
		return nil, err
	}
	return cat, nil
}

type CategoryListResult struct {
	Total int64       `json:"total"`
	List  []*Category `json:"list"`
}

// ListRoot 根类目
func (s *CategoryService) ListRoot(ctx context.Context) ([]*Category, error) {
	list, err := s.repo.ListRoot(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

// ListChildren 子类目
func (s *CategoryService) ListChildren(ctx context.Context, parentID int64) ([]*Category, error) {
	list, err := s.repo.ListByParent(ctx, parentID)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

// ListByLevel 按层级查询
func (s *CategoryService) ListByLevel(ctx context.Context, level int8) ([]*Category, error) {
	list, err := s.repo.ListByLevel(ctx, level)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

// ListAll 所有类目
func (s *CategoryService) ListAll(ctx context.Context) ([]*Category, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

// Update 更新类目（不修改 parent_id/level/path）
func (s *CategoryService) Update(ctx context.Context, id int64, req *UpdateCategoryReq) (*Category, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryNotFound
		}
		return nil, err
	}

	if req.Name != nil && *req.Name != cat.Name {
		existing, err := s.repo.FindByName(ctx, *req.Name, cat.ParentID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, errcode.ErrCategoryNameExists
		}
		cat.Name = *req.Name
	}
	if req.IconURL != nil {
		cat.IconURL = *req.IconURL
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		cat.Status = *req.Status
	}

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// Delete 删除类目（有子节点时禁止删除）
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrCategoryNotFound
		}
		return err
	}
	count, err := s.repo.CountByParentID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errcode.ErrCategoryHasChildren
	}
	return s.repo.Delete(ctx, id)
}
