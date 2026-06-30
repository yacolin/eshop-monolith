package product

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

type SpuService struct {
	repo         IspuRepository
	categoryRepo IcategoryRepository
	brandRepo    IbrandRepository
	attrRepo     IattributeRepository
	db           *gorm.DB
	rdb          *redis.Client
	sf           singleflight.Group

	localCache  *spuLocalCache
	bloomFilter *spuBloomFilter
	hotCounter  *spuHotKeyCounter
}

func NewSpuService(
	repo IspuRepository,
	categoryRepo IcategoryRepository,
	brandRepo IbrandRepository,
	attrRepo IattributeRepository,
	db *gorm.DB,
	rdb *redis.Client,
) *SpuService {
	return &SpuService{
		repo:         repo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		attrRepo:     attrRepo,
		db:           db,
		rdb:          rdb,

		localCache:  newSPULocalCache(),
		bloomFilter: newSPUBloomFilter(),
		hotCounter:  newSPUHotKeyCounter(),
	}
}

// Create 创建商品（事务内写 SPU + SKU + Description + ProductAttribute）
func (s *SpuService) Create(ctx context.Context, req *CreateSPUReq) (*SPU, error) {
	// 校验类目
	if _, err := s.categoryRepo.FindByID(ctx, req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryNotFound
		}
		return nil, err
	}
	// 校验品牌
	if req.BrandID > 0 {
		if _, err := s.brandRepo.FindByID(ctx, req.BrandID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errcode.ErrBrandNotFound
			}
			return nil, err
		}
	}
	// 校验 SKU 编码重复（请求内）
	seen := map[string]bool{}
	for _, sku := range req.SKUs {
		if seen[sku.SkuCode] {
			return nil, fmt.Errorf("%w: %s", errcode.ErrSKUCodeExists, sku.SkuCode)
		}
		seen[sku.SkuCode] = true
	}
	// 校验 SKU 编码重复（数据库）
	for _, sku := range req.SKUs {
		existing, err := s.repo.FindSKUByCode(ctx, sku.SkuCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("%w: %s", errcode.ErrSKUCodeExists, sku.SkuCode)
		}
	}
	// 校验属性
	if len(req.Attributes) > 0 {
		attrIDs := map[int64]bool{}
		for _, pa := range req.Attributes {
			if attrIDs[pa.AttributeID] {
				return nil, errcode.ErrProductAttrDuplicate
			}
			attrIDs[pa.AttributeID] = true
			if _, err := s.attrRepo.FindByID(ctx, pa.AttributeID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errcode.ErrAttributeNotFound
				}
				return nil, err
			}
		}
	}

	// 序列化 images
	var imagesJSON string
	if len(req.Images) > 0 {
		b, _ := json.Marshal(req.Images)
		imagesJSON = string(b)
	}

	// 事务创建
	var spu *SPU
	err := s.db.Transaction(func(tx *gorm.DB) error {
		spu = &SPU{
			Name:        req.Name,
			Subtitle:    req.Subtitle,
			CategoryID:  req.CategoryID,
			BrandID:     req.BrandID,
			Unit:        req.Unit,
			MainImage:   req.MainImage,
			Images:      imagesJSON,
			VideoURL:    req.VideoURL,
			SortOrder:   req.SortOrder,
			CreatedBy:   req.CreatedBy,
			Status:      0, // draft
		}
		if req.Description != "" || req.MobileDesc != "" {
			spu.HasDescription = 1
		}
		if err := s.repo.CreateSPUWithTx(tx, spu); err != nil {
			return err
		}

		// 批量创建 SKU
		for _, sku := range req.SKUs {
			specJSON := "{}"
			if sku.Spec != nil {
				b, _ := json.Marshal(sku.Spec)
				specJSON = string(b)
			}
			skuModel := &SKU{
				ProductID:     spu.ID,
				SkuCode:       sku.SkuCode,
				Barcode:       sku.Barcode,
				Spec:          specJSON,
				Price:         sku.Price,
				MarketPrice:   sku.MarketPrice,
				CostPrice:     sku.CostPrice,
				Weight:        sku.Weight,
				Volume:        sku.Volume,
				Length:        sku.Length,
				Width:         sku.Width,
				Height:        sku.Height,
				MinPurchaseQty: sku.MinPurchase,
				MaxPurchaseQty: sku.MaxPurchase,
				Image:         sku.Image,
				Status:        1,
			}
			if skuModel.MinPurchaseQty <= 0 {
				skuModel.MinPurchaseQty = 1
			}
			if err := s.repo.CreateSKUWithTx(tx, skuModel); err != nil {
				return err
			}
		}

		// 创建图文详情
		if req.Description != "" || req.MobileDesc != "" {
			desc := &Description{
				ProductID:         spu.ID,
				Description:       req.Description,
				MobileDescription: req.MobileDesc,
			}
			if err := s.repo.CreateDescriptionWithTx(tx, desc); err != nil {
				return err
			}
		}

		// 创建属性值
		for _, pa := range req.Attributes {
			paModel := &ProductAttribute{
				ProductID:   spu.ID,
				AttributeID: pa.AttributeID,
				Value:       pa.Value,
			}
			if err := s.repo.CreateProductAttrWithTx(tx, paModel); err != nil {
				return err
			}
		}

		// 更新聚合价格
		if err := s.repo.UpdateSPUPriceWithTx(tx, spu.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if spu != nil {
		s.bloomFilter.add(spu.ID)
		s.localCache.setSingle(spu.ID, spu)
		if s.rdb != nil {
			delSPUEntity(context.Background(), s.rdb, spu.ID)
			delAllSPUListCache(context.Background(), s.rdb)
			delayedDeleteSPU(context.Background(), s.rdb, spu.ID)
		}
	}
	return spu, nil
}

// GetByID 获取商品详情，多级缓存: L1 → Bloom Filter → L2 Redis → DB
func (s *SpuService) GetByID(ctx context.Context, id int64) (*SPU, error) {
	// L1 本地缓存
	if spu, ok := s.localCache.getSingle(id); ok {
		s.hotCounter.increment(id)
		return spu, nil
	}

	// Bloom Filter 快速拦截（无假阴性，返回 false 则 ID 一定不存在）
	if !s.bloomFilter.mayExist(id) {
		return nil, errcode.ErrSPUNotFound
	}

	// L2 Redis
	if s.rdb != nil {
		if spu, err := getSPUEntity(ctx, s.rdb, id); err == nil {
			s.localCache.setSingle(id, spu)
			s.hotCounter.increment(id)
			return spu, nil
		}
	}

	// DB 兜底
	spu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrSPUNotFound
		}
		return nil, err
	}

	// 回填 L2 + L1
	s.bloomFilter.add(id)
	s.localCache.setSingle(id, spu)
	if s.rdb != nil {
		if err := setSPUEntity(ctx, s.rdb, spu); err != nil {
			_ = err
		}
	}
	return spu, nil
}

type cursorInfo struct {
	ID int64
}

func (s *SpuService) List(ctx context.Context, req *SPUListReq) (*SPUListResult, error) {
	req.Normalize()
	var cursor cursorInfo
	if req.Cursor != "" {
		if b, err := base64.StdEncoding.DecodeString(req.Cursor); err == nil {
			cursor.ID, _ = strconv.ParseInt(string(b), 10, 64)
		}
	}
	useZSET := s.rdb != nil && req.Name == "" && req.PriceMin == 0 && req.PriceMax == 0
	if useZSET {
		result, err := s.listFromZSET(ctx, req, cursor)
		if err == nil {
			return result, nil
		}
	}
	return s.listFromDB(ctx, req, cursor)
}

func (s *SpuService) listFromZSET(ctx context.Context, req *SPUListReq, cursor cursorInfo) (*SPUListResult, error) {
	key := cacheKeySPUListIDs(req.CategoryID, req.BrandID, req.Status)
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {
		exists, err := s.rdb.Exists(ctx, key).Result()
		if err != nil || exists == 0 {
			var all []SPU
			db := s.db.WithContext(ctx).Model(&SPU{}).Select("id")
			if req.CategoryID != nil {
				db = db.Where("category_id IN (SELECT id FROM sp_categories WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM sp_categories WHERE id = ?), ?, '/%'))", *req.CategoryID, *req.CategoryID, *req.CategoryID)
			}
			if req.BrandID != nil {
				db = db.Where("brand_id = ?", *req.BrandID)
			}
			if req.Status != nil {
				db = db.Where("status = ?", *req.Status)
			}
			if err := db.Order("id DESC").Find(&all).Error; err != nil {
				return nil, err
			}
			if len(all) == 0 {
				return &SPUListResult{List: []*SPU{}}, nil
			}
			if err := setSPUListIDs(ctx, s.rdb, req.CategoryID, req.BrandID, req.Status, all); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	if v != nil {
		return v.(*SPUListResult), nil
	}
	var cursorRank int64 = -1
	if req.Cursor != "" && cursor.ID > 0 {
		rank, err := getSPUListRank(ctx, s.rdb, req.CategoryID, req.BrandID, req.Status, cursor.ID)
		if err != nil || rank < 0 {
			return nil, fmt.Errorf("cursor not in cache")
		}
		cursorRank = rank + 1
	}
	ids, err := getSPUListIDs(ctx, s.rdb, req.CategoryID, req.BrandID, req.Status, cursorRank, req.Size+1)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return &SPUListResult{List: []*SPU{}}, nil
	}
	return s.buildListResult(ctx, ids, req.Size)
}

func (s *SpuService) listFromDB(ctx context.Context, req *SPUListReq, cursor cursorInfo) (*SPUListResult, error) {
	ids, err := s.repo.ListIDs(ctx, req.Name, req.CategoryID, req.BrandID, req.Status, req.PriceMin, req.PriceMax, req.Size+1, cursor.ID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return &SPUListResult{List: []*SPU{}}, nil
	}
	return s.buildListResult(ctx, ids, req.Size)
}

func (s *SpuService) buildListResult(ctx context.Context, ids []int64, size int) (*SPUListResult, error) {
	hasMore := len(ids) > size
	if hasMore {
		ids = ids[:size]
	}
	var list []*SPU
	if s.rdb != nil {
		hit, miss, err := batchFetchSPUEntities(ctx, s.rdb, ids)
		if err != nil {
			return s.fetchFullFromDB(ctx, ids, hasMore)
		}
		if len(miss) > 0 {
			var dbSPUs []SPU
			if err := s.db.WithContext(ctx).Where("id IN ?", miss).Find(&dbSPUs).Error; err != nil {
				return nil, err
			}
			for i := range dbSPUs {
				hit[dbSPUs[i].ID] = &dbSPUs[i]
				s.bloomFilter.add(dbSPUs[i].ID)
				s.localCache.setSingle(dbSPUs[i].ID, &dbSPUs[i])
				if err := setSPUEntity(ctx, s.rdb, &dbSPUs[i]); err != nil {
					_ = err
				}
			}
		}
		list = make([]*SPU, 0, len(ids))
		for _, id := range ids {
			if spu, ok := hit[id]; ok {
				list = append(list, spu)
			}
		}
	} else {
		return s.fetchFullFromDB(ctx, ids, hasMore)
	}
	cursor := ""
	if len(list) > 0 {
		last := list[len(list)-1]
		cursor = base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(last.ID, 10)))
	}
	return &SPUListResult{List: list, Cursor: cursor, HasMore: hasMore}, nil
}

// WarmupCache 全量预热 SPU 到 Bloom Filter + L2 + L1
func (s *SpuService) WarmupCache(ctx context.Context) (int, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return 0, err
	}
	if len(all) == 0 {
		return 0, nil
	}

	ids := make([]int64, len(all))
	for i := range all {
		ids[i] = all[i].ID
	}
	s.bloomFilter.addAll(ids)

	if s.rdb != nil {
		pipe := s.rdb.Pipeline()
		for i := range all {
			data, err := sonic.Marshal(&all[i])
			if err != nil {
				continue
			}
			pipe.Set(ctx, cacheKeySPU(all[i].ID), data, spuEntityTTL)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			_ = err
		}
	}

	for i := range all {
		s.localCache.warmupSingle(all[i].ID, &all[i])
	}
	return len(all), nil
}

func (s *SpuService) fetchFullFromDB(ctx context.Context, ids []int64, hasMore bool) (*SPUListResult, error) {
	var spus []SPU
	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&spus).Error; err != nil {
		return nil, err
	}
	spuMap := make(map[int64]*SPU, len(spus))
	for i := range spus {
		spuMap[spus[i].ID] = &spus[i]
	}
	list := make([]*SPU, 0, len(ids))
	for _, id := range ids {
		if spu, ok := spuMap[id]; ok {
			list = append(list, spu)
		}
	}
	cursor := ""
	if len(list) > 0 {
		last := list[len(list)-1]
		cursor = base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(last.ID, 10)))
	}
	return &SPUListResult{List: list, Cursor: cursor, HasMore: hasMore}, nil
}

// Update 更新商品基本信息
func (s *SpuService) Update(ctx context.Context, id int64, req *UpdateSPUReq) (*SPU, error) {
	spu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrSPUNotFound
		}
		return nil, err
	}
	if req.Name != nil {
		spu.Name = *req.Name
	}
	if req.Subtitle != nil {
		spu.Subtitle = *req.Subtitle
	}
	if req.Unit != nil {
		spu.Unit = *req.Unit
	}
	if req.MainImage != nil {
		spu.MainImage = *req.MainImage
	}
	if req.Images != nil {
		b, _ := json.Marshal(*req.Images)
		spu.Images = string(b)
	}
	if req.VideoURL != nil {
		spu.VideoURL = *req.VideoURL
	}
	if req.SortOrder != nil {
		spu.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		spu.Status = *req.Status
	}
	if req.UpdatedBy != nil {
		spu.UpdatedBy = *req.UpdatedBy
	}
	s.localCache.removeSingle(id)
	if s.rdb != nil {
		delSPUEntity(context.Background(), s.rdb, id)
		delAllSPUListCache(context.Background(), s.rdb)
	}
	if err := s.repo.Update(ctx, spu); err != nil {
		return nil, err
	}
	if s.rdb != nil {
		delayedDeleteSPU(context.Background(), s.rdb, id)
	}
	return spu, nil
}

// Delete 删除商品（软删除）
func (s *SpuService) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrSPUNotFound
		}
		return err
	}
	s.localCache.removeSingle(id)
	if s.rdb != nil {
		delSPUEntity(context.Background(), s.rdb, id)
		delAllSPUListCache(context.Background(), s.rdb)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if s.rdb != nil {
		delayedDeleteSPU(context.Background(), s.rdb, id)
	}
	return nil
}
