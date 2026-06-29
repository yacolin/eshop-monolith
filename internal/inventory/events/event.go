package events

// InventoryLowEvent 库存不足事件（保留兼容旧模块引用）
type InventoryLowEvent struct {
	SkuID       int64
	ProductID   int64
	WarehouseID int64
	Quantity    int64
	Threshold   int64
}
