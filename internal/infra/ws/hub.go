package ws

import (
	"context"
	"sync"

	"eshop-monolith/pkg/logger"
)

// Hub WebSocket 连接管理器，按 user_id 分组管理客户端
type Hub struct {
	mu         sync.RWMutex
	clients    map[int64]map[*Client]bool // user_id → 该用户的连接集合
	register   chan *Client               // 注册连接
	unregister chan *Client               // 注销连接

	// 关闭控制
	shutdownCh chan struct{}
	closed     bool
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		shutdownCh: make(chan struct{}),
	}
}

// Run 启动 Hub 主循环（在独立 goroutine 中运行）
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			count := len(h.clients[client.UserID])
			h.mu.Unlock()
			logger.Info("WebSocket 客户端已连接",
				"user_id", client.UserID,
				"同一用户连接数", count)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				if _, exists := h.clients[client.UserID][client]; exists {
					delete(h.clients[client.UserID], client)
					close(client.Send)
					if len(h.clients[client.UserID]) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			logger.Info("WebSocket 客户端已断开",
				"user_id", client.UserID)

		case <-h.shutdownCh:
			// 关闭所有连接
			h.mu.Lock()
			for userID, conns := range h.clients {
				for client := range conns {
					close(client.Send)
					client.Conn.Close()
				}
				delete(h.clients, userID)
			}
			h.closed = true
			h.mu.Unlock()
			logger.Info("WebSocket Hub 已关闭所有连接")
			return
		}
	}
}

// SendToUser 向指定用户发送消息（JSON bytes）
func (h *Hub) SendToUser(userID int64, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
				// 发送缓冲区满，说明该客户端消费过慢，断开
				logger.Warn("WebSocket 客户端发送缓冲区满，断开",
					"user_id", userID)
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// SendToUsers 向多个用户发送消息
func (h *Hub) SendToUsers(userIDs []int64, data []byte) {
	for _, uid := range userIDs {
		h.SendToUser(uid, data)
	}
}

// Broadcast 向所有连接的客户端广播消息
func (h *Hub) Broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.clients {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
				logger.Warn("WebSocket 广播：客户端缓冲区满，断开",
					"user_id", client.UserID)
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// GetOnlineCount 获取在线用户数（连接数）
func (h *Hub) GetOnlineCount() (userCount int, connCount int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connCount = 0
	for _, clients := range h.clients {
		connCount += len(clients)
	}
	return len(h.clients), connCount
}

// IsOnline 判断用户是否在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

// HandleEvent EventBus 事件处理器，将事件推送给对应用户
// 事件结构体需包含 UserID int64 或 CustomerID string 字段
func (h *Hub) HandleEvent(event interface{}) {
	// 尝试从不同类型的事件中提取用户 ID
	userID := extractUserID(event)
	if userID <= 0 {
		return
	}

	// 序列化为推送消息
	msg := NewPushMessage(event)
	if msg == nil {
		return
	}

	msgJSON, err := msg.Marshal()
	if err != nil {
		logger.Error("WebSocket 消息序列化失败", "error", err)
		return
	}

	h.SendToUser(userID, msgJSON)
}

// Shutdown 优雅关闭 Hub
func (h *Hub) Shutdown(ctx context.Context) error {
	select {
	case h.shutdownCh <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
