package trade

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentNo     string         `gorm:"type:varchar(32);not null;uniqueIndex:uk_payment_no" json:"payment_no"`
	OrderNo       string         `gorm:"type:varchar(32);not null;index:idx_order_no" json:"order_no"`
	OrderID       int64          `gorm:"not null;index:idx_order_id" json:"order_id"`
	OrderType     string         `gorm:"type:varchar(20);not null;default:'order'" json:"order_type"`
	Amount        int64          `gorm:"not null" json:"amount"`
	Currency      string         `gorm:"type:varchar(10);not null;default:'CNY'" json:"currency"`
	PaymentMethod string         `gorm:"type:varchar(32);not null" json:"payment_method"`
	Channel       string         `gorm:"type:varchar(32);default:''" json:"channel"`
	TransactionID string         `gorm:"type:varchar(128);default:'';index:idx_transaction_id" json:"transaction_id"`
	Status        string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	FailureReason string         `gorm:"type:varchar(500);default:''" json:"failure_reason"`
	PaidAt        *time.Time     `gorm:"type:datetime" json:"paid_at"`
	CreatedAt     time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (Payment) TableName() string { return "tx_payments" }

// ── Refund ───────────────────────────────────────

type Refund struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	RefundNo            string         `gorm:"type:varchar(32);not null;uniqueIndex:uk_refund_no" json:"refund_no"`
	PaymentNo           string         `gorm:"type:varchar(32);not null;index:idx_payment_no" json:"payment_no"`
	OrderNo             string         `gorm:"type:varchar(32);not null;index:idx_order_no" json:"order_no"`
	OrderID             int64          `gorm:"not null;index:idx_order_id" json:"order_id"`
	Amount              int64          `gorm:"not null" json:"amount"`
	Reason              string         `gorm:"type:varchar(500);default:''" json:"reason"`
	Status              string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ChannelTransactionID string        `gorm:"type:varchar(128);default:''" json:"channel_transaction_id"`
	FailureReason       string         `gorm:"type:varchar(500);default:''" json:"failure_reason"`
	AppliedAt           *time.Time     `gorm:"type:datetime" json:"applied_at"`
	SuccessAt           *time.Time     `gorm:"type:datetime" json:"success_at"`
	CreatedAt           time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (Refund) TableName() string { return "tx_refunds" }
