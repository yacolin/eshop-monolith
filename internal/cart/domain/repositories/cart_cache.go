package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/internal/cart/api/dto"
	cartModels "eshop-monolith/internal/cart/domain/models"
	"eshop-monolith/internal/infra/repository/models"
)

const (
	cartCacheTTL = 7 * 24 * time.Hour
)

// CachedCartRepository 购物车缓存仓储，装饰模式包裹 CartRepository
// 读路径：cache-aside，先查 Redis ⇒ 未命中则回源 DB ⇒ 回填缓存
// 写路径：先写 DB ⇒ 删除 data 缓存（索引 key 不删，下次读自动刷新）
type CachedCartRepository struct {
	inner IcartRepository
	rdb   *redis.Client
	db    *gorm.DB // 仅用于 DeleteItem 时查找 cart_id
}

// NewCachedCartRepository 创建购物车缓存仓储
// rdb 为 nil 时跳过缓存，直接委托给 inner，方便本地开发/测试
func NewCachedCartRepository(inner IcartRepository, rdb *redis.Client, db *gorm.DB) IcartRepository {
	if rdb == nil {
		return inner
	}
	return &CachedCartRepository{
		inner: inner,
		rdb:   rdb,
		db:    db,
	}
}

func cartDataKey(cartID int64) string         { return fmt.Sprintf("cart:data:%d", cartID) }
func cartUserKey(userID int64) string          { return fmt.Sprintf("cart:user:%d", userID) }
func cartSessKey(sessionID string) string      { return fmt.Sprintf("cart:sess:%s", sessionID) }

// getCartFromCache 从 Redis 读取缓存的购物车完整数据
func (r *CachedCartRepository) getCartFromCache(ctx context.Context, cartID int64) (*cartModels.Cart, error) {
	data, err := r.rdb.Get(ctx, cartDataKey(cartID)).Bytes()
	if err != nil {
		return nil, err
	}
	var po models.CartPO
	if err := json.Unmarshal(data, &po); err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// setCartCache 将购物车完整数据写入 Redis，同时建立 user/session → cart_id 索引
func (r *CachedCartRepository) setCartCache(ctx context.Context, cart *cartModels.Cart) {
	po := models.CartFromDomain(cart)
	data, err := json.Marshal(po)
	if err != nil {
		return
	}

	pipe := r.rdb.Pipeline()
	pipe.Set(ctx, cartDataKey(cart.ID), data, cartCacheTTL)
	if cart.UserID > 0 {
		pipe.Set(ctx, cartUserKey(cart.UserID), cart.ID, cartCacheTTL)
	}
	if cart.SessionID != "" {
		pipe.Set(ctx, cartSessKey(cart.SessionID), cart.ID, cartCacheTTL)
	}
	pipe.Exec(ctx)
}

// evictCartCache 删除 data 缓存（索引 key 保留，指向的 cart_id 不变）
func (r *CachedCartRepository) evictCartCache(ctx context.Context, cartID int64) {
	r.rdb.Del(ctx, cartDataKey(cartID))
}

// GetByUserID 根据用户 ID 获取购物车（缓存支持）
func (r *CachedCartRepository) GetByUserID(ctx context.Context, userID int64) (*cartModels.Cart, error) {
	cartID, err := r.rdb.Get(ctx, cartUserKey(userID)).Int64()
	if err == nil {
		cart, err := r.getCartFromCache(ctx, cartID)
		if err == nil {
			return cart, nil
		}
	}

	cart, err := r.inner.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	r.setCartCache(ctx, cart)
	return cart, nil
}

// GetBySessionID 根据会话 ID 获取购物车（缓存支持）
func (r *CachedCartRepository) GetBySessionID(ctx context.Context, sessionID string) (*cartModels.Cart, error) {
	cartID, err := r.rdb.Get(ctx, cartSessKey(sessionID)).Int64()
	if err == nil {
		cart, err := r.getCartFromCache(ctx, cartID)
		if err == nil {
			return cart, nil
		}
	}

	cart, err := r.inner.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	r.setCartCache(ctx, cart)
	return cart, nil
}

// Create 创建购物车（不缓存，新购物车无商品）
func (r *CachedCartRepository) Create(ctx context.Context, cart *cartModels.Cart) error {
	return r.inner.Create(ctx, cart)
}

// Update 更新购物车
func (r *CachedCartRepository) Update(ctx context.Context, cart *cartModels.Cart) error {
	if err := r.inner.Update(ctx, cart); err != nil {
		return err
	}
	r.evictCartCache(ctx, cart.ID)
	return nil
}

// Delete 删除购物车
func (r *CachedCartRepository) Delete(ctx context.Context, id int64) error {
	if err := r.inner.Delete(ctx, id); err != nil {
		return err
	}
	r.evictCartCache(ctx, id)
	return nil
}

// ListByQuery 分页查询购物车列表（穿透到 DB，不适合缓存）
func (r *CachedCartRepository) ListByQuery(ctx context.Context, q dto.CartListQuery, offset, limit int) ([]cartModels.Cart, error) {
	return r.inner.ListByQuery(ctx, q, offset, limit)
}

// CountByQuery 统计购物车数量（穿透到 DB）
func (r *CachedCartRepository) CountByQuery(ctx context.Context, q dto.CartListQuery) (int64, error) {
	return r.inner.CountByQuery(ctx, q)
}

// AddItem 添加购物车项
func (r *CachedCartRepository) AddItem(ctx context.Context, item *cartModels.CartItem) error {
	if err := r.inner.AddItem(ctx, item); err != nil {
		return err
	}
	r.evictCartCache(ctx, item.CartID)
	return nil
}

// UpdateItem 更新购物车项
func (r *CachedCartRepository) UpdateItem(ctx context.Context, item *cartModels.CartItem) error {
	if err := r.inner.UpdateItem(ctx, item); err != nil {
		return err
	}
	r.evictCartCache(ctx, item.CartID)
	return nil
}

// DeleteItem 删除购物车项
func (r *CachedCartRepository) DeleteItem(ctx context.Context, id int64) error {
	// 先查找 cart_id 用于缓存失效
	var cartID int64
	if err := r.db.WithContext(ctx).Table("cart_items").
		Select("cart_id").Where("id = ?", id).Scan(&cartID).Error; err != nil {
		return err
	}

	if err := r.inner.DeleteItem(ctx, id); err != nil {
		return err
	}
	r.evictCartCache(ctx, cartID)
	return nil
}

// GetItemByCartAndSku 根据购物车和SKU获取购物车项（优先走缓存）
func (r *CachedCartRepository) GetItemByCartAndSku(ctx context.Context, cartID, skuID int64) (*cartModels.CartItem, error) {
	cart, err := r.getCartFromCache(ctx, cartID)
	if err == nil {
		for _, item := range cart.Items {
			if item.SkuID == skuID {
				return &item, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}

	return r.inner.GetItemByCartAndSku(ctx, cartID, skuID)
}
