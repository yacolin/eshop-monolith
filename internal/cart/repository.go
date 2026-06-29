package cart

import (
	"context"

	"gorm.io/gorm"
)

type IcartRepository interface {
	FindOrCreate(ctx context.Context, userID int64, sessionID string) (*Cart, error)
	FindByUserID(ctx context.Context, userID int64) (*Cart, error)
	FindBySessionID(ctx context.Context, sessionID string) (*Cart, error)
	FindItem(ctx context.Context, cartID, skuID int64) (*CartItem, error)
	AddOrUpdateItem(tx *gorm.DB, item *CartItem) error
	RemoveItem(tx *gorm.DB, cartID, skuID int64) error
	ClearItems(tx *gorm.DB, cartID int64) error
	ListItems(ctx context.Context, cartID int64) ([]CartItem, error)
	UpdateSummary(tx *gorm.DB, cartID int64) error
	DeleteCart(ctx context.Context, id int64) error
}

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) IcartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) FindOrCreate(ctx context.Context, userID int64, sessionID string) (*Cart, error) {
	var cart Cart
	// 先按 user_id 查找
	if userID > 0 {
		err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error
		if err == nil {
			return &cart, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	// 再按 session_id 查找
	if sessionID != "" {
		err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&cart).Error
		if err == nil {
			return &cart, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	// 都不存在则创建
	cart = Cart{
		UserID:    userID,
		SessionID: sessionID,
	}
	err := r.db.WithContext(ctx).Create(&cart).Error
	return &cart, err
}

func (r *CartRepository) FindByUserID(ctx context.Context, userID int64) (*Cart, error) {
	var cart Cart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *CartRepository) FindBySessionID(ctx context.Context, sessionID string) (*Cart, error) {
	var cart Cart
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&cart).Error
	return &cart, err
}

func (r *CartRepository) FindItem(ctx context.Context, cartID, skuID int64) (*CartItem, error) {
	var item CartItem
	err := r.db.WithContext(ctx).Where("cart_id = ? AND sku_id = ?", cartID, skuID).First(&item).Error
	return &item, err
}

func (r *CartRepository) AddOrUpdateItem(tx *gorm.DB, item *CartItem) error {
	// 尝试更新已有记录数量
	res := tx.Where("cart_id = ? AND sku_id = ?", item.CartID, item.SkuID).
		Assign(CartItem{Quantity: item.Quantity, Price: item.Price}).
		FirstOrCreate(item)
	return res.Error
}

func (r *CartRepository) RemoveItem(tx *gorm.DB, cartID, skuID int64) error {
	return tx.Where("cart_id = ? AND sku_id = ?", cartID, skuID).Delete(&CartItem{}).Error
}

func (r *CartRepository) ClearItems(tx *gorm.DB, cartID int64) error {
	return tx.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error
}

func (r *CartRepository) ListItems(ctx context.Context, cartID int64) ([]CartItem, error) {
	var items []CartItem
	err := r.db.WithContext(ctx).Where("cart_id = ?", cartID).Order("created_at ASC").Find(&items).Error
	return items, err
}

func (r *CartRepository) UpdateSummary(tx *gorm.DB, cartID int64) error {
	return tx.Exec(`
		UPDATE tx_carts
		SET item_count = (SELECT COUNT(*) FROM tx_cart_items WHERE cart_id = ?),
		    total_amount = (SELECT COALESCE(SUM(price * quantity), 0) FROM tx_cart_items WHERE cart_id = ?)
		WHERE id = ?`, cartID, cartID, cartID).Error
}

func (r *CartRepository) DeleteCart(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Cart{}).Error
}
