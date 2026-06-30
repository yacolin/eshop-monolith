---
name: "perf-opt"
description: "项目性能优化模式指南，涵盖三级缓存架构(L1→BloomFilter→L2→DB)、缓存治理(穿透/击穿/雪崩防护)、Keyset游标分页、降级策略等。适用于需要性能优化、缓存接入、分页改造的场景。"
---

# 性能优化指南

## 适用范围

本指南适用于重构后的新模块（Flat Go Structure）中需要性能优化的场景。常规 CRUD 业务不需要参考本指南。

## 1. 三级缓存架构

SPU 等核心读接口采用 L1(本地LRU) → BloomFilter → L2(Redis) → DB 三级缓存方案：

```
请求 → L1 Local LRU  ← 命中直接返回
         ↓ 未命中
      Bloom Filter  ← 拦截不存在 ID，防止穿透到 Redis
         ↓ 通过
      L2 Redis  ← ZSET 查 ID 列表 / String 查实体
         ↓ 未命中
      DB  ← 回填 L2 + L1
```

### 1.1 L1 本地缓存

每个 Service 内嵌一个本地 LRU 缓存实例，使用 `hashicorp/golang-lru/v2`：

```go
type productLocalCache struct {
    mu        sync.RWMutex
    single    *lru.Cache[int64, *cacheEntry[*CachedProductItem]]
    list      *lru.Cache[string, *cacheEntry[*SPUListResult]]
    singleTTL time.Duration
    listTTL   time.Duration
}
```

- `singleCache`：实体级别的本地缓存（如单个 SPU 详情）
- `listCache`：列表级别的本地缓存（如分页列表结果）
- 线程安全：`sync.RWMutex` 读写锁

写入时 TTL 加入 ±20% 随机偏移，防止批量过期导致缓存雪崩：

```go
func jitteredTTL(base time.Duration, jitter float64) time.Duration {
    delta := time.Duration(float64(base) * jitter)
    return base + time.Duration(rand.Int63n(int64(delta*2+1))) - delta
}
```

统一在 `cache.go` 文件中按实体组织缓存函数：

```
cache.go 结构：
├── Key Builders        — cacheKey{Entity}(params) 构建 Redis key
├── {Entity} Cache      — get/set/del 三件套
├── Warmup              — 全量数据预热到 Redis + L1
├── Batch Operations    — MGet 批量获取
├── Bloom Filter        — 缓存穿透防护
└── Delayed Double Delete — 写操作后的延迟双删
```

### 1.2 Bloom Filter

入口级拦截，阻挡对不存在 ID 的请求穿透到 Redis：

```go
import "github.com/bits-and-blooms/bloom/v3"

type productBloomFilter struct {
    filter *bloom.BloomFilter
}

// 估算容量 10 万，允许 1% 误判率
bloomN = 100000
bloomP = 0.01
```

- 每次回源 DB 查询到数据后在 Bloom Filter 中标记 ID
- `mayExist()` 快速判断 ID 是否存在（可能误判，但不会漏判）
- 不存在 ID 直接返回 404，不穿透 Redis/DB
- 写操作（Create/Delete）后同步更新 Bloom Filter
- 空值标记缓存：对不存在的数据缓存 `__EMPTY__` 占位符 30s，阻止恶意穿透

### 1.3 L2 Redis 缓存

统一 Redis 缓存函数：

- 使用 `redis.UniversalClient` 接口（兼容单机/集群）
- 序列化使用 `github.com/bytedance/sonic`（高性能 JSON）
- TTL 统一在文件顶部常量定义：`spuEntityTTL = 10 * time.Minute`
- 批量获取提供 `batchFetch{Entity}Entities`，使用 `MGet` 一次查询

列表缓存使用 ZSET + 实体分离：

```go
ZSET key:  spu:list:ids:cat={cid}:brand={bid}:status={s}
实体 key:  spu:{id}
```

- ZSET score = SortOrder，member = 补零 ID（`%020d`），用于 ZRevRange 游标翻页
- 实体用独立 key 存储，列表只存 ID 列表 + 排序信息
- Lua 脚本实现一次往返完成 ZRANGE + ZCARD + MGET：

```lua
-- zrangeMGetScript: 一次 Redis 网络往返完成列表页查询
local ids = redis.call("ZREVRANGE", KEYS[1], ARGV[1], ARGV[2])
local total = redis.call("ZCARD", KEYS[1])
local keys = {}
for i, id in ipairs(ids) do
    keys[i] = ARGV[4] .. id
end
local values = redis.call("MGET", unpack(keys))
return {total, values}
```

### 1.4 缓存预热

服务启动时或缓存全量失效后，提供 Warmup 方法：

```go
func (s *SpuService) WarmupCache(ctx context.Context) (int, error) {
    spus, err := s.repo.FindAll(ctx)
    // 1. Pipeline 写 Redis ZSET + 实体
    pipe := s.rdb.Pipeline()
    for _, spu := range spus {
        pipe.Set(ctx, cacheKeySPU(spu.ID), data, spuEntityTTL)
        pipe.ZAdd(ctx, cacheKeySPUListIDs(...), redis.Z{Score: ..., Member: ...})
    }
    pipe.Exec(ctx)
    // 2. 预热 L1
    s.localCache.warmup(items)
    // 3. 重建 Bloom Filter
    s.bloomFilter.addAll(ids)
}
```

## 2. 缓存治理

### 2.1 缓存穿透防护 — Bloom Filter + 空值标记

- Bloom Filter 拦截对不存在 ID 的请求
- 空值 30s 缓存 `__EMPTY__`，应对批量不存在 ID 的恶意穿透
- 二者同时生效：Bloom Filter 快速拦截，空值缓存兜底

### 2.2 缓存击穿防护 — Singleflight

热点 key 过期时，仅单协程回源加载，防止并发请求打穿 DB：

```go
import "golang.org/x/sync/singleflight"

type SpuService struct {
    sf singleflight.Group
}

// ZSET 缓存重建时防并发击穿
v, err, _ := s.sf.Do(cacheKey, func() (interface{}, error) {
    // 查 DB 并回填 ZSET
    return nil, nil
})
```

### 2.3 缓存雪崩防护 — TTL Jitter

- L1 本地缓存写入时 TTL 加入 ±20% 随机偏移
- Redis 缓存 TTL 使用固定值（不影响集群分布）

### 2.4 热点 Key 识别

统计 10s 窗口内各 key 访问频次，阈值 1000 次/s 标记为热点：

```go
type hotKeyCounter struct {
    mu       sync.Mutex
    counters map[int64]*hotKeyEntry  // id → {count, start}
}

const (
    hotKeyThreshold = 1000
    hotKeyWindow    = 10 * time.Second
)
```

## 3. Keyset 游标分页

大表列表（如 SPU）不使用 offset 分页，改用 keyset 游标分页避免深度翻页性能问题。

### DTO 定义

```go
type SPUListReq struct {
    Size   int    `form:"size,default=10" binding:"gte=1,lte=1000"`
    Cursor string `form:"cursor"`
}

type SPUListResult struct {
    List    []*SPU `json:"list"`
    Cursor  string `json:"cursor"`
    HasMore bool   `json:"has_more"`
}
```

- `size` 控制每页条数，不传默认 10，上限 1000
- `cursor` 由上次返回结果中的 `cursor` 字段提供，首次请求不传
- `HasMore` 标记是否还有下一页

### 数据库查询

```go
// 取 size+1 条判断是否有下一页
ids, err := s.repo.ListIDs(ctx, ..., req.Size+1, cursor.SortOrder, cursor.ID)
hasMore := len(ids) > req.Size

// Repo 层使用 (sort_order, id) 复合条件定位游标
db.Where("(sort_order < ? OR (sort_order = ? AND id < ?))", cursorSortOrder, cursorSortOrder, cursorID)
db.Order("sort_order DESC, id DESC").Limit(limit)
```

## 4. Brand 全量缓存 + 内存过滤

数据量小的实体（如 Brand）使用全量缓存 + 内存过滤分页：

```go
// 1. 尝试读缓存
if s.rdb != nil {
    cached, err := getBrandAllCache(ctx, s.rdb)
    if err == nil { allBrands = cached }
}
// 2. 缓存未命中 → 查 DB 并回填
if allBrands == nil {
    allBrands, err = s.repo.FindAll(ctx)
    if s.rdb != nil { setBrandAllCache(ctx, s.rdb, allBrands) }
}
// 3. 内存过滤 + 排序 + 分页
filtered := make([]Brand, 0, len(allBrands))
for _, b := range allBrands {
    if req.Name != "" && !strings.Contains(b.Name, req.Name) { continue }
}
```

## 5. 延迟双删

写操作先删缓存，写 DB，500ms 后再删一次，解决并发读写导致的缓存一致性问题：

```go
// Service 层写操作模式
if s.rdb != nil {
    del{Entity}Cache(context.Background(), s.rdb)       // 先删缓存
}
if err := s.repo.Update(ctx, entity); err != nil {
    return err
}
if s.rdb != nil {
    delayedDelete{Entity}(context.Background(), s.rdb)  // 500ms 后再删
}
```

```go
func delayedDeleteSPU(ctx context.Context, rdb redis.UniversalClient, id int64) {
    go func() {
        time.Sleep(500 * time.Millisecond)
        delSPUEntity(context.Background(), rdb, id)
        delAllSPUListCache(context.Background(), rdb)
    }()
}
```

## 6. 降级策略

Redis 不可用时静默降级查 DB，不阻塞主流程：

- `s.rdb != nil` 作为缓存开关，构造函数传 nil 可禁用缓存
- 缓存读取失败直接走 DB
- 缓存回填失败忽略错误

```go
if s.rdb != nil {
    cached, err := getCache(ctx, s.rdb, id)
    if err == nil { return cached, nil }
}
spu, err := s.repo.FindByID(ctx, id)
if s.rdb != nil {
    _ = setCache(ctx, s.rdb, spu) // 忽略回填错误
}
```

## 检查清单

### 缓存架构
- [ ] 缓存函数是否统一在 `cache.go` 中，按实体组织
- [ ] L1 本地 LRU 缓存是否有合理的容量和 TTL
- [ ] TTL 是否带 ±20% 随机 jitter 防雪崩
- [ ] 是否实现 Bloom Filter 拦截不存在 ID 的穿透
- [ ] ZSET 缓存重建是否使用 `singleflight.Group` 防击穿
- [ ] 是否有热点 key 识别机制
- [ ] 是否有 Warmup 预热方法

### 列表查询
- [ ] 大表列表是否使用 keyset 游标分页（非 offset）
- [ ] 游标分页是取 `size+1` 条判断 `HasMore`

### 一致性
- [ ] Redis 不可用时是否静默降级查 DB
- [ ] 写操作是否有延迟双删（先删缓存→写 DB→500ms 后再删）

### 缓存 Key 规范
- [ ] 缓存 key 是否使用统一命名前缀（如 `spu:`、`brand:`）
- [ ] TTL 是否在文件顶部常量定义
- [ ] 列表缓存是否使用 ZSET 存储 ID 列表 + 独立实体 key
