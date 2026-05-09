package models

import (
	"time"

	payDomain "eshop-monolith/internal/payment/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// RefundPO 退款记录持久化对象
type RefundPO struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	PaymentID     int64          `gorm:"not null;index"`
	OrderID       int64          `gorm:"not null;index"`
	RefundAmount  int64          `gorm:"not null"`
	RefundReason  string         `gorm:"type:text"`
	TransactionID string         `gorm:"type:varchar(255);index"`
	Status        string         `gorm:"type:varchar(20);not null;default:'pending'"`
	FailureReason string         `gorm:"type:text"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Payment       *PaymentPO     `gorm:"foreignKey:PaymentID"`
	Order         *PaymentOrderPO `gorm:"foreignKey:OrderID"`
}

func (RefundPO) TableName() string { return "refunds" }

func (po *RefundPO) ToDomain() *payDomain.Refund {
	var payment *payDomain.Payment
	if po.Payment != nil {
		payment = po.Payment.ToDomain()
	}
	var payOrder *payDomain.Order
	if po.Order != nil {
		payOrder = po.Order.ToDomain()
	}
	return &payDomain.Refund{
		ID:            po.ID,
		PaymentID:     po.PaymentID,
		OrderID:       po.OrderID,
		RefundAmount:  po.RefundAmount,
		RefundReason:  po.RefundReason,
		TransactionID: po.TransactionID,
		Status:        po.Status,
		FailureReason: po.FailureReason,
		CreatedAt:     utils.Timestamp(po.CreatedAt),
		UpdatedAt:     utils.Timestamp(po.UpdatedAt),
		Payment:       payment,
		Order:         payOrder,
	}
}

func RefundFromDomain(r *payDomain.Refund) *RefundPO {
	return &RefundPO{
		ID:            r.ID,
		PaymentID:     r.PaymentID,
		OrderID:       r.OrderID,
		RefundAmount:  r.RefundAmount,
		RefundReason:  r.RefundReason,
		TransactionID: r.TransactionID,
		Status:        r.Status,
		FailureReason: r.FailureReason,
		CreatedAt:     time.Time(r.CreatedAt),
		UpdatedAt:     time.Time(r.UpdatedAt),
	}
}
