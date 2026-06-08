package ws

import (
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWSRoutes(v1 *gin.RouterGroup, hub *Hub) {
	wsHandler := NewHandler(hub)

	v1.GET("/ws", wsHandler.Upgrade)
	v1.GET("/ws/stats", wsHandler.GetOnlineStats)

	v1.POST("/ws/reconnect", wsHandler.Reconnect)
	v1.GET("/ws/session", wsHandler.GetUserSession)

	// 需要认证的 WebSocket 路由
	auth := v1.Group("/ws")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/test/push", wsHandler.PushTestMessage)
	}
}