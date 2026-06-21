---
name: "code-standard"
description: "检查项目Go代码是否符合项目的代码规范和最佳实践。当用户提交代码前、进行代码审查、项目重构、询问规范问题、检查命名或代码风格时调用。也适用于用户说'规范检查'、'代码审查'、'check style'、'lint'、'代码质量'等场景。"
---

# 代码规范检查

## 技能信息

- **技能名称**: 代码规范检查
- **技能ID**: code-standard
- **描述**: 检查项目代码是否符合项目的代码规范和最佳实践
- **版本**: 2.0.0
- **作者**: System

## 适用场景

- 代码提交前的规范检查
- 代码审查过程中的规范验证
- 项目重构时的规范统一

## 代码规范

### 1. 项目结构

- **cmd目录**: 存放各服务的入口文件
- **internal目录**: 按业务模块划分
  - **{module}/api/**: HTTP 处理层
    - `handlers/` — Gin handler，调用 service，通过 `c.Error(err)` 传递错误
    - `dto/` — 请求/响应 DTO，请求用 `<Action><Entity>DTO` 命名
    - `routes/` — 路由注册
  - **{module}/service/**: 业务逻辑编排，事务管理，事件发布
  - **{module}/domain/**: 领域层
    - `models/` — 每个实体独立文件，GORM model
    - `repositories/` — 接口定义 + GORM 实现（同一文件）
  - **{module}/events/**: 领域事件结构体
- **internal/infra/**: 基础设施层
  - `repository/db.go` — `Repositories` 聚合，AutoMigrate 注册
  - `repository/models/` — PO 持久化对象，含 `ToDomain()`/`FromDomain()` 转换
  - `router/router.go` — 路由总入口
  - `eventbus/` — 内部事件总线，模块间异步通信
- **pkg/**: 公共代码
  - `config/`、`errcode/`、`response/`、`middleware/`、`logger/`、`utils/`、`query/`
- **configs/**: 配置文件
- **scripts/**: 脚本文件（RBAC 种子数据、数据迁移等）

### 2. 文件职责（单一职责）

- **一个文件只包含一个领域实体**的完整定义（model/repository/handler 同理）
- 例外：短小的辅助方法可以放同一文件
- **Repository 文件**: 接口 + 实现放在**同一个文件**，不要拆分 `interface` 和 `impl`
- 命名示例：`coupon_repo.go`（内含 `IcouponRepository` + `CouponRepository`），`user_coupon_repo.go`

### 3. 命名规范

- **包名**: 小写字母，简短清晰
- **结构体/接口名**: 驼峰，首字母大写
- **接口命名**: `I` + 小写首字母 + 资源名 + `Repository`，如 `IcouponRepository`、`IproductRepository`、`IorderRepository`
- **方法/变量名**: 驼峰，首字母小写
- **常量名**: 全大写，下划线分隔：`CouponTypeFixed`、`OrderStatusPending`
- **JSON标签**: 蛇形命名（小写+下划线）：`json:"coupon_type"`、`json:"min_amount"`
- **数据库表名**: 蛇形复数：`orders`、`order_items`、`user_coupons`、`promotion_products`
- **GORM 字段标签**: `gorm:"type:bigint;not null;index;comment:备注"`

### 4. DTO 命名约定

| 用途 | 命名格式 | 示例 |
|------|---------|------|
| 创建请求 | `Create<Entity>DTO` | `CreateOrderDTO`、`CreateCouponDTO` |
| 更新请求 | `Update<Entity>DTO` | `UpdateOrderDTO`、`UpdateCouponDTO` |
| 响应 | `<Entity>Response` | `OrderResponse`、`CouponResponse` |
| 列表结果 | `<Entity>ListResult` | `OrderListResult`、`CouponListResult` |
| 列表查询 | `<Entity>ListQuery` | `OrderListQuery`、`CouponListQuery` |

### 5. import 导入顺序

固定分三组，组间空行分隔：
1. **标准库**（`context`、`errors`、`strconv`、`time` 等）
2. **第三方库**（`gin`、`gorm`、`redis` 等）
3. **内部包**（`eshop-monolith/internal/...`、`eshop-monolith/pkg/...`）

### 6. 代码风格

- **缩进**: 使用 tab 缩进（Go 标准）
- **换行**: 每行不超过 120 字符
- **括号**: 左括号不单独占一行（Go 强制）
- **注释**: 使用中文注释，清晰说明代码功能
- **空行**: 函数间使用空行分隔，逻辑块间使用空行分隔

### 7. 错误处理

- **Handler 层**: 使用 `c.Error(err)` 传递错误，不直接返回 JSON
- **中间件层**: `ErrorHandler` 自动分类处理：
  - `*errcode.BizError` → 业务错误，返回对应 code + message
  - `validator.ValidationErrors` → 422 + 字段级错误详情
  - `gorm.ErrRecordNotFound` → **自动转为 404**
  - 其他 → 500 系统错误
- **Service 层**: 不关心 `gorm.ErrRecordNotFound`，返回原始 error，由 middleware 统一处理
- **Repository 层**: 不吞掉 `gorm.ErrRecordNotFound`，直接返回原始 error

### 8. 数据库操作

- **ORM框架**: 使用 gorm
- **主键**: 全表统一 `int64 autoIncrement`
- **金额存储**: 统一以「分」为单位（int64），避免浮点精度问题
- **时间字段**: `CreatedAt`、`UpdatedAt` 使用 `utils.Timestamp` 类型，JSON 序列化为毫秒级时间戳
- **软删除**: 嵌入 `gorm.DeletedAt`
- **表名**: 每个 Model 实现 `TableName()` 返回蛇形复数表名
- **外键关联**: 使用 `gorm:"foreignKey:..."` 标签

### 9. API 设计

- **路由结构**: RESTful，URL 含版本号 `/api/v1/`
- **请求参数**: 使用 DTO 结构体绑定
- **响应格式**: 统一 `response.Success(c, data)` 或 `response.Error(c, err)`
- **分页**: 嵌入 `query.Pagination`，返回 `query.ListResult[T]`
- **API 文档**: 使用 swaggo，@Tags 使用**英文小写复数**（`orders`、`coupons`、`promotions`、`products`）

### 10. Swagger 注解规范

```go
// ListCoupons 优惠券模板列表
// @Summary 优惠券模板列表
// @Description 分页查询优惠券模板列表
// @Tags coupons                        // ← 英文小写复数
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=dto.CouponListResult}
// @Router /api/v1/coupons [get]
```

- @Tags 必须用**英文**，与 route 分组名一致
- @Success 引用 `response.Response{data=...}`
- @Param 路径参数用 `path`，查询参数用 `query`，请求体用 `body`

### 11. Service 层规范

- 方法签名：`(ctx context.Context, ...) (*models.Xxx, error)`
- 事务内操作使用 `tx *gorm.DB` 传参，不直接操作 `s.db`
- 事件在事务**成功提交后**发布（`s.bus.Publish(...)` 在 `Transaction` block 外）
- 跨模块依赖通过 Service 构造函数注入，不直接引用其他模块的 Repository

### 12. Repository 层规范

- **文件结构**: 接口 + 实现放**同一个文件**，不拆分 `*_impl.go`
- **接口命名**: `I` + 小写首字母 + 资源名 + `Repository`（`IcouponRepository`）
- **实现结构体**: `CouponRepository`
- **构造方法**: `NewCouponRepository(db *gorm.DB) IcouponRepository`
- **事务方法**: 以 `WithTx` 后缀，接收 `*gorm.DB` 参数（`CreateWithTx(tx, ...)`）
- **非事务方法**: 接收 `context.Context`，使用 `r.db.WithContext(ctx)`
- **FindByID**: 不处理 `gorm.ErrRecordNotFound`，直接返回原始 error

### 13. Handler 层规范

- **结构体**: `XxxHandler`，组合对应 Service
- **固定模式**:
  1. 声明 DTO，`c.ShouldBindJSON(&req)` 或 `c.ShouldBindQuery(&req)`
  2. 错误走 `c.Error(err)`
  3. 成功走 `response.Success(c, data)`
- 用户 ID 从 JWT 上下文 `c.Get("user_id")` 获取
- 路径参数使用 `utils.ParseIntParam(c, "id")`

### 14. 性能优化

- **分页查询**: 支持分页和筛选
- **连接池**: 使用数据库连接池
- **缓存**: 合理使用缓存
- **并发**: 合理使用 goroutine

### 15. 安全规范

- **密码存储**: 使用 bcrypt 加密（`utils.CryptPassword`）
- **JWT验证**: 使用 `middleware.JWTAuth()` 中间件
- **输入验证**: 对所有输入进行 `binding` 标签校验
- **SQL注入**: 使用参数化查询（GORM 自动处理）

### 16. 测试规范

- **单元测试**: 为关键函数编写单元测试
- **集成测试**: 为服务间交互编写集成测试
- **测试覆盖率**: 保持较高的测试覆盖率

### 17. 层依赖方向

```text
API (handler + dto) → Service → Domain (model + repo interface) ← Repository (impl)
                          ↑
                    EventBus (模块间解耦)
```

- **依赖方向严格单向**: API → Service → Domain ← Repository
- **Repository 层实现 Domain 层定义的接口，依赖方向反转**

### 18. 模块结构模板

添加新模块时的标准目录：
```
internal/{module}/
├── api/
│   ├── dto/           # 请求/响应 DTO
│   ├── handlers/      # Gin handler
│   └── routes/        # 路由注册
├── service/            # 业务逻辑
├── domain/
│   ├── models/         # 领域模型（每实体一个文件）
│   └── repositories/   # 接口 + 实现（每 repo 一个文件）
└── events/             # 领域事件
```

## 检查规则

### 1. 项目结构检查

- [x] 检查项目是否遵循标准目录结构
- [x] 检查模块目录结构是否一致（api/service/domain/events）
- [x] 检查公共代码是否放在 pkg 目录
- [x] 检查 Repository 文件是否接口+实现在同一文件

### 2. 文件职责检查

- [x] 每个文件是否只包含一个领域实体
- [x] Repository 是否不拆分 `_impl.go`
- [x] Model 是否单个实体单个文件

### 3. 命名规范检查

- [x] 检查接口命名是否 `I` + 小写首字母（`IcouponRepository`）
- [x] 检查 DTO 命名是否符合约定（`CreateXxxDTO`、`XxxResponse`、`XxxListResult`）
- [x] 检查 JSON 标签是否使用蛇形命名法
- [x] 检查数据库表名是否使用蛇形复数
- [x] 检查 Swagger @Tags 是否英文小写复数

### 4. 导入顺序检查

- [x] 标准库 → 第三方 → 内部包 三组顺序
- [x] 组间有空行分隔

### 5. 错误处理检查

- [x] Repository 是否不吞 `gorm.ErrRecordNotFound`
- [x] Handler 是否使用 `c.Error(err)` 而非直接返回
- [x] 重要错误是否记录日志

### 6. 数据库操作检查

- [x] 主键类型是否正确（业务表 int64，用户表 UUID）
- [x] 金额是否以「分」为单位（int64）
- [x] 时间字段是否使用 `utils.Timestamp`
- [x] 是否实现 `TableName()` 方法
- [x] 是否使用软删除

### 7. API 设计检查

- [x] 路由是否遵循 RESTful 风格
- [x] URL 是否包含版本号 `/api/v1/`
- [x] 是否使用 DTO 绑定参数
- [x] 是否使用统一响应格式
- [x] 是否使用 swaggo 生成文档
- [x] @Tags 是否英文

### 8. Service 层检查

- [x] 方法签名是否 `(ctx context.Context, ...)`
- [x] 事件发布是否在事务成功之后
- [x] 跨模块依赖是否通过构造函数注入

## 工具集成

### 1. 代码检查工具

- **gofmt**: 代码格式化
- **go vet**: 静态分析可疑结构
- **staticcheck**: 高级静态分析
- **goimports**: 自动整理导入顺序
- **go test**: 运行测试

## 常见问题

### 1. Repository 文件拆分

**问题**: 接口和实现拆在不同文件（`xxx_repo.go` + `xxx_repo_impl.go`）

**解决方案**: 合并到同一个文件 `xxx_repo.go`，接口在上，实现在下。

### 2. 命名不一致

**问题**: 变量名、函数名、结构体名命名风格不一致

**解决方案**: 统一使用驼峰命名法，遵循项目规范

### 3. 错误处理不当

**问题**: Repository 层吞掉 gorm.ErrRecordNotFound，Handler 层直接返回 JSON

**解决方案**: Repository 返回原始 error，middleware 统一处理；Handler 使用 `c.Error(err)`

### 4. DTO 后缀混乱

**问题**: 同时出现 `CreateXxxReq`、`CreateXxxDTO`、`CreateXxxRequest`

**解决方案**: 统一使用 `CreateXxxDTO` / `UpdateXxxDTO`

### 5. Swagger Tags 中英文混用

**问题**: 部分 @Tags 用中文，部分用英文

**解决方案**: 统一使用英文小写复数（`orders`、`coupons`、`promotions`）

### 6. 金额存储使用 float64

**问题**: 使用 float64 或 decimal 存储金额，存在精度问题

**解决方案**: 统一以「分」为单位用 int64

## 示例代码

### Model 定义示例

```go
package models

import "eshop-monolith/pkg/utils"

type CouponType string

const (
	CouponTypeFixed      CouponType = "fixed"
	CouponTypePercentage CouponType = "percentage"
)

type Coupon struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	CouponType  CouponType      `json:"coupon_type"`
	Value       int64           `json:"value"`      // 面值，单位：分
	MinAmount   int64           `json:"min_amount"` // 最低消费，单位：分
	TotalStock  int             `json:"total_stock"`
	RemainStock int             `json:"remain_stock"`
	StartTime   utils.Timestamp `json:"start_time"`
	EndTime     utils.Timestamp `json:"end_time"`

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (Coupon) TableName() string { return "coupons" }
```

### Repository 示例（一个文件内接口+实现）

```go
package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"
	"gorm.io/gorm"
)

// IcouponRepository 优惠券模板仓储接口
type IcouponRepository interface {
	Create(ctx context.Context, coupon *models.Coupon) error
	FindByID(ctx context.Context, id int64) (*models.Coupon, error)
	CreateWithTx(tx *gorm.DB, coupon *models.Coupon) error
}

// CouponRepository 优惠券模板仓储实现
type CouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) IcouponRepository {
	return &CouponRepository{db: db}
}

func (r *CouponRepository) FindByID(ctx context.Context, id int64) (*models.Coupon, error) {
	var po repoModels.CouponPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		return nil, err // 不处理 gorm.ErrRecordNotFound
	}
	return po.ToDomain(), nil
}
```

### HTTP Handler 示例

```go
// ListCoupons 优惠券模板列表
// @Summary 优惠券模板列表
// @Description 分页查询优惠券模板列表
// @Tags coupons
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=dto.CouponListResult}
// @Router /api/v1/coupons [get]
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	// ...
	list, total, err := h.couponService.ListCoupons(c, int(page), int(pageSize))
	if err != nil {
		c.Error(err)  // 交给 ErrorHandler 处理
		return
	}
	response.Success(c, dto.CouponListResult{Total: total, List: items})
}
```

### Service 层示例

```go
func (s *CouponService) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	coupon.RemainStock = coupon.TotalStock
	coupon.Status = models.CouponStatusActive
	return s.couponRepo.Create(ctx, coupon)
}
```
