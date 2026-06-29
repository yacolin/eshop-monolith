package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

// ── 外部依赖接口 ─────────────────────────────────

type SkuInfo interface {
	GetID() int64
	GetProductID() int64
	GetSkuCode() string
	GetPrice() int64
	GetImage() string
	GetSpecJSON() string
}

type SkuService interface {
	FindByID(ctx context.Context, skuID int64) (SkuInfo, error)
}

type InventoryService interface {
	Lock(ctx context.Context, skuID int64, quantity int) error
	Unlock(ctx context.Context, skuID int64, quantity int) error
}

// ── OrderService ─────────────────────────────────

type OrderService struct {
	repo      IorderRepository
	skuSvc    SkuService
	invSvc    InventoryService
	db        *gorm.DB
}

func NewOrderService(repo IorderRepository, skuSvc SkuService, invSvc InventoryService, db *gorm.DB) *OrderService {
	return &OrderService{repo: repo, skuSvc: skuSvc, invSvc: invSvc, db: db}
}

type OrderListResult struct {
	Total int64    `json:"total"`
	List  []*Order `json:"list"`
}

// Create 创建订单（事务：写订单 → 写明细 → 写日志 → 锁库存）
func (s *OrderService) Create(ctx context.Context, userID int64, req *CreateOrderReq) (*Order, error) {
	// 查询所有 SKU 信息
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

	// 计算金额
	var totalAmount int64
	for _, item := range items {
		totalAmount += item.GetPrice() * int64(item.Quantity)
	}
	payAmount := totalAmount

	// 生成订单号
	orderNo := fmt.Sprintf("ORD%d%d", time.Now().UnixMilli(), userID%10000)

	order := &Order{
		OrderNo:       orderNo,
		UserID:        userID,
		TotalAmount:   totalAmount,
		PayAmount:     payAmount,
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

	// 事务创建
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单
		if err := s.repo.CreateOrder(ctx, order); err != nil {
			return err
		}
		// 2. 创建订单明细
		for _, item := range items {
			oi := &OrderItem{
				OrderID:     order.ID,
				OrderNo:     orderNo,
				SkuID:       item.GetID(),
				ProductID:   item.GetProductID(),
				SkuCode:     item.GetSkuCode(),
				ProductName: "",
				SkuSpec:     item.GetSpecJSON(),
				Image:       item.GetImage(),
				Price:       item.GetPrice(),
				Quantity:    item.Quantity,
				Subtotal:    item.GetPrice() * int64(item.Quantity),
			}
			if err := s.repo.CreateItem(ctx, oi); err != nil {
				return err
			}
		}
		// 3. 创建订单日志
		log := &OrderLog{
			OrderID: order.ID,
			OrderNo: orderNo,
			ToStatus: "pending",
			Note:    "订单创建",
		}
		if err := s.repo.CreateLog(ctx, log); err != nil {
			return err
		}
		// 4. 锁库存
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
	order, err := s.repo.FindByOrderNo(ctx, orderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
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

// UpdateStatus 更新订单状态（含状态机校验）
func (s *OrderService) UpdateStatus(ctx context.Context, orderNo string, req *UpdateOrderStatusReq) (*Order, error) {
	order, err := s.repo.FindByOrderNo(ctx, orderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrOrderNotFound
		}
		return nil, err
	}

	// 状态机校验
	if !isValidTransition(order.Status, req.Status) {
		return nil, errcode.ErrInvalidOrderStatus
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status": req.Status,
	}
	switch req.Status {
	case "paid":
		updates["paid_at"] = now
		updates["payment_status"] = "paid"
	case "shipped":
		updates["shipped_at"] = now
	case "delivered":
		updates["delivered_at"] = now
	case "cancelled":
		updates["closed_at"] = now
		updates["payment_status"] = "refunded"
	case "completed":
		updates["completed_at"] = now
	}

	if err := s.db.Model(order).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 记录日志
	s.repo.CreateLog(ctx, &OrderLog{
		OrderID:  order.ID,
		OrderNo:  orderNo,
		FromStatus: order.Status,
		ToStatus: req.Status,
		Note:     req.Note,
	})

	order.Status = req.Status
	return order, nil
}

// 状态机转换规则
var validTransitions = map[string][]string{
	"pending":     {"paid", "cancelled"},
	"paid":        {"shipped", "cancelled"},
	"shipped":     {"delivered"},
	"delivered":   {"completed"},
	"completed":   {},
	"cancelled":   {},
	"closed":      {},
}

func isValidTransition(from, to string) bool {
	targets, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}
