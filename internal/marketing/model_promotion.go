package marketing

import (
	"time"

	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Promotion struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PromoName     string         `gorm:"type:varchar(100);not null" json:"promo_name"`
	PromoType     int8           `gorm:"not null" json:"promo_type"`
	PromoCode     string         `gorm:"type:varchar(50);uniqueIndex:idx_code" json:"promo_code"`
	StartTime     time.Time      `gorm:"type:datetime;not null;index:idx_time" json:"start_time"`
	EndTime       time.Time      `gorm:"type:datetime;not null;index:idx_time" json:"end_time"`
	TotalQuantity int            `gorm:"not null;default:0" json:"total_quantity"`
	PerUserLimit  int            `gorm:"not null;default:0" json:"per_user_limit"`
	UsedQuantity  int            `gorm:"not null;default:0" json:"used_quantity"`
	RuleID        *int64         `gorm:"index" json:"rule_id"`
	Status        int8           `gorm:"not null;default:1;index:idx_status" json:"status"`
	CreatedBy     *int64         `json:"created_by"`
	UpdatedBy     *int64         `json:"updated_by"`
	CreatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Promotion) TableName() string { return "mkt_promotions" }

// ── PromotionRule ────────────────────────────────

type PromotionRule struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PromotionID   int64          `gorm:"not null;index:idx_promotion" json:"promotion_id"`
	RuleName      string         `gorm:"type:varchar(100);default:''" json:"rule_name"`
	ConditionType int8           `gorm:"not null" json:"condition_type"`
	ConditionValue float64       `gorm:"type:decimal(10,2)" json:"condition_value"`
	BenefitType   int8           `gorm:"not null" json:"benefit_type"`
	BenefitValue  float64        `gorm:"type:decimal(10,2)" json:"benefit_value"`
	IsStackable   int8           `gorm:"not null;default:0;index:idx_stack" json:"is_stackable"`
	StackPriority int            `gorm:"not null;default:0;index:idx_stack" json:"stack_priority"`
	CreatedBy     *int64         `json:"created_by"`
	UpdatedBy     *int64         `json:"updated_by"`
	CreatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PromotionRule) TableName() string { return "mkt_promotion_rules" }

// ── PromotionProduct ─────────────────────────────

type PromotionProduct struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PromotionID   int64          `gorm:"not null;index:idx_promotion;uniqueIndex:uk_promo_product" json:"promotion_id"`
	ProductType   int8           `gorm:"not null;uniqueIndex:uk_promo_product" json:"product_type"`
	ProductID     *int64         `gorm:"uniqueIndex:uk_promo_product" json:"product_id"`
	CategoryID    *int64         `gorm:"uniqueIndex:uk_promo_product" json:"category_id"`
	CreatedBy     *int64         `json:"created_by"`
	CreatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PromotionProduct) TableName() string { return "mkt_promotion_products" }
