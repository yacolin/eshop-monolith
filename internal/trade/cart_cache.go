package trade

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

const (
	cartCacheTTL = 2 * time.Hour
)

// ── Key Builders ──

func cacheKeyCartItems(userID int64) string { return fmt.Sprintf("cart:items:%d", userID) }
func cacheKeyCartMeta(userID int64) string  { return fmt.Sprintf("cart:meta:%d", userID) }

// ── CartItem Cache ──

type cartMetaCache struct {
	ItemCount   int   `json:"item_count"`
	TotalAmount int64 `json:"total_amount"`
}

func getCartCache(ctx context.Context, rdb redis.UniversalClient, userID int64) (*cartMetaCache, map[int64]*CartItem, error) {
	pipe := rdb.Pipeline()
	metaCmd := pipe.Get(ctx, cacheKeyCartMeta(userID))
	itemsCmd := pipe.HGetAll(ctx, cacheKeyCartItems(userID))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Meta
	metaData, err := metaCmd.Bytes()
	if err != nil {
		return nil, nil, err
	}
	var meta cartMetaCache
	if err := sonic.Unmarshal(metaData, &meta); err != nil {
		return nil, nil, err
	}

	// Items
	items := make(map[int64]*CartItem, len(itemsCmd.Val()))
	for skuIDStr, itemJSON := range itemsCmd.Val() {
		skuID, _ := strconv.ParseInt(skuIDStr, 10, 64)
		var item CartItem
		if err := sonic.Unmarshal([]byte(itemJSON), &item); err != nil {
			continue
		}
		items[skuID] = &item
	}
	return &meta, items, nil
}

func setCartCache(ctx context.Context, rdb redis.UniversalClient, userID int64, meta *cartMetaCache, items []CartItem) error {
	pipe := rdb.Pipeline()

	metaData, _ := sonic.Marshal(meta)
	pipe.Set(ctx, cacheKeyCartMeta(userID), metaData, cartCacheTTL)

	itemKey := cacheKeyCartItems(userID)
	pipe.Del(ctx, itemKey)
	for i := range items {
		itemData, _ := sonic.Marshal(&items[i])
		pipe.HSet(ctx, itemKey, strconv.FormatInt(items[i].SkuID, 10), itemData)
	}
	pipe.Expire(ctx, itemKey, cartCacheTTL)

	_, err := pipe.Exec(ctx)
	return err
}

func delCartCache(ctx context.Context, rdb redis.UniversalClient, userID int64) {
	rdb.Del(ctx, cacheKeyCartItems(userID), cacheKeyCartMeta(userID))
}
