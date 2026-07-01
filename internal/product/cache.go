package product

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"eshop-monolith/pkg/logger"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/bytedance/sonic"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
)

const (
	spuEntityTTL       = 10 * time.Minute
	spuListTTL         = 10 * time.Minute
	brandAllTTL        = 10 * time.Minute
	brandEntityTTL     = 10 * time.Minute
	categoryAllTTL     = 10 * time.Minute
	categoryEntityTTL  = 10 * time.Minute
	delayedDeleteDelay = 500 * time.Millisecond

	// L1 local cache
	spuLocalCacheSize      = 8192
	spuLocalCacheTTL       = 60 * time.Second
	spuLocalJitter         = 0.2
	brandLocalCacheSize    = 512
	categoryLocalCacheSize = 1024
	localCacheTTL          = 60 * time.Second
	localJitter            = 0.2

	// Bloom Filter
	bloomN = 100000
	bloomP = 0.01

	// Hot key
	hotKeyThreshold = 1000
	hotKeyWindow    = 10 * time.Second

	emptyPlaceholder = "__EMPTY__"
	emptyCacheTTL    = 30 * time.Second
)

// ── Key Builders ──

func cacheKeySPU(id int64) string { return fmt.Sprintf("spu:%d", id) }
func cacheKeyBrandAll() string    { return "brand:all" }
func cacheKeySPUListIDs(categoryID, brandID *int64, status *int8) string {
	cid := "0"
	if categoryID != nil {
		cid = strconv.FormatInt(*categoryID, 10)
	}
	bid := "0"
	if brandID != nil {
		bid = strconv.FormatInt(*brandID, 10)
	}
	st := ""
	if status != nil {
		st = strconv.Itoa(int(*status))
	}
	return fmt.Sprintf("spu:list:ids:cat=%s:brand=%s:status=%s", cid, bid, st)
}
func spuZSETMember(id int64) string { return fmt.Sprintf("%020d", id) }

// ── SPU Entity Cache ──

func getSPUEntity(ctx context.Context, rdb redis.UniversalClient, id int64) (*SPU, error) {
	data, err := rdb.Get(ctx, cacheKeySPU(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var spu SPU
	if err := sonic.Unmarshal(data, &spu); err != nil {
		return nil, err
	}
	return &spu, nil
}

func setSPUEntity(ctx context.Context, rdb redis.UniversalClient, spu *SPU) error {
	data, err := sonic.Marshal(spu)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeySPU(spu.ID), data, spuEntityTTL).Err()
}

func delSPUEntity(ctx context.Context, rdb redis.UniversalClient, id int64) {
	rdb.Del(ctx, cacheKeySPU(id))
}

// batchFetchSPUEntities 批量获取 SPU 实体缓存，返回 (命中map, 未命中id列表)
func batchFetchSPUEntities(ctx context.Context, rdb redis.UniversalClient, ids []int64) (map[int64]*SPU, []int64, error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = cacheKeySPU(id)
	}
	results, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, ids, err
	}
	hit := make(map[int64]*SPU, len(ids))
	var miss []int64
	for i, val := range results {
		if val == nil {
			miss = append(miss, ids[i])
			continue
		}
		s, ok := val.(string)
		if !ok {
			miss = append(miss, ids[i])
			continue
		}
		var spu SPU
		if err := sonic.Unmarshal([]byte(s), &spu); err != nil {
			miss = append(miss, ids[i])
			continue
		}
		hit[ids[i]] = &spu
	}
	return hit, miss, nil
}

// ── SPU List ZSET Cache ──

// getSPUListIDs 从 ZSET 获取指定游标后的 ID 列表（id ASC）。cursorRank=-1 表示第一页。
func getSPUListIDs(ctx context.Context, rdb redis.UniversalClient, categoryID, brandID *int64, status *int8, cursorRank int64, count int) ([]int64, error) {
	key := cacheKeySPUListIDs(categoryID, brandID, status)
	start := int64(0)
	if cursorRank >= 0 {
		start = cursorRank
	}
	members, err := rdb.ZRange(ctx, key, start, start+int64(count)-1).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(members))
	for _, m := range members {
		id, _ := strconv.ParseInt(strings.TrimSpace(m), 10, 64)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// setSPUListIDs 将 SPU 列表 ID 写入 ZSET，按 id ASC 排序
func setSPUListIDs(ctx context.Context, rdb redis.UniversalClient, categoryID, brandID *int64, status *int8, spus []SPU) error {
	key := cacheKeySPUListIDs(categoryID, brandID, status)
	pairs := make([]redis.Z, 0, len(spus))
	for i := range spus {
		pairs = append(pairs, redis.Z{
			Score:  float64(spus[i].ID),
			Member: spuZSETMember(spus[i].ID),
		})
	}
	if err := rdb.ZAdd(ctx, key, pairs...).Err(); err != nil {
		return err
	}
	return rdb.Expire(ctx, key, spuListTTL).Err()
}

// getSPUListRank 获取 SPU ID 在 ZSET 中的排名（ZRANK），-1 表示不存在
func getSPUListRank(ctx context.Context, rdb redis.UniversalClient, categoryID, brandID *int64, status *int8, id int64) (int64, error) {
	key := cacheKeySPUListIDs(categoryID, brandID, status)
	return rdb.ZRank(ctx, key, spuZSETMember(id)).Result()
}

// delAllSPUListCache 删除所有 SPU 列表缓存 ZSET
func delAllSPUListCache(ctx context.Context, rdb redis.UniversalClient) {
	iter := rdb.Scan(ctx, 0, "spu:list:ids:*", 0).Iterator()
	for iter.Next(ctx) {
		rdb.Del(ctx, iter.Val())
	}
	// SCAN 可能因连接断开或超时失败，忽略错误（延迟双删的二次删除会兜底）
	_ = iter.Err()
}

// ── Brand Cache ──

func getBrandAllCache(ctx context.Context, rdb redis.UniversalClient) ([]Brand, error) {
	data, err := rdb.Get(ctx, cacheKeyBrandAll()).Bytes()
	if err != nil {
		return nil, err
	}
	var brands []Brand
	if err := sonic.Unmarshal(data, &brands); err != nil {
		return nil, err
	}
	return brands, nil
}

func setBrandAllCache(ctx context.Context, rdb redis.UniversalClient, brands []Brand) error {
	data, err := sonic.Marshal(brands)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyBrandAll(), data, brandAllTTL).Err()
}

func delBrandAllCache(ctx context.Context, rdb redis.UniversalClient) {
	rdb.Del(ctx, cacheKeyBrandAll())
}

// ── Brand Entity Cache ──

func cacheKeyBrand(id int64) string { return fmt.Sprintf("brand:%d", id) }

func getBrandEntity(ctx context.Context, rdb redis.UniversalClient, id int64) (*Brand, error) {
	data, err := rdb.Get(ctx, cacheKeyBrand(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var b Brand
	if err := sonic.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func setBrandEntity(ctx context.Context, rdb redis.UniversalClient, b *Brand) error {
	data, err := sonic.Marshal(b)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyBrand(b.ID), data, brandEntityTTL).Err()
}

func delBrandEntity(ctx context.Context, rdb redis.UniversalClient, id int64) {
	rdb.Del(ctx, cacheKeyBrand(id))
}

func delayedDeleteBrandEntity(ctx context.Context, rdb redis.UniversalClient, id int64) {
	go func() {
		time.Sleep(delayedDeleteDelay)
		delBrandEntity(context.Background(), rdb, id)
	}()
}

// ── Category Cache ──

func cacheKeyCategory(id int64) string { return fmt.Sprintf("category:%d", id) }
func cacheKeyCategoryAll() string      { return "category:all" }

func getCategoryAllCache(ctx context.Context, rdb redis.UniversalClient) ([]Category, error) {
	data, err := rdb.Get(ctx, cacheKeyCategoryAll()).Bytes()
	if err != nil {
		return nil, err
	}
	var cats []Category
	if err := sonic.Unmarshal(data, &cats); err != nil {
		return nil, err
	}
	return cats, nil
}

func setCategoryAllCache(ctx context.Context, rdb redis.UniversalClient, cats []Category) error {
	data, err := sonic.Marshal(cats)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyCategoryAll(), data, categoryAllTTL).Err()
}

func delCategoryAllCache(ctx context.Context, rdb redis.UniversalClient) {
	rdb.Del(ctx, cacheKeyCategoryAll())
}

func getCategoryEntity(ctx context.Context, rdb redis.UniversalClient, id int64) (*Category, error) {
	data, err := rdb.Get(ctx, cacheKeyCategory(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var c Category
	if err := sonic.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func setCategoryEntity(ctx context.Context, rdb redis.UniversalClient, c *Category) error {
	data, err := sonic.Marshal(c)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyCategory(c.ID), data, categoryEntityTTL).Err()
}

func delCategoryEntity(ctx context.Context, rdb redis.UniversalClient, id int64) {
	rdb.Del(ctx, cacheKeyCategory(id))
}

func delayedDeleteCategoryAll(ctx context.Context, rdb redis.UniversalClient) {
	go func() {
		time.Sleep(delayedDeleteDelay)
		delCategoryAllCache(context.Background(), rdb)
	}()
}

func delayedDeleteCategoryEntity(ctx context.Context, rdb redis.UniversalClient, id int64) {
	go func() {
		time.Sleep(delayedDeleteDelay)
		delCategoryEntity(context.Background(), rdb, id)
	}()
}

// ── TTL Jitter ──

func jitteredTTL(base time.Duration, jitter float64) time.Duration {
	delta := time.Duration(float64(base) * jitter)
	return base + time.Duration(rand.Int63n(int64(delta*2+1))) - delta
}

// ── L1 SPU Local Cache ──

type cacheEntry[T any] struct {
	item      T
	expiresAt time.Time
}

type spuLocalCache struct {
	mu        sync.RWMutex
	single    *lru.Cache[int64, *cacheEntry[*SPU]]
	singleTTL time.Duration
}

func newSPULocalCache() *spuLocalCache {
	single, _ := lru.New[int64, *cacheEntry[*SPU]](spuLocalCacheSize)
	return &spuLocalCache{
		single:    single,
		singleTTL: spuLocalCacheTTL,
	}
}

func (c *spuLocalCache) getSingle(id int64) (*SPU, bool) {
	c.mu.RLock()
	entry, ok := c.single.Get(id)
	c.mu.RUnlock()
	if !ok || entry.expiresAt.Before(time.Now()) {
		if ok {
			c.mu.Lock()
			c.single.Remove(id)
			c.mu.Unlock()
		}
		return nil, false
	}
	return entry.item, true
}

func (c *spuLocalCache) setSingle(id int64, spu *SPU) {
	c.mu.Lock()
	c.single.Add(id, &cacheEntry[*SPU]{
		item:      spu,
		expiresAt: time.Now().Add(jitteredTTL(c.singleTTL, spuLocalJitter)),
	})
	c.mu.Unlock()
}

func (c *spuLocalCache) removeSingle(id int64) {
	c.mu.Lock()
	c.single.Remove(id)
	c.mu.Unlock()
}

// warmupSingle 预热单条 SPU 到 L1（不设过期，由 TTL jitter 统一控制）
func (c *spuLocalCache) warmupSingle(id int64, spu *SPU) {
	c.mu.Lock()
	c.single.Add(id, &cacheEntry[*SPU]{
		item:      spu,
		expiresAt: time.Now().Add(jitteredTTL(c.singleTTL, spuLocalJitter)),
	})
	c.mu.Unlock()
}

func (c *spuLocalCache) clear() {
	c.mu.Lock()
	c.single.Purge()
	c.mu.Unlock()
}

// ── Generic Local Cache (for small datasets) ──

type simpleLocalCache[T any] struct {
	mu    sync.RWMutex
	cache *lru.Cache[int64, *cacheEntry[T]]
	ttl   time.Duration
	size  int
}

func newSimpleLocalCache[T any](size int, ttl time.Duration) *simpleLocalCache[T] {
	c, _ := lru.New[int64, *cacheEntry[T]](size)
	return &simpleLocalCache[T]{cache: c, ttl: ttl, size: size}
}

func (c *simpleLocalCache[T]) get(id int64) (T, bool) {
	c.mu.RLock()
	entry, ok := c.cache.Get(id)
	c.mu.RUnlock()
	if !ok || entry.expiresAt.Before(time.Now()) {
		if ok {
			c.mu.Lock()
			c.cache.Remove(id)
			c.mu.Unlock()
		}
		var zero T
		return zero, false
	}
	return entry.item, true
}

func (c *simpleLocalCache[T]) set(id int64, item T) {
	c.mu.Lock()
	c.cache.Add(id, &cacheEntry[T]{
		item:      item,
		expiresAt: time.Now().Add(jitteredTTL(c.ttl, localJitter)),
	})
	c.mu.Unlock()
}

func (c *simpleLocalCache[T]) remove(id int64) {
	c.mu.Lock()
	c.cache.Remove(id)
	c.mu.Unlock()
}

func (c *simpleLocalCache[T]) clear() {
	c.mu.Lock()
	c.cache.Purge()
	c.mu.Unlock()
}

// ── SPU Bloom Filter ──

type spuBloomFilter struct {
	mu     sync.RWMutex
	filter *bloom.BloomFilter
	count  int64 // 已添加的 ID 数量，为 0 时跳过检查
}

func newSPUBloomFilter() *spuBloomFilter {
	return &spuBloomFilter{
		filter: bloom.NewWithEstimates(bloomN, bloomP),
	}
}

func (b *spuBloomFilter) add(id int64) {
	b.mu.Lock()
	b.filter.AddString(strconv.FormatInt(id, 10))
	b.count++
	b.mu.Unlock()
}

func (b *spuBloomFilter) addAll(ids []int64) {
	b.mu.Lock()
	for _, id := range ids {
		b.filter.AddString(strconv.FormatInt(id, 10))
	}
	b.count += int64(len(ids))
	b.mu.Unlock()
}

// mayExist 返回 id 是否可能存在。count==0（未预热）时始终返回 true，不拦截。
func (b *spuBloomFilter) mayExist(id int64) bool {
	b.mu.RLock()
	c := b.count
	if c == 0 {
		b.mu.RUnlock()
		return true // 未预热，放行
	}
	ok := b.filter.TestString(strconv.FormatInt(id, 10))
	b.mu.RUnlock()
	return ok
}

func (b *spuBloomFilter) clear() {
	b.mu.Lock()
	b.filter = bloom.NewWithEstimates(bloomN, bloomP)
	b.count = 0
	b.mu.Unlock()
}

// ── Hot Key Counter ──

type hotKeyEntry struct {
	count int64
	start time.Time
}

type spuHotKeyCounter struct {
	mu       sync.Mutex
	counters map[int64]*hotKeyEntry
}

func newSPUHotKeyCounter() *spuHotKeyCounter {
	return &spuHotKeyCounter{counters: make(map[int64]*hotKeyEntry)}
}

func (h *spuHotKeyCounter) increment(id int64) bool {
	h.mu.Lock()
	entry, ok := h.counters[id]
	now := time.Now()
	if !ok || now.Sub(entry.start) > hotKeyWindow {
		entry = &hotKeyEntry{start: now, count: 0}
		h.counters[id] = entry
	}
	entry.count++
	hot := entry.count >= hotKeyThreshold
	h.mu.Unlock()
	return hot
}

func (h *spuHotKeyCounter) reset(id int64) {
	h.mu.Lock()
	delete(h.counters, id)
	h.mu.Unlock()
}

// ── Delayed Double Delete ──

func delayedDeleteSPU(ctx context.Context, rdb redis.UniversalClient, id int64) {
	go func() {
		time.Sleep(delayedDeleteDelay)
		logger.Info("spu cache invalidated", "id", id)
		bg := context.Background()
		delSPUEntity(bg, rdb, id)
		delAllSPUListCache(bg, rdb)
	}()
}

func delayedDeleteBrand(ctx context.Context, rdb redis.UniversalClient) {
	go func() {
		time.Sleep(delayedDeleteDelay)
		delBrandAllCache(context.Background(), rdb)
	}()
}
