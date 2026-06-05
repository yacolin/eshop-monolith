package service

import (
	"context"
	"fmt"
	"strconv"

	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/notification/domain/models"
	"eshop-monolith/internal/notification/domain/repositories"
	inventoryEvents "eshop-monolith/internal/inventory/events"
	orderEvents "eshop-monolith/internal/order/events"
	paymentEvents "eshop-monolith/internal/payment/events"
	flashEvents "eshop-monolith/internal/flashsale/events"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/query"
)

// NotificationService 通知服务
type NotificationService struct {
	repo *repositories.NotificationRepository
}

// NewNotificationService 创建通知服务
func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateNotification 创建通知
func (s *NotificationService) CreateNotification(ctx context.Context, userID int64, title, content string, notifType models.NotificationType) (*models.Notification, error) {
	n := &models.Notification{
		UserID:  userID,
		Title:   title,
		Content: content,
		Type:    notifType,
		IsRead:  false,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("创建通知失败: %w", err)
	}
	return n, nil
}

// ListNotifications 获取用户通知列表（分页）
func (s *NotificationService) ListNotifications(ctx context.Context, userID int64, page, size int) (*query.ListResult[*models.Notification], error) {
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
		return nil, fmt.Errorf("查询通知列表失败: %w", err)
	}

	result := make([]*models.Notification, len(list))
	for i := range list {
		result[i] = &list[i]
	}

	return &query.ListResult[*models.Notification]{
		Total: total,
		List:  result,
	}, nil
}

// GetUnreadCount 获取未读通知数
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

// MarkAsRead 标记通知为已读
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

// MarkAllAsRead 标记所有通知为已读
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	return s.repo.Delete(ctx, notificationID, userID)
}

// ---------------------------------------------------------------------------
// EventBus 事件处理器 —— 通过闭包注入 bus 完成订阅
// ---------------------------------------------------------------------------

// RegisterHandlers 注册通知模块关心的事件
func (s *NotificationService) RegisterHandlers(bus *eventbus.Bus) {
	// 订单事件
	bus.Subscribe("order.OrderPaidEvent", s.handleOrderPaid)
	bus.Subscribe("order.OrderShippedEvent", s.handleOrderShipped)
	bus.Subscribe("order.OrderDeliveredEvent", s.handleOrderDelivered)
	bus.Subscribe("order.OrderCancelledEvent", s.handleOrderCancelled)

	// 闪购事件
	bus.Subscribe("flashsale.FlashOrderCreatedEvent", s.handleFlashOrderCreated)
	bus.Subscribe("flashsale.FlashOrderPaidEvent", s.handleFlashOrderPaid)
	bus.Subscribe("flashsale.FlashOrderCancelledEvent", s.handleFlashOrderCancelled)

	// 退款事件
	bus.Subscribe("payment.RefundCreatedEvent", s.handleRefundCreated)
	bus.Subscribe("payment.RefundFailedEvent", s.handleRefundFailed)

	// 库存预警
	bus.Subscribe("inventory.InventoryLowEvent", s.handleInventoryLow)
}

// handleOrderPaid 订单支付成功
func (s *NotificationService) handleOrderPaid(event interface{}) {
	e, ok := event.(orderEvents.OrderPaidEvent)
	if !ok {
		return
	}
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), userID,
		"订单支付成功",
		fmt.Sprintf("您的订单 #%d 已支付成功，金额 %.2f 元", e.OrderID, float64(e.TotalAmount)/100),
		models.NotificationTypeOrder)
	if err != nil {
		logger.Error("创建订单支付通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleOrderShipped 订单发货
func (s *NotificationService) handleOrderShipped(event interface{}) {
	e, ok := event.(orderEvents.OrderShippedEvent)
	if !ok {
		return
	}
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), userID,
		"订单已发货",
		fmt.Sprintf("您的订单 #%d 已发货，请耐心等待收货", e.OrderID),
		models.NotificationTypeOrder)
	if err != nil {
		logger.Error("创建订单发货通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleOrderDelivered 订单已签收
func (s *NotificationService) handleOrderDelivered(event interface{}) {
	e, ok := event.(orderEvents.OrderDeliveredEvent)
	if !ok {
		return
	}
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), userID,
		"订单已签收",
		fmt.Sprintf("您的订单 #%d 已签收，感谢您的购买！", e.OrderID),
		models.NotificationTypeOrder)
	if err != nil {
		logger.Error("创建订单签收通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleOrderCancelled 订单已取消
func (s *NotificationService) handleOrderCancelled(event interface{}) {
	e, ok := event.(orderEvents.OrderCancelledEvent)
	if !ok {
		return
	}
	userID := parseCustomerID(e.CustomerID)
	if userID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), userID,
		"订单已取消",
		fmt.Sprintf("您的订单 #%d 已取消", e.OrderID),
		models.NotificationTypeOrder)
	if err != nil {
		logger.Error("创建订单取消通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleFlashOrderCreated 闪购订单创建
func (s *NotificationService) handleFlashOrderCreated(event interface{}) {
	e, ok := event.(flashEvents.FlashOrderCreatedEvent)
	if !ok {
		return
	}
	if e.UserID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), e.UserID,
		"恭喜！抢购成功",
		fmt.Sprintf("您已成功抢购商品，订单 #%d 待支付", e.OrderID),
		models.NotificationTypeFlash)
	if err != nil {
		logger.Error("创建闪购通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleFlashOrderPaid 闪购订单支付成功
func (s *NotificationService) handleFlashOrderPaid(event interface{}) {
	e, ok := event.(flashEvents.FlashOrderPaidEvent)
	if !ok {
		return
	}
	if e.UserID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), e.UserID,
		"闪购订单支付成功",
		fmt.Sprintf("您的闪购订单 #%d 已支付成功，等待发货", e.OrderID),
		models.NotificationTypeFlash)
	if err != nil {
		logger.Error("创建闪购支付通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleFlashOrderCancelled 闪购订单取消
func (s *NotificationService) handleFlashOrderCancelled(event interface{}) {
	e, ok := event.(flashEvents.FlashOrderCancelledEvent)
	if !ok {
		return
	}
	if e.UserID == 0 {
		return
	}
	_, err := s.CreateNotification(context.Background(), e.UserID,
		"闪购订单已取消",
		fmt.Sprintf("您的闪购订单 #%d 已取消", e.OrderID),
		models.NotificationTypeFlash)
	if err != nil {
		logger.Error("创建闪购取消通知失败", "order_id", e.OrderID, "error", err)
	}
}

// handleRefundCreated 退款已受理
func (s *NotificationService) handleRefundCreated(event interface{}) {
	// 退款事件中没有用户 ID，需要查订单；简化处理：记录日志
	logger.Info("退款已受理", "event", event)
}

// handleRefundFailed 退款失败
func (s *NotificationService) handleRefundFailed(event interface{}) {
	e, ok := event.(paymentEvents.RefundFailedEvent)
	if !ok {
		return
	}
	logger.Warn("退款失败", "refund_id", e.RefundID, "order_id", e.OrderID, "reason", e.FailureReason)
}

// handleInventoryLow 库存预警（推送给管理员 userID=0 表示全体管理员）
func (s *NotificationService) handleInventoryLow(event interface{}) {
	e, ok := event.(inventoryEvents.InventoryLowEvent)
	if !ok {
		return
	}
	_, err := s.CreateNotification(context.Background(), 0,
		"库存预警",
		fmt.Sprintf("商品 %s 库存不足，当前库存 %d，阈值 %d", e.ProductID, e.Quantity, e.Threshold),
		models.NotificationTypeSystem)
	if err != nil {
		logger.Error("创建库存预警通知失败", "product_id", e.ProductID, "error", err)
	}
}

// parseCustomerID 将订单中的 CustomerID (string) 转为 int64
func parseCustomerID(customerID string) int64 {
	id, err := strconv.ParseInt(customerID, 10, 64)
	if err != nil {
		logger.Warn("解析 CustomerID 失败", "customer_id", customerID)
		return 0
	}
	return id
}
