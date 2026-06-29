package trade

import (
	"context"
)

// ── 外部依赖接口 ─────────────────────────────────

type SkuInfo interface {
	GetID() int64
	GetProductID() int64
	GetSkuCode() string
	GetPrice() int64
	GetImage() string
	GetSpecJSON() string
}

type SkuProvider interface {
	FindByID(ctx context.Context, skuID int64) (SkuInfo, error)
}

type InventoryService interface {
	Lock(ctx context.Context, skuID int64, quantity int) error
	Unlock(ctx context.Context, skuID int64, quantity int) error
}

// ── CartService ──────────────────────────────────

