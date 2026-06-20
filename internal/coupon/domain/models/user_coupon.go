package models

import (
	"time"

	"eshop-monolith/pkg/utils"
)

// UserCouponStatus 用户优惠券状态
type UserCouponStatus string

const (
	UserCouponStatusUnused UserCouponStatus = "unused"  // 未使用
	UserCouponStatusUsed   UserCouponStatus = "used"    // 已使用
	UserCouponStatusExpired UserCouponStatus = "expired" // 已过期
)

// UserCoupon 用户领取的优惠券
type UserCoupon struct {
	ID         int64            `json:"id"`
	UserID     int64            `json:"user_id"`     // 用户ID
	CouponID   int64            `json:"coupon_id"`   // 优惠券模板ID
	OrderNo    string           `json:"order_no"`    // 使用的订单号
	CouponCode string           `json:"coupon_code"` // 优惠券码
	Status     UserCouponStatus `json:"status"`      // 状态
	UsedAt     *time.Time       `json:"used_at"`     // 使用时间
	ExpireAt   time.Time        `json:"expire_at"`   // 过期时间

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

// TableName 表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}
