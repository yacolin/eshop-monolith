# eshop-monolith

一个完整的电商单体应用系统 Demo，展示了 Go + Gin + GORM 的分层架构和工程最佳实践。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7 + Bloom Filter + 本地 LRU Cache
- **消息队列**: RabbitMQ
- **配置管理**: Viper
- **日志**: Zap + Lumberjack
- **序列化**: Sonic (bytedance)
- **WebSocket**: Gorilla WebSocket
- **监控**: Prometheus
- **文档**: Swagger
- **限流**: Rate limiting middleware

## 核心特性

- ✅ **8 大业务模块**: product / trade / user / inventory / marketing / review / base / dashboard
- ✅ **扁平化模块结构**: 每个模块 handler/model/repo/service/dto 同包，无冗余嵌套
- ✅ **多级缓存体系**: L1 本地 LRU → Bloom Filter 拦截 → L2 Redis → DB 兜底
- ✅ **Keyset 游标分页**: 基于 base64 编码游标的高性能分页，避免 offset 性能衰减
- ✅ **统一认证授权**: JWT 认证 + RBAC 权限控制
- ✅ **消息队列集成**: RabbitMQ 实现模块间异步通信 + WebSocket 实时推送
- ✅ **WebSocket 实时推送**: 在线通知、断线重连、增量同步
- ✅ **Prometheus 监控**: 业务指标暴露
- ✅ **Swagger 文档**: 自动生成 API 文档
- ✅ **统一错误处理**: 错误分类中间件，业务错误码 + 参数校验 + 404 自动处理

## 项目结构

```
eshop-monolith/
├── cmd/
│   └── server/
│       └── main.go                      # 应用入口
├── internal/
│   ├── product/                         # 商品模块（品牌、类目、属性、SPU/SKU）
│   ├── trade/                           # 交易模块（购物车、订单、支付/退款）
│   ├── user/                            # 用户模块（用户、认证、RBAC、地址）
│   ├── inventory/                       # 库存模块（预占、释放、扣减、补货）
│   ├── marketing/                       # 营销模块（促销、优惠券、秒杀）
│   ├── review/                          # 评价模块（商品评价、审核、评分汇总）
│   ├── base/                            # 基础模块（通知推送）
│   ├── dashboard/                       # 仪表盘模块（运营数据汇总）
│   └── infra/                           # 基础设施层
│       ├── repository/                  # 数据库/Redis 初始化 + 仓储聚合
│       ├── router/                      # 统一路由注册
│       ├── rabbitmq/                    # RabbitMQ 客户端 + 消费者
│       └── ws/                          # WebSocket Hub + 会话管理
├── pkg/                                 # 通用工具包
│   ├── config/                          # 配置管理（Viper）
│   ├── errcode/                         # 业务错误码
│   ├── logger/                          # 结构化日志（Zap）
│   ├── middleware/                      # HTTP 中间件
│   │   ├── errorhandler.go              # 全局错误处理
│   │   └── jwtauth.go                   # JWT 认证
│   ├── query/                           # 分页/排序查询
│   ├── response/                        # 统一 JSON 响应
│   └── utils/                           # 工具函数
├── configs/
│   └── config.yaml                      # 配置文件
├── docs/                                # Swagger + 项目文档
├── scripts/                             # SQL 脚本 + 工具脚本
├── sample_data/                         # 示例数据集
├── tests/                               # 集成测试
├── docker-compose.yml
├── Makefile
└── go.mod
```

## 模块一览

| 模块 | 说明 | API 标签 |
|------|------|----------|
| **product** | 商品 SPU/SKU 管理、品牌、类目、属性规格 | `products`, `brands`, `categories`, `attributes` |
| **trade** | 购物车、订单、支付/退款 | `carts`, `orders`, `payments`, `refunds` |
| **user** | 用户注册登录、资料管理、地址管理、RBAC 权限 | `users`, `auth`, `roles`, `permissions`, `addresses` |
| **inventory** | 库存管理、预占、扣减、补货、流水 | `inventories` |
| **marketing** | 促销活动、优惠券领取核销、秒杀抢购 | `promotions`, `coupons`, `flash` |
| **review** | 商品评价创建、审核、评分汇总 | `reviews` |
| **base** | 站内通知推送、系统通知 | `notifications` |
| **dashboard** | 运营仪表盘数据汇总 | `dashboard` |

## 架构设计

### 扁平化模块架构

```
Handler → Service → Repository → DB
              ↑
         External Dependency（跨模块接口）
```

每个业务模块采用**扁平包结构**，文件按职责命名，不嵌套子目录：

```
internal/product/
├── model_spu.go            # 领域模型 + TableName
├── model_sku.go
├── model_category.go
├── repo_spu.go             # Repository 接口 + GORM 实现（同一文件）
├── repo_category.go
├── service_spu.go          # 业务逻辑
├── service_brand.go
├── handler_spu.go          # Gin Handler + 路由注册
├── handler_brand.go
├── dto_spu.go              # 请求/响应 DTO
└── cache.go                # 缓存层（L1/L2/BloomFilter/热 Key 计数）
```

### 跨模块依赖

跨模块调用通过**接口依赖反转**实现，而非直接导入对方包：

```go
// trade/interfaces.go — 定义 trade 需要的外部依赖
type SkuProvider interface {
    FindByID(ctx context.Context, skuID int64) (SkuInfo, error)
}
type InventoryService interface {
    Lock(ctx context.Context, skuID int64, quantity int) error
    Unlock(ctx context.Context, skuID int64, quantity int) error
}
```

### 依赖方向

```
product → DB
   ↑
trade  → product (via interface) + inventory (via interface)
   ↑
review → trade (via interface) + user (via interface)
```

### 多级缓存架构（Product SPU）

```
请求到达
  │
  ├── L1 Local LRU Cache (8K 条目, 60s TTL + jitter)
  │     └── 命中 → 返回
  │
  ├── Bloom Filter（快速拦截不存在的 ID）
  │     └── 不存在 → 直接返回 404（无 DB 查询）
  │
  ├── L2 Redis (10min TTL)
  │     └── 命中 → 回填 L1 → 返回
  │
  └── DB 兜底
        └── 回填 L2 + L1 + Bloom Filter → 返回
```

- 列表查询使用 **Redis ZSET** 缓存 ID 列表，支持游标分页
- 热 Key 计数窗口 10s，超阈值触发本地缓存强化
- 写操作使用**延迟双删**策略保证缓存一致性
- 启动时异步预热全量数据到 Bloom Filter + L2 + L1

### 错误处理

```go
// handler 中通过 c.Error(err) 传递错误
func (h *SpuHandler) GetByID(c *gin.Context) {
    result, err := h.svc.GetByID(c, id)
    if err != nil {
        c.Error(err)  // → ErrorHandler 中间件自动分类处理
        return
    }
    response.Success(c, result)
}
```

ErrorHandler 自动分类：
- `*errcode.BizError` → 业务错误，返回 code + message
- `validator.ValidationErrors` → 422 + 字段级错误详情
- `gorm.ErrRecordNotFound` → 自动转为 404
- 其他 → 500 系统错误

## 本地运行

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 7+
- RabbitMQ（可选，用于消息推送）
- Make (可选)

### 运行步骤

```bash
# 1. 启动依赖服务（MySQL + Redis + RabbitMQ）
#    或使用 Docker Compose：
docker-compose up -d

# 2. 初始化数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS eshop_db;"

# 3. 下载依赖
go mod download

# 4. 运行应用（表结构由 AutoMigrate 自动创建）
go run ./cmd/server

# 或使用 Make
make start
```

### 生成 Swagger 文档

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go --output docs
```

启动后访问: **http://localhost:8080/swagger/index.html**

## API 端点

### 商品 Product

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/products` | 商品列表（Keyset 游标分页） |
| GET | `/api/v1/products/:id` | 商品详情（含 SKU + 属性 + 图文描述） |
| POST | `/api/v1/products` | 创建商品（事务内写 SPU+SKU+Description+Attribute） |
| PUT | `/api/v1/products/:id` | 更新商品 |
| DELETE | `/api/v1/products/:id` | 删除商品 |
| GET | `/api/v1/categories` | 类目树（层级聚合子类目品牌） |
| GET | `/api/v1/brands` | 品牌列表 |
| GET | `/api/v1/attributes` | 属性列表 |
| GET | `/api/v1/category-brands` | 类目品牌关联 |
| GET | `/api/v1/skus` | SKU 列表（含库存信息） |

### 交易 Trade

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/carts` | 获取购物车 |
| POST | `/api/v1/orders` | 创建订单 |
| GET | `/api/v1/orders` | 订单列表 |
| GET | `/api/v1/orders/:id` | 订单详情（含订单项） |
| POST | `/api/v1/payments` | 创建支付 |
| POST | `/api/v1/refunds` | 创建退款 |

### 用户 User

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/users/register` | 用户注册 |
| POST | `/api/v1/users/login` | 用户登录 |
| GET | `/api/v1/users/profile` | 获取用户资料 |
| PUT | `/api/v1/users/profile` | 更新用户资料 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |
| GET/POST | `/api/v1/roles` | 角色管理 |
| GET/POST | `/api/v1/permissions` | 权限管理 |
| GET/POST/PUT/DELETE | `/api/v1/addresses` | 地址管理 |

### 库存 Inventory

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/inventories/stock` | 查询库存 |
| POST | `/api/v1/inventories/lock` | 预占库存 |
| POST | `/api/v1/inventories/unlock` | 释放库存 |
| POST | `/api/v1/inventories/deduct` | 扣减库存 |
| POST | `/api/v1/inventories/restock` | 入库/补货 |
| GET | `/api/v1/inventories/logs` | 库存变更流水 |

### 营销 Marketing

| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST | `/api/v1/promotions` | 促销管理 |
| POST | `/api/v1/coupons/claim` | 领取优惠券 |
| POST | `/api/v1/coupons/use` | 使用优惠券 |
| POST | `/api/v1/flash/buy` | 秒杀抢购 |
| POST | `/api/v1/flash/confirm` | 确认秒杀订单 |

### 评价 Review

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/reviews` | 创建评价 |
| GET | `/api/v1/reviews/me` | 我的评价 |
| GET | `/api/v1/products/:id/reviews` | 商品评价列表 |
| GET | `/api/v1/products/:id/rating` | 商品评分汇总 |
| PATCH | `/api/v1/admin/reviews/:id/moderate` | 审核评价 |

### 通知 Notification

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/notifications` | 通知列表 |
| GET | `/api/v1/notifications/unread` | 未读通知数 |
| PUT | `/api/v1/notifications/:id/read` | 标记已读 |

### 仪表盘 Dashboard

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/dashboard/stats` | 运营数据汇总 |

### WebSocket

| 路径 | 说明 |
|------|------|
| GET | `/api/v1/ws?token=xxx&last_seq=0` | WebSocket 连接 |
| GET | `/api/v1/ws/stats` | 在线统计 |
| POST | `/api/v1/ws/reconnect` | 断线重连 |
| GET | `/api/v1/ws/session` | 会话信息 |

### 系统

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查（含预热状态） |
| GET | `/metrics` | Prometheus 监控指标 |
| GET | `/swagger/*any` | Swagger UI |

## 关键设计

### Keyset 游标分页

使用 base64 编码的游标替代传统 offset/limit：

```
请求: GET /api/v1/products?cursor=MTIz&size=10
响应: { "list": [...], "cursor": "MTM0", "has_more": true }
```

- 避免大 offset 时的性能衰减
- 列表 ID 缓存于 Redis ZSET，按 score 排序支持游标定位

### 缓存一致性（延迟双删）

写操作时：
1. 删除 L1 缓存
2. 删除 L2 Redis 缓存
3. 执行数据库写操作
4. 500ms 后再次异步删除 L2（解决主从延迟读脏）

### 秒杀

秒杀流程按步骤拆分，防止库存超卖：
1. `POST /flash/buy` — 预占资格
2. `POST /flash/confirm` — 确认订单（正式扣减）

## 并发优化（已实施）

以下 7 个方向的并发优化已全部落地，覆盖 Dashboard 查询、列表 COUNT+LIMIT、商品详情、通知落库、WebSocket 广播、缓存预热、DTO 复用：

| # | 方向 | 关键变更 | Commit |
|---|------|---------|--------|
| 1 | Dashboard 并行聚合 | errgroup 并行执行 7 个统计查询；InvalidateCache 接入订单创建 | d64c273 |
| 2 | 列表 COUNT+LIMIT 并发 | 泛型 ConcurrentCountList[T]，Order/Inventory 列表 COUNT+LIMIT 并发 | 28e4b34 |
| 3 | 商品详情 Fan-Out | SPU/SKU/属性/描述四路 errgroup 并发获取后合并 | 773489e |
| 4 | 通知 Worker Pool | 4 worker goroutine + buffered channel 异步落库，降级同步写 | 5d4133a |
| 5 | WS 广播 Fan-Out | snapshot 快照 + errgroup 并发发送，读锁仅用于复制不阻塞 | c286189 |
| 6 | 批量预热 Pipeline | DB Bloom Filter/Redis Pipeline/L1 三阶段 errgroup 并行 | f3bd696 |
| 7 | sync.Pool 复用 DTO | SPUDetailResponse 对象池化复用，减少高频路径 GC | be07444 |

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 执行集成测试
go test ./tests/...
```

## 配置说明

配置通过 `configs/config.yaml` 加载，支持环境变量覆盖：

```bash
MYSQL_PASSWORD=secret \
JWT_SECRET=production-secret \
RABBITMQ_HOST=192.168.1.100 \
go run ./cmd/server
```

## 生产建议

- 修改 JWT secret 和数据库密码
- 启用 HTTPS
- 设置 `server.mode: release`
- 日志格式切换为 `format: json`
- 调整连接池参数
- 配置 Prometheus 监控
- 部署 RabbitMQ 集群

## 许可证

MIT
