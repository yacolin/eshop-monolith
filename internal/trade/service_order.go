package trade

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"eshop-monolith/internal/infra/rabbitmq"
)

type OrderService struct {
	repo     IorderRepository
	skuSvc   SkuProvider
	invSvc   InventoryService
	db       *gorm.DB
	eventBus *rabbitmq.Client
}

func NewOrderService(repo IorderRepository, skuSvc SkuProvider, invSvc InventoryService, db *gorm.DB, eventBus *rabbitmq.Client) *OrderService {
	return &OrderService{repo: repo, skuSvc: skuSvc, invSvc: invSvc, db: db, eventBus: eventBus}
}

type OrderListResult struct {
	Total int64    `json:"total"`
	List  []*Order `json:"list"`
}

func (s *OrderService) Create(ctx context.Context, userID int64, req *CreateOrderReq) (*Order, error) {
	type skuWithQty struct {
		SkuInfo
		Quantity int
	}
	var items []skuWithQty
	for _, item := range req.Items {
		sku, err := s.skuSvc.FindByID(ctx, item.SkuID)
		if err != nil {
			return nil, fmt.Errorf("sku not found: %d", item.SkuID)
		}
		items = append(items, skuWithQty{SkuInfo: sku, Quantity: item.Quantity})
	}

	var totalAmount int64
	for _, item := range items {
		totalAmount += item.GetPrice() * int64(item.Quantity)
	}
	orderNo := fmt.Sprintf("ORD%d%d", time.Now().UnixMilli(), userID%10000)

	order := &Order{
		OrderNo:       orderNo,
		UserID:        userID,
		TotalAmount:   totalAmount,
		PayAmount:     totalAmount,
		Status:        "pending",
		PaymentStatus: "unpaid",
		Consignee:     req.Address.Consignee,
		Phone:         req.Address.Phone,
		Province:      req.Address.Province,
		City:          req.Address.City,
		District:      req.Address.District,
		DetailAddr:    req.Address.DetailAddr,
		ZipCode:       req.Address.ZipCode,
		CouponID:      req.CouponID,
		BuyerRemark:   req.BuyerRemark,
		Source:        req.Source,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateOrder(ctx, order); err != nil {
			return err
		}
		for _, item := range items {
			if err := s.repo.CreateItem(ctx, &OrderItem{
				OrderID:   order.ID,
				OrderNo:   orderNo,
				SkuID:     item.GetID(),
				ProductID: item.GetProductID(),
				SkuCode:   item.GetSkuCode(),
				SkuSpec:   item.GetSpecJSON(),
				Image:     item.GetImage(),
				Price:     item.GetPrice(),
				Quantity:  item.Quantity,
				Subtotal:  item.GetPrice() * int64(item.Quantity),
			}); err != nil {
				return err
			}
		}
		s.repo.CreateLog(ctx, &OrderLog{OrderID: order.ID, OrderNo: orderNo, ToStatus: "pending", Note: "订单创建"})
		for _, item := range items {
			if err := s.invSvc.Lock(ctx, item.GetID(), item.Quantity); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	return s.repo.FindByOrderNo(ctx, orderNo)
}

func (s *OrderService) List(ctx context.Context, req *OrderListReq) (*OrderListResult, error) {
	req.Normalize()
	list, total, err := s.repo.List(ctx, req.UserID, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*Order, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &OrderListResult{Total: total, List: items}, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderNo string, req *UpdateOrderStatusReq) (*Order, error) {
	order, err := s.repo.FindByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	if !isValidTransition(order.Status, req.Status) {
		return nil, fmt.Errorf("invalid status transition: %s -> %s", order.Status, req.Status)
	}
	now := time.Now()
	updates := map[string]interface{}{"status": req.Status}
	switch req.Status {
	case "paid":
		updates["paid_at"], updates["payment_status"] = now, "paid"
	case "shipped":
		updates["shipped_at"] = now
	case "delivered":
		updates["delivered_at"] = now
	case "cancelled":
		updates["closed_at"], updates["payment_status"] = now, "refunded"
	case "completed":
		updates["completed_at"] = now
	}
	s.db.Model(order).Updates(updates)
	s.repo.CreateLog(ctx, &OrderLog{OrderID: order.ID, OrderNo: orderNo, FromStatus: order.Status, ToStatus: req.Status, Note: req.Note})
	order.Status = req.Status

	s.publishOrderEvent(ctx, order, req.Status)
	return order, nil
}

func (s *OrderService) publishOrderEvent(ctx context.Context, order *Order, status string) {
	if s.eventBus == nil {
		return
	}
	switch status {
	case "paid":
		s.eventBus.Publish(ctx, OrderPaidEvent{
			CustomerID:  strconv.FormatInt(order.UserID, 10),
			OrderID:     order.ID,
			TotalAmount: order.PayAmount,
		})
	case "shipped":
		s.eventBus.Publish(ctx, OrderShippedEvent{
			CustomerID: strconv.FormatInt(order.UserID, 10),
			OrderID:    order.ID,
		})
	case "delivered":
		s.eventBus.Publish(ctx, OrderDeliveredEvent{
			CustomerID: strconv.FormatInt(order.UserID, 10),
			OrderID:    order.ID,
		})
	case "cancelled":
		s.eventBus.Publish(ctx, OrderCancelledEvent{
			CustomerID: strconv.FormatInt(order.UserID, 10),
			OrderID:    order.ID,
			UserID:     order.UserID,
		})
	}
}

var validTransitions = map[string][]string{
	"pending":   {"paid", "cancelled"},
	"paid":      {"shipped", "cancelled"},
	"shipped":   {"delivered"},
	"delivered": {"completed"},
}

func isValidTransition(from, to string) bool {
	for _, t := range validTransitions[from] {
		if t == to {
			return true
		}
	}
	return false
}

// ── PaymentService ───────────────────────────────

