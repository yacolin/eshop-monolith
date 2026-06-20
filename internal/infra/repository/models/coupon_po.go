package models

import (
	"time"

	domain "eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// CouponPO 优惠券模板持久化对象
type CouponPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(255);not null;comment:优惠券名称"`
	Description string         `gorm:"type:text;comment:优惠券描述"`
	CouponType  string         `gorm:"type:varchar(20);not null;index;comment:优惠券类型(fixed/percentage/voucher)"`
	Scope       string         `gorm:"type:varchar(20);not null;default:global;comment:适用范围(global/category/product)"`
	ScopeValue  string         `gorm:"type:varchar(500);comment:范围值"`
	Value       int64          `gorm:"type:bigint;not null;comment:面值(分或百分比*100)"`
	MinAmount   int64          `gorm:"type:bigint;not null;default:0;comment:最低消费金额(分)"`
	MaxDiscount int64          `gorm:"type:bigint;not null;default:0;comment:最大抵扣金额(分)"`
	TotalStock  int            `gorm:"type:int;not null;comment:发放总量"`
	RemainStock int            `gorm:"type:int;not null;comment:剩余数量"`
	UserLimit   int            `gorm:"type:int;not null;default:0;comment:每人限领数量"`
	StartTime   time.Time      `gorm:"type:timestamp;not null;index;comment:领取开始时间"`
	EndTime     time.Time      `gorm:"type:timestamp;not null;index;comment:领取结束时间"`
	ValidDays   int            `gorm:"type:int;not null;default:0;comment:领取后有效天数"`
	Status      string         `gorm:"type:varchar(20);not null;index;comment:状态(active/inactive/expired)"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (CouponPO) TableName() string { return "coupons" }

func (po *CouponPO) ToDomain() *domain.Coupon {
	return &domain.Coupon{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		CouponType:  domain.CouponType(po.CouponType),
		Scope:       domain.CouponScope(po.Scope),
		ScopeValue:  po.ScopeValue,
		Value:       po.Value,
		MinAmount:   po.MinAmount,
		MaxDiscount: po.MaxDiscount,
		TotalStock:  po.TotalStock,
		RemainStock: po.RemainStock,
		UserLimit:   po.UserLimit,
		StartTime:   utils.Timestamp(po.StartTime),
		EndTime:     utils.Timestamp(po.EndTime),
		ValidDays:   po.ValidDays,
		Status:      domain.CouponStatus(po.Status),
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func CouponFromDomain(c *domain.Coupon) *CouponPO {
	return &CouponPO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CouponType:  string(c.CouponType),
		Scope:       string(c.Scope),
		ScopeValue:  c.ScopeValue,
		Value:       c.Value,
		MinAmount:   c.MinAmount,
		MaxDiscount: c.MaxDiscount,
		TotalStock:  c.TotalStock,
		RemainStock: c.RemainStock,
		UserLimit:   c.UserLimit,
		StartTime:   time.Time(c.StartTime),
		EndTime:     time.Time(c.EndTime),
		ValidDays:   c.ValidDays,
		Status:      string(c.Status),
		CreatedAt:   time.Time(c.CreatedAt),
		UpdatedAt:   time.Time(c.UpdatedAt),
	}
}
