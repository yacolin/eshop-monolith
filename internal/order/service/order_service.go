package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/pkg/errcode"

	"eshop-monolith/internal/order/api/dto"
	"eshop-monolith/internal/order/domain/models"
	"eshop-monolith/internal/order/domain/repositories"
	"eshop-monolith/internal/order/events"

	"gorm.io/gorm"
)

type OrderService struct {
	db            *gorm.DB
	orderRepo     repositories.IorderRepository
	inventoryRepo repositories.InventoryForOrder
	bus           *eventbus.Bus
}

// NewOrderService 创建订单服务
func NewOrderService(db *gorm.DB, orderRepo repositories.IorderRepository, inventoryRepo repositories.InventoryForOrder, bus *eventbus.Bus) *OrderService {
	return &OrderService{
		db:            db,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		bus:           bus,
	}
}

// CreateOrder 创建订单（事务内完成库存预占 + 订单写入）
func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderDTO) (*models.Order, error) {
	var order *models.Order

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		totalAmount := int64(0)
		var orderItems []models.OrderItem

		for _, item := range req.Items {
			pid, err := strconv.ParseInt(item.ProductID, 10, 64)
			if err != nil {
				return errcode.ErrInvalidParams
			}
			if err := s.inventoryRepo.ReserveWithTx(tx, pid, item.Quantity); err != nil {
				return err
			}
			amount := item.UnitPrice * int64(item.Quantity)
			totalAmount += amount
			orderItems = append(orderItems, models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				Amount:    amount,
			})
		}

		order = &models.Order{
			CustomerID:  req.CustomerID,
			TotalAmount: totalAmount,
			Currency:    req.Currency,
			Status:      models.OrderStatusPending,
			Items:       orderItems,
		}

		return s.orderRepo.CreateWithTx(tx, order)
	})
	if err != nil {
		return nil, err
	}

	// 事务外发布事件
	if s.bus != nil {
		s.bus.Publish(events.OrderCreatedEvent{
			OrderID:     order.ID,
			CustomerID:  order.CustomerID,
			TotalAmount: order.TotalAmount,
			Currency:    order.Currency,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt.ToTime(),
		})
	}

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
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if req.Status != "" {
		if err := s.orderRepo.UpdateStatus(ctx, orderID, req.Status); err != nil {
			return nil, err
		}
		order.Status = req.Status
	}

	return order, nil
}

// CancelOrder 取消订单（事务内释放库存 + 更新状态）
func (s *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			pid, parseErr := strconv.ParseInt(item.ProductID, 10, 64)
			if parseErr != nil {
				return errcode.ErrInvalidParams
			}
			if releaseErr := s.inventoryRepo.ReleaseWithTx(tx, pid, item.Quantity); releaseErr != nil {
				return releaseErr
			}
		}
		return s.orderRepo.UpdateStatusWithTx(tx, orderID, models.OrderStatusCancelled)
	})
	if err != nil {
		return err
	}

	if s.bus != nil {
		s.bus.Publish(events.OrderCancelledEvent{
			OrderID:     order.ID,
			CustomerID:  order.CustomerID,
			TotalAmount: order.TotalAmount,
			CancelledAt: time.Now(),
		})
	}

	return nil
}

// HandlePaidSuccess 支付成功处理（事务内扣减库存+更新订单状态, 事务外发布事件）
func (s *OrderService) HandlePaidSuccess(ctx context.Context, orderID int64) error {
	if err := s.UpdateOrderStatus(ctx, orderID, models.OrderStatusPaid); err != nil {
		return err
	}

	// 事务外发布支付成功事件
	if s.bus != nil {
		order, err := s.orderRepo.FindByID(ctx, orderID)
		if err == nil {
			s.bus.Publish(events.OrderPaidEvent{
				OrderID:     order.ID,
				CustomerID:  order.CustomerID,
				TotalAmount: order.TotalAmount,
				Currency:    order.Currency,
				PaidAt:      time.Now(),
			})
		}
	}

	return nil
}

// UpdateOrderStatus 更新订单状态（事务内完成库存扣减/释放 + 状态更新）
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
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

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if status == models.OrderStatusPaid && order.Status == models.OrderStatusPending {
			for _, item := range order.Items {
				pid, parseErr := strconv.ParseInt(item.ProductID, 10, 64)
				if parseErr != nil {
					return errcode.ErrInvalidParams
				}
				if deductErr := s.inventoryRepo.DeductWithTx(tx, pid, item.Quantity); deductErr != nil {
					return deductErr
				}
			}
		} else if status == models.OrderStatusCancelled && order.Status == models.OrderStatusPending {
			for _, item := range order.Items {
				pid, parseErr := strconv.ParseInt(item.ProductID, 10, 64)
				if parseErr != nil {
					return errcode.ErrInvalidParams
				}
				if releaseErr := s.inventoryRepo.ReleaseWithTx(tx, pid, item.Quantity); releaseErr != nil {
					return releaseErr
				}
			}
		}
		return s.orderRepo.UpdateStatusWithTx(tx, orderID, status)
	})

	return err
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
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusCancelled && order.Status != models.OrderStatusDelivered {
		return errors.New("only cancelled or delivered orders can be deleted")
	}

	return s.orderRepo.Delete(ctx, orderID)
}

// GetOrdersByUserID 根据用户ID获取订单列表
func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64, page, pageSize int) ([]models.Order, int64, error) {
	return s.orderRepo.FindByUserID(ctx, userID, page, pageSize)
}
