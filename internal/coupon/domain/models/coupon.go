package models

import "eshop-monolith/pkg/utils"

// CouponType 优惠券类型
type CouponType string

const (
	CouponTypeFixed       CouponType = "fixed"       // 满减券
	CouponTypePercentage  CouponType = "percentage"  // 折扣券
	CouponTypeVoucher     CouponType = "voucher"     // 代金券（无门槛）
)

// CouponScope 优惠券适用范围
type CouponScope string

const (
	CouponScopeGlobal   CouponScope = "global"    // 全场通用
	CouponScopeCategory CouponScope = "category"  // 指定分类
	CouponScopeProduct  CouponScope = "product"   // 指定商品
)

// CouponStatus 优惠券状态
type CouponStatus string

const (
	CouponStatusActive   CouponStatus = "active"
	CouponStatusInactive CouponStatus = "inactive"
	CouponStatusExpired  CouponStatus = "expired"
)

// Coupon 优惠券模板
type Coupon struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`           // 优惠券名称
	Description   string           `json:"description"`    // 优惠券描述
	CouponType    CouponType       `json:"coupon_type"`    // 优惠券类型
	Scope         CouponScope      `json:"scope"`          // 适用范围
	ScopeValue    string           `json:"scope_value"`    // 范围值（分类ID/商品ID，逗号分隔）
	Value         int64            `json:"value"`          // 面值（分 或 百分比*100）
	MinAmount     int64            `json:"min_amount"`     // 最低消费金额（分），0 表示无限制
	MaxDiscount   int64            `json:"max_discount"`   // 最大抵扣金额（分），仅百分比券有效，0 不限制
	TotalStock    int              `json:"total_stock"`    // 发放总量
	RemainStock   int              `json:"remain_stock"`   // 剩余数量
	UserLimit     int              `json:"user_limit"`     // 每人限领数量，0 不限制
	StartTime     utils.Timestamp  `json:"start_time"`     // 领取开始时间
	EndTime       utils.Timestamp  `json:"end_time"`       // 领取结束时间
	ValidDays     int              `json:"valid_days"`     // 领取后有效天数，0 表示固定有效期
	Status        CouponStatus     `json:"status"`         // 状态

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 表名
func (Coupon) TableName() string {
	return "coupons"
}

// IsValid 检查优惠券模板是否有效
func (c *Coupon) IsValid() bool {
	return c.Status == CouponStatusActive && c.RemainStock > 0
}
