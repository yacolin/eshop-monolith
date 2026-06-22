package repositories

import (
	"context"

	"gorm.io/gorm"

	"eshop-monolith/internal/cart/api/dto"
	cartModels "eshop-monolith/internal/cart/domain/models"
	"eshop-monolith/internal/infra/repository/models"
)

// IcartRepository 购物车仓储接口
type IcartRepository interface {
	// GetByUserID 根据用户ID获取购物车
	GetByUserID(ctx context.Context, userID int64) (*cartModels.Cart, error)

	// GetBySessionID 根据会话ID获取购物车
	GetBySessionID(ctx context.Context, sessionID string) (*cartModels.Cart, error)

	// Create 创建购物车
	Create(ctx context.Context, cart *cartModels.Cart) error

	// Update 更新购物车
	Update(ctx context.Context, cart *cartModels.Cart) error

	// Delete 删除购物车
	Delete(ctx context.Context, id int64) error

	// ListByQuery 根据查询参数获取购物车列表
	ListByQuery(ctx context.Context, q dto.CartListQuery, offset, limit int) ([]cartModels.Cart, error)

	// CountByQuery 根据查询参数统计购物车数量
	CountByQuery(ctx context.Context, q dto.CartListQuery) (int64, error)

	// AddItem 添加购物车项
	AddItem(ctx context.Context, item *cartModels.CartItem) error

	// UpdateItem 更新购物车项
	UpdateItem(ctx context.Context, item *cartModels.CartItem) error

	// DeleteItem 删除购物车项
	DeleteItem(ctx context.Context, id int64) error

	// GetItemByCartAndSku 根据购物车ID和SKU ID获取购物车项
	GetItemByCartAndSku(ctx context.Context, cartID, skuID int64) (*cartModels.CartItem, error)
}

// CartRepository 购物车仓储实现
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository 创建购物车仓储实例
func NewCartRepository(db *gorm.DB) IcartRepository {
	return &CartRepository{db: db}
}

// GetByUserID 根据用户ID获取购物车
func (r *CartRepository) GetByUserID(ctx context.Context, userID int64) (*cartModels.Cart, error) {
	var po models.CartPO
	err := r.db.WithContext(ctx).First(&po, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	// 加载购物车项
	err = r.db.WithContext(ctx).Find(&po.Items, "cart_id = ?", po.ID).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// GetBySessionID 根据会话ID获取购物车
func (r *CartRepository) GetBySessionID(ctx context.Context, sessionID string) (*cartModels.Cart, error) {
	var po models.CartPO
	err := r.db.WithContext(ctx).First(&po, "session_id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	// 加载购物车项
	err = r.db.WithContext(ctx).Find(&po.Items, "cart_id = ?", po.ID).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// Create 创建购物车
func (r *CartRepository) Create(ctx context.Context, cart *cartModels.Cart) error {
	po := models.CartFromDomain(cart)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	cart.ID = po.ID
	return nil
}

// Update 更新购物车
func (r *CartRepository) Update(ctx context.Context, cart *cartModels.Cart) error {
	po := models.CartFromDomain(cart)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除购物车
func (r *CartRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CartPO{}, id).Error
}

// ListByQuery 根据查询参数获取购物车列表
func (r *CartRepository) ListByQuery(ctx context.Context, q dto.CartListQuery, offset, limit int) ([]cartModels.Cart, error) {
	var pos []models.CartPO
	// 构建查询条件
	db := r.db.WithContext(ctx)
	if q.UserID > 0 {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.SessionID != "" {
		db = db.Where("session_id = ?", q.SessionID)
	}
	// 执行查询
	err := db.Offset(offset).Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	// 加载购物车项
	for i := range pos {
		err = r.db.WithContext(ctx).Find(&pos[i].Items, "cart_id = ?", pos[i].ID).Error
		if err != nil {
			return nil, err
		}
	}

	carts := make([]cartModels.Cart, len(pos))
	for i, po := range pos {
		carts[i] = *po.ToDomain()
	}
	return carts, nil
}

// CountByQuery 根据查询参数统计购物车数量
func (r *CartRepository) CountByQuery(ctx context.Context, q dto.CartListQuery) (int64, error) {
	var count int64
	// 构建查询条件
	db := r.db.WithContext(ctx).Model(&models.CartPO{})
	if q.UserID > 0 {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.SessionID != "" {
		db = db.Where("session_id = ?", q.SessionID)
	}
	// 执行计数
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// AddItem 添加购物车项
func (r *CartRepository) AddItem(ctx context.Context, item *cartModels.CartItem) error {
	po := models.CartItemFromDomain(item)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	item.ID = po.ID
	return nil
}

// UpdateItem 更新购物车项
func (r *CartRepository) UpdateItem(ctx context.Context, item *cartModels.CartItem) error {
	po := models.CartItemFromDomain(item)
	return r.db.WithContext(ctx).Save(po).Error
}

// DeleteItem 删除购物车项
func (r *CartRepository) DeleteItem(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CartItemPO{}, id).Error
}

// GetItemByCartAndSku 根据购物车ID和SKU ID获取购物车项
func (r *CartRepository) GetItemByCartAndSku(ctx context.Context, cartID, skuID int64) (*cartModels.CartItem, error) {
	var po models.CartItemPO
	err := r.db.WithContext(ctx).Where("cart_id = ? AND sku_id = ?", cartID, skuID).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}
