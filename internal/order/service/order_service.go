package service

import (
	"context"
	"errors"
	"time"

	"eshop-monolith/internal/eventbus"
	inventoryRepo "eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/pkg/errcode"

	"eshop-monolith/internal/order/api/dto"
	"eshop-monolith/internal/order/domain/models"
	"eshop-monolith/internal/order/domain/repositories"
	"eshop-monolith/internal/order/events"
)

type OrderService struct {
	orderRepo     repositories.IorderRepository
	inventoryRepo inventoryRepo.IinventoryRepository
	bus           *eventbus.Bus
}

// NewOrderService 创建订单服务
func NewOrderService(orderRepo repositories.IorderRepository, inventoryRepo inventoryRepo.IinventoryRepository, bus *eventbus.Bus) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		bus:           bus,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderDTO) (*models.Order, error) {
	// 计算订单总金额
	totalAmount := int64(0)
	var orderItems []models.OrderItem

	// 预占库存并创建订单项
	for _, item := range req.Items {
		// 预占库存
		if err := s.inventoryRepo.ReserveInventory(ctx, 0, item.Quantity); err != nil {
			return nil, err
		}

		// 计算单项金额
		amount := item.UnitPrice * int64(item.Quantity)
		totalAmount += amount

		// 创建订单项
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Amount:    amount,
		})
	}

	// 创建订单
	order := &models.Order{
		CustomerID:  req.CustomerID,
		TotalAmount: totalAmount,
		Currency:    req.Currency,
		Status:      models.OrderStatusPending,
		Items:       orderItems,
	}

	// 保存订单
	if err := s.orderRepo.Create(ctx, order); err != nil {
		// 释放已预占的库存
		for _, item := range req.Items {
			s.inventoryRepo.ReleaseInventory(ctx, 0, item.Quantity)
		}
		return nil, err
	}

	// 发布订单创建事件
	s.bus.Publish(events.OrderCreatedEvent{
		OrderID:     order.ID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Currency:    order.Currency,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.ToTime(),
	})

	return order, nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(ctx context.Context, orderID int64) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return order, nil
}

// UpdateOrder 更新订单
func (s *OrderService) UpdateOrder(ctx context.Context, orderID int64, req *dto.UpdateOrderDTO) (*models.Order, error) {
	// 获取订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 更新订单信息
	if req.Status != "" {
		if err := s.orderRepo.UpdateStatus(ctx, orderID, req.Status); err != nil {
			return nil, err
		}
		order.Status = req.Status
	}

	return order, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	// 获取订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	// 释放库存
	for _, item := range order.Items {
		if err := s.inventoryRepo.ReleaseInventory(ctx, 0, item.Quantity); err != nil {
			return err
		}
	}

	// 更新订单状态
	if err := s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
		return err
	}

	// 发布订单取消事件
	s.bus.Publish(events.OrderCancelledEvent{
		OrderID:     order.ID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		CancelledAt: time.Now(),
	})

	return nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	// 检查状态是否有效
	validStatuses := map[string]bool{
		models.OrderStatusPending:   true,
		models.OrderStatusPaid:      true,
		models.OrderStatusShipped:   true,
		models.OrderStatusDelivered: true,
		models.OrderStatusCancelled: true,
	}

	if !validStatuses[status] {
		return errors.New("invalid order status")
	}

	// 获取订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 处理状态转换逻辑
	if status == models.OrderStatusPaid && order.Status == models.OrderStatusPending {
		// 支付成功，扣减库存
		for _, item := range order.Items {
			if err := s.inventoryRepo.DeductInventory(ctx, 0, item.Quantity); err != nil {
				return err
			}
		}
	} else if status == models.OrderStatusCancelled && order.Status == models.OrderStatusPending {
		// 取消订单，释放库存
		for _, item := range order.Items {
			if err := s.inventoryRepo.ReleaseInventory(ctx, 0, item.Quantity); err != nil {
				return err
			}
		}
	}

	// 更新状态
	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}

// ListOrders 列出订单
func (s *OrderService) ListOrders(ctx context.Context, q dto.OrderListQuery) (*dto.OrderListResult, error) {
	offset := (q.Page - 1) * q.Size
	orders, err := s.orderRepo.ListByQuery(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.orderRepo.CountByQuery(ctx, q)
	if err != nil {
		return nil, err
	}

	return &dto.OrderListResult{
		List:  orders,
		Total: total,
	}, nil
}

// DeleteOrder 删除订单
func (s *OrderService) DeleteOrder(ctx context.Context, orderID int64) error {
	// 获取订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 检查订单状态（只有已取消或已完成的订单可以删除）
	if order.Status != models.OrderStatusCancelled && order.Status != models.OrderStatusDelivered {
		return errors.New("only cancelled or delivered orders can be deleted")
	}

	// 删除订单
	return s.orderRepo.Delete(ctx, orderID)
}

// GetOrdersByUserID 根据用户ID获取订单列表
func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64, page, pageSize int) ([]models.Order, int64, error) {
	return s.orderRepo.FindByUserID(ctx, userID, page, pageSize)
}
