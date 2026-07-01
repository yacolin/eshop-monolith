package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/logger"
)

type CategoryService struct {
	repo       IcategoryRepository
	rdb        *redis.Client
	allLocal   *simpleLocalCache[[]Category]
	entityLocal *simpleLocalCache[*Category]
}

func NewCategoryService(repo IcategoryRepository, rdb *redis.Client) *CategoryService {
	return &CategoryService{
		repo:        repo,
		rdb:         rdb,
		allLocal:    newSimpleLocalCache[[]Category](categoryLocalCacheSize, localCacheTTL),
		entityLocal: newSimpleLocalCache[*Category](categoryLocalCacheSize, localCacheTTL),
	}
}

const maxCategoryLevel = 3

// ── 缓存辅助 ──

func (s *CategoryService) getAllFromCache(ctx context.Context) ([]Category, error) {
	// L1
	if cached, ok := s.allLocal.get(0); ok {
		return cached, nil
	}
	// L2 Redis
	if s.rdb != nil {
		if cached, err := getCategoryAllCache(ctx, s.rdb); err == nil {
			s.allLocal.set(0, cached)
			return cached, nil
		}
	}
	// DB
	logger.Info("category cache miss, fallback to DB")
	all, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	if s.rdb != nil {
		setCategoryAllCache(ctx, s.rdb, all)
	}
	s.allLocal.set(0, all)
	return all, nil
}

// WarmupCache 预热类目缓存到 L1 + L2
func (s *CategoryService) WarmupCache(ctx context.Context) (int, error) {
	all, err := s.repo.ListAll(ctx)
	if err != nil {
		return 0, err
	}
	s.allLocal.set(0, all)
	if s.rdb != nil {
		if err := setCategoryAllCache(ctx, s.rdb, all); err != nil {
			return 0, err
		}
		for i := range all {
			setCategoryEntity(ctx, s.rdb, &all[i])
			s.entityLocal.set(all[i].ID, &all[i])
		}
	}
	return len(all), nil
}


func (s *CategoryService) invalidateCache(ctx context.Context) {
	s.allLocal.remove(0)
	s.entityLocal.clear()
	if s.rdb != nil {
		logger.Info("category cache invalidated")
		delCategoryAllCache(ctx, s.rdb)
	}
}

func (s *CategoryService) delayedInvalidate(ctx context.Context) {
	if s.rdb != nil {
		delayedDeleteCategoryAll(ctx, s.rdb)
	}
}

// ── Create ──

func (s *CategoryService) Create(ctx context.Context, req *CreateCategoryReq) (*Category, error) {
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

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	cat.Path = fmt.Sprintf("%s%d/", cat.Path, cat.ID)
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx)
	return cat, nil
}

// ── GetByID ──

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*Category, error) {
	// L1
	if cached, ok := s.entityLocal.get(id); ok {
		return cached, nil
	}
	// L2 Redis
	if s.rdb != nil {
		if cached, err := getCategoryEntity(ctx, s.rdb, id); err == nil {
			s.entityLocal.set(id, cached)
			return cached, nil
		}
	}
	// DB
		logger.Info("category entity cache miss, fallback to DB", "id", id)
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryNotFound
		}
		return nil, err
	}
	// 回填
	if s.rdb != nil {
		setCategoryEntity(ctx, s.rdb, cat)
	}
	s.entityLocal.set(id, cat)
	return cat, nil
}

// ── List ──

type CategoryListResult struct {
	Total int64       `json:"total"`
	List  []*Category `json:"list"`
}

func (s *CategoryService) List(ctx context.Context, req *CategoryListReq) (*CategoryListResult, error) {
	req.Normalize()
	list, total, err := s.repo.List(ctx, req.Name, req.Status, req.Level, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &CategoryListResult{Total: total, List: items}, nil
}

// ── 从全量缓存派生的查询 ──

// ListRoot 根类目
func (s *CategoryService) ListRoot(ctx context.Context) ([]*Category, error) {
	all, err := s.getAllFromCache(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, 0)
	for i := range all {
		if all[i].ParentID == 0 {
			items = append(items, &all[i])
		}
	}
	return items, nil
}

// ListChildren 子类目
func (s *CategoryService) ListChildren(ctx context.Context, parentID int64) ([]*Category, error) {
	all, err := s.getAllFromCache(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, 0)
	for i := range all {
		if all[i].ParentID == parentID {
			items = append(items, &all[i])
		}
	}
	return items, nil
}

// ListByLevel 按层级查询
func (s *CategoryService) ListByLevel(ctx context.Context, level int8) ([]*Category, error) {
	all, err := s.getAllFromCache(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, 0)
	for i := range all {
		if all[i].Level == level {
			items = append(items, &all[i])
		}
	}
	return items, nil
}

// ListAll 所有类目
func (s *CategoryService) ListAll(ctx context.Context) ([]*Category, error) {
	all, err := s.getAllFromCache(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*Category, len(all))
	for i := range all {
		items[i] = &all[i]
	}
	return items, nil
}

// ── Update ──

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

	s.invalidateCache(ctx)
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	s.delayedInvalidate(ctx)
	return cat, nil
}

// ── Delete ──

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
	s.invalidateCache(ctx)
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.delayedInvalidate(ctx)
	return nil
}
