package eventbus

import (
	"eshop-monolith/pkg/logger"
	"eshop-monolith/internal/user/events"
)

// RegisterUserHandlers 注册用户事件处理器
func RegisterUserHandlers(bus *Bus) {
	bus.Subscribe("user.UserRegisteredEvent", handleUserRegistered)
	bus.Subscribe("user.UserLoggedInEvent", handleUserLoggedIn)
	bus.Subscribe("user.UserProfileUpdatedEvent", handleUserProfileUpdated)
}

// handleUserRegistered 处理用户注册事件
func handleUserRegistered(event interface{}) {
	e, ok := event.(events.UserRegisteredEvent)
	if !ok {
		return
	}
	logger.Info("User registered", "user_id", e.UserID, "username", e.Username, "email", e.Email)
	// 这里可以添加发送欢迎邮件等逻辑
}

// handleUserLoggedIn 处理用户登录事件
func handleUserLoggedIn(event interface{}) {
	e, ok := event.(events.UserLoggedInEvent)
	if !ok {
		return
	}
	logger.Info("User logged in", "user_id", e.UserID, "username", e.Username, "ip", e.IP)
	// 这里可以添加登录记录等逻辑
}

// handleUserProfileUpdated 处理用户资料更新事件
func handleUserProfileUpdated(event interface{}) {
	e, ok := event.(events.UserProfileUpdatedEvent)
	if !ok {
		return
	}
	logger.Info("User profile updated", "user_id", e.UserID, "nickname", e.Nickname)
	// 这里可以添加资料更新通知等逻辑
}
