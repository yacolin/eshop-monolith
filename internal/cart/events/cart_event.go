package events

// CartItemAddedEvent 购物车项添加事件
type CartItemAddedEvent struct {
	CartID    int64 `json:"cart_id"`
	UserID    int64 `json:"user_id"`
	SkuID     int64 `json:"sku_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
}

func (e CartItemAddedEvent) RoutingKey() string { return "cart.item-added" }

// CartItemUpdatedEvent 购物车项更新事件
type CartItemUpdatedEvent struct {
	CartID      int64 `json:"cart_id"`
	UserID      int64 `json:"user_id"`
	ItemID      int64 `json:"item_id"`
	SkuID       int64 `json:"sku_id"`
	ProductID   int64 `json:"product_id"`
	OldQuantity int   `json:"old_quantity"`
	NewQuantity int   `json:"new_quantity"`
	Price       int64 `json:"price"`
}

func (e CartItemUpdatedEvent) RoutingKey() string { return "cart.item-updated" }

// CartItemRemovedEvent 购物车项移除事件
type CartItemRemovedEvent struct {
	CartID    int64 `json:"cart_id"`
	UserID    int64 `json:"user_id"`
	ItemID    int64 `json:"item_id"`
	SkuID     int64 `json:"sku_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
}

func (e CartItemRemovedEvent) RoutingKey() string { return "cart.item-removed" }

// CartClearedEvent 购物车清空事件
type CartClearedEvent struct {
	CartID int64 `json:"cart_id"`
	UserID int64 `json:"user_id"`
}

func (e CartClearedEvent) RoutingKey() string { return "cart.cleared" }
