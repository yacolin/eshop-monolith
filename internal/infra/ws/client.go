package ws

import (
	"time"

	"eshop-monolith/pkg/logger"

	"github.com/gorilla/websocket"
)

const (
	// 写入超时
	writeWait = 10 * time.Second

	// 心跳：服务端 Ping 间隔
	pingPeriod = 30 * time.Second

	// 允许的 Pong 等待时间（必须大于 pingPeriod）
	pongWait = 60 * time.Second

	// 发送通道缓冲区大小
	sendBufferSize = 64
)

// Client 代表一个 WebSocket 连接
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID int64
}

// NewClient 创建客户端
func NewClient(hub *Hub, conn *websocket.Conn, userID int64) *Client {
	return &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, sendBufferSize),
		UserID: userID,
	}
}

// ReadPump 从 WebSocket 连接读取消息（goroutine）
// 目前 WebSocket 主要用作服务端推送，客户端消息仅用于心跳响应
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512) // 只读取控制消息，限制数据量
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Warn("WebSocket 读取异常",
					"user_id", c.UserID,
					"error", err)
			}
			break
		}
		// 客户端消息现阶段不处理（纯推送模式）
	}
}

// WritePump 从 send channel 写入 WebSocket 连接（goroutine）
// 包含心跳 Ping
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了 send channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将 send channel 中积压的消息一并写入（减少网络包数量）
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
