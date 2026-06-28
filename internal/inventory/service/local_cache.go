package service

import (
	"math/rand"
	"sync"
	"time"

	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/pkg/query"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	localCacheSize     = 16384
	localCacheTTL      = 120 * time.Second
	localCacheTTLJitter = 0.2
	listCacheSize      = 2048
	listCacheTTL       = 60 * time.Second
	listCacheTTLJitter = 0.2
)

func jitteredTTL(base time.Duration, jitter float64) time.Duration {
	delta := time.Duration(float64(base) * jitter)
	return base + time.Duration(rand.Int63n(int64(delta*2+1)))-delta
}

type cacheEntry[T any] struct {
	item      T
	expiresAt time.Time
}

func newCacheEntry[T any](item T, ttl time.Duration) *cacheEntry[T] {
	return &cacheEntry[T]{item: item, expiresAt: time.Now().Add(ttl)}
}

func (e *cacheEntry[T]) isExpired() bool {
	return time.Now().After(e.expiresAt)
}

type productLocalCache struct {
	mu       sync.RWMutex
	single   *lru.Cache[int64, *cacheEntry[*dto.CachedProductItem]]
	list     *lru.Cache[string, *cacheEntry[*query.ListResult[dto.CachedProductItem]]]
	singleTTL time.Duration
	listTTL   time.Duration
}

func newProductLocalCache() *productLocalCache {
	single, _ := lru.New[int64, *cacheEntry[*dto.CachedProductItem]](localCacheSize)
	list, _ := lru.New[string, *cacheEntry[*query.ListResult[dto.CachedProductItem]]](listCacheSize)
	return &productLocalCache{
		single:    single,
		list:      list,
		singleTTL: localCacheTTL,
		listTTL:   listCacheTTL,
	}
}

func (c *productLocalCache) getSingle(id int64) (*dto.CachedProductItem, bool) {
	c.mu.RLock()
	entry, ok := c.single.Get(id)
	c.mu.RUnlock()
	if !ok || entry.isExpired() {
		if ok {
			c.mu.Lock()
			c.single.Remove(id)
			c.mu.Unlock()
		}
		return nil, false
	}
	return entry.item, true
}

func (c *productLocalCache) setSingle(id int64, item *dto.CachedProductItem) {
	c.mu.Lock()
	c.single.Add(id, newCacheEntry(item, jitteredTTL(c.singleTTL, localCacheTTLJitter)))
	c.mu.Unlock()
}

func (c *productLocalCache) getList(key string) (*query.ListResult[dto.CachedProductItem], bool) {
	c.mu.RLock()
	entry, ok := c.list.Get(key)
	c.mu.RUnlock()
	if !ok || entry.isExpired() {
		if ok {
			c.mu.Lock()
			c.list.Remove(key)
			c.mu.Unlock()
		}
		return nil, false
	}
	return entry.item, true
}

func (c *productLocalCache) setList(key string, result *query.ListResult[dto.CachedProductItem]) {
	c.mu.Lock()
	c.list.Add(key, newCacheEntry(result, jitteredTTL(c.listTTL, listCacheTTLJitter)))
	c.mu.Unlock()
}

func (c *productLocalCache) removeSingle(id int64) {
	c.mu.Lock()
	c.single.Remove(id)
	c.mu.Unlock()
}

func (c *productLocalCache) clear() {
	c.mu.Lock()
	c.single.Purge()
	c.list.Purge()
	c.mu.Unlock()
}

func (c *productLocalCache) warmup(items []dto.CachedProductItem) {
	c.mu.Lock()
	for i := range items {
		entry := &cacheEntry[*dto.CachedProductItem]{
			item:      &items[i],
			expiresAt: time.Now().Add(jitteredTTL(c.singleTTL, localCacheTTLJitter)),
		}
		c.single.Add(items[i].ID, entry)
	}
	c.mu.Unlock()
}