package trade

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ── 外部依赖接口 ─────────────────────────────────

type CartService struct {
	repo        IcartRepository
	skuProvider SkuProvider
	db          *gorm.DB
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
