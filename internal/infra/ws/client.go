package ws

import (
	"encoding/json"
	"time"

	"eshop-monolith/pkg/logger"

	"github.com/gorilla/websocket"
)

const (
	// 写入超时
	writeWait = 10 * time.Second

	// 心跳：客户端Ping间隔（30秒）
	pingPeriod = 30 * time.Second

	// 允许的Pong等待时间（必须大于pingPeriod，连续失败2次判定断开）
	pongWait = 90 * time.Second

	// 发送通道缓冲区大小
	sendBufferSize = 64

	// 连续心跳失败次数上限
	maxPingFailures = 2
)

// Client 代表一个 WebSocket 连接
type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	UserID    int64
	Username  string // 用户名（用于在线事件广播）
	Nickname  string // 昵称（用于在线事件广播）
	LastSeq   int64 // 客户端最后收到的消息序列号

	// 心跳相关
	pingFailCount int        // 连续Ping失败次数
	lastPingTime  time.Time  // 上次发送Ping的时间
	pingTicker    *time.Ticker
}

// NewClient 创建客户端
func NewClient(hub *Hub, conn *websocket.Conn, userID int64, lastSeq int64) *Client {
	return &Client{
		Hub:       hub,
		Conn:      conn,
		Send:      make(chan []byte, sendBufferSize),
		UserID:    userID,
		LastSeq:   lastSeq,
		pingTicker: time.NewTicker(pingPeriod),
	}
}

// ReadPump 从 WebSocket 连接读取消息（goroutine）
func (c *Client) ReadPump() {
	defer func() {
		c.pingTicker.Stop()
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.pingFailCount = 0 // 收到Pong，重置失败计数
		c.updateLastSeqFromPong(appData)
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Warn("WebSocket 读取异常",
					"user_id", c.UserID,
					"error", err)
			}
			break
		}

		c.handleClientMessage(message)
	}
}

// WritePump 从 send channel 写入 WebSocket 连接（goroutine）
func (c *Client) WritePump() {
	defer func() {
		c.pingTicker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-c.pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte(c.buildPingPayload())); err != nil {
				c.pingFailCount++
				if c.pingFailCount >= maxPingFailures {
					logger.Warn("WebSocket 客户端心跳失败次数达到上限，断开连接",
						"user_id", c.UserID,
						"fail_count", c.pingFailCount)
					return
				}
				logger.Warn("WebSocket 客户端心跳失败",
					"user_id", c.UserID,
					"fail_count", c.pingFailCount)
			} else {
				c.pingFailCount = 0
				c.lastPingTime = time.Now()
			}
		}
	}
}

// handleClientMessage 处理客户端发送的消息
func (c *Client) handleClientMessage(message []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		logger.Warn("WebSocket 解析客户端消息失败",
			"user_id", c.UserID,
			"error", err)
		return
	}

	switch msg.Type {
	case "ping":
		// 客户端Ping，回复Pong
		c.sendPong(msg.Seq)
	case "sync":
		// 客户端请求同步，上报lastSeq
		c.Hub.SyncMessages(c.UserID, msg.LastSeq)
	case "reconnect":
		// 客户端重连，提交lastSeq
		c.LastSeq = msg.LastSeq
		c.Hub.SyncMessages(c.UserID, msg.LastSeq)
	}
}

// sendPong 发送Pong响应
func (c *Client) sendPong(seq int64) {
	pong := &ClientMessage{
		Type:    "pong",
		Seq:     seq,
		LastSeq: c.LastSeq,
	}
	data, err := json.Marshal(pong)
	if err != nil {
		logger.Error("WebSocket 序列化Pong消息失败", "user_id", c.UserID, "error", err)
		return
	}

	select {
	case c.Send <- data:
	default:
		logger.Warn("WebSocket 发送Pong失败，缓冲区满", "user_id", c.UserID)
	}
}

// buildPingPayload 构建Ping消息负载（包含服务器当前时间）
func (c *Client) buildPingPayload() string {
	type pingPayload struct {
		Timestamp int64 `json:"timestamp"`
		LastSeq   int64 `json:"last_seq"`
	}
	payload := pingPayload{
		Timestamp: time.Now().UnixMilli(),
		LastSeq:   c.LastSeq,
	}
	data, _ := json.Marshal(payload)
	return string(data)
}

// updateLastSeqFromPong 从Pong响应中更新lastSeq
func (c *Client) updateLastSeqFromPong(appData string) {
	if appData == "" {
		return
	}
	var payload struct {
		LastSeq int64 `json:"last_seq"`
	}
	if err := json.Unmarshal([]byte(appData), &payload); err == nil && payload.LastSeq > c.LastSeq {
		c.LastSeq = payload.LastSeq
	}
}

// ClientMessage 客户端消息结构
type ClientMessage struct {
	Type    string `json:"type"`
	Seq     int64  `json:"seq"`
	LastSeq int64  `json:"last_seq"`
}