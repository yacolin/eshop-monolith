package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/internal/cart/api/dto"
	"eshop-monolith/internal/cart/domain/models"
	"eshop-monolith/internal/cart/domain/repositories"
	"eshop-monolith/internal/cart/events"
	"eshop-monolith/internal/infra/rabbitmq"
	invService "eshop-monolith/internal/inventory/service"
)

// CartService 购物车服务
type CartService struct {
	cartRepo         repositories.IcartRepository
	inventoryService *invService.InventoryService
	productService   *invService.ProductService
	skuService       *invService.SkuService
	rabbit           *rabbitmq.Client
}

// NewCartService 创建购物车服务实例
func NewCartService(cartRepo repositories.IcartRepository, inventoryService *invService.InventoryService, productService *invService.ProductService, skuService *invService.SkuService, rabbit *rabbitmq.Client) *CartService {
	return &CartService{
		cartRepo:         cartRepo,
		inventoryService: inventoryService,
		productService:   productService,
		skuService:       skuService,
		rabbit:           rabbit,
	}
}

// GetCart 获取购物车
func (s *CartService) GetCart(ctx context.Context, userID int64, sessionID string) (*dto.CartResponse, error) {
	var cart *models.Cart
	var err error

	// 优先根据用户ID获取购物车
	if userID > 0 {
		cart, err = s.cartRepo.GetByUserID(ctx, userID)
	} else if sessionID != "" {
		// 否则根据会话ID获取购物车
		cart, err = s.cartRepo.GetBySessionID(ctx, sessionID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 购物车不存在，创建一个新的
			newCart := &models.Cart{
				UserID:    userID,
				SessionID: sessionID,
			}
			err = s.cartRepo.Create(ctx, newCart)
			if err != nil {
				return nil, err
			}
			cart = newCart
		} else {
			return nil, err
		}
	}

	// 转换为响应格式
	return s.toCartResponse(ctx, cart)
}

// AddToCart 添加商品到购物车
func (s *CartService) AddToCart(ctx context.Context, userID int64, sessionID string, req *dto.AddToCartDTO) (*dto.CartResponse, error) {
	// 获取或创建购物车
	cart, err := s.getOrCreateCart(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 验证 SKU 是否存在
	sku, err := s.skuService.GetSku(ctx, req.SkuID)
	if err != nil {
		return nil, err
	}

	// 检查库存
	stock, err := s.inventoryService.GetInventoryBySkuID(ctx, req.SkuID)
	if err != nil {
		return nil, err
	}
	if stock.Quantity < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// 检查购物车中是否已有相同 SKU
	existingItem, err := s.cartRepo.GetItemByCartAndSku(ctx, cart.ID, req.SkuID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingItem != nil {
		// 已有该 SKU，更新数量
		existingItem.Quantity += req.Quantity
		err = s.cartRepo.UpdateItem(ctx, existingItem)
		if err != nil {
			return nil, err
		}

		// 发布购物车项更新事件
		s.rabbit.Publish(ctx,events.CartItemUpdatedEvent{
			CartID:      cart.ID,
			UserID:      cart.UserID,
			ItemID:      existingItem.ID,
			SkuID:       existingItem.SkuID,
			ProductID:   existingItem.ProductID,
			OldQuantity: existingItem.Quantity - req.Quantity,
			NewQuantity: existingItem.Quantity,
			Price:       existingItem.Price,
		})
	} else {
		// 新增购物车项
		newItem := &models.CartItem{
			CartID:    cart.ID,
			SkuID:     req.SkuID,
			ProductID: sku.ProductID,
			Quantity:  req.Quantity,
			Price:     sku.Price,
			SKU:       sku.SKUCode,
		}
		err = s.cartRepo.AddItem(ctx, newItem)
		if err != nil {
			return nil, err
		}

		// 发布购物车项添加事件
		s.rabbit.Publish(ctx,events.CartItemAddedEvent{
			CartID:    cart.ID,
			UserID:    cart.UserID,
			SkuID:     newItem.SkuID,
			ProductID: newItem.ProductID,
			Quantity:  newItem.Quantity,
			Price:     newItem.Price,
		})
	}

	// 重新获取购物车详情
	updatedCart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		updatedCart, err = s.cartRepo.GetBySessionID(ctx, sessionID)
		if err != nil {
			return nil, err
		}
	}

	// 转换为响应格式
	return s.toCartResponse(ctx, updatedCart)
}

// UpdateCartItem 更新购物车项
func (s *CartService) UpdateCartItem(ctx context.Context, userID int64, sessionID string, itemID int64, req *dto.UpdateCartItemDTO) (*dto.CartResponse, error) {
	// 获取购物车
	cart, err := s.getCart(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 查找购物车项
	var targetItem *models.CartItem
	for i := range cart.Items {
		if cart.Items[i].ID == itemID {
			targetItem = &cart.Items[i]
			break
		}
	}

	if targetItem == nil {
		return nil, errors.New("cart item not found")
	}

	// 检查库存
	stock, err := s.inventoryService.GetInventoryBySkuID(ctx, targetItem.SkuID)
	if err != nil {
		return nil, err
	}

	if stock.Quantity < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// 保存旧数量
	oldQuantity := targetItem.Quantity

	// 更新数量
	targetItem.Quantity = req.Quantity
	err = s.cartRepo.UpdateItem(ctx, targetItem)
	if err != nil {
		return nil, err
	}

	// 发布购物车项更新事件
	s.rabbit.Publish(ctx,events.CartItemUpdatedEvent{
		CartID:      cart.ID,
		UserID:      cart.UserID,
		ItemID:      itemID,
		SkuID:       targetItem.SkuID,
		ProductID:   targetItem.ProductID,
		OldQuantity: oldQuantity,
		NewQuantity: req.Quantity,
		Price:       targetItem.Price,
	})

	// 重新获取购物车详情
	updatedCart, err := s.getCart(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	return s.toCartResponse(ctx, updatedCart)
}

// RemoveCartItem 删除购物车项
func (s *CartService) RemoveCartItem(ctx context.Context, userID int64, sessionID string, itemID int64) (*dto.CartResponse, error) {
	// 获取购物车
	cart, err := s.getCart(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 查找购物车项
	var targetItem *models.CartItem
	for i := range cart.Items {
		if cart.Items[i].ID == itemID {
			targetItem = &cart.Items[i]
			break
		}
	}

	if targetItem == nil {
		return nil, errors.New("cart item not found")
	}

	// 保存要删除的项信息
	removedItem := *targetItem

	// 删除购物车项
	err = s.cartRepo.DeleteItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// 发布购物车项移除事件
	s.rabbit.Publish(ctx,events.CartItemRemovedEvent{
		CartID:    cart.ID,
		UserID:    cart.UserID,
		ItemID:    itemID,
		SkuID:     removedItem.SkuID,
		ProductID: removedItem.ProductID,
		Quantity:  removedItem.Quantity,
		Price:     removedItem.Price,
	})

	// 重新获取购物车详情
	updatedCart, err := s.getCart(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	return s.toCartResponse(ctx, updatedCart)
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(ctx context.Context, userID int64, sessionID string) error {
	// 获取购物车
	cart, err := s.getCart(ctx, userID, sessionID)
	if err != nil {
		return err
	}

	// 删除所有购物车项
	for _, item := range cart.Items {
		err = s.cartRepo.DeleteItem(ctx, item.ID)
		if err != nil {
			return err
		}
	}

	// 发布购物车清空事件
	s.rabbit.Publish(ctx,events.CartClearedEvent{
		CartID: cart.ID,
		UserID: cart.UserID,
	})

	return nil
}

// getOrCreateCart 获取或创建购物车
func (s *CartService) getOrCreateCart(ctx context.Context, userID int64, sessionID string) (*models.Cart, error) {
	var cart *models.Cart
	var err error

	// 优先根据用户ID获取购物车
	if userID > 0 {
		cart, err = s.cartRepo.GetByUserID(ctx, userID)
	} else if sessionID != "" {
		// 否则根据会话ID获取购物车
		cart, err = s.cartRepo.GetBySessionID(ctx, sessionID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 购物车不存在，创建一个新的
			newCart := &models.Cart{
				UserID:    userID,
				SessionID: sessionID,
			}
			err = s.cartRepo.Create(ctx, newCart)
			if err != nil {
				return nil, err
			}
			return newCart, nil
		} else {
			return nil, err
		}
	}

	return cart, nil
}

// getCart 获取购物车
func (s *CartService) getCart(ctx context.Context, userID int64, sessionID string) (*models.Cart, error) {
	var cart *models.Cart
	var err error

	// 优先根据用户ID获取购物车
	if userID > 0 {
		cart, err = s.cartRepo.GetByUserID(ctx, userID)
	} else if sessionID != "" {
		// 否则根据会话ID获取购物车
		cart, err = s.cartRepo.GetBySessionID(ctx, sessionID)
	}

	if err != nil {
		return nil, err
	}

	return cart, nil
}

// toCartResponse 将购物车模型转换为响应格式
func (s *CartService) toCartResponse(ctx context.Context, cart *models.Cart) (*dto.CartResponse, error) {
	response := &dto.CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		Items:      make([]dto.CartItemResponse, 0, len(cart.Items)),
		TotalItems: 0,
		TotalPrice: 0,
	}

	for _, item := range cart.Items {
		sku, err := s.skuService.GetSku(ctx, item.SkuID)
		if err != nil {
			continue
		}
		product, err := s.productService.GetProductByID(ctx, sku.ProductID)
		if err != nil {
			continue
		}
		stock, err := s.inventoryService.GetInventoryBySkuID(ctx, item.SkuID)
		if err != nil {
			continue
		}
		itemResponse := dto.CartItemResponse{
			ID:          item.ID,
			SkuID:       item.SkuID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			Price:       item.Price,
			SKU:         item.SKU,
			ProductName: product.Name,
			Stock:       stock.Quantity,
		}
		response.Items = append(response.Items, itemResponse)
		response.TotalItems += item.Quantity
		response.TotalPrice += item.Price * int64(item.Quantity)
	}
	return response, nil
}
