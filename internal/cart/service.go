package cart

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// SkuInfo SKU 查询接口（由 product 模块注入实现，避免直接依赖）
type SkuInfo interface {
	GetID() int64
	GetProductID() int64
	GetPrice() int64
	GetImage() string
	GetSpecJSON() string
}

type SkuProvider interface {
	FindByID(ctx context.Context, skuID int64) (SkuInfo, error)
}

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

// AddItem 添加商品到购物车，已存在则累加数量
func (s *CartService) AddItem(ctx context.Context, userID int64, sessionID string, req *AddItemReq) (*CartResponse, error) {
	// 查询 SKU 信息
	sku, err := s.skuProvider.FindByID(ctx, req.SkuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sku not found: %d", req.SkuID)
		}
		return nil, err
	}

	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 尝试读取已有的 sku_spec
	specJSON := sku.GetSpecJSON()
	if specJSON == "" {
		specJSON = "{}"
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		item := &CartItem{
			CartID:      cart.ID,
			SkuID:       sku.GetID(),
			ProductID:   sku.GetProductID(),
			ProductName: "", // product name would need separate lookup
			SkuSpec:     specJSON,
			Image:       sku.GetImage(),
			Price:       sku.GetPrice(),
			Quantity:    req.Quantity,
		}
		// FirstOrCreate 会更新已有记录的 quantity += 1
		var existing CartItem
		err := tx.Where("cart_id = ? AND sku_id = ?", cart.ID, sku.GetID()).First(&existing).Error
		if err == nil {
			// 已存在，累加数量
			item.Quantity = existing.Quantity + req.Quantity
			item.ID = existing.ID
			item.ProductName = existing.ProductName
			return tx.Model(&existing).Updates(map[string]interface{}{
				"quantity": item.Quantity,
				"price":    item.Price,
			}).Error
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		return tx.Create(item).Error
	})
	if err != nil {
		return nil, err
	}

	// 更新购物车汇总
	s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.UpdateSummary(tx, cart.ID)
	})

	return s.GetCart(ctx, userID, sessionID)
}

func (s *CartService) UpdateQuantity(ctx context.Context, userID int64, sessionID string, req *UpdateItemReq) (*CartResponse, error) {
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if req.Quantity <= 0 {
			return s.repo.RemoveItem(tx, cart.ID, req.SkuID)
		}
		item, err := s.repo.FindItem(ctx, cart.ID, req.SkuID)
		if err != nil {
			return err
		}
		return tx.Model(item).Update("quantity", req.Quantity).Error
	})
	if err != nil {
		return nil, err
	}
	s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.UpdateSummary(tx, cart.ID)
	})
	return s.GetCart(ctx, userID, sessionID)
}

func (s *CartService) RemoveItem(ctx context.Context, userID int64, sessionID string, skuID int64) (*CartResponse, error) {
	cart, err := s.repo.FindOrCreate(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.RemoveItem(tx, cart.ID, skuID)
	})
	s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.UpdateSummary(tx, cart.ID)
	})
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
