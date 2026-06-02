package eventbus

// RegisterHandlers 注册事件处理器
func RegisterHandlers(bus *Bus) {
	// 注册订单事件处理器
	RegisterOrderHandlers(bus)

	// 注册库存事件处理器
	RegisterInventoryHandlers(bus)

	// 注册用户事件处理器
	RegisterUserHandlers(bus)

	// 注册支付事件处理器
	RegisterPaymentHandlers(bus)

	// 注册购物车事件处理器
	RegisterCartHandlers(bus)

	// 注册闪购事件处理器
	RegisterFlashHandlers(bus)
}
