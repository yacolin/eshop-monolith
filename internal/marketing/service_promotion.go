package marketing

import (
	"context"
	"errors"
	"time"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/logger"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type PromotionService struct {
	repo IpromotionRepository
	db   *gorm.DB
	rdb  *redis.Client

	localCache  *promotionLocalCache
	bloomFilter *promotionBloomFilter
}

func NewPromotionService(repo IpromotionRepository, db *gorm.DB, rdb *redis.Client) *PromotionService {
	return &PromotionService{
		repo: repo,
		db:   db,
		rdb:  rdb,

		localCache:  newPromotionLocalCache(),
		bloomFilter: newPromotionBloomFilter(),
	}
}

type PromotionListResult struct {
	Total int64        `json:"total"`
	List  []*Promotion `json:"list"`
}

// spuBrief SPU 概要信息（促销商品展示用）
type spuBrief struct {
	ID        int64  `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	Subtitle  string `gorm:"column:subtitle"`
	MainImage string `gorm:"column:main_image"`
	Unit      string `gorm:"column:unit"`
	MinPrice  int64  `gorm:"column:min_price"`
	MaxPrice  int64  `gorm:"column:max_price"`
	SalesCount int   `gorm:"column:sales_count"`
	Status    int8   `gorm:"column:status"`
}

// Create 创建促销（含规则和商品范围）
func (s *PromotionService) Create(ctx context.Context, req *CreatePromotionReq) (*Promotion, error) {
	start, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
	end, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)

	p := &Promotion{
		PromoName:     req.PromoName,
		PromoType:     req.PromoType,
		PromoCode:     req.PromoCode,
		StartTime:     start,
		EndTime:       end,
		TotalQuantity: req.TotalQuantity,
		PerUserLimit:  req.PerUserLimit,
		Status:        1,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		// Create rule
		rule := &PromotionRule{
			PromotionID:    p.ID,
			RuleName:       req.RuleName,
			ConditionType:  req.ConditionType,
			ConditionValue: req.ConditionValue,
			BenefitType:    req.BenefitType,
			BenefitValue:   req.BenefitValue,
			IsStackable:    req.IsStackable,
			StackPriority:  req.StackPriority,
		}
		if err := tx.Create(rule).Error; err != nil {
			return err
		}
		// Set products
		if len(req.ProductIDs) > 0 {
			products := make([]PromotionProduct, len(req.ProductIDs))
			for i, pid := range req.ProductIDs {
				products[i] = PromotionProduct{
					PromotionID: p.ID,
					ProductType: 3,
					ProductID:   &pid,
				}
			}
			return tx.Create(&products).Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 回填缓存
	s.bloomFilter.add(p.ID)
	s.localCache.setSingle(p.ID, p)
	if s.rdb != nil {
		_ = setPromotionEntity(context.Background(), s.rdb, p)
	}
	return p, nil
}

// GetByID 获取促销，三级缓存: L1 → BloomFilter → L2 Redis → DB
func (s *PromotionService) GetByID(ctx context.Context, id int64) (*Promotion, error) {
	// L1 本地缓存
	if p, ok := s.localCache.getSingle(id); ok {
		return p, nil
	}

	// Bloom Filter 快速拦截
	if !s.bloomFilter.mayExist(id) {
		return nil, errcode.ErrPromotionNotFound
	}

	// L2 Redis
	if s.rdb != nil {
		if p, err := getPromotionEntity(ctx, s.rdb, id); err == nil {
			s.localCache.setSingle(id, p)
			return p, nil
		}
	}

	logger.Info("promotion cache miss, fallback to DB", "id", id)
	// DB 兜底
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}

	// 回填 L2 + L1
	s.bloomFilter.add(id)
	s.localCache.setSingle(id, p)
	if s.rdb != nil {
		_ = setPromotionEntity(ctx, s.rdb, p)
	}
	return p, nil
}

// GetDetailByID 获取促销详情（含规则、商品范围），Fan-Out 并发获取
func (s *PromotionService) GetDetailByID(ctx context.Context, id int64) (*PromotionDetailResponse, error) {
	// L2 缓存一次读完整 detail（不含 L1，detail 体积较大走 L2 即可）
	if s.rdb != nil {
		if d, err := getPromotionDetail(ctx, s.rdb, id); err == nil {
			return d, nil
		}
	}

	g, egCtx := errgroup.WithContext(ctx)

	var p *Promotion
	g.Go(func() error {
		var err error
		p, err = s.GetByID(egCtx, id)
		return err
	})

	var rule *PromotionRule
	g.Go(func() error {
		var err error
		rule, err = s.repo.FindRuleByPromotionID(egCtx, id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})

	var products []PromotionProduct
	g.Go(func() error {
		var err error
		products, err = s.repo.ListProducts(egCtx, id)
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// 批量查询 SPU 信息
	spuIDs := make([]int64, 0, len(products))
	for _, prod := range products {
		if prod.ProductType == 3 && prod.ProductID != nil {
			spuIDs = append(spuIDs, *prod.ProductID)
		}
	}
	spuMap := make(map[int64]spuBrief, len(spuIDs))
	if len(spuIDs) > 0 {
		var rows []spuBrief
		if err := s.db.WithContext(ctx).Table("sp_products").
			Where("id IN ?", spuIDs).
			Scan(&rows).Error; err == nil {
			for i := range rows {
				spuMap[rows[i].ID] = rows[i]
			}
		}
	}

	// 组装商品列表
	items := make([]PromotionProductItem, 0, len(products))
	for _, prod := range products {
		item := PromotionProductItem{
			ID:          prod.ID,
			ProductType: prod.ProductType,
			ProductID:   prod.ProductID,
			CategoryID:  prod.CategoryID,
		}
		if prod.ProductType == 3 && prod.ProductID != nil {
			if brief, ok := spuMap[*prod.ProductID]; ok {
				item.SpuName = brief.Name
				item.Subtitle = brief.Subtitle
				item.MainImage = brief.MainImage
				item.Unit = brief.Unit
				item.MinPrice = brief.MinPrice
				item.MaxPrice = brief.MaxPrice
				item.SalesCount = brief.SalesCount
				item.SpuStatus = brief.Status
			}
		}
		items = append(items, item)
	}
	if items == nil {
		items = []PromotionProductItem{}
	}

	resp := &PromotionDetailResponse{
		Promotion: p,
		Rule:      rule,
		Products:  items,
	}

	// 回填 L2
	if s.rdb != nil {
		_ = setPromotionDetail(ctx, s.rdb, resp)
	}
	return resp, nil
}

// WarmupCache 全量预热促销到 Bloom Filter + L2 + L1
func (s *PromotionService) WarmupCache(ctx context.Context) (int, error) {
	var all []Promotion
	if err := s.db.WithContext(ctx).Model(&Promotion{}).Find(&all).Error; err != nil {
		return 0, err
	}
	if len(all) == 0 {
		return 0, nil
	}

	g, _ := errgroup.WithContext(ctx)

	// Bloom Filter
	g.Go(func() error {
		ids := make([]int64, len(all))
		for i := range all {
			ids[i] = all[i].ID
		}
		s.bloomFilter.addAll(ids)
		return nil
	})

	// L1 本地缓存
	g.Go(func() error {
		for i := range all {
			s.localCache.setSingle(all[i].ID, &all[i])
		}
		return nil
	})

	// L2 Redis 实体缓存 + 清理旧 detail 缓存
	if s.rdb != nil {
		g.Go(func() error {
			pipe := s.rdb.Pipeline()
			for i := range all {
				data, marshalErr := sonic.Marshal(&all[i])
				if marshalErr != nil {
					continue
				}
				pipe.Set(ctx, cacheKeyPromotion(all[i].ID), data, promotionEntityTTL)
				pipe.Del(ctx, cacheKeyPromotionDetail(all[i].ID))
			}
			_, err := pipe.Exec(ctx)
			return err
		})
	}

	return len(all), g.Wait()
}

func (s *PromotionService) List(ctx context.Context, req *PromotionListReq) (*PromotionListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	list, total, err := s.repo.List(ctx, req.Status, req.PromoType, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*Promotion, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &PromotionListResult{Total: total, List: items}, nil
}

func (s *PromotionService) Update(ctx context.Context, id int64, req *UpdatePromotionReq) (*Promotion, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	if req.PromoName != nil {
		p.PromoName = *req.PromoName
	}
	if req.StartTime != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		p.StartTime = t
	}
	if req.EndTime != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		p.EndTime = t
	}
	if req.TotalQuantity != nil {
		p.TotalQuantity = *req.TotalQuantity
	}
	if req.PerUserLimit != nil {
		p.PerUserLimit = *req.PerUserLimit
	}
	if req.Status != nil {
		p.Status = *req.Status
	}

	// 先删缓存，再写 DB
	s.localCache.removeSingle(id)
	if s.rdb != nil {
		delPromotionEntity(context.Background(), s.rdb, id)
		delPromotionDetail(context.Background(), s.rdb, id)
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	if s.rdb != nil {
		delayedDeletePromotion(context.Background(), s.rdb, id)
	}

	// Update rule
	if req.RuleName != nil || req.ConditionType != nil || req.BenefitType != nil {
		rule, err := s.repo.FindRuleByPromotionID(ctx, id)
		if err == nil {
			if req.RuleName != nil {
				rule.RuleName = *req.RuleName
			}
			if req.ConditionType != nil {
				rule.ConditionType = *req.ConditionType
			}
			if req.ConditionValue != nil {
				rule.ConditionValue = *req.ConditionValue
			}
			if req.BenefitType != nil {
				rule.BenefitType = *req.BenefitType
			}
			if req.BenefitValue != nil {
				rule.BenefitValue = *req.BenefitValue
			}
			if req.IsStackable != nil {
				rule.IsStackable = *req.IsStackable
			}
			if req.StackPriority != nil {
				rule.StackPriority = *req.StackPriority
			}
			s.repo.UpdateRule(ctx, rule)
		}
	}
	return p, nil
}

func (s *PromotionService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrPromotionNotFound
		}
		return err
	}
	// 先删缓存
	s.localCache.removeSingle(id)
	if s.rdb != nil {
		delPromotionEntity(context.Background(), s.rdb, id)
		delPromotionDetail(context.Background(), s.rdb, id)
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("promotion_id = ?", id).Delete(&PromotionRule{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&Promotion{}).Error
	})
}
