package ws

import (
	"net/http"
	"strconv"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（生产环境应配置具体域名）
	},
}

// Handler WebSocket 请求处理器
type Handler struct {
	hub *Hub
}

// NewHandler 创建 WebSocket 处理器
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// Upgrade 升级 HTTP 连接为 WebSocket
// 通过查询参数 token 进行 JWT 鉴权
// URL: /api/v1/ws?token=xxx
func (h *Handler) Upgrade(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	// 解析 JWT token
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	// 提取 user_id
	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrade 失败时不能调用 c.Error（响应已写入）
		return
	}

	// 创建客户端并注册到 Hub
	client := NewClient(h.hub, conn, userID)
	h.hub.register <- client

	// 启动读写 goroutine
	go client.WritePump()
	go client.ReadPump()
}

// GetOnlineStats 获取在线统计（调试/管理用）
func (h *Handler) GetOnlineStats(c *gin.Context) {
	userCount, connCount := h.hub.GetOnlineCount()
	response.Success(c, gin.H{
		"online_users": userCount,
		"connections":  connCount,
	})
}

// extractUserIDFromClaims 从 JWT claims 中提取 user_id
func extractUserIDFromClaims(claims map[string]interface{}) (int64, error) {
	v, ok := claims["user_id"]
	if !ok {
		return 0, errcode.ErrUnauthorized
	}

	switch id := v.(type) {
	case float64:
		return int64(id), nil
	case int64:
		return id, nil
	case int:
		return int64(id), nil
	case string:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return 0, errcode.ErrUnauthorized
		}
		return n, nil
	default:
		return 0, errcode.ErrUnauthorized
	}
}
