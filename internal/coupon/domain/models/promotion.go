package models

import "eshop-monolith/pkg/utils"

// PromotionType 促销活动类型
type PromotionType string

const (
	PromotionTypeTimeDiscount PromotionType = "time_discount" // 限时折扣
	PromotionTypeFullReduce   PromotionType = "full_reduce"   // 满减活动
)

// PromotionStatus 促销活动状态
type PromotionStatus string

const (
	PromotionStatusPending   PromotionStatus = "pending"   // 待开始
	PromotionStatusActive    PromotionStatus = "active"    // 进行中
	PromotionStatusFinished  PromotionStatus = "finished"  // 已结束
	PromotionStatusCancelled PromotionStatus = "cancelled" // 已取消
)

// Promotion 促销活动
type Promotion struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`        // 活动名称
	Description string           `json:"description"` // 活动描述
	PromoType   PromotionType    `json:"promo_type"`  // 活动类型
	Scope       string           `json:"scope"`       // 适用范围（all/category/product）
	ScopeValue  string           `json:"scope_value"` // 范围值
	Rule        string           `json:"rule"`        // 规则JSON（折扣率/满减条件）
	StartTime   utils.Timestamp  `json:"start_time"`  // 开始时间
	EndTime     utils.Timestamp  `json:"end_time"`    // 结束时间
	Status      PromotionStatus  `json:"status"`      // 状态
	SortOrder   int              `json:"sort_order"`  // 排序权重

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 表名
func (Promotion) TableName() string {
	return "promotions"
}

// PromotionProduct 促销关联商品
type PromotionProduct struct {
	ID          int64           `json:"id"`
	PromotionID int64           `json:"promotion_id"`
	ProductID   int64           `json:"product_id"`
	Discount    int64           `json:"discount"`     // 折扣价或折扣率（分/百分比）
	CreatedAt   utils.Timestamp `json:"created_at"`
}

// TableName 表名
func (PromotionProduct) TableName() string {
	return "promotion_products"
}
