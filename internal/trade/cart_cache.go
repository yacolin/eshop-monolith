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
	cartCacheTTL = 24 * time.Hour
)

func cartItemsKey(userID int64) string { return fmt.Sprintf("cart:%d:items", userID) }
func cartMetaKey(userID int64) string  { return fmt.Sprintf("cart:%d:meta", userID) }

type cartMeta struct {
	ItemCount   int   `json:"item_count"`
	TotalAmount int64 `json:"total_amount"`
}

type cachedCartItem struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	SkuSpec     string `json:"sku_spec"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
}

func readCartFromRedis(ctx context.Context, rdb redis.UniversalClient, userID int64) (*CartResponse, error) {
	pipe := rdb.Pipeline()
	metaCmd := pipe.Get(ctx, cartMetaKey(userID))
	itemsCmd := pipe.HGetAll(ctx, cartItemsKey(userID))
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	var meta cartMeta
	metaData, err := metaCmd.Bytes()
	if err != nil {
		return nil, err
	}
	if err := sonic.Unmarshal(metaData, &meta); err != nil {
		return nil, err
	}

	items := make([]CartItemResponse, 0, len(itemsCmd.Val()))
	for skuIDStr, itemJSON := range itemsCmd.Val() {
		var ci cachedCartItem
		if err := sonic.Unmarshal([]byte(itemJSON), &ci); err != nil {
			continue
		}
		skuID, _ := strconv.ParseInt(skuIDStr, 10, 64)
		items = append(items, CartItemResponse{
			SkuID:       skuID,
			ProductID:   ci.ProductID,
			ProductName: ci.ProductName,
			SkuSpec:     ci.SkuSpec,
			Image:       ci.Image,
			Price:       ci.Price,
			Quantity:    ci.Quantity,
			Subtotal:    ci.Price * int64(ci.Quantity),
		})
	}

	return &CartResponse{
		ItemCount:   meta.ItemCount,
		TotalAmount: meta.TotalAmount,
		Items:       items,
	}, nil
}

func writeCartToRedis(ctx context.Context, rdb redis.UniversalClient, userID int64, resp *CartResponse) error {
	pipe := rdb.Pipeline()
	metaData, _ := sonic.Marshal(cartMeta{ItemCount: resp.ItemCount, TotalAmount: resp.TotalAmount})
	pipe.Set(ctx, cartMetaKey(userID), metaData, cartCacheTTL)
	pipe.Del(ctx, cartItemsKey(userID))
	for _, item := range resp.Items {
		data, _ := sonic.Marshal(cachedCartItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			SkuSpec:     item.SkuSpec,
			Image:       item.Image,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
		pipe.HSet(ctx, cartItemsKey(userID), strconv.FormatInt(item.SkuID, 10), data)
	}
	pipe.Expire(ctx, cartItemsKey(userID), cartCacheTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func delCartFromRedis(ctx context.Context, rdb redis.UniversalClient, userID int64) {
	rdb.Del(ctx, cartItemsKey(userID), cartMetaKey(userID))
}
