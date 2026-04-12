# eshop-monolith

一个完整的电商单体应用系统 Demo，展示了清晰的分层架构和最佳实践。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **配置管理**: Viper
- **日志**: Zap

## 核心特性

- ✅ **清晰分层架构**: API → Service → Domain → Repository 四层架构
- ✅ **模块化设计**: 订单、库存、用户三大模块低耦合高内聚
- ✅ **统一认证授权**: JWT 认证 + RBAC 权限控制
- ✅ **本地事务保证**: 基于 MySQL 的 ACID 事务
- ✅ **幂等性控制**: 基于 Redis 的请求幂等性保证
- ✅ **事件驱动**: 内部事件总线实现模块间解耦
- ✅ **可观测性**: 结构化日志 + 请求追踪

## 项目亮点

### 1. 清晰的架构分层

```
API 层 (Handlers) → 处理 HTTP 请求/响应
    ↓
Service 层 → 业务逻辑编排
    ↓
Domain 层 → 领域模型 + 仓储接口
    ↓
Repository 层 → 数据持久化
```

### 2. 模块化设计

- **订单模块**: 订单管理、状态流转、事件发布
- **库存模块**: 产品管理、库存管理、预占/释放
- **用户模块**: 用户管理、JWT 认证、RBAC 权限

### 3. 本地事务保证

使用 MySQL 数据库事务保证数据一致性，避免了分布式事务的复杂性：

```go
// 示例：创建订单时自动扣减库存
func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建订单
        // 2. 扣减库存
        // 3. 发布订单创建事件
        return nil
    })
}
```

### 4. 事件驱动架构

通过内部事件总线实现模块间解耦：

- 订单创建 → 发送通知
- 库存不足 → 触发补货
- 订单完成 → 更新统计

## 适用场景

本 Demo 适合用于：

- 📚 学习 Go 语言分层架构设计
- 🎓 面试项目展示（展示架构设计能力）
- 💼 中小型项目快速开发模板
- 🚀 从零开始的 Go Web 项目脚手架

## 目录结构

```
eshop-monolith/
├── cmd/
│   └── server/
│       └── main.go                    # 应用入口
├── internal/
│   ├── api/                           # API 层
│   │   └── routes/
│   │       └── router.go              # 统一路由注册
│   ├── domain/                        # 领域层
│   │   └── shared/                    # 共享领域
│   │       ├── errors.go              # 领域错误
│   │       ├── models.go              # 共享模型
│   │       └── value_objects.go       # 值对象
│   ├── eventbus/                      # 内部事件总线
│   │   ├── bus.go                     # 事件总线实现
│   │   ├── handlers.go                # 事件处理器聚合
│   │   ├── inventory_handlers.go      # 库存事件处理器
│   │   ├── order_handlers.go          # 订单事件处理器
│   │   ├── payment_handlers.go        # 支付事件处理器
│   │   └── user_handlers.go           # 用户事件处理器
│   ├── inventory/                     # 库存模块
│   │   ├── api/
│   │   │   ├── dto/
│   │   │   │   ├── category_dto.go    # 分类DTO
│   │   │   │   ├── inventory_dto.go   # 库存DTO
│   │   │   │   └── product_dto.go     # 产品DTO
│   │   │   ├── handlers/
│   │   │   │   ├── category_handler.go # 分类处理器
│   │   │   │   ├── inventory_handler.go # 库存处理器
│   │   │   │   └── product_handler.go   # 产品处理器
│   │   │   └── routes/
│   │   │       ├── category_routes.go  # 分类路由
│   │   │       ├── inventory_routes.go # 库存路由
│   │   │       └── product_routes.go   # 产品路由
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   │   ├── category.go        # 分类模型
│   │   │   │   ├── inventory.go       # 库存模型
│   │   │   │   └── product.go         # 产品模型
│   │   │   └── repositories/
│   │   │       ├── category_repo.go    # 分类仓储
│   │   │       ├── inventory_repo.go   # 库存仓储
│   │   │       └── product_repo.go     # 产品仓储
│   │   ├── events/
│   │   │   ├── category_event.go       # 分类事件
│   │   │   ├── inventory_event.go      # 库存事件
│   │   │   └── product_event.go        # 产品事件
│   │   └── service/
│   │       ├── category_service.go     # 分类服务
│   │       ├── inventory_service.go    # 库存服务
│   │       └── product_service.go      # 产品服务
│   ├── order/                         # 订单模块
│   │   ├── api/
│   │   │   ├── dto/
│   │   │   │   └── order_dto.go       # 订单DTO
│   │   │   ├── handlers/
│   │   │   │   └── order_handler.go   # 订单处理器
│   │   │   └── routes/
│   │   │       └── order_routes.go     # 订单路由
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   │   └── order.go           # 订单模型
│   │   │   └── repositories/
│   │   │       └── order_repo.go       # 订单仓储
│   │   ├── events/
│   │   │   └── order_event.go          # 订单事件
│   │   └── service/
│   │       └── order_service.go        # 订单服务
│   ├── payment/                       # 支付模块
│   │   ├── api/
│   │   │   ├── dto/
│   │   │   │   └── payment_dto.go      # 支付DTO
│   │   │   ├── handlers/
│   │   │   │   └── payment_handler.go # 支付处理器
│   │   │   └── routes/
│   │   │       └── payment_routes.go   # 支付路由
│   │   ├── domain/
│   │   │   ├── models/
│   │   │   │   └── payment.go         # 支付模型
│   │   │   └── repositories/
│   │   │       ├── payment_method_repo.go # 支付方式仓储
│   │   │       ├── payment_repo.go     # 支付仓储
│   │   │       └── refund_repo.go      # 退款仓储
│   │   ├── events/
│   │   │   └── payment_event.go        # 支付事件
│   │   └── service/
│   │       └── payment_service.go      # 支付服务
│   ├── pkg/                           # 内部公共包
│   │   ├── config/                    # 配置管理
│   │   │   └── config.go              # 配置实现
│   │   ├── errcode/                   # 错误码
│   │   │   └── errcode.go             # 错误码定义
│   │   ├── logger/                    # 日志工具
│   │   │   └── logger.go              # 日志实现
│   │   ├── middleware/                # 中间件
│   │   │   ├── errorhandler.go        # 错误处理中间件
│   │   │   ├── jwtauth.go             # JWT 认证中间件
│   │   │   └── rbac.go                # RBAC 权限中间件
│   │   ├── query/                     # 查询工具
│   │   │   └── query.go               # 查询实现
│   │   ├── response/                  # 统一响应
│   │   │   └── response.go            # 响应实现
│   │   └── utils/                     # 工具函数
│   │       ├── cryptopwd.go           # 密码加密
│   │       ├── parseIntParam.go       # 参数解析
│   │       └── timestamp.go           # 时间戳工具
│   ├── repository/                    # 仓储实现层
│   │   └── db.go                      # 数据库连接
│   └── user/                          # 用户模块
│       ├── api/
│       │   ├── dto/
│       │   │   ├── permission_dto.go  # 权限DTO
│       │   │   ├── userIdentity_dto.go # 用户身份DTO
│       │   │   ├── userInfo_dto.go     # 用户信息DTO
│       │   │   └── user_dto.go         # 用户DTO
│       │   ├── handlers/
│       │   │   ├── auth_handler.go     # 认证处理器
│       │   │   ├── permission_handler.go # 权限处理器
│       │   │   ├── role_handler.go     # 角色处理器
│       │   │   └── user_handler.go     # 用户处理器
│       │   └── routes/
│       │       ├── auth_routes.go      # 认证路由
│       │       ├── permission_routes.go # 权限路由
│       │       ├── role_routes.go      # 角色路由
│       │       └── user_routes.go      # 用户路由
│       ├── domain/
│       │   ├── auth/
│       │   │   └── provider.go        # 认证提供商
│       │   ├── models/
│       │   │   ├── auth_token.go      # 认证令牌模型
│       │   │   ├── permission.go       # 权限模型
│       │   │   ├── role.go             # 角色模型
│       │   │   ├── user.go             # 用户模型
│       │   │   ├── user_identity.go    # 用户身份模型
│       │   │   └── user_info.go        # 用户信息模型
│       │   └── repositories/
│       │       ├── authToken_repo.go   # 认证令牌仓储
│       │       ├── permission_repo.go  # 权限仓储
│       │       ├── role_repo.go        # 角色仓储
│       │       ├── userIdentity_repo.go # 用户身份仓储
│       │       ├── userInfo_repo.go    # 用户信息仓储
│       │       └── user_repo.go        # 用户仓储
│       ├── events/
│       │   └── user_event.go           # 用户事件
│       └── service/
│           ├── auth_service.go         # 认证服务
│           ├── permission_service.go   # 权限服务
│           ├── token_service.go        # 令牌服务
│           └── user_service.go         # 用户服务
├── configs/
│   └── config.yaml                    # 统一配置文件
├── docs/
│   ├── API.md                         # API 文档
│   ├── DEPLOYMENT.md                  # 部署指南
│   └── DEVELOPMENT.md                 # 开发指南
├── scripts/
│   ├── init.sql                       # 数据库初始化
│   ├── seed.sql                       # 测试数据
│   └── permissions.sql                # 权限数据
├── test/
│   ├── unit/                          # 单元测试
│   ├── integration/                   # 集成测试
│   └── e2e/                           # 端到端测试
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 架构说明

### 分层职责

| 层级          | 目录            | 职责                                | 依赖                     |
| ------------- | --------------- | ----------------------------------- | ------------------------ |
| API 层        | `*/api/`        | HTTP 请求处理、参数验证、响应格式化 | Service 层               |
| Service 层    | `*/service/`    | 业务逻辑编排、事务管理、事件发布    | Domain 层、Repository 层 |
| Domain 层     | `*/domain/`     | 领域模型定义、业务规则、仓储接口    | 无（核心层）             |
| Repository 层 | `*/repository/` | 数据持久化、SQL 操作                | Domain 层                |
| EventBus      | `eventbus/`     | 模块间事件通信                      | Service 层               |

### 依赖方向

```
API 层 → Service 层 → Domain 层 ← Repository 层
                      ↑
                 EventBus
```

## 当前实现

### 1. 订单模块

**API 端点**:

- `POST /api/v1/orders` 创建订单（本地事务）
- `GET /api/v1/orders` 订单列表（支持分页、筛选）
- `GET /api/v1/orders/:id` 订单详情
- `PUT /api/v1/orders/:id` 更新订单状态
- `DELETE /api/v1/orders/:id` 取消订单
- `POST /api/v1/orders/:id/cancel` 取消订单
- `PATCH /api/v1/orders/:id/status` 更新订单状态

**核心功能**:

- ✅ 订单 CRUD
- ✅ 本地事务保证（订单创建 + 库存扣减）
- ✅ 订单状态机管理
- ✅ 领域事件发布

### 2. 库存模块

**API 端点**:

- `POST /api/v1/products` 创建产品
- `GET /api/v1/products` 产品列表
- `GET /api/v1/products/:id` 产品详情
- `PUT /api/v1/products/:id` 更新产品
- `DELETE /api/v1/products/:id` 删除产品
- `POST /api/v1/inventories` 创建库存
- `GET /api/v1/inventories` 库存列表
- `POST /api/v1/inventories/reserve` 预占库存
- `POST /api/v1/inventories/release` 释放库存
- `POST /api/v1/categories` 创建分类
- `GET /api/v1/categories` 分类列表
- `GET /api/v1/categories/:id` 分类详情
- `PUT /api/v1/categories/:id` 更新分类
- `DELETE /api/v1/categories/:id` 删除分类

**核心功能**:

- ✅ 产品管理（CRUD）
- ✅ 库存管理（预占、释放、调整）
- ✅ 库存检查
- ✅ 库存预警
- ✅ 分类管理（CRUD）

### 3. 用户模块

**API 端点**:

- `POST /api/v1/users/register` 用户注册
- `POST /api/v1/users/login` 用户登录
- `GET /api/v1/users/profile` 获取用户资料
- `PUT /api/v1/users/profile` 更新用户资料
- `POST /api/v1/auth/refresh` 刷新 Token
- `GET /api/v1/roles` 角色管理
- `GET /api/v1/permissions` 权限管理
- `POST /api/v1/permissions/check` 权限检查
- `GET /api/v1/users` 用户列表
- `GET /api/v1/users/:user_id` 获取用户详情
- `GET /api/v1/users/:user_id/roles` 获取用户角色
- `POST /api/v1/users/:user_id/roles` 分配角色给用户
- `DELETE /api/v1/users/:user_id/roles/:role_id` 从用户移除角色

**核心功能**:

- ✅ 用户注册与登录
- ✅ JWT 认证
- ✅ 密码加密（bcrypt）
- ✅ Token 刷新机制
- ✅ RBAC 权限控制
- ✅ 用户信息管理

### 4. 支付模块

**API 端点**:

- `POST /api/v1/payments` 创建支付
- `GET /api/v1/payments` 支付列表（支持分页、筛选）
- `GET /api/v1/payments/:id` 支付详情
- `PATCH /api/v1/payments/:id/status` 更新支付状态
- `GET /api/v1/orders/payment/:order_id` 根据订单ID获取支付
- `POST /api/v1/refunds` 创建退款
- `GET /api/v1/refunds` 退款列表（支持分页、筛选）
- `PATCH /api/v1/refunds/:id/status` 更新退款状态
- `GET /api/v1/payment-methods` 获取支付方式列表

**核心功能**:

- ✅ 支付管理（创建、查询、更新状态）
- ✅ 退款管理（创建、查询、更新状态）
- ✅ 支付方式管理
- ✅ 与订单系统集成
- ✅ 支付状态变更触发订单状态更新
- ✅ 领域事件发布

**API 端点**:

- `POST /api/v1/users/register` 用户注册
- `POST /api/v1/users/login` 用户登录
- `GET /api/v1/users/profile` 获取用户资料
- `PUT /api/v1/users/profile` 更新用户资料
- `POST /api/v1/auth/refresh` 刷新 Token
- `GET /api/v1/roles` 角色管理
- `GET /api/v1/permissions` 权限管理
- `POST /api/v1/permissions/check` 权限检查

**核心功能**:

- ✅ 用户注册与登录
- ✅ JWT 认证
- ✅ 密码加密（bcrypt）
- ✅ Token 刷新机制
- ✅ RBAC 权限控制

## 本地运行

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 7+
- Make (可选)

### 运行步骤

#### 1. 启动依赖服务

```bash
# 本地安装并启动 MySQL 和 Redis
# MySQL: 确保 3306 端口可用
# Redis: 确保 6379 端口可用
```

#### 2. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE eshop_db;"

# 执行初始化脚本
mysql -u root -p eshop_db < scripts/init.sql

# 导入测试数据
mysql -u root -p eshop_db < scripts/seed.sql

# 导入权限数据
mysql -u root -p eshop_db < scripts/permissions.sql
```

#### 3. 配置环境

```bash
# 配置文件已存在于 configs/config.yaml
# 可根据需要修改数据库密码等配置
# 例如：修改 MySQL 密码
# vim configs/config.yaml
```

#### 4. 运行应用

```bash
# 下载依赖
go mod download

# 运行应用
go run ./cmd/server

# 或使用 make
make run
```

### 使用 Makefile

```bash
# 查看所有命令
make help

# 安装依赖
make deps

# 运行测试
make test

# 构建应用
make build

# 运行应用
make run

# 清理
make clean
```

## 配置说明

### 配置文件 `configs/config.yaml`

```yaml
server:
  port: 8080
  mode: debug # debug, release, test
  read_timeout: 30s
  write_timeout: 30s

mysql:
  host: localhost
  port: 3306
  username: root
  password: root
  database: eshop_db
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10

jwt:
  secret: "your-secret-key-change-in-production"
  expire_hours: 24
  refresh_expire_hours: 168 # 7 days

log:
  level: info # debug, info, warn, error
  format: json # json, text
  output: stdout # stdout, file
  file_path: logs/app.log

rate_limit:
  enabled: true
  requests_per_second: 100
  burst: 200

cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:8080"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Origin"
    - "Content-Type"
    - "Authorization"
```

### 环境变量覆盖

所有配置都支持环境变量覆盖（使用 `.` 替换为 `_`）：

```bash
SERVER_PORT=8081 \
MYSQL_PASSWORD=secret \
JWT_SECRET=production-secret \
go run ./cmd/server
```

## API 使用示例

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

### 创建产品

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {access_token}" \
  -d '{
    "name": "iPhone 15",
    "description": "Apple iPhone 15 Pro Max",
    "price": 999900,
    "sku": "IPHONE15-001"
  }'
```

### 创建订单（自动扣减库存）

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {access_token}" \
  -d '{
    "items": [
      {
        "product_id": "product-uuid",
        "quantity": 2,
        "unit_price": 999900
      }
    ],
    "shipping_address": {
      "province": "Guangdong",
      "city": "Shenzhen",
      "address": "Nanshan District"
    }
  }'
```

### 查询订单

```bash
curl -X GET http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer {access_token}" \
  -H "X-Idempotency-Key: unique-request-id"
```

## 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./test/unit/...

# 运行集成测试
go test ./test/integration/...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 性能测试

```bash
# 使用 wrk 进行压力测试
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/health

# 使用 ab 测试
ab -n 10000 -c 100 http://localhost:8080/api/v1/health
```

## 性能优化建议

### 数据库优化

- ✅ 添加合适的索引
- ✅ 使用连接池
- ✅ 避免 N+1 查询
- ✅ 使用预编译语句

### 缓存策略

- ✅ 热点数据缓存（用户信息、产品信息）
- ✅ 本地缓存 + Redis 二级缓存
- ✅ 缓存穿透、雪崩防护

### 并发控制

- ✅ 使用 sync.Pool 复用对象
- ✅ 限流中间件
- ✅ 熔断降级

## 从微服务迁移到单体的变化

### 移除的组件

- ❌ gRPC 客户端/服务端
- ❌ RabbitMQ 消息队列
- ❌ Saga 分布式事务协调器
- ❌ 服务发现和负载均衡
- ❌ 独立的服务配置文件

### 新增/改造的组件

- ✅ 内部事件总线（替代 MQ）
- ✅ 统一的事务管理器
- ✅ 模块间直接调用
- ✅ 统一的配置管理
- ✅ 共享的领域模型

### 架构对比

| 方面     | 微服务版本      | 单体版本      |
| -------- | --------------- | ------------- |
| 部署单元 | 3个独立服务     | 1个应用       |
| 事务     | Saga 最终一致性 | ACID 强一致性 |
| 通信     | HTTP/gRPC + MQ  | 直接方法调用  |
| 数据存储 | 3个独立数据库   | 1个共享数据库 |
| 复杂度   | 高              | 中            |
| 开发效率 | 低              | 高            |
| 运维成本 | 高              | 低            |

## 生产部署建议

### 1. 安全配置

```yaml
# 生产环境必须修改
jwt:
  secret: "使用强随机密钥（至少32位）"

mysql:
  password: "使用强密码"

# 启用 HTTPS
server:
  tls_enabled: true
  cert_file: /path/to/cert.pem
  key_file: /path/to/key.pem
```

### 2. 日志和监控

```yaml
# 使用结构化日志
log:
  format: json
  output: file
  file_path: /var/log/eshop/app.log

# 启用健康检查
health:
  enabled: true
  path: /health
```

### 3. 性能调优

```bash
# 设置 GOMAXPROCS
export GOMAXPROCS=8

# 设置 GC 百分比
export GOGC=100

# 使用优化编译
go build -ldflags="-s -w" -o app ./cmd/server
```

## 故障排查

### 应用无法启动

1. 检查配置文件：
   ```bash
   cat configs/config.yaml
   ```
2. 检查数据库连接：
   ```bash
   mysql -h localhost -u root -p -e "SELECT 1"
   ```
3. 检查 Redis 连接：
   ```bash
   redis-cli ping
   ```
4. 查看应用日志：
   ```bash
   # 直接查看应用输出日志
   ```

### 性能问题

1. 开启慢查询日志：
   ```sql
   SET GLOBAL slow_query_log = 'ON';
   SET GLOBAL long_query_time = 1;
   ```
2. 查看数据库连接数：
   ```sql
   SHOW PROCESSLIST;
   ```
3. 分析索引使用：
   ```sql
   EXPLAIN SELECT * FROM orders WHERE user_id = 'xxx';
   ```

### 常见错误

| 错误                        | 原因           | 解决方案               |
| --------------------------- | -------------- | ---------------------- |
| `connection refused`        | 数据库未启动   | 启动本地 MySQL 服务    |
| `invalid token`             | JWT 密钥不匹配 | 检查 `jwt.secret` 配置 |
| `duplicate key`             | 唯一约束冲突   | 检查业务逻辑           |
| `context deadline exceeded` | 超时           | 增加超时时间或优化查询 |

## 开发指南

### 添加新功能模块

1. 创建领域模型：
   ```go
   // internal/domain/payment/models.go
   type Payment struct {
       ID     string
       Amount int64
       Status string
   }
   ```
2. 定义仓储接口：
   ```go
   // internal/domain/payment/repository.go
   type PaymentRepository interface {
       Create(ctx context.Context, payment *Payment) error
       FindByID(ctx context.Context, id string) (*Payment, error)
   }
   ```
3. 实现仓储：
   ```go
   // internal/repository/payment_repo.go
   type paymentRepo struct {
       db *gorm.DB
   }
   ```
4. 创建服务：
   ```go
   // internal/service/payment_service.go
   type PaymentService struct {
       repo PaymentRepository
   }
   ```
5. 添加处理器：
   ```go
   // internal/api/handlers/payment_handler.go
   func (h *PaymentHandler) CreatePayment(c *gin.Context) {
       // 处理请求
   }
   ```
6. 注册路由：
   ```go
   // internal/api/routes/router.go
   func SetupRoutes(r *gin.Engine, handlers *Handlers) {
       paymentGroup := r.Group("/api/v1/payments")
       paymentGroup.POST("/", handlers.Payment.CreatePayment)
   }
   ```

### 代码规范

1. **命名规范**：
   - 包名：小写，单数
   - 文件名：snake_case
   - 变量/函数：camelCase
   - 导出类型/函数：PascalCase
2. **错误处理**：
   ```go
   if err != nil {
       return fmt.Errorf("failed to create order: %w", err)
   }
   ```
3. **日志规范**：
   ```go
   logger.Info("order created",
       zap.String("order_id", order.ID),
       zap.Int64("amount", order.TotalAmount),
   )
   ```

### Git 提交规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

示例：

```
feat(order): add order cancellation feature

- Add cancel order endpoint
- Implement order status validation
- Add refund logic for cancelled orders

Closes #123
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### Pull Request 要求

- ✅ 通过所有测试
- ✅ 更新相关文档
- ✅ 符合代码规范
- ✅ 添加必要的注释

## 许可证

本项目采用 MIT 许可证 - 查看 LICENSE 文件了解详情

## 联系方式

- 项目链接: <https://github.com/yourusername/eshop-monolith>
- 问题反馈: [Issues](https://github.com/yourusername/eshop-monolith/issues)

## 致谢

感谢所有贡献者的支持！

```

这个 README 完整描述了单体项目的架构、使用方法、API 示例、部署指南等，相比原微服务版本：

1. **简化了架构说明** - 突出分层架构和模块化设计
2. **移除了分布式特性** - 强调本地事务和内部事件总线
3. **统一了配置管理** - 单一配置文件
4. **提供了迁移对比** - 帮助理解两种架构的差异
5. **增加了更多实用内容** - 性能优化、故障排查、开发指南等
```
