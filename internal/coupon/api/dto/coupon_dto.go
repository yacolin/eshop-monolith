package dto

import (
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// CouponResponse 优惠券模板响应
type CouponResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CouponType  string           `json:"coupon_type"`
	Scope       string           `json:"scope"`
	ScopeValue  string           `json:"scope_value,omitempty"`
	Value       int64            `json:"value"`
	MinAmount   int64            `json:"min_amount"`
	MaxDiscount int64            `json:"max_discount"`
	TotalStock  int              `json:"total_stock"`
	RemainStock int              `json:"remain_stock"`
	UserLimit   int              `json:"user_limit"`
	StartTime   utils.Timestamp  `json:"start_time"`
	EndTime     utils.Timestamp  `json:"end_time"`
	ValidDays   int              `json:"valid_days"`
	Status      string           `json:"status"`
	CreatedAt   utils.Timestamp  `json:"created_at"`
	UpdatedAt   utils.Timestamp  `json:"updated_at"`
}

// CouponListResult 优惠券模板列表
type CouponListResult struct {
	Total int64            `json:"total"`
	List  []CouponResponse `json:"list"`
}

// CreateCouponReq 创建优惠券请求
type CreateCouponReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CouponType  string `json:"coupon_type" binding:"required,oneof=fixed percentage voucher"`
	Scope       string `json:"scope" binding:"required,oneof=global category product"`
	ScopeValue  string `json:"scope_value"`
	Value       int64  `json:"value" binding:"required,min=1"`
	MinAmount   int64  `json:"min_amount"`
	MaxDiscount int64  `json:"max_discount"`
	TotalStock  int    `json:"total_stock" binding:"required,min=1"`
	UserLimit   int    `json:"user_limit"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	ValidDays   int    `json:"valid_days"`
}

// UpdateCouponReq 更新优惠券请求
type UpdateCouponReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope" binding:"oneof=global category product"`
	ScopeValue  string `json:"scope_value"`
	Value       int64  `json:"value" binding:"min=0"`
	MinAmount   int64  `json:"min_amount"`
	MaxDiscount int64  `json:"max_discount"`
	UserLimit   int    `json:"user_limit"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	ValidDays   int    `json:"valid_days"`
	Status      string `json:"status" binding:"oneof=active inactive"`
}

// ClaimCouponReq 领取优惠券请求
type ClaimCouponReq struct {
	CouponID int64 `json:"coupon_id" binding:"required"`
}

// UseCouponReq 使用优惠券请求
type UseCouponReq struct {
	UserCouponID int64  `json:"user_coupon_id" binding:"required"`
	OrderNo      string `json:"order_no" binding:"required"`
	OrderAmount  int64  `json:"order_amount" binding:"required,min=1"`
}

// UserCouponResponse 用户优惠券响应
type UserCouponResponse struct {
	ID         int64            `json:"id"`
	UserID     int64            `json:"user_id"`
	CouponID   int64            `json:"coupon_id"`
	CouponCode string           `json:"coupon_code"`
	OrderNo    string           `json:"order_no,omitempty"`
	Status     string           `json:"status"`
	ExpireAt   utils.Timestamp  `json:"expire_at"`
	UsedAt     *utils.Timestamp `json:"used_at,omitempty"`
	CreatedAt  utils.Timestamp  `json:"created_at"`

	// 冗余展示优惠券信息
	CouponName        string `json:"coupon_name,omitempty"`
	CouponType        string `json:"coupon_type,omitempty"`
	CouponValue       int64  `json:"coupon_value,omitempty"`
	CouponMinAmount   int64  `json:"coupon_min_amount,omitempty"`
	CouponDescription string `json:"coupon_description,omitempty"`
}

// UserCouponListResult 用户优惠券列表
type UserCouponListResult struct {
	Total int64                `json:"total"`
	List  []UserCouponResponse `json:"list"`
}

// CouponListQuery 优惠券查询参数
type CouponListQuery struct {
	query.Pagination
	Status string `form:"status"`
}
