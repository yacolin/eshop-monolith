package trade

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IpaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	FindByPaymentNo(ctx context.Context, paymentNo string) (*Payment, error)
	FindByOrderNo(ctx context.Context, orderNo string) (*Payment, error)
	UpdateStatus(ctx context.Context, id int64, status, transactionID, failureReason string, paidAt *time.Time) error
	CreateLog(ctx context.Context, log *PaymentLog) error

	CreateRefund(ctx context.Context, refund *Refund) error
	FindRefundByNo(ctx context.Context, refundNo string) (*Refund, error)
	FindRefundsByPaymentNo(ctx context.Context, paymentNo string) ([]Refund, error)
	UpdateRefundStatus(ctx context.Context, id int64, status, channelTransactionID, failureReason string) error
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IpaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) FindByPaymentNo(ctx context.Context, paymentNo string) (*Payment, error) {
	var p Payment
	err := r.db.WithContext(ctx).Where("payment_no = ?", paymentNo).First(&p).Error
	return &p, err
}

func (r *PaymentRepository) FindByOrderNo(ctx context.Context, orderNo string) (*Payment, error) {
	var p Payment
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&p).Error
	return &p, err
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id int64, status, transactionID, failureReason string, paidAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&Payment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         status,
		"transaction_id": transactionID,
		"failure_reason": failureReason,
		"paid_at":        paidAt,
	}).Error
}

func (r *PaymentRepository) CreateLog(ctx context.Context, log *PaymentLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *PaymentRepository) CreateRefund(ctx context.Context, refund *Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

func (r *PaymentRepository) FindRefundByNo(ctx context.Context, refundNo string) (*Refund, error) {
	var rf Refund
	err := r.db.WithContext(ctx).Where("refund_no = ?", refundNo).First(&rf).Error
	return &rf, err
}

func (r *PaymentRepository) FindRefundsByPaymentNo(ctx context.Context, paymentNo string) ([]Refund, error) {
	var list []Refund
	err := r.db.WithContext(ctx).Where("payment_no = ?", paymentNo).Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *PaymentRepository) UpdateRefundStatus(ctx context.Context, id int64, status, channelTransactionID, failureReason string) error {
	return r.db.WithContext(ctx).Model(&Refund{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":                 status,
		"channel_transaction_id": channelTransactionID,
		"failure_reason":         failureReason,
	}).Error
}
