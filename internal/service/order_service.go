package service

import (
	"context"
	"fmt"

	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/domain/order"
	"eshop-monolith/internal/eventbus"
)

// OrderService 订单服务
type OrderService struct {
	orderRepo     order.Repository
	inventoryRepo inventory.Repository
	bus           *eventbus.Bus
}

// NewOrderService 创建订单服务
func NewOrderService(orderRepo order.Repository, inventoryRepo inventory.Repository, bus *eventbus.Bus) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		bus:           bus,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID          int64                 `json:"user_id"`
	Items           []OrderItemRequest    `json:"items"`
	ShippingAddress order.ShippingAddress `json:"shipping_address"`
}

// OrderItemRequest 订单商品请求
type OrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*order.Order, error) {
	// 计算总金额
	var totalAmount int64
	orderItems := make([]order.OrderItem, 0, len(req.Items))

	// 检查并预占库存
	for _, item := range req.Items {
		// 预占库存
		if err := s.inventoryRepo.ReserveInventory(ctx, item.ProductID, item.Quantity); err != nil {
			return nil, err
		}

		// 计算总金额
		totalAmount += item.UnitPrice * int64(item.Quantity)

		// 创建订单商品
		orderItems = append(orderItems, order.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	// 创建订单
	newOrder := &order.Order{
		UserID:          req.UserID,
		TotalAmount:     totalAmount,
		Status:          order.OrderStatusPending,
		ShippingAddress: req.ShippingAddress,
		Items:           orderItems,
	}

	// 保存订单
	if err := s.orderRepo.Create(ctx, newOrder); err != nil {
		// 释放库存
		for _, item := range req.Items {
			_ = s.inventoryRepo.ReleaseInventory(ctx, item.ProductID, item.Quantity)
		}
		return nil, err
	}

	// 发布订单创建事件
	s.bus.Publish(order.OrderCreatedEvent{
		OrderID:     fmt.Sprintf("%d", newOrder.ID),
		UserID:      fmt.Sprintf("%d", newOrder.UserID),
		TotalAmount: newOrder.TotalAmount,
	})

	return newOrder, nil
}

// GetOrderByID 根据ID获取订单
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*order.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}

// ListOrdersByUserID 根据用户ID列出订单
func (s *OrderService) ListOrdersByUserID(ctx context.Context, userID string, page, pageSize int) ([]order.Order, int64, error) {
	return s.orderRepo.FindByUserID(ctx, userID, page, pageSize)
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id, status string) error {
	// 获取订单
	existingOrder, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	oldStatus := existingOrder.Status

	// 更新状态
	if err := s.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	// 发布订单状态变更事件
	s.bus.Publish(order.OrderStatusChangedEvent{
		OrderID:   id,
		OldStatus: oldStatus,
		NewStatus: status,
		UserID:    fmt.Sprintf("%d", existingOrder.UserID),
	})

	// 如果订单取消，释放库存
	if status == order.OrderStatusCancelled {
		for _, item := range existingOrder.Items {
			_ = s.inventoryRepo.ReleaseInventory(ctx, item.ProductID, item.Quantity)
		}

		// 发布订单取消事件
		s.bus.Publish(order.OrderCancelledEvent{
			OrderID:     id,
			UserID:      fmt.Sprintf("%d", existingOrder.UserID),
			TotalAmount: existingOrder.TotalAmount,
		})
	}

	return nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	return s.UpdateOrderStatus(ctx, id, order.OrderStatusCancelled)
}
