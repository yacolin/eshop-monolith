package events

// InventoryReservedEvent 库存预占事件
type InventoryReservedEvent struct {
	SkuID    string `json:"sku_id"`
	Quantity int    `json:"quantity"`
}

func (e InventoryReservedEvent) RoutingKey() string { return "inventory.reserved" }

// InventoryReleasedEvent 库存释放事件
type InventoryReleasedEvent struct {
	SkuID    string `json:"sku_id"`
	Quantity int    `json:"quantity"`
}

func (e InventoryReleasedEvent) RoutingKey() string { return "inventory.released" }

// InventoryLowEvent 库存不足事件
type InventoryLowEvent struct {
	SkuID     string `json:"sku_id"`
	Quantity  int    `json:"quantity"`
	Threshold int    `json:"threshold"`
}

func (e InventoryLowEvent) RoutingKey() string { return "inventory.low" }
