package trade

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eshop-monolith/pkg/errcode"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CartService struct {
	repo        IcartRepository
	skuProvider SkuProvider
	db          *gorm.DB
	rdb         *redis.Client

	mu          sync.Mutex
	pendingSync map[int64]*time.Timer // userID → debounce timer
}

func NewCartService(repo IcartRepository, skuProvider SkuProvider, db *gorm.DB, rdb *redis.Client) *CartService {
	return &CartService{
		repo:        repo,
		skuProvider: skuProvider,
		db:          db,
		rdb:         rdb,
		pendingSync: make(map[int64]*time.Timer),
	}
}

// ── Redis CRUD（主存） ─────────────────────────────

func (s *CartService) GetCart(ctx context.Context, userID int64, sessionID string) (*CartResponse, error) {
	if s.rdb != nil && userID > 0 {
		resp, err := readCartFromRedis(ctx, s.rdb, userID)
		if err == nil {
			return resp, nil
		}
	}
	// DB 兜底 + 回填 Redis
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	resp := buildCartResp(cart, items)
	if s.rdb != nil && userID > 0 {
		_ = writeCartToRedis(ctx, s.rdb, userID, resp)
	}
	return resp, nil
}

func (s *CartService) AddItem(ctx context.Context, userID int64, sessionID string, req *AddItemReq) (*CartResponse, error) {
	sku, err := s.skuProvider.FindByID(ctx, req.SkuID)
	if err != nil {
		return nil, fmt.Errorf("sku not found: %d", req.SkuID)
	}
	// 读取当前 Redis 购物车，追加或更新商品
	resp, err := s.loadOrCreateCartRedis(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 计算添加后的总数量
	existingQty := 0
	for _, item := range resp.Items {
		if item.SkuID == req.SkuID {
			existingQty = item.Quantity
			break
		}
	}
	if existingQty+req.Quantity > int(sku.GetAvailableQuantity()) {
		return nil, errcode.ErrInsufficientStock
	}

	found := false
	for i, item := range resp.Items {
		if item.SkuID == req.SkuID {
			resp.Items[i].Quantity += req.Quantity
			resp.Items[i].Price = sku.GetPrice()
			resp.Items[i].Subtotal = resp.Items[i].Price * int64(resp.Items[i].Quantity)
			found = true
			break
		}
	}
	if !found {
		resp.Items = append(resp.Items, CartItemResponse{
			SkuID:       sku.GetID(),
			ProductID:   sku.GetProductID(),
			ProductName: sku.GetProductName(),
			SkuSpec:     sku.GetSpecJSON(),
			Image:       sku.GetImage(),
			Price:       sku.GetPrice(),
			Quantity:    req.Quantity,
			Subtotal:    sku.GetPrice() * int64(req.Quantity),
		})
	}

	resp.recalc()
	if s.rdb != nil && userID > 0 {
		_ = writeCartToRedis(ctx, s.rdb, userID, resp)
	}
	s.debounceSync(userID)
	return resp, nil
}

func (s *CartService) UpdateQuantity(ctx context.Context, userID int64, sessionID string, req *UpdateItemReq) (*CartResponse, error) {
	resp, err := s.loadOrCreateCartRedis(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	if req.Quantity <= 0 {
		for i, item := range resp.Items {
			if item.SkuID == req.SkuID {
				resp.Items = append(resp.Items[:i], resp.Items[i+1:]...)
				break
			}
		}
	} else {
		for i, item := range resp.Items {
			if item.SkuID == req.SkuID {
				resp.Items[i].Quantity = req.Quantity
				resp.Items[i].Subtotal = item.Price * int64(req.Quantity)
				break
			}
		}
	}

	resp.recalc()
	if s.rdb != nil && userID > 0 {
		_ = writeCartToRedis(ctx, s.rdb, userID, resp)
	}
	s.debounceSync(userID)
	return resp, nil
}

func (s *CartService) RemoveItem(ctx context.Context, userID int64, sessionID string, skuID int64) (*CartResponse, error) {
	resp, err := s.loadOrCreateCartRedis(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	for i, item := range resp.Items {
		if item.SkuID == skuID {
			resp.Items = append(resp.Items[:i], resp.Items[i+1:]...)
			break
		}
	}

	resp.recalc()
	if s.rdb != nil && userID > 0 {
		_ = writeCartToRedis(ctx, s.rdb, userID, resp)
	}
	s.debounceSync(userID)
	return resp, nil
}

func (s *CartService) ClearCart(ctx context.Context, userID int64, sessionID string) error {
	if s.rdb != nil && userID > 0 {
		delCartFromRedis(ctx, s.rdb, userID)
	}
	// 同步清 DB
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

// loadOrCreateCartRedis 读 Redis 购物车，不存在则建空车
func (s *CartService) loadOrCreateCartRedis(ctx context.Context, userID int64, sessionID string) (*CartResponse, error) {
	if s.rdb != nil && userID > 0 {
		resp, err := readCartFromRedis(ctx, s.rdb, userID)
		if err == nil {
			return resp, nil
		}
	}
	// Redis miss → 从 DB 读
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	return buildCartResp(cart, items), nil
}

// ── 异步落库 ────────────────────────────────────

// debounceSync 防抖延迟同步：5s 内无新写入再落库
func (s *CartService) debounceSync(userID int64) {
	if s.rdb == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.pendingSync[userID]; ok {
		t.Stop()
	}
	s.pendingSync[userID] = time.AfterFunc(5*time.Second, func() {
		s.syncToDB(context.Background(), userID)
	})
}

func (s *CartService) syncToDB(ctx context.Context, userID int64) {
	defer func() {
		s.mu.Lock()
		delete(s.pendingSync, userID)
		s.mu.Unlock()
	}()

	resp, err := readCartFromRedis(ctx, s.rdb, userID)
	if err != nil {
		return
	}

	cart, err := s.repo.FindOrCreate(ctx, userID, "")
	if err != nil {
		return
	}
	_ = s.db.Transaction(func(tx *gorm.DB) error {
		_ = s.repo.ClearItems(tx, cart.ID)
		for _, item := range resp.Items {
			ci := &CartItem{
				CartID:      cart.ID,
				SkuID:       item.SkuID,
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				Image:       item.Image,
				Price:       item.Price,
				Quantity:    item.Quantity,
			}
			if item.SkuSpec != "" {
				ci.SkuSpec = item.SkuSpec
			}
			if err := tx.Create(ci).Error; err != nil {
				return err
			}
		}
		return nil
	})
	_ = s.repo.UpdateSummary(s.db, cart.ID)
}

// ── Helper ────────────────────────────────────────

func buildCartResp(cart *Cart, items []CartItem) *CartResponse {
	resp := &CartResponse{
		ID:    cart.ID,
		Items: make([]CartItemResponse, 0, len(items)),
	}
	for _, item := range items {
		resp.Items = append(resp.Items, CartItemResponse{
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
	resp.recalc()
	return resp
}

// recalc 从 items 重新计算 item_count / total_amount
func (r *CartResponse) recalc() {
	r.ItemCount = len(r.Items)
	r.TotalAmount = 0
	for _, item := range r.Items {
		r.TotalAmount += item.Subtotal
	}
}
