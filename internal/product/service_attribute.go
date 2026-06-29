package product

import (
	"context"
	"errors"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type AttributeService struct {
	repo       IattributeRepository
	categoryRepo IcategoryRepository
}

func NewAttributeService(repo IattributeRepository, categoryRepo IcategoryRepository) *AttributeService {
	return &AttributeService{repo: repo, categoryRepo: categoryRepo}
}

func (s *AttributeService) Create(ctx context.Context, req *CreateAttributeReq) (*Attribute, error) {
	// 校验类目存在
	if _, err := s.categoryRepo.FindByID(ctx, req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryParentNotFound
		}
		return nil, err
	}

	attr := &Attribute{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		InputType:  req.InputType,
		Values:     req.Values,
		Unit:       req.Unit,
		Required:   req.Required,
		Searchable: req.Searchable,
		IsSkuSpec:  req.IsSkuSpec,
		SortOrder:  req.SortOrder,
		Status:     1,
	}
	if err := s.repo.Create(ctx, attr); err != nil {
		return nil, err
	}
	return attr, nil
}

func (s *AttributeService) GetByID(ctx context.Context, id int64) (*Attribute, error) {
	attr, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrAttributeNotFound
		}
		return nil, err
	}
	return attr, nil
}

func (s *AttributeService) ListByCategory(ctx context.Context, categoryID int64) ([]*Attribute, error) {
	list, err := s.repo.ListByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	items := make([]*Attribute, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

func (s *AttributeService) ListSearchable(ctx context.Context, categoryID int64) ([]*Attribute, error) {
	list, err := s.repo.ListSearchable(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	items := make([]*Attribute, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

func (s *AttributeService) ListSkuSpec(ctx context.Context, categoryID int64) ([]*Attribute, error) {
	list, err := s.repo.ListSkuSpec(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	items := make([]*Attribute, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

func (s *AttributeService) ListAll(ctx context.Context) ([]*Attribute, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Attribute, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return items, nil
}

func (s *AttributeService) Update(ctx context.Context, id int64, req *UpdateAttributeReq) (*Attribute, error) {
	attr, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrAttributeNotFound
		}
		return nil, err
	}
	if req.Name != nil {
		attr.Name = *req.Name
	}
	if req.InputType != nil {
		attr.InputType = *req.InputType
	}
	if req.Values != nil {
		attr.Values = *req.Values
	}
	if req.Unit != nil {
		attr.Unit = *req.Unit
	}
	if req.Required != nil {
		attr.Required = *req.Required
	}
	if req.Searchable != nil {
		attr.Searchable = *req.Searchable
	}
	if req.IsSkuSpec != nil {
		attr.IsSkuSpec = *req.IsSkuSpec
	}
	if req.SortOrder != nil {
		attr.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		attr.Status = *req.Status
	}
	if err := s.repo.Update(ctx, attr); err != nil {
		return nil, err
	}
	return attr, nil
}

func (s *AttributeService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrAttributeNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
