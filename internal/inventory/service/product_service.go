package service

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"eshop-monolith/internal/infra/domain/shared"
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/query"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	productCacheZSet     = "product:zset"
	productInfoPrefix    = "product:info:"
	productCategoryZSet  = "product:zset:category:" // + categoryID
	cachedProductListTTL = time.Hour
	hotKeyThreshold      = 1000
	hotKeyWindow         = 10 * time.Second
	emptyPlaceholder     = "__EMPTY__"
	emptyCacheTTL        = 30 * time.Second
)

var zrangeMGetScript = redis.NewScript(`
local ids
if ARGV[3] == "desc" then
	ids = redis.call("ZREVRANGE", KEYS[1], ARGV[1], ARGV[2])
else
	ids = redis.call("ZRANGE", KEYS[1], ARGV[1], ARGV[2])
end
local total = redis.call("ZCARD", KEYS[1])
if #ids == 0 then return {total, {}} end
local keys = {}
for i, id in ipairs(ids) do
	keys[i] = ARGV[4] .. id
end
local values = redis.call("MGET", unpack(keys))
return {total, values}
`)

var zrangeByScoreScript = redis.NewScript(`
	local ids = redis.call("ZRANGEBYSCORE", KEYS[1], "(" .. ARGV[1], "+inf", "LIMIT", 0, ARGV[2])
	local total = redis.call("ZCARD", KEYS[1])
	if #ids == 0 then return {total, {}, 0} end
	local keys = {}
	for i, id in ipairs(ids) do
		keys[i] = ARGV[3] .. id
	end
	local values = redis.call("MGET", unpack(keys))
	return {total, values, ids[#ids]}
`)

func categoryZSetKey(categoryID int64) string {
	return productCategoryZSet + strconv.FormatInt(categoryID, 10)
}

type hotKeyCounter struct {
	mu       sync.Mutex
	counters map[int64]*hotKeyEntry
}

type hotKeyEntry struct {
	count int64
	start time.Time
}

func newHotKeyCounter() *hotKeyCounter {
	return &hotKeyCounter{counters: make(map[int64]*hotKeyEntry)}
}

func (h *hotKeyCounter) increment(id int64) bool {
	h.mu.Lock()
	entry, ok := h.counters[id]
	now := time.Now()
	if !ok || now.Sub(entry.start) > hotKeyWindow {
		entry = &hotKeyEntry{start: now}
		h.counters[id] = entry
	}
	entry.count++
	hot := entry.count >= hotKeyThreshold
	h.mu.Unlock()
	return hot
}

func (h *hotKeyCounter) reset(id int64) {
	h.mu.Lock()
	delete(h.counters, id)
	h.mu.Unlock()
}

type ProductService struct {
	repo          repositories.IproductRepository
	inventoryRepo repositories.IinventoryRepository
	skuRepo       repositories.IskuRepository
	rabbit        *rabbitmq.Client
	db            *gorm.DB
	rdb           *redis.Client
	localCache    *productLocalCache
	bloomFilter   *productBloomFilter
	singleGroup   singleflight.Group
	hotCounter    *hotKeyCounter
	metrics       *cacheMetrics
}

func NewProductService(
	repo repositories.IproductRepository,
	inventoryRepo repositories.IinventoryRepository,
	skuRepo repositories.IskuRepository,
	rabbit *rabbitmq.Client,
	db *gorm.DB,
	rdb *redis.Client,
) *ProductService {
	return &ProductService{
		repo:          repo,
		inventoryRepo: inventoryRepo,
		skuRepo:       skuRepo,
		rabbit:        rabbit,
		db:            db,
		rdb:           rdb,
		localCache:    newProductLocalCache(),
		bloomFilter:   newProductBloomFilter(),
		hotCounter:    newHotKeyCounter(),
		metrics:       getCacheMetrics(),
	}
}

func (s *ProductService) evictProductCache(id int64) {
	s.localCache.removeSingle(id)
	s.rdb.Del(context.Background(), productInfoPrefix+strconv.FormatInt(id, 10))
	s.rdb.ZRem(context.Background(), productCacheZSet, id)
	s.metrics.incEviction()
}

func (s *ProductService) delayedDoubleDelete(id int64) {
	s.evictProductCache(id)
	time.AfterFunc(500*time.Millisecond, func() {
		s.evictProductCache(id)
	})
}

func (s *ProductService) WarmupProductCache(ctx context.Context) (int, error) {
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return 0, err
	}

	items := make([]dto.CachedProductItem, 0, len(products))
	ids := make([]int64, 0, len(products))
	pipe := s.rdb.Pipeline()
	ctxBg := context.Background()

	pipe.Del(ctxBg, productCacheZSet)

	for _, p := range products {
		item := dto.CachedProductItem{
			ID:       p.ID,
			Name:     p.Name,
			MinPrice: p.MinPrice,
		}
		items = append(items, item)
		ids = append(ids, p.ID)
		data, err := sonic.Marshal(item)
		if err != nil {
			return 0, err
		}
		pipe.Set(ctxBg, productInfoPrefix+strconv.FormatInt(p.ID, 10), data, 0)
		pipe.ZAdd(ctxBg, productCacheZSet, redis.Z{
			Score:  float64(p.ID),
			Member: p.ID,
		})
	}

	// 加载所有产品-分类关系，构建每个分类的 ZSET
	var rels []struct {
		ProductID  int64 `gorm:"column:product_id"`
		CategoryID int64 `gorm:"column:category_id"`
	}
	if err := s.db.WithContext(ctxBg).Table("product_categories").Find(&rels).Error; err != nil {
		return 0, err
	}

	// 按 category_id 分组
	catMap := make(map[int64][]int64)
	for _, rel := range rels {
		catMap[rel.CategoryID] = append(catMap[rel.CategoryID], rel.ProductID)
	}
	for catID, productIDs := range catMap {
		key := categoryZSetKey(catID)
		pipe.Del(ctxBg, key)
		for _, pid := range productIDs {
			pipe.ZAdd(ctxBg, key, redis.Z{Score: float64(pid), Member: pid})
		}
	}

	if _, err = pipe.Exec(ctxBg); err != nil {
		return 0, err
	}

	s.bloomFilter.clear()
	s.bloomFilter.addAll(ids)
	s.localCache.warmup(items)

	// 预加载热门列表页（前 3 页）到 L1 缓存，避免启动后首次请求 miss
	s.warmupListPages(items)

	return len(products), nil
}

// warmupListPages 预加载热门列表页到 L1 本地缓存
func (s *ProductService) warmupListPages(items []dto.CachedProductItem) {
	const pageSize = 10
	total := int64(len(items))
	for page := 1; page <= 3; page++ {
		start := (page - 1) * pageSize
		if start >= len(items) {
			break
		}
		end := start + pageSize
		if end > len(items) {
			end = len(items)
		}

		// ASC 页
		ascList := make([]dto.CachedProductItem, end-start)
		copy(ascList, items[start:end])
		ascKey := "list:" + strconv.Itoa(page) + ":" + strconv.Itoa(pageSize) + ":asc"
		s.localCache.setList(ascKey, &query.ListResult[dto.CachedProductItem]{
			List:  ascList,
			Total: total,
		})

		// DESC 页（从尾部取，倒序）
		descStart := len(items) - end
		if descStart < 0 {
			descStart = 0
		}
		descEnd := len(items) - start
		if descStart < descEnd {
			descItems := make([]dto.CachedProductItem, descEnd-descStart)
			for i, j := 0, descEnd-1; j >= descStart; i, j = i+1, j-1 {
				descItems[i] = items[j]
			}
			descKey := "list:" + strconv.Itoa(page) + ":" + strconv.Itoa(pageSize) + ":desc"
			s.localCache.setList(descKey, &query.ListResult[dto.CachedProductItem]{
				List:  descItems,
				Total: total,
			})
		}
	}
}

func (s *ProductService) ListCachedProducts(ctx context.Context, q dto.ProductListQuery) (*query.ListResult[dto.CachedProductItem], error) {
	var catPart string
	zsetKey := productCacheZSet
	if q.CategoryID != nil {
		catPart = "cat_" + strconv.FormatInt(*q.CategoryID, 10) + ":"
		zsetKey = categoryZSetKey(*q.CategoryID)
	}
	cacheKey := "list:" + catPart + strconv.FormatInt(int64(q.Page), 10) + ":" + strconv.FormatInt(int64(q.Size), 10) + ":" + q.Order
	if result, ok := s.localCache.getList(cacheKey); ok {
		return result, nil
	}

	ctxBg := context.Background()
	offset := int64((q.Page - 1) * q.Size)
	stop := offset + int64(q.Size) - 1

	order := "asc"
	if q.Order == "desc" {
		order = "desc"
	}

	result, err := zrangeMGetScript.Run(ctxBg, s.rdb, []string{zsetKey}, offset, stop, order, productInfoPrefix).Result()
	if err != nil {
		return nil, err
	}

	values, ok := result.([]any)
	if !ok || len(values) != 2 {
		return &query.ListResult[dto.CachedProductItem]{Total: 0}, nil
	}

	total, ok := values[0].(int64)
	if !ok {
		if totalFloat, ok := values[0].(float64); ok {
			total = int64(totalFloat)
		} else {
			return nil, nil
		}
	}

	items, _ := values[1].([]any)
	products := make([]dto.CachedProductItem, 0, len(items))
	for _, data := range items {
		if data == nil {
			continue
		}
		var p dto.CachedProductItem
		if err := sonic.Unmarshal([]byte(data.(string)), &p); err != nil {
			continue
		}
		products = append(products, p)
	}

	listResult := &query.ListResult[dto.CachedProductItem]{List: products, Total: total}
	s.localCache.setList(cacheKey, listResult)
	return listResult, nil
}

// ListCachedProductsByCursor 基于游标从缓存查询产品列表（深分页优化）
// 支持 category_id 筛选，通过分类 ZSET 实现
func (s *ProductService) ListCachedProductsByCursor(ctx context.Context, q dto.ProductCacheCursorQuery) (*dto.ProductCacheCursorResult, error) {
	ctxBg := context.Background()
	limit := q.Size + 1

	zsetKey := productCacheZSet
	if q.CategoryID != nil {
		zsetKey = categoryZSetKey(*q.CategoryID)
	}

	result, err := zrangeByScoreScript.Run(ctxBg, s.rdb, []string{zsetKey}, q.Cursor, limit, productInfoPrefix).Result()
	if err != nil {
		return nil, err
	}

	values, ok := result.([]any)
	if !ok || len(values) != 3 {
		return &dto.ProductCacheCursorResult{}, nil
	}

	items, _ := values[1].([]any)
	products := make([]dto.CachedProductItem, 0, len(items))
	for _, data := range items {
		if data == nil {
			continue
		}
		var p dto.CachedProductItem
		if err := sonic.Unmarshal([]byte(data.(string)), &p); err != nil {
			continue
		}
		products = append(products, p)
	}

	res := &dto.ProductCacheCursorResult{}
	if len(products) > q.Size {
		res.List = products[:q.Size]
		res.NextCursor = products[q.Size-1].ID
		res.HasMore = true
	} else {
		res.List = products
		res.NextCursor = 0
		res.HasMore = false
	}
	return res, nil
}

func (s *ProductService) GetCachedProductByID(ctx context.Context, id int64) (*dto.CachedProductItem, error) {
	if !s.bloomFilter.mayExist(id) {
		s.metrics.incBloomReject()
		return nil, errcode.ErrProductNotFound
	}

	if item, ok := s.localCache.getSingle(id); ok {
		s.metrics.incL1Hit()
		s.hotCounter.increment(id)
		return item, nil
	}
	s.metrics.incL1Miss()

	sfKey := "product:" + strconv.FormatInt(id, 10)
	v, err, shared := s.singleGroup.Do(sfKey, func() (any, error) {
		var dbItem *dto.CachedProductItem

		start := time.Now()
		data, err := s.rdb.Get(context.Background(), productInfoPrefix+strconv.FormatInt(id, 10)).Bytes()
		if err == redis.Nil {
			s.metrics.incL2Miss()
			s.rdb.Set(context.Background(), productInfoPrefix+strconv.FormatInt(id, 10), emptyPlaceholder, emptyCacheTTL)
			s.metrics.observeDuration("l2", time.Since(start).Seconds())
			return nil, errcode.ErrProductNotFound
		}
		if err != nil {
			start = time.Now()
			s.metrics.incDBFallback()
			product, dbErr := s.repo.FindByID(context.Background(), id)
			if dbErr != nil {
				return nil, errcode.ErrProductNotFound
			}
			dbItem = &dto.CachedProductItem{
				ID:       product.ID,
				Name:     product.Name,
				MinPrice: product.MinPrice,
			}
			s.localCache.setSingle(id, dbItem)
			s.metrics.observeDuration("db", time.Since(start).Seconds())
			return dbItem, nil
		}

		s.metrics.incL2Hit()
		s.metrics.observeDuration("l2", time.Since(start).Seconds())
		if string(data) == emptyPlaceholder {
			return nil, errcode.ErrProductNotFound
		}

		var item dto.CachedProductItem
		if err := sonic.Unmarshal(data, &item); err != nil {
			return nil, err
		}

		s.localCache.setSingle(id, &item)
		return &item, nil
	})

	if shared {
		s.metrics.incSingleflightDedup()
	}

	if err != nil {
		return nil, err
	}
	return v.(*dto.CachedProductItem), nil
}

func (s *ProductService) GetProductWithInventory(ctx context.Context, id int64) (*dto.ProductDetailDTO, error) {
	type prodResult struct {
		product *models.Product
		err     error
	}
	type invResult struct {
		inventory *models.Inventory
		err       error
	}

	prodCh := make(chan prodResult, 1)
	invCh := make(chan invResult, 1)

	go func() {
		p, err := s.repo.FindByID(ctx, id)
		prodCh <- prodResult{product: p, err: err}
	}()

	go func() {
		inv, err := s.inventoryRepo.FindInventoryBySkuID(ctx, id)
		invCh <- invResult{inventory: inv, err: err}
	}()

	pr := <-prodCh
	if pr.err != nil {
		return nil, pr.err
	}

	ir := <-invCh
	detail := &dto.ProductDetailDTO{
		ID:          pr.product.ID,
		Name:        pr.product.Name,
		Description: pr.product.Description,
		MinPrice:    pr.product.MinPrice,
		CreatedAt:   pr.product.CreatedAt,
		UpdatedAt:   pr.product.UpdatedAt,
	}
	if ir.err == nil && ir.inventory != nil {
		detail.Quantity = ir.inventory.Quantity
		detail.Status = ir.inventory.Status
		detail.Reserved = ir.inventory.Reserved
		detail.Threshold = ir.inventory.Threshold
	}

	// 补充分类信息
	detail.Categories, _ = s.getCategoryInfo(ctx, id)

	return detail, nil
}

// getCategoryInfo 查询产品的所有分类
func (s *ProductService) getCategoryInfo(ctx context.Context, productID int64) ([]dto.ProductCategoryBrief, error) {
	var categories []dto.ProductCategoryBrief
	err := s.db.WithContext(ctx).
		Table("product_categories").
		Select("categories.id, categories.name").
		Joins("JOIN categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Order("categories.id ASC").
		Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetProductWithCategory 获取产品详情（含分类信息）
func (s *ProductService) GetProductWithCategory(ctx context.Context, id int64) (*dto.ProductWithCategoryDTO, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &dto.ProductWithCategoryDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		MinPrice:    product.MinPrice,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	// 一次额外查询补全分类信息
	result.Categories, _ = s.getCategoryInfo(ctx, id)

	return result, nil
}

// ListProductsWithCategory 列出产品（含分类信息，子查询单次 LEFT JOIN）
func (s *ProductService) ListProductsWithCategory(ctx context.Context, q dto.ProductListQuery) (*query.ListResult[dto.ProductWithCategoryDTO], error) {
	offset := (q.Page - 1) * q.Size

	rows, err := s.repo.ListProductsEnriched(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountProducts(ctx, q)
	if err != nil {
		return nil, err
	}

	// 按 product_id 分组，合并同一产品的多个分类
	productMap := make(map[int64]*dto.ProductWithCategoryDTO, len(rows))
	var order []int64
	for _, row := range rows {
		if _, ok := productMap[row.ID]; !ok {
			productMap[row.ID] = &dto.ProductWithCategoryDTO{
				ID:          row.ID,
				Name:        row.Name,
				Description: row.Description,
				MinPrice:    row.MinPrice,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			}
			order = append(order, row.ID)
		}
		if row.CategoryID != nil {
			productMap[row.ID].Categories = append(productMap[row.ID].Categories, dto.ProductCategoryBrief{
				ID:   *row.CategoryID,
				Name: *row.CategoryName,
			})
		}
	}

	enrichedList := make([]dto.ProductWithCategoryDTO, len(order))
	for i, id := range order {
		enrichedList[i] = *productMap[id]
	}

	return &query.ListResult[dto.ProductWithCategoryDTO]{
		List:  enrichedList,
		Total: total,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*models.Product, error) {
	newProduct := &models.Product{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newProduct).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return errcode.ErrDuplicateSKU
			}
			return err
		}
		for _, categoryID := range req.CategoryIDs {
			pc := &shared.ProductCategory{
				ProductID:  newProduct.ID,
				CategoryID: categoryID,
			}
			if err := tx.Create(pc).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.bloomFilter.add(newProduct.ID)

	// 添加到每个分类的 ZSET
	for _, catID := range req.CategoryIDs {
		s.rdb.ZAdd(context.Background(), categoryZSetKey(catID), redis.Z{Score: float64(newProduct.ID), Member: newProduct.ID})
	}

	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0]
	}
	s.rabbit.Publish(ctx, events.ProductCreatedEvent{
		ProductID:  newProduct.ID,
		Name:       newProduct.Name,
		CategoryID: categoryIDValue,
	})

	return newProduct, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.FindByID(ctx, id)
}

// GetSkusWithInventory 获取产品的 SKU 列表并批量注入库存信息（单次 LEFT JOIN）
func (s *ProductService) GetSkusWithInventory(ctx context.Context, productID int64) ([]dto.SkuDetailResponse, error) {
	rows, err := s.skuRepo.FindByProductIDWithInventory(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.SkuDetailResponse, len(rows))
	for i, row := range rows {
		availableQuantity := 0
		inventoryStatus := string(models.InventoryStatusOutOfStock)
		if row.Quantity != nil {
			availableQuantity = *row.Quantity - *row.Reserved
			if row.Status != nil {
				inventoryStatus = *row.Status
			}
		}
		sku := &models.Sku{
			ID:        row.ID,
			ProductID: row.ProductID,
			Name:      row.Name,
			Price:     row.Price,
			SKUCode:   row.SKUCode,
			Image:     row.Image,
		}
		if row.Spec != "" && row.Spec != "null" {
			_ = sonic.Unmarshal([]byte(row.Spec), &sku.Spec)
		}
		result[i] = dto.SkuDetailToResponse(sku, availableQuantity, inventoryStatus)
	}
	return result, nil
}

func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	return s.repo.FindBySKU(ctx, sku)
}

func (s *ProductService) GetProductWithCategories(ctx context.Context, productID int64) (*models.Product, []models.Category, error) {
	prod, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, nil, err
	}

	var categories []models.Category
	if err := s.db.WithContext(ctx).Table("categories").
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	return prod, categories, nil
}

func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]models.Product, int64, error) {
	return s.repo.ListByCategory(ctx, categoryID, page, pageSize)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductDTO) (*models.Product, error) {
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}

	// 记录旧的分类 ID，用于更新后清理旧的分类 ZSET
	var oldCategoryIDs []int64
	if req.CategoryIDs != nil {
		s.db.WithContext(ctx).Table("product_categories").Where("product_id = ?", id).Pluck("category_id", &oldCategoryIDs)
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(existingProduct).Error; err != nil {
			return err
		}
		if req.CategoryIDs != nil {
			if err := tx.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
				return err
			}
			for _, categoryID := range req.CategoryIDs {
				pc := &shared.ProductCategory{
					ProductID:  existingProduct.ID,
					CategoryID: categoryID,
				}
				if err := tx.Create(pc).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 更新分类 ZSET：从旧的移除，添加到新的
	if req.CategoryIDs != nil {
		for _, catID := range oldCategoryIDs {
			s.rdb.ZRem(context.Background(), categoryZSetKey(catID), id)
		}
		for _, catID := range req.CategoryIDs {
			s.rdb.ZAdd(context.Background(), categoryZSetKey(catID), redis.Z{Score: float64(id), Member: id})
		}
	}

	s.delayedDoubleDelete(id)

	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0]
	}
	s.rabbit.Publish(ctx, events.ProductUpdatedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		CategoryID: categoryIDValue,
	})

	return existingProduct, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 查询产品所属分类，用于后续清理分类 ZSET
	var catIDs []int64
	s.db.WithContext(ctx).Table("product_categories").Where("product_id = ?", id).Pluck("category_id", &catIDs)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 从每个分类 ZSET 中移除
	for _, catID := range catIDs {
		s.rdb.ZRem(context.Background(), categoryZSetKey(catID), id)
	}

	s.delayedDoubleDelete(id)

	s.rabbit.Publish(ctx, events.ProductDeletedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		CategoryID: 0,
	})

	return nil
}

func (s *ProductService) ListAllProducts(ctx context.Context, page, pageSize int) ([]models.Product, int64, error) {
	return s.repo.ListAll(ctx, page, pageSize)
}

func (s *ProductService) ListProducts(ctx context.Context, q dto.ProductListQuery) (*dto.ProductListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, total, err := s.repo.ListProductsWithTotal(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResult{
		List:  list,
		Total: total,
	}, nil
}

// ListProductsByCursor 基于游标的产品列表查询（深分页优化）
func (s *ProductService) ListProductsByCursor(ctx context.Context, q dto.ProductCursorQuery) (*dto.ProductCursorResult, error) {
	// 多查一条判断是否还有更多数据
	limit := q.Size + 1
	list, err := s.repo.ListProductsByCursor(ctx, q, limit)
	if err != nil {
		return nil, err
	}

	result := &dto.ProductCursorResult{}
	if len(list) > q.Size {
		// 有更多数据，截断多查的那条
		result.List = list[:q.Size]
		result.NextCursor = list[q.Size-1].ID
		result.HasMore = true
	} else {
		result.List = list
		result.NextCursor = 0
		result.HasMore = false
	}
	return result, nil
}
