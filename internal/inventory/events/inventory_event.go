package events

// InventoryReservedEvent 库存预占事件
type InventoryReservedEvent struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// InventoryReleasedEvent 库存释放事件
type InventoryReleasedEvent struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// InventoryLowEvent 库存不足事件
type InventoryLowEvent struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Threshold int    `json:"threshold"`
}
