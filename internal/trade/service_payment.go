package trade

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ── 外部依赖接口 ─────────────────────────────────

type PaymentService struct {
	repo IpaymentRepository
	db   *gorm.DB
}

func NewPaymentService(repo IpaymentRepository, db *gorm.DB) *PaymentService {
	return &PaymentService{repo: repo, db: db}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentReq) (*Payment, error) {
	payment := &Payment{
		PaymentNo:     fmt.Sprintf("PAY%d", time.Now().UnixNano()),
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
	s.repo.CreateLog(ctx, &PaymentLog{PaymentID: payment.ID, PaymentNo: payment.PaymentNo, Action: "create", Status: "pending"})
	return payment, nil
}

func (s *PaymentService) HandleCallback(ctx context.Context, req *PaymentCallbackReq) (*Payment, error) {
	payment, err := s.repo.FindByPaymentNo(ctx, req.PaymentNo)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	updates := map[string]interface{}{"status": req.Status, "transaction_id": req.TransactionID, "failure_reason": req.FailureReason}
	if req.Status == "success" {
		updates["paid_at"] = &now
	}
	s.db.Model(payment).Updates(updates)
	payment.Status = req.Status
	payment.TransactionID = req.TransactionID
	s.repo.CreateLog(ctx, &PaymentLog{
		PaymentID: payment.ID, PaymentNo: payment.PaymentNo, Channel: req.Channel,
		TransactionID: req.TransactionID, Action: "pay_callback", RequestBody: req.RawBody, Status: req.Status,
	})
	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, paymentNo string) (*Payment, error) {
	return s.repo.FindByPaymentNo(ctx, paymentNo)
}

func (s *PaymentService) CreateRefund(ctx context.Context, req *CreateRefundReq) (*Refund, error) {
	payment, err := s.repo.FindByPaymentNo(ctx, req.PaymentNo)
	if err != nil {
		return nil, err
	}
	refund := &Refund{
		RefundNo:    fmt.Sprintf("REF%d", time.Now().UnixNano()),
		PaymentNo:   payment.PaymentNo,
		OrderNo:     payment.OrderNo,
		OrderID:     payment.OrderID,
		Amount:      req.Amount,
		Reason:      req.Reason,
		Status:      "pending",
		AppliedAt:   timePtr(time.Now()),
	}
	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return nil, err
	}
	return refund, nil
}

func timePtr(t time.Time) *time.Time { return &t }

