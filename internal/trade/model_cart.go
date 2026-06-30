package trade

import (
	"time"

	"eshop-monolith/pkg/utils"
)

type Cart struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64          `gorm:"not null;default:0;index:idx_user_id" json:"user_id"`
	SessionID   string         `gorm:"type:varchar(64);default:'';index:idx_session_id" json:"session_id"`
	ItemCount   int            `gorm:"not null;default:0" json:"item_count"`
	TotalAmount int64          `gorm:"not null;default:0" json:"total_amount"`
	ExpiredAt   *time.Time     `gorm:"type:datetime;index:idx_expired_at" json:"expired_at"`
	CreatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Cart) TableName() string { return "tx_carts" }

// ── CartItem ─────────────────────────────────────

type CartItem struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID      int64      `gorm:"not null;index:idx_cart_id" json:"cart_id"`
	SkuID       int64      `gorm:"not null;index:idx_sku_id" json:"sku_id"`
	ProductID   int64      `gorm:"not null;default:0" json:"product_id"`
	ProductName string     `gorm:"type:varchar(200);not null;default:''" json:"product_name"`
	SkuSpec     string     `gorm:"type:json" json:"sku_spec,omitempty"`
	Image       string     `gorm:"type:varchar(512);default:''" json:"image"`
	Price       int64      `gorm:"not null" json:"price"`
	Quantity    int        `gorm:"not null;default:1" json:"quantity"`
	CreatedAt utils.Timestamp  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (CartItem) TableName() string { return "tx_cart_items" }
