# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Run the application
go run ./cmd/server

# Build the binary
go build -o bin/server ./cmd/server

# Download dependencies
go mod download

# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Tidy dependencies
go mod tidy
```

## Project Architecture

电商单体应用, 使用 Go + Gin + GORM + MySQL + Redis。

### Layer Dependency

```
API (handler + dto) → Service → Domain (model + repo interface) ← Repository (impl)
                          ↑
                    EventBus (模块间解耦)
```

- **依赖方向严格单向**: API → Service → Domain ← Repository
- **Repository 层实现 Domain 层定义的接口, 依赖方向反转**
- **EventBus 内部事件总线**替代消息队列, 实现模块间异步通信

### Module Structure (`internal/{module}/`)

每个业务模块（order、inventory、user、payment、cart）独立目录:

| 目录 | 职责 |
|------|------|
| `handler_*.go` | Gin handler, 调用 service, 通过 `c.Error(err)` 传递错误 |
| `dto_*.go` | 请求/响应 DTO, 使用 `query.Pagination` 嵌入分页 |
| `service_*.go` | 业务逻辑编排, 事务管理, 事件发布 |
| `model_*.go` | GORM model, 主键用 `int64 autoIncrement` |
| `repo_*.go` | Repository 接口 + GORM 实现（同一文件） |
| `events.go` | 领域事件结构体 |

> 注: 模块采用扁平结构, 文件按 `{kind}_{entity}.go` 命名, 不拆分 api/service/domain 子目录

### Cross-Module Patterns

- **`internal/infra/repository/db.go`**: `Repositories` 聚合所有模块的 repo 实现, 在 `main.go` 中初始化后注入路由；`InitDB()` 只负责连接, 不做自动迁移
- **表结构管理**: 不启用 AutoMigrate（会改动线上数据库）, 建表用 `docs/schema.sql` 手动执行；模型变更后用 `go run scripts/gen_schema.go > docs/schema.sql` 重新生成 DDL
- **`internal/infra/eventbus/`**: 订阅/发布模式, 按模块拆分 handler 文件
- **`internal/infra/router/router.go`**: 路由总入口, 所有模块路由在此注册

### Shared Packages (`internal/pkg/`)

| 包 | 用途 |
|----|------|
| `config/` | Viper 从 `configs/config.yaml` 加载, 环境变量覆盖 |
| `errcode/` | 业务错误码 `BizError{Code, Message}` |
| `response/` | 统一 JSON 响应: `Success`, `BizError`, `SysError`, `BindError` |
| `middleware/` | `ErrorHandler`(全局panic恢复+错误分类), `JWTAuth`, `RequirePermission`/`RequireRole` |
| `logger/` | Zap SugaredLogger, 结构化日志带请求上下文 |
| `utils/` | `Timestamp`(毫秒JSON), `ParseIntParam`, `ParseToken`, `CryptPassword` |
| `query/` | 分页(`Pagination`), 排序, `ListResult[T]` 泛型列表 |

### Key Patterns

- **错误处理**: handler 中调用 `c.Error(err)` → `ErrorHandler` 中间件自动分类:
  - `*errcode.BizError` → 业务错误, 返回对应 code + message
  - `validator.ValidationErrors` → 422 + 字段级错误详情
  - `gorm.ErrRecordNotFound` → **自动转为 404**, Repository/Service 层不处理
  - 其他 → 500 系统错误
- **Repository 文件**: 接口 + 实现放在**同一文件**, 不拆分 `_impl.go`; 接口命名 `I`+小写首字母(`IcouponRepository`, `IorderRepository`)
- **DTO 命名**: 请求用 `CreateXxxDTO` / `UpdateXxxDTO`, 响应用 `XxxResponse`, 列表用 `XxxListResult`
- **导入顺序**: 标准库 → 第三方 → 内部包, 三组间空行分隔
- **金额存储**: 统一以「分」为单位(int64), 避免浮点精度问题
- **时间字段**: 使用 `utils.Timestamp` 类型, JSON 序列化为毫秒级时间戳
- **表名**: GORM `TableName()` 返回带模块前缀的蛇形复数 (`tx_orders`、`usr_users`、`sp_products`)
- **软删除**: `gorm.DeletedAt` 嵌入 model
- **密码**: bcrypt (`utils.CryptPassword`)
- **主键**: 全表统一 `int64 autoIncrement`, 无 UUID 主键
- **Swagger**: @Tags 使用**英文小写复数** (`orders`、`coupons`、`promotions`)

### Entry Point

`cmd/server/main.go`:
1. `config.Load()` 加载配置
2. `repository.InitDB()` + `repository.InitRedis()`
3. `repository.NewRepositories()` 聚合所有 repo
4. `routes.SetupRouter()` 注册路由, 创建 EventBus
5. 启动 HTTP 服务 + 优雅关闭

### 添加新模块步骤

1. `internal/{module}/model_{entity}.go` — 定义 GORM model（每实体一个文件, 含 `TableName()`）
2. `internal/{module}/repo_{entity}.go` — 定义 repo 接口 + GORM 实现（同一文件）
3. `internal/{module}/service_{entity}.go` — 业务逻辑
4. `internal/{module}/dto_{entity}.go` + `handler_{entity}.go` — DTO + Handler
5. `internal/{module}/routes.go` — 路由注册
6. `internal/infra/router/router.go` — 挂载路由, 如需与订单系统集成则传递 Service
7. 重新生成建表脚本: `go run scripts/gen_schema.go > docs/schema.sql`（新增表时手动执行）
