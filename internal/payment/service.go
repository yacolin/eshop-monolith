package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type PaymentService struct {
	repo IpaymentRepository
	db   *gorm.DB
}

func NewPaymentService(repo IpaymentRepository, db *gorm.DB) *PaymentService {
	return &PaymentService{repo: repo, db: db}
}

// CreatePayment 创建支付单
func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentReq) (*Payment, error) {
	paymentNo := fmt.Sprintf("PAY%d", time.Now().UnixNano())
	payment := &Payment{
		PaymentNo:     paymentNo,
		OrderNo:       req.OrderNo,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Channel:       req.Channel,
		Status:        "pending",
	}
	if req.OrderType != "" {
		payment.OrderType = req.OrderType
	}
	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}
	s.repo.CreateLog(ctx, &PaymentLog{
		PaymentID: payment.ID,
		PaymentNo: payment.PaymentNo,
		Action:    "create",
		Status:    "pending",
	})
	return payment, nil
}

// HandleCallback 处理支付回调
func (s *PaymentService) HandleCallback(ctx context.Context, req *PaymentCallbackReq) (*Payment, error) {
	payment, err := s.repo.FindByPaymentNo(ctx, req.PaymentNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPaymentNotFound
		}
		return nil, err
	}
	now := time.Now()
	status := req.Status
	if status == "success" {
		payment.PaidAt = &now
	} else {
		status = "failed"
	}

	// 使用 db 更新，不通过 repo（避免时间字段问题）
	err = s.db.Model(payment).Updates(map[string]interface{}{
		"status":         status,
		"transaction_id": req.TransactionID,
		"failure_reason": req.FailureReason,
		"paid_at":        payment.PaidAt,
	}).Error
	if err != nil {
		return nil, err
	}
	payment.Status = status
	payment.TransactionID = req.TransactionID

	s.repo.CreateLog(ctx, &PaymentLog{
		PaymentID:     payment.ID,
		PaymentNo:     payment.PaymentNo,
		Channel:       req.Channel,
		TransactionID: req.TransactionID,
		Action:        "pay_callback",
		RequestBody:   req.RawBody,
		Status:        status,
	})
	return payment, nil
}

// GetPayment 查询支付单
func (s *PaymentService) GetPayment(ctx context.Context, paymentNo string) (*Payment, error) {
	payment, err := s.repo.FindByPaymentNo(ctx, paymentNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPaymentNotFound
		}
		return nil, err
	}
	return payment, nil
}

// GetPaymentByOrder 按订单查支付
func (s *PaymentService) GetPaymentByOrder(ctx context.Context, orderNo string) (*Payment, error) {
	payment, err := s.repo.FindByOrderNo(ctx, orderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPaymentNotFound
		}
		return nil, err
	}
	return payment, nil
}

// CreateRefund 创建退款
func (s *PaymentService) CreateRefund(ctx context.Context, req *CreateRefundReq) (*Refund, error) {
	// 查找支付单以获取 payment_no
	payment, err := s.repo.FindByPaymentNo(ctx, req.PaymentNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPaymentNotFound
		}
		return nil, err
	}

	refundNo := fmt.Sprintf("REF%d", time.Now().UnixNano())
	refund := &Refund{
		RefundNo:  refundNo,
		PaymentNo: payment.PaymentNo,
		OrderNo:   payment.OrderNo,
		OrderID:   payment.OrderID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    "pending",
		AppliedAt: timePtr(time.Now()),
	}
	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return nil, err
	}
	return refund, nil
}

func timePtr(t time.Time) *time.Time { return &t }

// HandleRefundCallback 处理退款回调
func (s *PaymentService) HandleRefundCallback(ctx context.Context, req *RefundCallbackReq) (*Refund, error) {
	refund, err := s.repo.FindRefundByNo(ctx, req.RefundNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrRefundNotFound
		}
		return nil, err
	}
	now := time.Now()
	err = s.db.Model(refund).Updates(map[string]interface{}{
		"status":                  req.Status,
		"channel_transaction_id": req.ChannelTransactionID,
		"failure_reason":         req.FailureReason,
		"success_at":             now,
	}).Error
	if err != nil {
		return nil, err
	}
	refund.Status = req.Status
	return refund, nil
}
