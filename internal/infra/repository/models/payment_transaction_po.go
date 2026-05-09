package models

import (
	"time"

	payDomain "eshop-monolith/internal/payment/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// PaymentTransactionPO 支付交易记录持久化对象
type PaymentTransactionPO struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	PaymentID     int64          `gorm:"not null;index"`
	TransactionID string         `gorm:"type:varchar(255);not null;index"`
	Amount        int64          `gorm:"not null"`
	Type          string         `gorm:"type:varchar(20);not null"`
	Status        string         `gorm:"type:varchar(20);not null"`
	ResponseData  string         `gorm:"type:json"`
	ErrorData     string         `gorm:"type:json"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Payment       *PaymentPO     `gorm:"foreignKey:PaymentID"`
}

func (PaymentTransactionPO) TableName() string { return "payment_transactions" }

func (po *PaymentTransactionPO) ToDomain() *payDomain.PaymentTransaction {
	return &payDomain.PaymentTransaction{
		ID:            po.ID,
		PaymentID:     po.PaymentID,
		TransactionID: po.TransactionID,
		Amount:        po.Amount,
		Type:          po.Type,
		Status:        po.Status,
		ResponseData:  po.ResponseData,
		ErrorData:     po.ErrorData,
		CreatedAt:     utils.Timestamp(po.CreatedAt),
		UpdatedAt:     utils.Timestamp(po.UpdatedAt),
	}
}

func PaymentTransactionFromDomain(t *payDomain.PaymentTransaction) *PaymentTransactionPO {
	return &PaymentTransactionPO{
		ID:            t.ID,
		PaymentID:     t.PaymentID,
		TransactionID: t.TransactionID,
		Amount:        t.Amount,
		Type:          t.Type,
		Status:        t.Status,
		ResponseData:  t.ResponseData,
		ErrorData:     t.ErrorData,
		CreatedAt:     time.Time(t.CreatedAt),
		UpdatedAt:     time.Time(t.UpdatedAt),
	}
}
