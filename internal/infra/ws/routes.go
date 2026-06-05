package ws

import (
	"github.com/gin-gonic/gin"
)

// RegisterWSRoutes 注册 WebSocket 路由
func RegisterWSRoutes(v1 *gin.RouterGroup, hub *Hub) {
	wsHandler := NewHandler(hub)

	// WebSocket 连接端点（JWT 鉴权通过 token 查询参数）
	v1.GET("/ws", wsHandler.Upgrade)

	// WebSocket 在线统计（调试/管理用）
	v1.GET("/ws/stats", wsHandler.GetOnlineStats)
}
