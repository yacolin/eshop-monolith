package base

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/trade"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/utils"
)

type NotificationService struct {
	repo InotificationRepository
}

func NewNotificationService(repo InotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

// CreateNotification 同步创建通知
func (s *NotificationService) CreateNotification(ctx context.Context, userID int64, title, content string, channel, category int8) (*Notification, error) {
	now := time.Now()
	n := &Notification{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Channel:   channel,
		Category:  category,
		IsRead:    false,
		CreatedAt: utils.Timestamp(now),
		UpdatedAt: utils.Timestamp(now),
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("create notification failed: %w", err)
	}
	return n, nil
}

// CreateSystemNotification 根据模板代码或直接参数创建系统通知
func (s *NotificationService) CreateSystemNotification(ctx context.Context, userID int64, templateCode, title, content string) (*Notification, error) {
	if templateCode != "" {
		tmpl, err := s.repo.GetTemplateByCode(ctx, templateCode)
		if err != nil {
			return nil, fmt.Errorf("resolve template %q: %w", templateCode, err)
		}
		if title == "" {
			title = tmpl.TitleTemplate
		}
		if content == "" {
			content = tmpl.ContentTemplate
		}
	}
	if title == "" || content == "" {
		return nil, fmt.Errorf("title and content are required when template_code is not provided")
	}
	return s.CreateNotification(ctx, userID, title, content, ChannelInApp, CategorySystem)
}

func (s *NotificationService) ListTemplates(ctx context.Context) ([]NotificationTemplate, error) {
	return s.repo.FindActiveTemplates(ctx)
}

func (s *NotificationService) GetTemplateByCode(ctx context.Context, code string) (*NotificationTemplate, error) {
	return s.repo.GetTemplateByCode(ctx, code)
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID int64, req *NotificationListReq) (*NotificationListResult, error) {
	req.Normalize()
	list, total, err := s.repo.ListByUserID(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	respList := make([]*NotificationResp, len(list))
	for i := range list {
		respList[i] = toResp(&list[i])
	}
	return &NotificationListResult{Total: total, List: respList}, nil
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

func (s *NotificationService) handleRefundCreated(e trade.RefundCreatedEvent) {
	logger.Info("refund created", "event", e)
}

func (s *NotificationService) handleRefundFailed(e trade.RefundFailedEvent) {
	logger.Warn("refund failed", "refund_id", e.RefundID, "order_id", e.OrderID, "reason", e.FailureReason)
}

func (s *NotificationService) handleInventoryLow(e inventory.InventoryLowEvent) {
	s.CreateNotification(context.Background(), 0, "库存预警",
		fmt.Sprintf("商品 %d 库存不足，当前库存 %d，阈值 %d", e.SkuID, e.Quantity, e.Threshold), ChannelInApp, CategorySystem)
}
