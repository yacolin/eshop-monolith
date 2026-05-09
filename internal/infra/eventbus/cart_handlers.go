package eventbus

import (
	"eshop-monolith/internal/cart/events"
	"eshop-monolith/pkg/logger"
)

// RegisterCartHandlers 注册购物车事件处理器
func RegisterCartHandlers(bus *Bus) {
	bus.Subscribe("cart.CartItemAddedEvent", handleCartItemAdded)
	bus.Subscribe("cart.CartItemUpdatedEvent", handleCartItemUpdated)
	bus.Subscribe("cart.CartItemRemovedEvent", handleCartItemRemoved)
	bus.Subscribe("cart.CartClearedEvent", handleCartCleared)
}

// handleCartItemAdded 处理购物车项添加事件
func handleCartItemAdded(event interface{}) {
	e, ok := event.(events.CartItemAddedEvent)
	if !ok {
		return
	}
	logger.Info("Cart item added",
		"cart_id", e.CartID,
		"user_id", e.UserID,
		"product_id", e.ProductID,
		"quantity", e.Quantity,
		"price", e.Price,
	)
	// 这里可以添加具体的业务逻辑，例如：
	// 1. 记录购物车操作日志
	// 2. 触发相关的业务流程
	// 3. 与其他模块进行集成
}

// handleCartItemUpdated 处理购物车项更新事件
func handleCartItemUpdated(event interface{}) {
	e, ok := event.(events.CartItemUpdatedEvent)
	if !ok {
		return
	}
	logger.Info("Cart item updated",
		"cart_id", e.CartID,
		"user_id", e.UserID,
		"item_id", e.ItemID,
		"product_id", e.ProductID,
		"old_quantity", e.OldQuantity,
		"new_quantity", e.NewQuantity,
		"price", e.Price,
	)
	// 这里可以添加具体的业务逻辑
}

// handleCartItemRemoved 处理购物车项移除事件
func handleCartItemRemoved(event interface{}) {
	e, ok := event.(events.CartItemRemovedEvent)
	if !ok {
		return
	}
	logger.Info("Cart item removed",
		"cart_id", e.CartID,
		"user_id", e.UserID,
		"item_id", e.ItemID,
		"product_id", e.ProductID,
		"quantity", e.Quantity,
		"price", e.Price,
	)
	// 这里可以添加具体的业务逻辑
}

// handleCartCleared 处理购物车清空事件
func handleCartCleared(event interface{}) {
	e, ok := event.(events.CartClearedEvent)
	if !ok {
		return
	}
	logger.Info("Cart cleared",
		"cart_id", e.CartID,
		"user_id", e.UserID,
	)
	// 这里可以添加具体的业务逻辑
}
