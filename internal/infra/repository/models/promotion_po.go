package models

import (
	"time"

	domain "eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// PromotionPO 促销活动持久化对象
type PromotionPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(255);not null;comment:活动名称"`
	Description string         `gorm:"type:text;comment:活动描述"`
	PromoType   string         `gorm:"type:varchar(20);not null;index;comment:活动类型(time_discount/full_reduce)"`
	Scope       string         `gorm:"type:varchar(20);not null;default:all;comment:适用范围(all/category/product)"`
	ScopeValue  string         `gorm:"type:varchar(500);comment:范围值"`
	Rule        string         `gorm:"type:text;not null;comment:规则JSON"`
	StartTime   time.Time      `gorm:"type:timestamp;not null;index;comment:开始时间"`
	EndTime     time.Time      `gorm:"type:timestamp;not null;index;comment:结束时间"`
	Status      string         `gorm:"type:varchar(20);not null;index;comment:状态(pending/active/finished/cancelled)"`
	SortOrder   int            `gorm:"type:int;not null;default:0;comment:排序权重"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (PromotionPO) TableName() string { return "promotions" }

func (po *PromotionPO) ToDomain() *domain.Promotion {
	return &domain.Promotion{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		PromoType:   domain.PromotionType(po.PromoType),
		Scope:       po.Scope,
		ScopeValue:  po.ScopeValue,
		Rule:        po.Rule,
		StartTime:   utils.Timestamp(po.StartTime),
		EndTime:     utils.Timestamp(po.EndTime),
		Status:      domain.PromotionStatus(po.Status),
		SortOrder:   po.SortOrder,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func PromotionFromDomain(p *domain.Promotion) *PromotionPO {
	return &PromotionPO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PromoType:   string(p.PromoType),
		Scope:       p.Scope,
		ScopeValue:  p.ScopeValue,
		Rule:        p.Rule,
		StartTime:   time.Time(p.StartTime),
		EndTime:     time.Time(p.EndTime),
		Status:      string(p.Status),
		SortOrder:   p.SortOrder,
		CreatedAt:   time.Time(p.CreatedAt),
		UpdatedAt:   time.Time(p.UpdatedAt),
	}
}

// PromotionProductPO 促销商品关联持久化对象
type PromotionProductPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	PromotionID int64          `gorm:"type:bigint;not null;index;comment:促销活动ID"`
	ProductID   int64          `gorm:"type:bigint;not null;index;comment:商品ID"`
	Discount    int64          `gorm:"type:bigint;not null;comment:折扣价或折扣率"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (PromotionProductPO) TableName() string { return "promotion_products" }

func (po *PromotionProductPO) ToDomain() *domain.PromotionProduct {
	return &domain.PromotionProduct{
		ID:          po.ID,
		PromotionID: po.PromotionID,
		ProductID:   po.ProductID,
		Discount:    po.Discount,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
	}
}

func PromotionProductFromDomain(pp *domain.PromotionProduct) *PromotionProductPO {
	return &PromotionProductPO{
		ID:          pp.ID,
		PromotionID: pp.PromotionID,
		ProductID:   pp.ProductID,
		Discount:    pp.Discount,
		CreatedAt:   time.Time(pp.CreatedAt),
	}
}
