package models

import (
	"time"

	payDomain "eshop-monolith/internal/payment/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// PaymentMethodPO 支付方式持久化对象
type PaymentMethodPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex"`
	Name        string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	Config      string         `gorm:"type:json"`
	Status      int            `gorm:"type:tinyint;default:1"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (PaymentMethodPO) TableName() string { return "payment_methods" }

func (po *PaymentMethodPO) ToDomain() *payDomain.PaymentMethod {
	return &payDomain.PaymentMethod{
		ID:          po.ID,
		Code:        po.Code,
		Name:        po.Name,
		Description: po.Description,
		Config:      po.Config,
		Status:      po.Status,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func PaymentMethodFromDomain(p *payDomain.PaymentMethod) *PaymentMethodPO {
	return &PaymentMethodPO{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		Config:      p.Config,
		Status:      p.Status,
		CreatedAt:   time.Time(p.CreatedAt),
		UpdatedAt:   time.Time(p.UpdatedAt),
	}
}
