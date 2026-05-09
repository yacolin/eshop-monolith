package models

import (
	"time"

	payDomain "eshop-monolith/internal/payment/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// PaymentPO 支付记录持久化对象
type PaymentPO struct {
	ID            int64                  `gorm:"primaryKey;autoIncrement"`
	OrderID       int64                  `gorm:"not null;index"`
	Amount        int64                  `gorm:"not null"`
	Currency      string                 `gorm:"type:varchar(10);not null;default:'CNY'"`
	PaymentMethod string                 `gorm:"type:varchar(50);not null"`
	TransactionID string                 `gorm:"type:varchar(255);index"`
	Status        string                 `gorm:"type:varchar(20);not null;default:'pending'"`
	FailureReason string                 `gorm:"type:text"`
	Metadata      string                 `gorm:"type:json"`
	PaidAt        *time.Time             `gorm:"type:timestamp"`
	CreatedAt     time.Time              `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt     time.Time              `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt     gorm.DeletedAt         `gorm:"index"`
	Order         *PaymentOrderPO        `gorm:"foreignKey:OrderID"`
	Transactions  []PaymentTransactionPO `gorm:"foreignKey:PaymentID"`
	Refunds       []RefundPO             `gorm:"foreignKey:PaymentID"`
}

func (PaymentPO) TableName() string { return "payments" }

func (po *PaymentPO) ToDomain() *payDomain.Payment {
	var paidAt *utils.Timestamp
	if po.PaidAt != nil {
		ts := utils.Timestamp(*po.PaidAt)
		paidAt = &ts
	}
	var payOrder *payDomain.Order
	if po.Order != nil {
		payOrder = po.Order.ToDomain()
	}
	transactions := make([]payDomain.PaymentTransaction, len(po.Transactions))
	for i, t := range po.Transactions {
		transactions[i] = *t.ToDomain()
	}
	refunds := make([]payDomain.Refund, len(po.Refunds))
	for i, r := range po.Refunds {
		refunds[i] = *r.ToDomain()
	}
	return &payDomain.Payment{
		ID:            po.ID,
		OrderID:       po.OrderID,
		Amount:        po.Amount,
		Currency:      po.Currency,
		PaymentMethod: po.PaymentMethod,
		TransactionID: po.TransactionID,
		Status:        po.Status,
		FailureReason: po.FailureReason,
		Metadata:      po.Metadata,
		PaidAt:        paidAt,
		CreatedAt:     utils.Timestamp(po.CreatedAt),
		UpdatedAt:     utils.Timestamp(po.UpdatedAt),
		Order:         payOrder,
		Transactions:  transactions,
		Refunds:       refunds,
	}
}

func PaymentFromDomain(p *payDomain.Payment) *PaymentPO {
	var paidAt *time.Time
	if p.PaidAt != nil {
		t := time.Time(*p.PaidAt)
		paidAt = &t
	}
	var payOrder *PaymentOrderPO
	if p.Order != nil {
		payOrder = PaymentOrderFromDomain(p.Order)
	}
	transactions := make([]PaymentTransactionPO, len(p.Transactions))
	for i, t := range p.Transactions {
		transactions[i] = *PaymentTransactionFromDomain(&t)
	}
	refunds := make([]RefundPO, len(p.Refunds))
	for i, r := range p.Refunds {
		refunds[i] = *RefundFromDomain(&r)
	}
	return &PaymentPO{
		ID:            p.ID,
		OrderID:       p.OrderID,
		Amount:        p.Amount,
		Currency:      p.Currency,
		PaymentMethod: p.PaymentMethod,
		TransactionID: p.TransactionID,
		Status:        p.Status,
		FailureReason: p.FailureReason,
		Metadata:      p.Metadata,
		PaidAt:        paidAt,
		CreatedAt:     time.Time(p.CreatedAt),
		UpdatedAt:     time.Time(p.UpdatedAt),
		Order:         payOrder,
		Transactions:  transactions,
		Refunds:       refunds,
	}
}
