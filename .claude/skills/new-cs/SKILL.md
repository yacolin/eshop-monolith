---
name: "new-cs"
description: "检查重构后项目代码是否符合新的扁平化代码规范（Flat Go Structure）。基于商品中心 product 模块的 brand 参考实现提炼。适用于 code review、新模块开发、规范统一。"
---

# 新代码规范（Flat Go Structure）

## 适用范围

本规范适用于商品中心（`internal/product/`）、库存中心（`internal/inventory/`）等**重构后的新模块**。旧模块仍遵循旧规范。

## 核心原则

1. **一个模块一个 package** — 不拆分子 package，所有文件在模块根目录平铺
2. **model 即 PO** — 无 `ToDomain()`/`FromDomain()` 双重模型转换
3. **文件按实体前缀命名** — `model_spu.go` `repo_spu.go` `service_spu.go`
4. **DTO/Handler 按实体拆分** — `dto_{entity}.go`、`handler_{entity}.go`

## 1. 目录结构

```
internal/{module}/           # 一个 package
├── dto_{entity}.go          #   请求/响应 DTO
├── handler_{entity}.go      #   Handler + 路由注册
├── model_{entity}.go        #   GORM 模型（同时也是 PO）
├── repo_{entity}.go         #   接口 + 实现（同一文件）
└── service_{entity}.go      #   业务逻辑
```

### 示例：product 模块

```
internal/product/
├── dto_brand.go             # CreateBrandReq, UpdateBrandReq, BrandListReq
├── dto_category.go          # CreateCategoryReq, UpdateCategoryReq
├── handler_brand.go         # BrandHandler + RegisterBrandRoutes()
├── handler_category.go      # CategoryHandler + RegisterCategoryRoutes()
├── model_brand.go           # Brand → sp_brands
├── model_spu.go             # SPU → sp_products
├── model_sku.go             # SKU → sp_skus
├── model_category.go        # Category → sp_categories
├── model_attribute.go       # Attribute → sp_attributes
├── repo_brand.go            # IbrandRepository + BrandRepository
├── repo_spu.go              # IspuRepository + SpuRepository
├── repo_category.go         # IcategoryRepository + CategoryRepository
├── service_brand.go         # BrandService
├── service_spu.go           # SpuService
└── service_category.go      # CategoryService
```

## 2. 文件命名规范

| 前缀 | 含义 | 示例 |
|------|------|------|
| `model_` | GORM 模型（同时也是 PO） | `model_spu.go`, `model_sku.go` |
| `repo_` | Repository 接口 + 实现 | `repo_spu.go`, `repo_category.go` |
| `service_` | 业务逻辑 | `service_spu.go`, `service_brand.go` |
| `dto` | 请求/响应 DTO | `dto_brand.go`, `dto_category.go` |
| `handler` | Handler + 路由注册 | `handler_brand.go`, `handler_category.go` |

### 规则

- 一个文件只包含一个实体的定义（model/repo/service 各一个文件）
- Repository 接口 + 实现在**同一个文件**，不拆分
- `dto_` 和 `handler_` 都按实体拆分（`dto_brand.go`、`handler_category.go`）
- 文件名使用蛇形命名：`model_product_attribute.go`

## 3. 包命名

- 包名 = 模块名，简短小写：`package product`、`package inventory`
- **不建子 package**：没有 `product/model/`、`product/repository/`
- 模块间调用直接 import 完整路径：`import "eshop-monolith/internal/product"`

## 4. 模型定义

### 规范

- 模型同时也是 PO，GORM 标签直接标注字段
- 不经过 `infra/repository/models/` 的 PO 层
- 实现 `TableName()` 返回表名
- 不使用 `utils.Timestamp`，使用标准 `time.Time`
- JSON 序列化时间戳在 DTO 转换时处理
- 主键统一 `int64 autoIncrement`
- 软删除使用 `gorm.DeletedAt`
- 金额以「分」为单位（int64）

### 示例

```go
package product

import (
    "time"
    "gorm.io/gorm"
)

type Brand struct {
    ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string         `gorm:"type:varchar(100);not null;uniqueIndex:uk_name" json:"name"`
    EnglishName string         `gorm:"type:varchar(100);default:''" json:"english_name"`
    LogoURL     string         `gorm:"type:varchar(512);default:''" json:"logo_url"`
    FirstLetter string         `gorm:"type:char(1);default:'';index:idx_first_letter" json:"first_letter"`
    SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
    Status      int8           `gorm:"not null;default:1;index:idx_status" json:"status"`
    Description string         `gorm:"type:text" json:"description"`
    CreatedAt   time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (Brand) TableName() string { return "sp_brands" }
```

## 5. DTO 命名

### 规范

| 用途 | 命名格式 | 示例 |
|------|---------|------|
| 创建请求 | `Create{Entity}Req` | `CreateBrandReq` |
| 更新请求 | `Update{Entity}Req` | `UpdateBrandReq` |
| 列表查询 | `{Entity}ListReq` | `BrandListReq` |
| 列表结果 | `{Entity}ListResult` | `BrandListResult` |

- 每个实体独立的 DTO 文件，如 `dto_brand.go`、`dto_category.go`
- JSON 标签使用蛇形命名：`json:"first_letter"`
- 表单参数使用 `form` 标签：`form:"first_letter"`
- 校验使用 `binding` 标签：`binding:"required,max=100"`
- 可选更新的字段用指针类型：`Name *string`（区分传空和不传）

### 示例

```go
package product

import "eshop-monolith/pkg/query"

type CreateBrandReq struct {
    Name        string `json:"name" binding:"required,max=100"`
    EnglishName string `json:"english_name" binding:"max=100"`
    LogoURL     string `json:"logo_url" binding:"max=512"`
    FirstLetter string `json:"first_letter" binding:"omitempty,len=1,alpha"`
    SortOrder   *int   `json:"sort_order"`
    Status      *int8  `json:"status" binding:"omitempty,oneof=0 1"`
    Description string `json:"description"`
}

type UpdateBrandReq struct {
    Name        *string `json:"name" binding:"omitempty,max=100"`
    EnglishName *string `json:"english_name" binding:"omitempty,max=100"`
    LogoURL     *string `json:"logo_url" binding:"omitempty,max=512"`
    FirstLetter *string `json:"first_letter" binding:"omitempty,len=1,alpha"`
    SortOrder   *int    `json:"sort_order"`
    Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
    Description *string `json:"description"`
}

type BrandListReq struct {
    query.Pagination
    Name        string `form:"name"`
    FirstLetter string `form:"first_letter"`
    Status      *int8  `form:"status"`
}
```

## 6. Repository 规范

### 规范

- 接口命名：`I` + 小写首字母 + 实体名 + `Repository`（`IbrandRepository`）
- 实现结构体：`{Entity}Repository`（`BrandRepository`）
- 构造方法：`New{Entity}Repository(db *gorm.DB) I{entity}Repository`
- 接口 + 实现在**同一文件**
- 方法返回 model 指针：`*Brand`
- 不处理 `gorm.ErrRecordNotFound`，直接返回原始 error
- 分页查询返回 `([]Entity, int64, error)`

### 示例

```go
package product

import (
    "context"
    "gorm.io/gorm"
)

type IbrandRepository interface {
    Create(ctx context.Context, brand *Brand) error
    FindByID(ctx context.Context, id int64) (*Brand, error)
    FindByName(ctx context.Context, name string) (*Brand, error)
    List(ctx context.Context, name, firstLetter string, status *int8, page, size int) ([]Brand, int64, error)
    Update(ctx context.Context, brand *Brand) error
    Delete(ctx context.Context, id int64) error
}

type BrandRepository struct {
    db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) IbrandRepository {
    return &BrandRepository{db: db}
}

func (r *BrandRepository) Create(ctx context.Context, brand *Brand) error {
    return r.db.WithContext(ctx).Create(brand).Error
}
```

### 事务方法

涉及多表写入时，使用事务方法，接收 `*gorm.DB` 参数：

```go
type IspuRepository interface {
    CreateWithTx(tx *gorm.DB, spu *SPU) error
    CreateSkuWithTx(tx *gorm.DB, sku *SKU) error
}
```

事务在 Service 层用 `gorm.DB.Transaction()` 编排。

## 7. Service 规范

### 规范

- 方法签名：`(ctx context.Context, req *XxxReq) (*Model, error)`
- 统一定义列表返回类型在 service 文件中：

```go
type BrandListResult struct {
    Total int64    `json:"total"`
    List  []*Brand `json:"list"`
}
```

- 业务校验返回 `errcode.BizError`
- 名称唯一校验在 Create/Update 中处理
- 不需要的 error 类型用 `errors.Is(err, gorm.ErrRecordNotFound)` 判断

## 8. Handler 规范

### 规范

- `dto_` 和 `handler_` 都按实体拆分
- Handler 结构体组合对应 Service
- 路由注册写在 handler 文件末尾，函数名 `Register{Module}Routes`
- 路由注册直接创建 repo → service → handler，不经过 `infra/repository/db.go`
- 不通过 `Repositories` 聚合注册新模块的 repo
- 旧模块的 `infra/repository/db.go` 和 `infra/router/router.go` 只做桥接

### 固定模式

```go
// Create 创建品牌
// @Summary 创建品牌
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body CreateBrandReq true "品牌信息"
// @Success 200 {object} response.Response{data=Brand}
// @Router /api/v1/brands [post]
func (h *BrandHandler) Create(c *gin.Context) {
    var req CreateBrandReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err)
        return
    }
    result, err := h.svc.Create(c, &req)
    if err != nil {
        c.Error(err)
        return
    }
    response.Success(c, result)
}
```

### 路由注册（写在 `handler_brand.go` 末尾）

```go
// ── Routes ────────────────────────────────────────

func RegisterBrandRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
    repo := NewBrandRepository(db)
    svc := NewBrandService(repo)
    h := NewBrandHandler(svc)

    brands := v1.Group("/brands")
    {
        brands.GET("", h.List)
        brands.GET("/:id", h.GetByID)
    }
    auth := v1.Group("/brands")
    auth.Use(middleware.JWTAuth())
    {
        auth.POST("", h.Create)
        auth.PUT("/:id", h.Update)
        auth.DELETE("/:id", h.Delete)
    }
}
```

## 9. 导入顺序

标准库 → 第三方库 → 内部包，三组间空行分隔：

```go
import (
    "context"
    "errors"

    "gorm.io/gorm"

    "eshop-monolith/pkg/errcode"
)
```

## 10. 错误处理

- Handler 层使用 `c.Error(err)` 传递错误，不直接返回 JSON
- `ErrorHandler` 中间件自动分类：
  - `*errcode.BizError` → 业务错误，返回对应 code + message
  - `validator.ValidationErrors` → 422 + 字段级错误详情
  - `gorm.ErrRecordNotFound` → 自动转为 404
  - 其他 → 500 系统错误
- Repository 层不处理 `gorm.ErrRecordNotFound`
- Service 层需要区分 NotFound 时用 `errors.Is(err, gorm.ErrRecordNotFound)` 判断

## 11. AutoMigrate

新模块的表已通过 SQL 手动创建，不在 `AutoMigrate` 中注册。后续 schema 变更通过 DDL 管理，不依赖 GORM AutoMigrate。

## 12. API 规范

- 路由结构化 RESTful，URL 含版本号 `/api/v1/`
- 响应格式统一 `response.Success(c, data)`
- Swagger @Tags 使用英文小写复数（`brands`、`categories`、`products`）
- 时间戳在 JSON 中序列化为毫秒级（在 DTO 响应时手动转换 `CreatedAt.UnixMilli()`）
- Path 参数使用 `utils.ParseIntParam(c, "id")`

## 检查清单

### 项目结构
- [ ] 模块内文件是否全部平铺，无子目录
- [ ] 是否使用 `{type}_{entity}.go` 前缀命名
- [ ] `dto_` 和 `handler_` 是否都按实体拆分

### 模型
- [ ] 模型是否同时是 PO，无 `ToDomain()`/`FromDomain()`
- [ ] 是否使用 `time.Time` 而非 `utils.Timestamp`
- [ ] 表名是否为 `sp_` 前缀
- [ ] 主键是否为 `int64 autoIncrement`
- [ ] 软删除是否使用 `gorm.DeletedAt`
- [ ] 金额是否以「分」为单位（int64）

### Repository
- [ ] 接口命名是否 `I` + 小写首字母（`IbrandRepository`）
- [ ] 接口和实现是否在同一文件
- [ ] 是否不处理 `gorm.ErrRecordNotFound`

### Service
- [ ] 方法签名是否 `(ctx context.Context, req *XxxReq)`
- [ ] 可选更新字段是否使用指针类型
- [ ] 列表返回类型是否定义在 service 文件中

### Handler
- [ ] handler 和路由注册是否在同一文件
- [ ] 是否使用 `c.Error(err)` 传递错误
- [ ] DTO binding 标签是否完整
- [ ] @Tags 是否英文小写复数

### 路由
- [ ] 路由注册函数是否直接创建 repo→service→handler
- [ ] URL 是否包含 `/api/v1/` 版本号
