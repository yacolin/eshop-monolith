package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	flashEvents "eshop-monolith/internal/flashsale/events"
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/trade"
	"eshop-monolith/pkg/logger"
)

type NotificationService struct {
	repo InotificationRepository
}

func NewNotificationService(repo InotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID int64, title, content string, channel, category int8) (*Notification, error) {
	n := &Notification{
		UserID:   userID,
		Title:    title,
		Content:  content,
		Channel:  channel,
		Category: category,
		IsRead:   false,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("create notification failed: %w", err)
	}
	return n, nil
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID int64, page, size int) ([]*Notification, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	list, total, err := s.repo.ListByUserID(ctx, userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Notification, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result, total, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	return s.repo.Delete(ctx, notificationID, userID)
}

// HandleMessage RabbitMQ 消息处理器
func (s *NotificationService) HandleMessage(msg rabbitmq.Message) error {
	switch msg.RoutingKey {
	case "order.paid":
		var e trade.OrderPaidEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleOrderPaid(e)
	case "order.shipped":
		var e trade.OrderShippedEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleOrderShipped(e)
	case "order.delivered":
		var e trade.OrderDeliveredEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleOrderDelivered(e)
	case "order.cancelled":
		var e trade.OrderCancelledEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleOrderCancelled(e)
	case "flash-order.created":
		var e flashEvents.FlashOrderCreatedEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleFlashOrderCreated(e)
	case "flash-order.paid":
		var e flashEvents.FlashOrderPaidEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleFlashOrderPaid(e)
	case "flash-order.cancelled":
		var e flashEvents.FlashOrderCancelledEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleFlashOrderCancelled(e)
	case "payment.refund.created":
		var e trade.RefundCreatedEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleRefundCreated(e)
	case "payment.refund.failed":
		var e trade.RefundFailedEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleRefundFailed(e)
	case "inventory.low":
		var e inventory.InventoryLowEvent
		json.Unmarshal(msg.Payload, &e)
		s.handleInventoryLow(e)
	}
	return nil
}

func parseCustomerID(customerID string) int64 {
	id, err := strconv.ParseInt(customerID, 10, 64)
	if err != nil {
		logger.Warn("parse CustomerID failed", "customer_id", customerID)
		return 0
	}
	return id
}

func (s *NotificationService) handleOrderPaid(e trade.OrderPaidEvent) {
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	s.CreateNotification(context.Background(), userID, "订单支付成功",
		fmt.Sprintf("您的订单 #%d 已支付成功", e.OrderID), ChannelInApp, CategoryOrder)
}

func (s *NotificationService) handleOrderShipped(e trade.OrderShippedEvent) {
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	s.CreateNotification(context.Background(), userID, "订单已发货",
		fmt.Sprintf("您的订单 #%d 已发货", e.OrderID), ChannelInApp, CategoryOrder)
}

func (s *NotificationService) handleOrderDelivered(e trade.OrderDeliveredEvent) {
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	s.CreateNotification(context.Background(), userID, "订单已签收",
		fmt.Sprintf("您的订单 #%d 已签收", e.OrderID), ChannelInApp, CategoryOrder)
}

func (s *NotificationService) handleOrderCancelled(e trade.OrderCancelledEvent) {
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	s.CreateNotification(context.Background(), userID, "订单已取消",
		fmt.Sprintf("您的订单 #%d 已取消", e.OrderID), ChannelInApp, CategoryOrder)
}

func (s *NotificationService) handleFlashOrderCreated(e flashEvents.FlashOrderCreatedEvent) {
	if e.UserID == 0 {
		return
	}
	s.CreateNotification(context.Background(), e.UserID, "抢购成功",
		fmt.Sprintf("您已成功抢购商品，订单 #%d 待支付", e.OrderID), ChannelInApp, CategoryMarketing)
}

func (s *NotificationService) handleFlashOrderPaid(e flashEvents.FlashOrderPaidEvent) {
	if e.UserID == 0 {
		return
	}
	s.CreateNotification(context.Background(), e.UserID, "闪购订单支付成功",
		fmt.Sprintf("闪购订单 #%d 已支付", e.OrderID), ChannelInApp, CategoryMarketing)
}

func (s *NotificationService) handleFlashOrderCancelled(e flashEvents.FlashOrderCancelledEvent) {
	if e.UserID == 0 {
		return
	}
	s.CreateNotification(context.Background(), e.UserID, "闪购订单已取消",
		fmt.Sprintf("闪购订单 #%d 已取消", e.OrderID), ChannelInApp, CategoryMarketing)
}

func (s *NotificationService) handleRefundCreated(e trade.RefundCreatedEvent) {
	logger.Info("refund created", "event", e)
}

func (s *NotificationService) handleRefundFailed(e trade.RefundFailedEvent) {
	logger.Warn("refund failed", "refund_id", e.RefundID, "order_id", e.OrderID, "reason", e.FailureReason)
}

func (s *NotificationService) handleInventoryLow(e inventory.InventoryLowEvent) {
	s.CreateNotification(context.Background(), 0, "库存预警",
		fmt.Sprintf("商品 %s 库存不足，当前库存 %d，阈值 %d", e.SkuID, e.Quantity, e.Threshold), ChannelInApp, CategorySystem)
}
