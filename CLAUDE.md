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
| `api/handlers/` | Gin handler, 调用 service, 通过 `c.Error(err)` 传递错误 |
| `api/dto/` | 请求/响应 DTO, 使用 `query.Pagination` 嵌入分页 |
| `api/routes/` | 路由注册 |
| `service/` | 业务逻辑编排, 事务管理, 事件发布 |
| `domain/models/` | GORM model, 主键用 `int64 autoIncrement` |
| `domain/repositories/` | Repository 接口定义 |
| `events/` | 领域事件结构体 |

### Cross-Module Patterns

- **`internal/repository/db.go`**: `Repositories` 聚合所有模块的 repo 实现, 在 `main.go` 中初始化后注入路由
- **`internal/eventbus/`**: 订阅/发布模式, 按模块拆分 handler 文件
- **`internal/domain/shared/`**: 跨模块共享的领域模型（如 `ProductCategory`）和领域错误

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

- **错误处理**: handler 中调用 `c.Error(err)` → `ErrorHandler` 中间件自动分类(BizError/Validation/System)并返回对应响应
- **金额存储**: 统一以「分」为单位(int64), 避免浮点精度问题
- **时间字段**: 使用 `utils.Timestamp` 类型, JSON 序列化为毫秒级时间戳
- **表名**: GORM `TableName()` 返回蛇形复数 (orders, order_items)
- **软删除**: `gorm.DeletedAt` 嵌入 model
- **密码**: bcrypt (`utils.CryptPassword`)
- **自动迁移**: `repository.InitDB` 中 `db.AutoMigrate()` 所有 model

### Entry Point

`cmd/server/main.go`:
1. `config.Load()` 加载配置
2. `repository.InitDB()` + `repository.InitRedis()`
3. `repository.NewRepositories()` 聚合所有 repo
4. `routes.SetupRouter()` 注册路由, 创建 EventBus
5. 启动 HTTP 服务 + 优雅关闭

### 添加新模块步骤

1. `internal/{module}/domain/models/` — 定义 model
2. `internal/{module}/domain/repositories/` — 定义 repo 接口
3. `internal/repository/db.go` — 实现 repo, 加入 `Repositories` struct, 注册 AutoMigrate
4. `internal/{module}/service/` — 业务逻辑
5. `internal/{module}/api/handlers/` + `api/dto/` — HTTP 层
6. `internal/{module}/api/routes/` — 路由注册
7. `internal/api/routes/router.go` 中挂载路由
