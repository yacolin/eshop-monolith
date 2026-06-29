package inventory

// InventoryLowEvent 库存不足事件
type InventoryLowEvent struct {
	SkuID       int64 `json:"sku_id"`
	ProductID   int64 `json:"product_id"`
	WarehouseID int64 `json:"warehouse_id"`
	Quantity    int64 `json:"quantity"`
	Threshold   int64 `json:"threshold"`
}

func (InventoryLowEvent) RoutingKey() string { return "inventory.low" }
