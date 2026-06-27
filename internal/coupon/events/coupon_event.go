package events

import "time"

// CouponIssuedEvent 优惠券发放事件
type CouponIssuedEvent struct {
	UserCouponID int64     `json:"user_coupon_id"`
	UserID       int64     `json:"user_id"`
	CouponID     int64     `json:"coupon_id"`
	CouponCode   string    `json:"coupon_code"`
	IssuedAt     time.Time `json:"issued_at"`
}

func (e CouponIssuedEvent) RoutingKey() string { return "coupon.issued" }

// CouponUsedEvent 优惠券使用事件
type CouponUsedEvent struct {
	UserCouponID int64     `json:"user_coupon_id"`
	UserID       int64     `json:"user_id"`
	CouponID     int64     `json:"coupon_id"`
	OrderNo      string    `json:"order_no"`
	DiscountAmt  int64     `json:"discount_amt"` // 实际抵扣金额，单位：分
	UsedAt       time.Time `json:"used_at"`
}

func (e CouponUsedEvent) RoutingKey() string { return "coupon.used" }

// CouponExpiredEvent 优惠券过期事件
type CouponExpiredEvent struct {
	UserCouponID int64     `json:"user_coupon_id"`
	UserID       int64     `json:"user_id"`
	CouponID     int64     `json:"coupon_id"`
	CouponCode   string    `json:"coupon_code"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (e CouponExpiredEvent) RoutingKey() string { return "coupon.expired" }
