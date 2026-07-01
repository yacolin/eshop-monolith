package marketing

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"eshop-monolith/pkg/logger"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/bytedance/sonic"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
)

const (
	promotionEntityTTL = 10 * time.Minute
	promotionDetailTTL = 5 * time.Minute
	delayedDeleteDelay = 500 * time.Millisecond

	// L1 local cache
	promotionLocalCacheSize = 1024
	promotionLocalCacheTTL  = 60 * time.Second
	localJitter             = 0.2

	// Bloom Filter
	promotionBloomN = 50000
	promotionBloomP = 0.01

	emptyPlaceholder = "__EMPTY__"
	emptyCacheTTL    = 30 * time.Second
)

// ── Key Builders ──

func cacheKeyPromotion(id int64) string { return fmt.Sprintf("promotion:%d", id) }
func cacheKeyPromotionDetail(id int64) string { return fmt.Sprintf("promotion:detail:%d", id) }

// ── Promotion Entity Cache ──

func getPromotionEntity(ctx context.Context, rdb redis.UniversalClient, id int64) (*Promotion, error) {
	data, err := rdb.Get(ctx, cacheKeyPromotion(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var p Promotion
	if err := sonic.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func setPromotionEntity(ctx context.Context, rdb redis.UniversalClient, p *Promotion) error {
	data, err := sonic.Marshal(p)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyPromotion(p.ID), data, promotionEntityTTL).Err()
}

func delPromotionEntity(ctx context.Context, rdb redis.UniversalClient, id int64) {
	rdb.Del(ctx, cacheKeyPromotion(id))
}

// ── Promotion Detail Cache ──

func getPromotionDetail(ctx context.Context, rdb redis.UniversalClient, id int64) (*PromotionDetailResponse, error) {
	data, err := rdb.Get(ctx, cacheKeyPromotionDetail(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var d PromotionDetailResponse
	if err := sonic.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func setPromotionDetail(ctx context.Context, rdb redis.UniversalClient, d *PromotionDetailResponse) error {
	data, err := sonic.Marshal(d)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyPromotionDetail(d.ID), data, promotionDetailTTL).Err()
}

func delPromotionDetail(ctx context.Context, rdb redis.UniversalClient, id int64) {
	rdb.Del(ctx, cacheKeyPromotionDetail(id))
}

// ── TTL Jitter ──

func jitteredTTL(base time.Duration, jitter float64) time.Duration {
	delta := time.Duration(float64(base) * jitter)
	return base + time.Duration(rand.Int63n(int64(delta*2+1))) - delta
}

// ── L1 Local Cache ──

type cacheEntry[T any] struct {
	item      T
	expiresAt time.Time
}

type promotionLocalCache struct {
	mu        sync.RWMutex
	single    *lru.Cache[int64, *cacheEntry[*Promotion]]
	singleTTL time.Duration
}

func newPromotionLocalCache() *promotionLocalCache {
	single, _ := lru.New[int64, *cacheEntry[*Promotion]](promotionLocalCacheSize)
	return &promotionLocalCache{
		single:    single,
		singleTTL: promotionLocalCacheTTL,
	}
}

func (c *promotionLocalCache) getSingle(id int64) (*Promotion, bool) {
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

func (c *promotionLocalCache) setSingle(id int64, p *Promotion) {
	c.mu.Lock()
	c.single.Add(id, &cacheEntry[*Promotion]{
		item:      p,
		expiresAt: time.Now().Add(jitteredTTL(c.singleTTL, localJitter)),
	})
	c.mu.Unlock()
}

func (c *promotionLocalCache) removeSingle(id int64) {
	c.mu.Lock()
	c.single.Remove(id)
	c.mu.Unlock()
}

func (c *promotionLocalCache) clear() {
	c.mu.Lock()
	c.single.Purge()
	c.mu.Unlock()
}

// ── Bloom Filter ──

type promotionBloomFilter struct {
	mu     sync.RWMutex
	filter *bloom.BloomFilter
	count  int64
}

func newPromotionBloomFilter() *promotionBloomFilter {
	return &promotionBloomFilter{
		filter: bloom.NewWithEstimates(promotionBloomN, promotionBloomP),
	}
}

func (b *promotionBloomFilter) add(id int64) {
	b.mu.Lock()
	b.filter.AddString(strconv.FormatInt(id, 10))
	b.count++
	b.mu.Unlock()
}

func (b *promotionBloomFilter) addAll(ids []int64) {
	b.mu.Lock()
	for _, id := range ids {
		b.filter.AddString(strconv.FormatInt(id, 10))
	}
	b.count += int64(len(ids))
	b.mu.Unlock()
}

func (b *promotionBloomFilter) mayExist(id int64) bool {
	b.mu.RLock()
	c := b.count
	if c == 0 {
		b.mu.RUnlock()
		return true
	}
	ok := b.filter.TestString(strconv.FormatInt(id, 10))
	b.mu.RUnlock()
	return ok
}

func (b *promotionBloomFilter) clear() {
	b.mu.Lock()
	b.filter = bloom.NewWithEstimates(promotionBloomN, promotionBloomP)
	b.count = 0
	b.mu.Unlock()
}

// ── Delayed Double Delete ──

func delayedDeletePromotion(ctx context.Context, rdb redis.UniversalClient, id int64) {
	go func() {
		time.Sleep(delayedDeleteDelay)
		logger.Info("promotion cache invalidated", "id", id)
		bg := context.Background()
		delPromotionEntity(bg, rdb, id)
		delPromotionDetail(bg, rdb, id)
	}()
}
