package trade

import (
	"context"
	"fmt"
	"time"

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

type SkuProvider interface {
	FindByID(ctx context.Context, skuID int64) (SkuInfo, error)
}

type InventoryService interface {
	Lock(ctx context.Context, skuID int64, quantity int) error
	Unlock(ctx context.Context, skuID int64, quantity int) error
}

// ── CartService ──────────────────────────────────

type CartService struct {
	repo       IcartRepository
	skuProvider SkuProvider
	db         *gorm.DB
}

func NewCartService(repo IcartRepository, skuProvider SkuProvider, db *gorm.DB) *CartService {
	return &CartService{repo: repo, skuProvider: skuProvider, db: db}
}

func (s *CartService) GetCart(ctx context.Context, userID int64, sessionID string) (*CartResponse, error) {
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	resp := &CartResponse{
		ID:          cart.ID,
		ItemCount:   cart.ItemCount,
		TotalAmount: cart.TotalAmount,
		Items:       make([]CartItemResponse, 0, len(items)),
	}
	for _, item := range items {
		resp.Items = append(resp.Items, CartItemResponse{
			ID:          item.ID,
			SkuID:       item.SkuID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			SkuSpec:     item.SkuSpec,
			Image:       item.Image,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Price * int64(item.Quantity),
		})
	}
	return resp, nil
}

func (s *CartService) AddItem(ctx context.Context, userID int64, sessionID string, req *AddItemReq) (*CartResponse, error) {
	sku, err := s.skuProvider.FindByID(ctx, req.SkuID)
	if err != nil {
		return nil, fmt.Errorf("sku not found: %d", req.SkuID)
	}
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var existing CartItem
		err := tx.Where("cart_id = ? AND sku_id = ?", cart.ID, sku.GetID()).First(&existing).Error
		if err == nil {
			existing.Quantity += req.Quantity
			existing.Price = sku.GetPrice()
			return tx.Model(&existing).Updates(map[string]interface{}{
				"quantity": existing.Quantity,
				"price":    existing.Price,
			}).Error
		}
		return tx.Create(&CartItem{
			CartID:    cart.ID,
			SkuID:     sku.GetID(),
			ProductID: sku.GetProductID(),
			SkuSpec:   sku.GetSpecJSON(),
			Image:     sku.GetImage(),
			Price:     sku.GetPrice(),
			Quantity:  req.Quantity,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	s.db.Transaction(func(tx *gorm.DB) error { return s.repo.UpdateSummary(tx, cart.ID) })
	return s.GetCart(ctx, userID, sessionID)
}

func (s *CartService) UpdateQuantity(ctx context.Context, userID int64, sessionID string, req *UpdateItemReq) (*CartResponse, error) {
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	s.db.Transaction(func(tx *gorm.DB) error {
		if req.Quantity <= 0 {
			return s.repo.RemoveItem(tx, cart.ID, req.SkuID)
		}
		return tx.Model(&CartItem{}).Where("cart_id = ? AND sku_id = ?", cart.ID, req.SkuID).Update("quantity", req.Quantity).Error
	})
	s.db.Transaction(func(tx *gorm.DB) error { return s.repo.UpdateSummary(tx, cart.ID) })
	return s.GetCart(ctx, userID, sessionID)
}

func (s *CartService) RemoveItem(ctx context.Context, userID int64, sessionID string, skuID int64) (*CartResponse, error) {
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	s.db.Transaction(func(tx *gorm.DB) error { return s.repo.RemoveItem(tx, cart.ID, skuID) })
	s.db.Transaction(func(tx *gorm.DB) error { return s.repo.UpdateSummary(tx, cart.ID) })
	return s.GetCart(ctx, userID, sessionID)
}

func (s *CartService) ClearCart(ctx context.Context, userID int64, sessionID string) error {
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.ClearItems(tx, cart.ID); err != nil {
			return err
		}
		return s.repo.UpdateSummary(tx, cart.ID)
	})
}

// ── OrderService ─────────────────────────────────

type OrderService struct {
	repo     IorderRepository
	skuSvc   SkuProvider
	invSvc   InventoryService
	db       *gorm.DB
}

func NewOrderService(repo IorderRepository, skuSvc SkuProvider, invSvc InventoryService, db *gorm.DB) *OrderService {
	return &OrderService{repo: repo, skuSvc: skuSvc, invSvc: invSvc, db: db}
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
	return order, nil
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

type PaymentService struct {
	repo IpaymentRepository
	db   *gorm.DB
}

func NewPaymentService(repo IpaymentRepository, db *gorm.DB) *PaymentService {
	return &PaymentService{repo: repo, db: db}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentReq) (*Payment, error) {
	payment := &Payment{
		PaymentNo:     fmt.Sprintf("PAY%d", time.Now().UnixNano()),
		OrderNo:       req.OrderNo,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Channel:       req.Channel,
		Status:        "pending",
	}
	if req.OrderType != "" {
		payment.OrderType = req.OrderType
	}
	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}
	s.repo.CreateLog(ctx, &PaymentLog{PaymentID: payment.ID, PaymentNo: payment.PaymentNo, Action: "create", Status: "pending"})
	return payment, nil
}

func (s *PaymentService) HandleCallback(ctx context.Context, req *PaymentCallbackReq) (*Payment, error) {
	payment, err := s.repo.FindByPaymentNo(ctx, req.PaymentNo)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	updates := map[string]interface{}{"status": req.Status, "transaction_id": req.TransactionID, "failure_reason": req.FailureReason}
	if req.Status == "success" {
		updates["paid_at"] = &now
	}
	s.db.Model(payment).Updates(updates)
	payment.Status = req.Status
	payment.TransactionID = req.TransactionID
	s.repo.CreateLog(ctx, &PaymentLog{
		PaymentID: payment.ID, PaymentNo: payment.PaymentNo, Channel: req.Channel,
		TransactionID: req.TransactionID, Action: "pay_callback", RequestBody: req.RawBody, Status: req.Status,
	})
	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, paymentNo string) (*Payment, error) {
	return s.repo.FindByPaymentNo(ctx, paymentNo)
}

func (s *PaymentService) CreateRefund(ctx context.Context, req *CreateRefundReq) (*Refund, error) {
	payment, err := s.repo.FindByPaymentNo(ctx, req.PaymentNo)
	if err != nil {
		return nil, err
	}
	refund := &Refund{
		RefundNo:    fmt.Sprintf("REF%d", time.Now().UnixNano()),
		PaymentNo:   payment.PaymentNo,
		OrderNo:     payment.OrderNo,
		OrderID:     payment.OrderID,
		Amount:      req.Amount,
		Reason:      req.Reason,
		Status:      "pending",
		AppliedAt:   timePtr(time.Now()),
	}
	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return nil, err
	}
	return refund, nil
}

func timePtr(t time.Time) *time.Time { return &t }
