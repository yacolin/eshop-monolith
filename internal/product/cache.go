package product

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

const (
	spuEntityTTL       = 10 * time.Minute
	spuListTTL         = 10 * time.Minute
	brandAllTTL        = 10 * time.Minute
	delayedDeleteDelay = 500 * time.Millisecond
)

// ── Key Builders ──

func cacheKeySPU(id int64) string                       { return fmt.Sprintf("spu:%d", id) }
func cacheKeyBrandAll() string                          { return "brand:all" }
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

// getSPUListIDs 从 ZSET 获取指定游标后的 ID 列表。cursorRank=-1 表示第一页。
func getSPUListIDs(ctx context.Context, rdb redis.UniversalClient, categoryID, brandID *int64, status *int8, cursorRank int64, count int) ([]int64, error) {
	key := cacheKeySPUListIDs(categoryID, brandID, status)
	start := int64(0)
	if cursorRank >= 0 {
		start = cursorRank
	}
	members, err := rdb.ZRevRange(ctx, key, start, start+int64(count)-1).Result()
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

// setSPUListIDs 将 SPU 列表 ID 写入 ZSET
func setSPUListIDs(ctx context.Context, rdb redis.UniversalClient, categoryID, brandID *int64, status *int8, spus []SPU) error {
	key := cacheKeySPUListIDs(categoryID, brandID, status)
	pairs := make([]redis.Z, 0, len(spus))
	for i := range spus {
		pairs = append(pairs, redis.Z{
			Score:  float64(spus[i].SortOrder),
			Member: spuZSETMember(spus[i].ID),
		})
	}
	if err := rdb.ZAdd(ctx, key, pairs...).Err(); err != nil {
		return err
	}
	return rdb.Expire(ctx, key, spuListTTL).Err()
}

// getSPUListRank 获取 SPU ID 在 ZSET 中的排名（ZREVRANK），-1 表示不存在
func getSPUListRank(ctx context.Context, rdb redis.UniversalClient, categoryID, brandID *int64, status *int8, id int64) (int64, error) {
	key := cacheKeySPUListIDs(categoryID, brandID, status)
	return rdb.ZRevRank(ctx, key, spuZSETMember(id)).Result()
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

// ── Delayed Double Delete ──

func delayedDeleteSPU(ctx context.Context, rdb redis.UniversalClient, id int64) {
	go func() {
		time.Sleep(delayedDeleteDelay)
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
