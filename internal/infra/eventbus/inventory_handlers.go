package eventbus

import (
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/logger"
)

// RegisterInventoryHandlers 注册库存事件处理器
func RegisterInventoryHandlers(bus *Bus) {
	bus.Subscribe("inventory.InventoryReservedEvent", handleInventoryReserved)
	bus.Subscribe("inventory.InventoryReleasedEvent", handleInventoryReleased)
	bus.Subscribe("inventory.InventoryLowEvent", handleInventoryLow)
}

// handleInventoryReserved 处理库存预占事件
func handleInventoryReserved(event interface{}) {
	e, ok := event.(events.InventoryReservedEvent)
	if !ok {
		return
	}
	logger.Info("Inventory reserved", "product_id", e.ProductID, "quantity", e.Quantity)
	// 这里可以添加库存预警检查等逻辑
}

// handleInventoryReleased 处理库存释放事件
func handleInventoryReleased(event interface{}) {
	e, ok := event.(events.InventoryReleasedEvent)
	if !ok {
		return
	}
	logger.Info("Inventory released", "product_id", e.ProductID, "quantity", e.Quantity)
	// 这里可以添加库存更新逻辑
}

// handleInventoryLow 处理库存不足事件
func handleInventoryLow(event interface{}) {
	e, ok := event.(events.InventoryLowEvent)
	if !ok {
		return
	}
	logger.Warn("Inventory low", "product_id", e.ProductID, "quantity", e.Quantity, "threshold", e.Threshold)
	// 这里可以添加补货提醒等逻辑
}
