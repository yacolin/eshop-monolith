package query

import (
	"eshop-monolith/internal/pkg/config"

	"gorm.io/gorm"
)

// 通用列表结果，使用 Go 泛型以适配不同实体类型
type ListResult[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

// 通用分页
type Pagination struct {
	Page int `form:"page,default=1" binding:"gte=1"`          // 页码，最小 1
	Size int `form:"size,default=10" binding:"gte=1,lte=100"` // 每页条数，范围 1..100
}

// Normalize 将保证 Page/Size 在合理范围内（对外部输入做最后防护）
func (p *Pagination) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	// use configured defaults and max if available
	cfg := config.Get()
	def := 10
	max := 100
	if cfg != nil {
		if cfg.Pagination.DefaultSize > 0 {
			def = cfg.Pagination.DefaultSize
		}
		if cfg.Pagination.MaxSize > 0 {
			max = cfg.Pagination.MaxSize
		}
	}
	if p.Size <= 0 {
		p.Size = def
	}
	if p.Size > max {
		p.Size = max
	}
}

// ApplyOrder 应用排序到 GORM 查询
// sortBy: 排序字段（如 "id", "created_at"）
// order: 排序方向（"asc" 或 "desc"）
// defaultOrder: 默认排序（如 "id asc"）
func ApplyOrder(db *gorm.DB, sortBy, order, defaultOrder string) *gorm.DB {
	ord := order
	if ord != "asc" && ord != "desc" {
		ord = "asc"
	}

	o := defaultOrder
	if sortBy != "" {
		o = sortBy + " " + ord
	}
	return db.Order(o)
}
