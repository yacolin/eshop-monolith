package product

import (
	"context"
	"errors"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type BrandService struct {
	repo IbrandRepository
}

func NewBrandService(repo IbrandRepository) *BrandService {
	return &BrandService{repo: repo}
}

func (s *BrandService) Create(ctx context.Context, req *CreateBrandReq) (*Brand, error) {
	existing, err := s.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errcode.ErrBrandNameExists
	}
	brand := &Brand{
		Name:        req.Name,
		EnglishName: req.EnglishName,
		LogoURL:     req.LogoURL,
		FirstLetter: req.FirstLetter,
		Description: req.Description,
	}
	if req.SortOrder != nil {
		brand.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		brand.Status = *req.Status
	}
	if err := s.repo.Create(ctx, brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (s *BrandService) GetByID(ctx context.Context, id int64) (*Brand, error) {
	brand, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrBrandNotFound
		}
		return nil, err
	}
	return brand, nil
}

type BrandListResult struct {
	Total int64    `json:"total"`
	List  []*Brand `json:"list"`
}

func (s *BrandService) List(ctx context.Context, req *BrandListReq) (*BrandListResult, error) {
	req.Normalize()
	list, total, err := s.repo.List(ctx, req.Name, req.FirstLetter, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*Brand, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &BrandListResult{Total: total, List: items}, nil
}

func (s *BrandService) Update(ctx context.Context, id int64, req *UpdateBrandReq) (*Brand, error) {
	brand, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrBrandNotFound
		}
		return nil, err
	}
	if req.Name != nil && *req.Name != brand.Name {
		existing, err := s.repo.FindByName(ctx, *req.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, errcode.ErrBrandNameExists
		}
		brand.Name = *req.Name
	}
	if req.EnglishName != nil {
		brand.EnglishName = *req.EnglishName
	}
	if req.LogoURL != nil {
		brand.LogoURL = *req.LogoURL
	}
	if req.FirstLetter != nil {
		brand.FirstLetter = *req.FirstLetter
	}
	if req.SortOrder != nil {
		brand.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		brand.Status = *req.Status
	}
	if req.Description != nil {
		brand.Description = *req.Description
	}
	if err := s.repo.Update(ctx, brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (s *BrandService) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrBrandNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
