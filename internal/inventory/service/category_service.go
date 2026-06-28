package service

import (
	"context"
	"strconv"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const (
	categoryCacheZSet     = "category:zset"
	categoryNonRootZSet   = "category:non-root:zset"
	categoryInfoPrefix    = "category:info:"
)

// CategoryService 分类服务
type CategoryService struct {
	repo        repositories.IcategoryRepository
	catAttr     repositories.IcategoryAttributeRepository
	rabbit      *rabbitmq.Client
	rdb         *redis.Client
	singleGroup singleflight.Group
}

// NewCategoryService 创建分类服务
func NewCategoryService(
	repo repositories.IcategoryRepository,
	catAttr repositories.IcategoryAttributeRepository,
	rabbit *rabbitmq.Client,
	rdb *redis.Client,
) *CategoryService {
	return &CategoryService{
		repo:    repo,
		catAttr: catAttr,
		rabbit:  rabbit,
		rdb:     rdb,
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
	s.rabbit.Publish(ctx,events.CategoryCreatedEvent{
		CategoryID: newCategory.ID,
		Name:       newCategory.Name,
		ParentID:   newCategory.ParentID,
	})

	// 写入缓存
	cacheItem := dto.CachedCategoryItem{
		ID:   newCategory.ID,
		Name: newCategory.Name,
	}
	if data, err := sonic.Marshal(cacheItem); err == nil {
		s.rdb.Set(context.Background(), categoryInfoPrefix+strconv.FormatInt(newCategory.ID, 10), data, 0)
		s.rdb.ZAdd(context.Background(), categoryCacheZSet, redis.Z{
			Score:  float64(newCategory.ID),
			Member: newCategory.ID,
		})
		if newCategory.ParentID != nil {
			s.rdb.ZAdd(context.Background(), categoryNonRootZSet, redis.Z{
				Score:  float64(newCategory.ID),
				Member: newCategory.ID,
			})
		}
	}

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

// ListNonRootCategories 从缓存列出非根分类
func (s *CategoryService) ListNonRootCategories(ctx context.Context) (*dto.CachedCategoryListResult, error) {
	ctxBg := context.Background()

	ids, err := s.rdb.ZRange(ctxBg, categoryNonRootZSet, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return &dto.CachedCategoryListResult{}, nil
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = categoryInfoPrefix + id
	}

	values, err := s.rdb.MGet(ctxBg, keys...).Result()
	if err != nil {
		return nil, err
	}

	items := make([]dto.CachedCategoryItem, 0, len(values))
	for _, data := range values {
		if data == nil {
			continue
		}
		var item dto.CachedCategoryItem
		if err := sonic.Unmarshal([]byte(data.(string)), &item); err != nil {
			continue
		}
		items = append(items, item)
	}

	return &dto.CachedCategoryListResult{
		Total: len(items),
		List:  items,
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
	s.rabbit.Publish(ctx,events.CategoryUpdatedEvent{
		CategoryID: existingCategory.ID,
		Name:       existingCategory.Name,
		ParentID:   existingCategory.ParentID,
	})

	// 更新缓存
	cacheItem := dto.CachedCategoryItem{
		ID:   existingCategory.ID,
		Name: existingCategory.Name,
	}
	if data, err := sonic.Marshal(cacheItem); err == nil {
		s.rdb.Set(context.Background(), categoryInfoPrefix+strconv.FormatInt(existingCategory.ID, 10), data, 0)
	}
	if existingCategory.ParentID != nil {
		s.rdb.ZAdd(context.Background(), categoryNonRootZSet, redis.Z{
			Score:  float64(existingCategory.ID),
			Member: existingCategory.ID,
		})
	} else {
		s.rdb.ZRem(context.Background(), categoryNonRootZSet, existingCategory.ID)
	}

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
	s.rabbit.Publish(ctx,events.CategoryDeletedEvent{
		CategoryID: existingCategory.ID,
		Name:       existingCategory.Name,
		ParentID:   existingCategory.ParentID,
	})

	// 删除缓存
	s.rdb.Del(context.Background(), categoryInfoPrefix+strconv.FormatInt(id, 10))
	s.rdb.ZRem(context.Background(), categoryCacheZSet, id)
	s.rdb.ZRem(context.Background(), categoryNonRootZSet, id)

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

// WarmupCategoryCache 全量预热分类到 Redis
func (s *CategoryService) WarmupCategoryCache(ctx context.Context) (int, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return 0, err
	}

	pipe := s.rdb.Pipeline()
	ctxBg := context.Background()

	pipe.Del(ctxBg, categoryCacheZSet)
	pipe.Del(ctxBg, categoryNonRootZSet)

	for _, c := range list {
		item := dto.CachedCategoryItem{
			ID:   c.ID,
			Name: c.Name,
		}
		data, err := sonic.Marshal(item)
		if err != nil {
			return 0, err
		}
		pipe.Set(ctxBg, categoryInfoPrefix+strconv.FormatInt(c.ID, 10), data, 0)
		pipe.ZAdd(ctxBg, categoryCacheZSet, redis.Z{
			Score:  float64(c.ID),
			Member: c.ID,
		})
		if c.ParentID != nil {
			pipe.ZAdd(ctxBg, categoryNonRootZSet, redis.Z{
				Score:  float64(c.ID),
				Member: c.ID,
			})
		}
	}

	if _, err = pipe.Exec(ctxBg); err != nil {
		return 0, err
	}

	return len(list), nil
}

// ListCachedCategories 从缓存列出全部分类
func (s *CategoryService) ListCachedCategories(ctx context.Context) ([]dto.CachedCategoryItem, error) {
	ctxBg := context.Background()

	ids, err := s.rdb.ZRange(ctxBg, categoryCacheZSet, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []dto.CachedCategoryItem{}, nil
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = categoryInfoPrefix + id
	}

	values, err := s.rdb.MGet(ctxBg, keys...).Result()
	if err != nil {
		return nil, err
	}

	items := make([]dto.CachedCategoryItem, 0, len(values))
	for _, data := range values {
		if data == nil {
			continue
		}
		var item dto.CachedCategoryItem
		if err := sonic.Unmarshal([]byte(data.(string)), &item); err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// GetCachedCategoryByID 从缓存查询单个分类
func (s *CategoryService) GetCachedCategoryByID(ctx context.Context, id int64) (*dto.CachedCategoryItem, error) {
	sfKey := "category:" + strconv.FormatInt(id, 10)
	v, err, _ := s.singleGroup.Do(sfKey, func() (any, error) {
		data, err := s.rdb.Get(context.Background(), categoryInfoPrefix+strconv.FormatInt(id, 10)).Bytes()
		if err == redis.Nil {
			return nil, errcode.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		var item dto.CachedCategoryItem
		if err := sonic.Unmarshal(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*dto.CachedCategoryItem), nil
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
