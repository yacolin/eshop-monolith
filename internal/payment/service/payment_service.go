package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/order/domain/repositories"
	"eshop-monolith/internal/payment/api/dto"
	"eshop-monolith/internal/payment/domain/models"
	paymentRepos "eshop-monolith/internal/payment/domain/repositories"
	"eshop-monolith/internal/payment/events"
	"eshop-monolith/internal/pkg/logger"
	"eshop-monolith/internal/pkg/utils"
)

// PaymentService 支付服务
type PaymentService struct {
	paymentRepo       paymentRepos.IPaymentRepository
	refundRepo        paymentRepos.IRefundRepository
	paymentMethodRepo paymentRepos.IPaymentMethodRepository
	orderRepo         repositories.IorderRepository
	bus               *eventbus.Bus
}

// NewPaymentService 创建支付服务实例
func NewPaymentService(
	paymentRepo paymentRepos.IPaymentRepository,
	refundRepo paymentRepos.IRefundRepository,
	paymentMethodRepo paymentRepos.IPaymentMethodRepository,
	orderRepo repositories.IorderRepository,
	bus *eventbus.Bus,
) *PaymentService {
	return &PaymentService{
		paymentRepo:       paymentRepo,
		refundRepo:        refundRepo,
		paymentMethodRepo: paymentMethodRepo,
		orderRepo:         orderRepo,
		bus:               bus,
	}
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*models.Payment, string, error) {
	// 验证支付方式是否存在
	_, err := s.paymentMethodRepo.GetByCode(ctx, req.PaymentMethod)
	if err != nil {
		return nil, "", fmt.Errorf("payment method not found: %w", err)
	}

	// 验证订单是否存在
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, "", fmt.Errorf("order not found: %w", err)
	}

	// 验证订单状态
	if order.Status != "pending" {
		return nil, "", errors.New("order status is not pending")
	}

	// 验证金额
	if req.Amount != order.TotalAmount {
		return nil, "", errors.New("payment amount does not match order amount")
	}

	// 创建支付记录
	payment := &models.Payment{
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Status:        "pending",
		Metadata:      req.Metadata,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, "", fmt.Errorf("create payment failed: %w", err)
	}

	// 发布支付创建事件
	s.bus.Publish(events.PaymentCreatedEvent{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
	})

	// 根据支付方式生成支付链接或二维码
	var paymentURL string
	switch req.PaymentMethod {
	case "alipay":
		// 生成支付宝支付链接
		paymentURL = fmt.Sprintf("https://openapi.alipay.com/gateway.do?order_id=%d&amount=%d", payment.ID, payment.Amount)
	case "wechat":
		// 生成微信支付二维码
		paymentURL = fmt.Sprintf("wechat://pay?order_id=%d&amount=%d", payment.ID, payment.Amount)
	case "bank":
		// 生成银行转账信息
		paymentURL = fmt.Sprintf("bank://transfer?order_id=%d&amount=%d", payment.ID, payment.Amount)
	case "cash":
		// 货到付款，不需要支付链接
		paymentURL = ""
	}

	return payment, paymentURL, nil
}

// GetPaymentByID 根据ID获取支付记录
func (s *PaymentService) GetPaymentByID(ctx context.Context, id int64) (*models.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

// GetPaymentByOrderID 根据订单ID获取支付记录
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID int64) (*models.Payment, error) {
	return s.paymentRepo.GetByOrderID(ctx, orderID)
}

// UpdatePaymentStatus 更新支付状态
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id int64, status string, transactionID string, failureReason string) error {
	// 获取支付记录
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	previousStatus := payment.Status

	// 更新支付状态
	payment.Status = status
	payment.TransactionID = transactionID
	payment.FailureReason = failureReason

	if status == "success" {
		now := utils.Timestamp(time.Now())
		payment.PaidAt = &now
	}

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("update payment status failed: %w", err)
	}

	// 发布支付状态更新事件
	s.bus.Publish(events.PaymentStatusUpdatedEvent{
		PaymentID:      payment.ID,
		OrderID:        payment.OrderID,
		Status:         payment.Status,
		PreviousStatus: previousStatus,
		TransactionID:  transactionID,
	})

	// 如果支付成功，更新订单状态
	if status == "success" {
		if err := s.orderRepo.UpdateStatus(ctx, payment.OrderID, "paid"); err != nil {
			logger.Error("update order status failed", "order_id", payment.OrderID, "error", err)
		}
	}

	// 如果支付失败，发布支付失败事件
	if status == "failed" {
		s.bus.Publish(events.PaymentFailedEvent{
			PaymentID:     payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentMethod: payment.PaymentMethod,
			FailureReason: failureReason,
		})
	}

	return nil
}

// ListPayments 获取支付列表
func (s *PaymentService) ListPayments(ctx context.Context, q dto.PaymentListQuery) (*dto.PaymentListResult, error) {
	offset := (q.Page - 1) * q.Size

	payments, err := s.paymentRepo.ListByQuery(ctx, q, offset, q.Size)
	if err != nil {
		return nil, fmt.Errorf("list payments failed: %w", err)
	}

	total, err := s.paymentRepo.CountByQuery(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("count payments failed: %w", err)
	}

	paymentResponses := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = dto.PaymentResponse{
			ID:            payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentMethod: payment.PaymentMethod,
			TransactionID: payment.TransactionID,
			Status:        payment.Status,
			CreatedAt:     payment.CreatedAt.ToTime().String(),
		}

		if payment.PaidAt != nil {
			paymentResponses[i].PaidAt = payment.PaidAt.ToTime().String()
		}
	}

	return &dto.PaymentListResult{
		List:  paymentResponses,
		Total: total,
	}, nil
}

// CreateRefund 创建退款
func (s *PaymentService) CreateRefund(ctx context.Context, req *dto.CreateRefundRequest) (*models.Refund, error) {
	// 获取支付记录
	payment, err := s.paymentRepo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// 验证订单是否匹配
	if payment.OrderID != req.OrderID {
		return nil, errors.New("payment and order do not match")
	}

	// 验证支付状态
	if payment.Status != "success" {
		return nil, errors.New("payment is not successful")
	}

	// 验证退款金额
	if req.RefundAmount > payment.Amount {
		return nil, errors.New("refund amount exceeds payment amount")
	}

	// 检查是否已经有退款记录
	refunds, err := s.refundRepo.GetByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("get refunds failed: %w", err)
	}

	// 计算已退款金额
	var refundedAmount int64
	for _, refund := range refunds {
		if refund.Status == "success" {
			refundedAmount += refund.RefundAmount
		}
	}

	// 验证剩余可退款金额
	if req.RefundAmount > payment.Amount-refundedAmount {
		return nil, errors.New("refund amount exceeds available amount")
	}

	// 创建退款记录
	refund := &models.Refund{
		PaymentID:    req.PaymentID,
		OrderID:      req.OrderID,
		RefundAmount: req.RefundAmount,
		RefundReason: req.RefundReason,
		Status:       "pending",
	}

	if err := s.refundRepo.Create(ctx, refund); err != nil {
		return nil, fmt.Errorf("create refund failed: %w", err)
	}

	// 发布退款创建事件
	s.bus.Publish(events.RefundCreatedEvent{
		RefundID:     refund.ID,
		PaymentID:    refund.PaymentID,
		OrderID:      refund.OrderID,
		RefundAmount: refund.RefundAmount,
		RefundReason: refund.RefundReason,
	})

	return refund, nil
}

// UpdateRefundStatus 更新退款状态
func (s *PaymentService) UpdateRefundStatus(ctx context.Context, id int64, status string, transactionID string, failureReason string) error {
	// 获取退款记录
	refund, err := s.refundRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("refund not found: %w", err)
	}

	previousStatus := refund.Status

	// 更新退款状态
	refund.Status = status
	refund.TransactionID = transactionID
	refund.FailureReason = failureReason

	if err := s.refundRepo.Update(ctx, refund); err != nil {
		return fmt.Errorf("update refund status failed: %w", err)
	}

	// 发布退款状态更新事件
	s.bus.Publish(events.RefundStatusUpdatedEvent{
		RefundID:       refund.ID,
		PaymentID:      refund.PaymentID,
		OrderID:        refund.OrderID,
		Status:         refund.Status,
		PreviousStatus: previousStatus,
		TransactionID:  transactionID,
	})

	// 如果退款成功，更新订单状态
	if status == "success" {
		// 检查是否所有金额都已退款
		payment, err := s.paymentRepo.GetByID(ctx, refund.PaymentID)
		if err == nil {
			refunds, err := s.refundRepo.GetByPaymentID(ctx, refund.PaymentID)
			if err == nil {
				var totalRefunded int64
				for _, r := range refunds {
					if r.Status == "success" {
						totalRefunded += r.RefundAmount
					}
				}

				// 如果所有金额都已退款，更新订单状态为 refunded
				if totalRefunded >= payment.Amount {
					if err := s.orderRepo.UpdateStatus(ctx, refund.OrderID, "refunded"); err != nil {
						logger.Error("update order status to refunded failed", "order_id", refund.OrderID, "error", err)
					}
				}
			}
		}
	}

	// 如果退款失败，发布退款失败事件
	if status == "failed" {
		s.bus.Publish(events.RefundFailedEvent{
			RefundID:      refund.ID,
			PaymentID:     refund.PaymentID,
			OrderID:       refund.OrderID,
			RefundAmount:  refund.RefundAmount,
			FailureReason: failureReason,
		})
	}

	return nil
}

// ListRefunds 获取退款列表
func (s *PaymentService) ListRefunds(ctx context.Context, q dto.RefundListQuery) (*dto.RefundListResponse, error) {
	offset := (q.Page - 1) * q.PageSize

	refunds, err := s.refundRepo.ListByQuery(ctx, q, offset, q.PageSize)
	if err != nil {
		return nil, fmt.Errorf("list refunds failed: %w", err)
	}

	total, err := s.refundRepo.CountByQuery(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("count refunds failed: %w", err)
	}

	refundResponses := make([]dto.RefundResponse, len(refunds))
	for i, refund := range refunds {
		refundResponses[i] = dto.RefundResponse{
			ID:            refund.ID,
			PaymentID:     refund.PaymentID,
			OrderID:       refund.OrderID,
			RefundAmount:  refund.RefundAmount,
			RefundReason:  refund.RefundReason,
			TransactionID: refund.TransactionID,
			Status:        refund.Status,
			CreatedAt:     refund.CreatedAt.ToTime().String(),
		}
	}

	return &dto.RefundListResponse{
		Refunds:  refundResponses,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

// ListPaymentMethods 获取支付方式列表
func (s *PaymentService) ListPaymentMethods(ctx context.Context, page, pageSize int) (*dto.PaymentMethodListResponse, error) {
	offset := (page - 1) * pageSize

	paymentMethods, total, err := s.paymentMethodRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list payment methods failed: %w", err)
	}

	paymentMethodResponses := make([]dto.PaymentMethodResponse, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		paymentMethodResponses[i] = dto.PaymentMethodResponse{
			ID:          paymentMethod.ID,
			Code:        paymentMethod.Code,
			Name:        paymentMethod.Name,
			Description: paymentMethod.Description,
			Status:      paymentMethod.Status,
			CreatedAt:   paymentMethod.CreatedAt.ToTime().String(),
		}
	}

	return &dto.PaymentMethodListResponse{
		PaymentMethods: paymentMethodResponses,
		Total:          total,
		Page:           page,
		PageSize:       pageSize,
	}, nil
}

// GetPaymentMethodByCode 根据编码获取支付方式
func (s *PaymentService) GetPaymentMethodByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	return s.paymentMethodRepo.GetByCode(ctx, code)
}
