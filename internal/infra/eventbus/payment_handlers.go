package eventbus

import (
	"eshop-monolith/pkg/logger"
	"eshop-monolith/internal/payment/events"
)

// RegisterPaymentHandlers 注册支付事件处理器
func RegisterPaymentHandlers(bus *Bus) {
	bus.Subscribe("payment.PaymentCreatedEvent", handlePaymentCreated)
	bus.Subscribe("payment.PaymentStatusUpdatedEvent", handlePaymentStatusUpdated)
	bus.Subscribe("payment.PaymentFailedEvent", handlePaymentFailed)
	bus.Subscribe("payment.RefundCreatedEvent", handleRefundCreated)
	bus.Subscribe("payment.RefundStatusUpdatedEvent", handleRefundStatusUpdated)
	bus.Subscribe("payment.RefundFailedEvent", handleRefundFailed)
}

// handlePaymentCreated 处理支付创建事件
func handlePaymentCreated(event interface{}) {
	e, ok := event.(events.PaymentCreatedEvent)
	if !ok {
		return
	}
	logger.Info("Payment created", "payment_id", e.PaymentID, "order_id", e.OrderID, "amount", e.Amount, "payment_method", e.PaymentMethod)
	// 这里可以添加支付创建通知等逻辑
}

// handlePaymentStatusUpdated 处理支付状态更新事件
func handlePaymentStatusUpdated(event interface{}) {
	e, ok := event.(events.PaymentStatusUpdatedEvent)
	if !ok {
		return
	}
	logger.Info("Payment status updated", "payment_id", e.PaymentID, "order_id", e.OrderID, "status", e.Status, "previous_status", e.PreviousStatus)
	// 这里可以添加支付状态更新通知等逻辑
}

// handlePaymentFailed 处理支付失败事件
func handlePaymentFailed(event interface{}) {
	e, ok := event.(events.PaymentFailedEvent)
	if !ok {
		return
	}
	logger.Error("Payment failed", "payment_id", e.PaymentID, "order_id", e.OrderID, "amount", e.Amount, "failure_reason", e.FailureReason)
	// 这里可以添加支付失败通知等逻辑
}

// handleRefundCreated 处理退款创建事件
func handleRefundCreated(event interface{}) {
	e, ok := event.(events.RefundCreatedEvent)
	if !ok {
		return
	}
	logger.Info("Refund created", "refund_id", e.RefundID, "payment_id", e.PaymentID, "order_id", e.OrderID, "refund_amount", e.RefundAmount)
	// 这里可以添加退款创建通知等逻辑
}

// handleRefundStatusUpdated 处理退款状态更新事件
func handleRefundStatusUpdated(event interface{}) {
	e, ok := event.(events.RefundStatusUpdatedEvent)
	if !ok {
		return
	}
	logger.Info("Refund status updated", "refund_id", e.RefundID, "payment_id", e.PaymentID, "order_id", e.OrderID, "status", e.Status, "previous_status", e.PreviousStatus)
	// 这里可以添加退款状态更新通知等逻辑
}

// handleRefundFailed 处理退款失败事件
func handleRefundFailed(event interface{}) {
	e, ok := event.(events.RefundFailedEvent)
	if !ok {
		return
	}
	logger.Error("Refund failed", "refund_id", e.RefundID, "payment_id", e.PaymentID, "order_id", e.OrderID, "refund_amount", e.RefundAmount, "failure_reason", e.FailureReason)
	// 这里可以添加退款失败通知等逻辑
}
