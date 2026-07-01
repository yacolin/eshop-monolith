package product

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/logger"
)

type BrandService struct {
	repo          IbrandRepository
	rdb           *redis.Client
	listLocal     *simpleLocalCache[[]Brand]
	entityLocal   *simpleLocalCache[*Brand]
}

func NewBrandService(repo IbrandRepository, rdb *redis.Client) *BrandService {
	return &BrandService{
		repo:        repo,
		rdb:         rdb,
		listLocal:   newSimpleLocalCache[[]Brand](brandLocalCacheSize, localCacheTTL),
		entityLocal: newSimpleLocalCache[*Brand](brandLocalCacheSize, localCacheTTL),
	}
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
	s.invalidateCache(ctx)
	return brand, nil
}

func (s *BrandService) GetByID(ctx context.Context, id int64) (*Brand, error) {
	// L1
	if cached, ok := s.entityLocal.get(id); ok {
		return cached, nil
	}
	// L2 Redis
	if s.rdb != nil {
		if cached, err := getBrandEntity(ctx, s.rdb, id); err == nil {
			s.entityLocal.set(id, cached)
			return cached, nil
		}
	}
	// DB
	logger.Info("brand cache miss, fallback to DB", "id", id)
	brand, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrBrandNotFound
		}
		return nil, err
	}
	// 回填
	if s.rdb != nil {
		setBrandEntity(ctx, s.rdb, brand)
	}
	s.entityLocal.set(id, brand)
	return brand, nil
}

type BrandListResult struct {
	Total int64    `json:"total"`
	List  []*Brand `json:"list"`
}

func (s *BrandService) List(ctx context.Context, req *BrandListReq) (*BrandListResult, error) {
	req.Normalize()
	var allBrands []Brand

	// L1
	if cached, ok := s.listLocal.get(0); ok {
		allBrands = cached
	}
	// L2 Redis
	if allBrands == nil && s.rdb != nil {
		if cached, err := getBrandAllCache(ctx, s.rdb); err == nil {
			allBrands = cached
			s.listLocal.set(0, cached)
		}
	}
	// DB
	logger.Info("brand list cache miss, fallback to DB")
	if allBrands == nil {
		var err error
		allBrands, err = s.repo.FindAll(ctx)
		if err != nil {
			return nil, err
		}
		if s.rdb != nil {
			setBrandAllCache(ctx, s.rdb, allBrands)
		}
		s.listLocal.set(0, allBrands)
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
	s.invalidateCache(ctx)
	if err := s.repo.Update(ctx, brand); err != nil {
		return nil, err
	}
	if s.rdb != nil {
		delayedDeleteBrand(ctx, s.rdb)
		delayedDeleteBrandEntity(ctx, s.rdb, id)
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
	s.invalidateCache(ctx)
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if s.rdb != nil {
		delayedDeleteBrand(ctx, s.rdb)
		delayedDeleteBrandEntity(ctx, s.rdb, id)
	}
	return nil
}


// WarmupCache 预热品牌缓存到 L1 + L2
func (s *BrandService) WarmupCache(ctx context.Context) (int, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return 0, err
	}
	s.listLocal.set(0, all)
	if s.rdb != nil {
		if err := setBrandAllCache(ctx, s.rdb, all); err != nil {
			return 0, err
		}
		for i := range all {
			setBrandEntity(ctx, s.rdb, &all[i])
			s.entityLocal.set(all[i].ID, &all[i])
		}
	}
	return len(all), nil
}

func (s *BrandService) invalidateCache(ctx context.Context) {
	s.listLocal.remove(0)
	s.entityLocal.clear()
	if s.rdb != nil {
		logger.Info("brand cache invalidated")
		delBrandAllCache(ctx, s.rdb)
	}
}
