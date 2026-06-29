# 商品中心 & 库存中心重构设计

日期：2026-06-29

## 背景

现有项目存在两个核心问题：

1. **目录结构嵌套过深**：`api/handlers/` `domain/models/` `domain/repositories/` 每模块 4-5 层目录，并伴随 PO/Domain 双重模型转换（`ToDomain()`/`FromDomain()`），样板代码多
2. **业务表设计不标准**：旧表缺少 SPU/SKU 分层、库存三字段模式、类目树 path+level 等商业标准设计

## 目标架构

```
internal/
├── product/                    # 商品中心（package product）
│   ├── dto.go                  #   所有 DTO
│   ├── handler.go              #   Handler + 路由注册
│   ├── model_brand.go          #   Brand → sp_brands
│   ├── model_category.go       #   Category → sp_categories
│   ├── model_category_brand.go #   CategoryBrand → sp_category_brands
│   ├── model_spu.go            #   SPU → sp_products
│   ├── model_sku.go            #   SKU → sp_skus
│   ├── model_description.go    #   Description → sp_product_descriptions
│   ├── model_attribute.go      #   Attribute → sp_attributes
│   ├── model_product_attr.go   #   ProductAttribute → sp_product_attributes
│   ├── repo_brand.go           #   Brand 仓储
│   ├── repo_category.go        #   Category 仓储（树操作）
│   ├── repo_spu.go             #   SPU+SKU 聚合仓储
│   ├── repo_attribute.go       #   Attribute 仓储
│   ├── service_brand.go        #   Brand 业务
│   ├── service_category.go     #   Category 业务
│   └── service_spu.go          #   SPU 创建（事务编排 SPU+SKU+Attr）
│
├── inventory/                  # 库存中心（package inventory，新扁平结构）
│   ├── dto.go
│   ├── handler.go
│   ├── model.go                #   Inventory → sp_inventories
│   ├── model_log.go            #   InventoryLog → sp_inventory_logs
│   ├── repository.go
│   └── service.go
│
├── order/              │ 旧模块保持原样，逐步替换
├── cart/               │
├── user/               │
└── ...                 │
```

### 文件命名规则

所有文件使用 `{type}_{entity}.go` 前缀命名，不建子目录：

| 前缀 | 含义 | 示例 |
|------|------|------|
| `model_` | GORM 领域模型（同时也是 PO，无双重模型转换） | `model_spu.go`, `model_sku.go` |
| `repo_` | Repository 接口 + GORM 实现（同一文件） | `repo_spu.go`, `repo_category.go` |
| `service_` | 业务逻辑 | `service_spu.go`, `service_brand.go` |
| `dto` | 所有模块请求/响应 DTO（单文件） | `dto.go` |
| `handler` | Handler + 路由注册（单文件） | `handler.go` |

### 与旧结构对比

| 对比项 | 旧结构 | 新结构 |
|--------|--------|--------|
| 目录深度 | 4-5 层子目录 | 全部平铺 |
| 模型 | `domain/models/` + `infra/repository/models/PO` 两套 | 单模型，GORM 标签直接标注 |
| 转换 | `ToDomain()` / `FromDomain()` | 无，模型即 PO |
| 路由 | `api/routes/` 子目录 | 在 `handler.go` 末尾 |
| 接口命名 | `IxxxRepository` | `IxxxRepository` |
| Repository 注册 | 通过 `infra/repository/db.go` 的 `Repositories` 聚合 | 在 route 中直接创建（不经过聚合） |
| 表名 | `products`, `categories` 等 | `sp_products`, `sp_categories` 等 |

## 重构阶段

### 阶段一：Category + Category-Brand

前置条件：无（基础数据表，不依赖其他服务）

工作内容：
- 创建 `model_category.go` + `model_category_brand.go`
- 创建 `repo_category.go`
- 创建 `service_category.go`（树操作：path 维护、level 校验、子节点查询）
- handler 中注册 CRUD + 树查询路由

### 阶段二：Attribute

前置条件：阶段一完成（Attribute 挂载在 Category 下）

工作内容：
- 创建 `model_attribute.go` + `model_product_attr.go`
- 创建 `repo_attribute.go`
- 创建 `service_attribute.go`（属性字典管理、类目属性关联）
- 注册路由

### 阶段三：SPU + SKU + Description + ProductAttribute

前置条件：阶段一、二完成（创建 SPU 时需要选择 Category、Attribute）

工作内容：
- 创建 `model_spu.go`、`model_sku.go`、`model_description.go`
- 创建 `repo_spu.go`（事务内同时写入 SPU + SKU + Description + ProductAttribute）
- 创建 `service_spu.go`（SPU 创建编排、SKU 价格聚合、状态管理）
- 注册路由

### 库存中心（后续独立阶段）

前置条件：Product Center 稳定后

工作内容：
- 扁平化重写 `internal/inventory/`
- 库存锁定/释放、补货、流水记录

## 数据迁移策略

| 阶段 | 迁移方案 |
|------|---------|
| Brand | 无旧数据，直接使用新表 |
| Category | 旧 `categories` 表有数据，新 `sp_categories` 需要导入或用双写过渡 |
| Attribute | 旧 `attributes` 表有数据，需迁移 |
| SPU/SKU | 旧 `products`/`skus` 表有数据，需迁移 |

## Brand 参考实现

Brand 已作为试点完成，文件清单：

- `internal/product/model_brand.go` — GORM 模型，`TableName()` → `sp_brands`
- `internal/product/repo_brand.go` — 接口 + 实现，基础 CRUD + List
- `internal/product/service_brand.go` — 名称唯一校验，CRUD 编排
- `internal/product/dto.go` — `CreateBrandReq`, `UpdateBrandReq`, `BrandListReq`
- `internal/product/handler.go` — Handler + Swagger + `RegisterBrandRoutes()`
