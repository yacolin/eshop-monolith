package eventbus

import (
	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/domain/order"
	"eshop-monolith/internal/domain/user"
	"eshop-monolith/internal/pkg/logger"
)

// RegisterHandlers 注册事件处理器
func RegisterHandlers(bus *Bus) {
	// 订单事件处理器
	bus.Subscribe("order.OrderCreatedEvent", handleOrderCreated)
	bus.Subscribe("order.OrderStatusChangedEvent", handleOrderStatusChanged)
	bus.Subscribe("order.OrderCancelledEvent", handleOrderCancelled)

	// 库存事件处理器
	bus.Subscribe("inventory.InventoryReservedEvent", handleInventoryReserved)
	bus.Subscribe("inventory.InventoryReleasedEvent", handleInventoryReleased)
	bus.Subscribe("inventory.InventoryLowEvent", handleInventoryLow)

	// 用户事件处理器
	bus.Subscribe("user.UserRegisteredEvent", handleUserRegistered)
	bus.Subscribe("user.UserLoggedInEvent", handleUserLoggedIn)
	bus.Subscribe("user.UserProfileUpdatedEvent", handleUserProfileUpdated)
}

// handleOrderCreated 处理订单创建事件
func handleOrderCreated(event interface{}) {
	e, ok := event.(order.OrderCreatedEvent)
	if !ok {
		return
	}
	logger.Info("Order created", "order_id", e.OrderID, "user_id", e.UserID, "total_amount", e.TotalAmount)
	// 这里可以添加发送通知、记录日志等逻辑
}

// handleOrderStatusChanged 处理订单状态变更事件
func handleOrderStatusChanged(event interface{}) {
	e, ok := event.(order.OrderStatusChangedEvent)
	if !ok {
		return
	}
	logger.Info("Order status changed", "order_id", e.OrderID, "old_status", e.OldStatus, "new_status", e.NewStatus)
	// 这里可以添加状态变更通知等逻辑
}

// handleOrderCancelled 处理订单取消事件
func handleOrderCancelled(event interface{}) {
	e, ok := event.(order.OrderCancelledEvent)
	if !ok {
		return
	}
	logger.Info("Order cancelled", "order_id", e.OrderID, "user_id", e.UserID)
	// 这里可以添加退款处理等逻辑
}

// handleInventoryReserved 处理库存预占事件
func handleInventoryReserved(event interface{}) {
	e, ok := event.(inventory.InventoryReservedEvent)
	if !ok {
		return
	}
	logger.Info("Inventory reserved", "product_id", e.ProductID, "quantity", e.Quantity)
	// 这里可以添加库存预警检查等逻辑
}

// handleInventoryReleased 处理库存释放事件
func handleInventoryReleased(event interface{}) {
	e, ok := event.(inventory.InventoryReleasedEvent)
	if !ok {
		return
	}
	logger.Info("Inventory released", "product_id", e.ProductID, "quantity", e.Quantity)
	// 这里可以添加库存更新逻辑
}

// handleInventoryLow 处理库存不足事件
func handleInventoryLow(event interface{}) {
	e, ok := event.(inventory.InventoryLowEvent)
	if !ok {
		return
	}
	logger.Warn("Inventory low", "product_id", e.ProductID, "quantity", e.Quantity, "threshold", e.Threshold)
	// 这里可以添加补货提醒等逻辑
}

// handleUserRegistered 处理用户注册事件
func handleUserRegistered(event interface{}) {
	e, ok := event.(user.UserRegisteredEvent)
	if !ok {
		return
	}
	logger.Info("User registered", "user_id", e.UserID, "username", e.Username, "email", e.Email)
	// 这里可以添加发送欢迎邮件等逻辑
}

// handleUserLoggedIn 处理用户登录事件
func handleUserLoggedIn(event interface{}) {
	e, ok := event.(user.UserLoggedInEvent)
	if !ok {
		return
	}
	logger.Info("User logged in", "user_id", e.UserID, "username", e.Username, "ip", e.IP)
	// 这里可以添加登录记录等逻辑
}

// handleUserProfileUpdated 处理用户资料更新事件
func handleUserProfileUpdated(event interface{}) {
	e, ok := event.(user.UserProfileUpdatedEvent)
	if !ok {
		return
	}
	logger.Info("User profile updated", "user_id", e.UserID, "username", e.Username)
	// 这里可以添加资料更新通知等逻辑
}