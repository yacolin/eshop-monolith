package models

import (
	"time"

	domain "eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// UserCouponPO 用户优惠券持久化对象
type UserCouponPO struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	UserID     int64          `gorm:"type:bigint;not null;index;comment:用户ID"`
	CouponID   int64          `gorm:"type:bigint;not null;index;comment:优惠券模板ID"`
	OrderNo    string         `gorm:"type:varchar(64);comment:使用的订单号"`
	CouponCode string         `gorm:"type:varchar(32);not null;uniqueIndex;comment:优惠券码"`
	Status     string         `gorm:"type:varchar(20);not null;index;comment:状态(unused/used/expired)"`
	UsedAt     *time.Time     `gorm:"type:timestamp;comment:使用时间"`
	ExpireAt   time.Time      `gorm:"type:timestamp;not null;index;comment:过期时间"`
	CreatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (UserCouponPO) TableName() string { return "user_coupons" }

func (po *UserCouponPO) ToDomain() *domain.UserCoupon {
	return &domain.UserCoupon{
		ID:         po.ID,
		UserID:     po.UserID,
		CouponID:   po.CouponID,
		OrderNo:    po.OrderNo,
		CouponCode: po.CouponCode,
		Status:     domain.UserCouponStatus(po.Status),
		UsedAt:     po.UsedAt,
		ExpireAt:   po.ExpireAt,
		CreatedAt:  utils.Timestamp(po.CreatedAt),
		UpdatedAt:  utils.Timestamp(po.UpdatedAt),
	}
}

func UserCouponFromDomain(uc *domain.UserCoupon) *UserCouponPO {
	return &UserCouponPO{
		ID:         uc.ID,
		UserID:     uc.UserID,
		CouponID:   uc.CouponID,
		OrderNo:    uc.OrderNo,
		CouponCode: uc.CouponCode,
		Status:     string(uc.Status),
		UsedAt:     uc.UsedAt,
		ExpireAt:   uc.ExpireAt,
		CreatedAt:  time.Time(uc.CreatedAt),
		UpdatedAt:  time.Time(uc.UpdatedAt),
	}
}
