package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	couponSvc "eshop-monolith/internal/coupon/service"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/pkg/errcode"

	"eshop-monolith/internal/order/api/dto"
	"eshop-monolith/internal/order/domain/models"
	"eshop-monolith/internal/order/domain/repositories"
	"eshop-monolith/internal/order/events"

	"gorm.io/gorm"
)

// generateOrderNo 生成全局唯一订单号
func generateOrderNo() string {
	now := time.Now().UnixMilli()
	r := rand.Intn(10000)
	return fmt.Sprintf("ORD%d%04d", now, r)
}

func orderToResponse(o *models.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = orderItemToResponse(&item, o.OrderNo)
	}
	return dto.OrderResponse{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		CustomerID:     o.CustomerID,
		TotalAmount:    o.TotalAmount,
		DiscountAmount: o.DiscountAmount,
		CouponID:       o.CouponID,
		Currency:       o.Currency,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
		Items:          items,
	}
}

func orderItemToResponse(item *models.OrderItem, orderNo string) dto.OrderItemResponse {
	return dto.OrderItemResponse{
		ID:        item.ID,
		OrderID:   item.OrderID,
		OrderNo:   orderNo,
		SkuID:     item.SkuID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		Amount:    item.Amount,
	}
}

type OrderService struct {
	db            *gorm.DB
	orderRepo     repositories.IorderRepository
	inventoryRepo repositories.InventoryForOrder
	skuForOrder   repositories.SkuForOrder
	couponService *couponSvc.CouponService
	bus           *eventbus.Bus
}

func NewOrderService(db *gorm.DB, orderRepo repositories.IorderRepository, inventoryRepo repositories.InventoryForOrder, skuForOrder repositories.SkuForOrder, bus *eventbus.Bus, couponService *couponSvc.CouponService) *OrderService {
	return &OrderService{
		db:            db,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		skuForOrder:   skuForOrder,
		couponService: couponService,
		bus:           bus,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderDTO) (*models.Order, error) {
	var order *models.Order
	orderNo := generateOrderNo()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		totalAmount := int64(0)
		var orderItems []models.OrderItem

		for _, item := range req.Items {
			productID, price, err := s.skuForOrder.GetSkuInfo(ctx, item.SkuID)
			if err != nil {
				return err
			}
			if err := s.inventoryRepo.ReserveWithTx(tx, item.SkuID, item.Quantity); err != nil {
				return err
			}
			amount := price * int64(item.Quantity)
			totalAmount += amount
			orderItems = append(orderItems, models.OrderItem{
				SkuID:     item.SkuID,
				ProductID: strconv.FormatInt(productID, 10),
				Quantity:  item.Quantity,
				UnitPrice: price,
				Amount:    amount,
			})
		}

		// 计算实际支付金额（扣除优惠券）
		discountAmount := int64(0)
		var couponID *int64

		if req.UserCouponID != nil && *req.UserCouponID > 0 && s.couponService != nil {
			userID, _ := strconv.ParseInt(req.CustomerID, 10, 64)
			if userID > 0 {
				discount, err := s.couponService.UseCoupon(ctx, *req.UserCouponID, userID, orderNo, totalAmount)
				if err != nil {
					return err
				}
				discountAmount = discount
				couponID = req.UserCouponID
			}
		}

		order = &models.Order{
			OrderNo:        orderNo,
			CustomerID:     req.CustomerID,
			TotalAmount:    totalAmount - discountAmount,
			DiscountAmount: discountAmount,
			CouponID:       couponID,
			Currency:       req.Currency,
			Status:         models.OrderStatusPending,
			Items:          orderItems,
		}

		return s.orderRepo.CreateWithTx(tx, order)
	})
	if err != nil {
		return nil, err
	}

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

func (s *OrderService) GetOrder(ctx context.Context, orderID int64) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return order, nil
}

// UpdateOrder 更新订单，含 status 变更时走完整库存逻辑
func (s *OrderService) UpdateOrder(ctx context.Context, orderID int64, req *dto.UpdateOrderDTO) (*models.Order, error) {
	if req.Status != "" {
		if err := s.UpdateOrderStatus(ctx, orderID, req.Status); err != nil {
			return nil, err
		}
	}

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	var order models.Order

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		o, err := s.orderRepo.FindByIDWithTx(tx, orderID)
		if err != nil {
			return err
		}
		order = *o

		if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusPaid {
			return errors.New("only pending or paid orders can be cancelled")
		}

		for _, item := range order.Items {
			switch order.Status {
			case models.OrderStatusPending:
				if releaseErr := s.inventoryRepo.ReleaseWithTx(tx, item.SkuID, item.Quantity); releaseErr != nil {
					return releaseErr
				}
			case models.OrderStatusPaid:
				if restoreErr := s.inventoryRepo.RestoreWithTx(tx, item.SkuID, item.Quantity); restoreErr != nil {
					return restoreErr
				}
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

func (s *OrderService) HandlePaidSuccess(ctx context.Context, orderID int64) error {
	if err := s.UpdateOrderStatus(ctx, orderID, models.OrderStatusPaid); err != nil {
		return err
	}

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

// UpdateOrderStatus 更新订单状态（事务内完成库存操作 + 状态更新）
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

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orderRepo.FindByIDWithTx(tx, orderID)
		if err != nil {
			return err
		}

		switch {
		case status == models.OrderStatusPaid && order.Status == models.OrderStatusPending:
			for _, item := range order.Items {
				if deductErr := s.inventoryRepo.DeductWithTx(tx, item.SkuID, item.Quantity); deductErr != nil {
					return deductErr
				}
			}

		case status == models.OrderStatusCancelled && order.Status == models.OrderStatusPending:
			for _, item := range order.Items {
				if releaseErr := s.inventoryRepo.ReleaseWithTx(tx, item.SkuID, item.Quantity); releaseErr != nil {
					return releaseErr
				}
			}

		case status == models.OrderStatusCancelled && order.Status == models.OrderStatusPaid:
			for _, item := range order.Items {
				if restoreErr := s.inventoryRepo.RestoreWithTx(tx, item.SkuID, item.Quantity); restoreErr != nil {
					return restoreErr
				}
			}
		}

		return s.orderRepo.UpdateStatusWithTx(tx, orderID, status)
	})
}

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

	result := &dto.OrderListResult{
		List:  make([]dto.OrderResponse, len(orders)),
		Total: total,
	}
	for i, o := range orders {
		result.List[i] = orderToResponse(&o)
	}
	return result, nil
}

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

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64, page, pageSize int) (*dto.OrderListResult, error) {
	orders, total, err := s.orderRepo.FindByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	result := &dto.OrderListResult{
		List:  make([]dto.OrderResponse, len(orders)),
		Total: total,
	}
	for i, o := range orders {
		result.List[i] = orderToResponse(&o)
	}
	return result, nil
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID int64, page, pageSize int) (*dto.OrderItemListResult, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, total, err := s.orderRepo.FindItemsByOrderID(ctx, orderID, page, pageSize)
	if err != nil {
		return nil, err
	}

	result := &dto.OrderItemListResult{
		List:  make([]dto.OrderItemResponse, len(items)),
		Total: total,
	}
	for i, item := range items {
		result.List[i] = orderItemToResponse(&item, order.OrderNo)
	}
	return result, nil
}

func (s *OrderService) ListAllOrderItems(ctx context.Context, q dto.OrderItemListQuery) (*dto.OrderItemListResult, error) {
	offset := (q.Page - 1) * q.Size
	items, total, err := s.orderRepo.ListAllItems(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	orderIDs := make([]int64, 0, len(items))
	seen := make(map[int64]bool)
	for _, item := range items {
		if !seen[item.OrderID] {
			orderIDs = append(orderIDs, item.OrderID)
			seen[item.OrderID] = true
		}
	}
	orderNoMap, _ := s.orderRepo.BatchGetOrderNo(ctx, orderIDs)

	result := &dto.OrderItemListResult{
		List:  make([]dto.OrderItemResponse, len(items)),
		Total: total,
	}
	for i, item := range items {
		result.List[i] = orderItemToResponse(&item, orderNoMap[item.OrderID])
	}
	return result, nil
}
