package product

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

type BrandService struct {
	repo IbrandRepository
	rdb  *redis.Client
}

func NewBrandService(repo IbrandRepository, rdb *redis.Client) *BrandService {
	return &BrandService{repo: repo, rdb: rdb}
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
	if s.rdb != nil {
		delBrandAllCache(context.Background(), s.rdb)
		delayedDeleteBrand(context.Background(), s.rdb)
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
	var allBrands []Brand
	if s.rdb != nil {
		cached, err := getBrandAllCache(ctx, s.rdb)
		if err == nil {
			allBrands = cached
		}
	}
	if allBrands == nil {
		var err error
		allBrands, err = s.repo.FindAll(ctx)
		if err != nil {
			return nil, err
		}
		if s.rdb != nil {
			if err := setBrandAllCache(ctx, s.rdb, allBrands); err != nil {
				_ = err
			}
		}
	}
	filtered := make([]Brand, 0, len(allBrands))
	for _, b := range allBrands {
		if req.Name != "" && !strings.Contains(b.Name, req.Name) {
			continue
		}
		if req.FirstLetter != "" && b.FirstLetter != req.FirstLetter {
			continue
		}
		if req.Status != nil && b.Status != *req.Status {
			continue
		}
		filtered = append(filtered, b)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].SortOrder != filtered[j].SortOrder {
			return filtered[i].SortOrder > filtered[j].SortOrder
		}
		return filtered[i].ID < filtered[j].ID
	})
	total := int64(len(filtered))
	offset := (req.Page - 1) * req.Size
	if offset >= len(filtered) {
		return &BrandListResult{Total: total, List: []*Brand{}}, nil
	}
	end := offset + req.Size
	if end > len(filtered) {
		end = len(filtered)
	}
	items := make([]*Brand, end-offset)
	for i := range items {
		items[i] = &filtered[offset+i]
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
	if s.rdb != nil {
		delBrandAllCache(context.Background(), s.rdb)
	}
	if err := s.repo.Update(ctx, brand); err != nil {
		return nil, err
	}
	if s.rdb != nil {
		delayedDeleteBrand(context.Background(), s.rdb)
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
	if s.rdb != nil {
		delBrandAllCache(context.Background(), s.rdb)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if s.rdb != nil {
		delayedDeleteBrand(context.Background(), s.rdb)
	}
	return nil
}
