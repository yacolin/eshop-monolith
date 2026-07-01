package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/pkg/logger"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

// UserInfoProvider 用户信息查询回调
type UserInfoProvider func(userID int64) (username string, nickname string, err error)

// Hub WebSocket 连接管理器，按 user_id 分组管理客户端
type Hub struct {
	mu         sync.RWMutex
	clients    map[int64]map[*Client]bool // user_id → 该用户的连接集合
	register   chan *Client               // 注册连接
	unregister chan *Client               // 注销连接

	msgCache   *MessageCache   // 消息缓存管理器
	sessionMgr *SessionManager // 会话管理器

	userInfoProvider UserInfoProvider // 用户信息查询回调
	globalSeq        atomic.Int64     // 全局顺序号（用于实时广播）

	shutdownCh chan struct{}
	closed     bool
}

// NewHub 创建 Hub
func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		msgCache:   NewMessageCache(redisClient),
		sessionMgr: NewSessionManager(redisClient),
		shutdownCh: make(chan struct{}),
	}
}

// SetUserInfoProvider 设置用户信息查询回调
func (h *Hub) SetUserInfoProvider(provider UserInfoProvider) {
	h.userInfoProvider = provider
}

// Run 启动 Hub 主循环（在独立 goroutine 中运行）
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case <-h.shutdownCh:
			h.handleShutdown()
			return
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	wasOffline := h.clients[client.UserID] == nil
	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*Client]bool)
	}
	h.clients[client.UserID][client] = true
	count := len(h.clients[client.UserID])
	h.mu.Unlock()

	logger.Info("WebSocket 客户端已连接",
		"user_id", client.UserID,
		"同一用户连接数", count,
		"last_seq", client.LastSeq)

	// 先发送欢迎消息（确认连接建立，无延迟）
	go h.sendWelcomeMessage(client)

	// 发送全量用户列表同步给当前客户端
	go h.sendUserSync(client)

	// 再广播用户在线和统计数据
	if wasOffline {
		go h.broadcastUserEvent(client, "online")
	}
	go h.broadcastStats()
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	isOffline := false
	if _, ok := h.clients[client.UserID]; ok {
		if _, exists := h.clients[client.UserID][client]; exists {
			delete(h.clients[client.UserID], client)
			close(client.Send)

			if len(h.clients[client.UserID]) == 0 {
				delete(h.clients, client.UserID)
				isOffline = true
			}
		}
	}
	h.mu.Unlock()

	if isOffline {
		go h.broadcastUserEvent(client, "offline")
	}
	go h.broadcastStats()

	go h.sessionMgr.UpdateLastSeq(client.UserID, client.LastSeq)

	logger.Info("WebSocket 客户端已断开",
		"user_id", client.UserID,
		"last_seq", client.LastSeq)
}

func (h *Hub) handleShutdown() {
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
}

// ── Snapshot helpers ──

// snapshotClients 返回所有客户端的快照（读锁释放后安全使用）
func (h *Hub) snapshotClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	snapshot := make([]*Client, 0, total)
	for _, clients := range h.clients {
		for c := range clients {
			snapshot = append(snapshot, c)
		}
	}
	return snapshot
}

// snapshotUserClients 返回指定用户的客户端快照
func (h *Hub) snapshotUserClients(userID int64) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := h.clients[userID]
	snapshot := make([]*Client, 0, len(clients))
	for c := range clients {
		snapshot = append(snapshot, c)
	}
	return snapshot
}

// snapshotTargetClients 收集指定用户列表下所有客户端的快照
func (h *Hub) snapshotTargetClients(userIDs []int64) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var snapshot []*Client
	for _, uid := range userIDs {
		if clients, ok := h.clients[uid]; ok {
			for c := range clients {
				snapshot = append(snapshot, c)
			}
		}
	}
	return snapshot
}

// sendToClient 向单个客户端发送数据，缓冲区满时断开连接
func (h *Hub) sendToClient(client *Client, data []byte) {
	select {
	case client.Send <- data:
		if id := extractSeqID(data); id > 0 {
			client.LastSeq = id
		}
	default:
		logger.Warn("WebSocket 客户端发送缓冲区满，断开", "user_id", client.UserID)
		select {
		case h.unregister <- client:
		default:
		}
	}
}

// SendToUser 向指定用户发送消息（JSON bytes）
func (h *Hub) SendToUser(userID int64, data []byte) {
	for _, client := range h.snapshotUserClients(userID) {
		h.sendToClient(client, data)
	}
}

// SendToUsers 向多个用户并发发送消息
func (h *Hub) SendToUsers(userIDs []int64, data []byte) {
	clients := h.snapshotTargetClients(userIDs)
	g, _ := errgroup.WithContext(context.Background())
	for _, client := range clients {
		client := client
		g.Go(func() error {
			h.sendToClient(client, data)
			return nil
		})
	}
	g.Wait()
}

// Broadcast 向所有连接的客户端广播消息
func (h *Hub) Broadcast(data []byte) {
	clients := h.snapshotClients()
	g, _ := errgroup.WithContext(context.Background())
	for _, client := range clients {
		client := client
		g.Go(func() error {
			select {
			case client.Send <- data:
			default:
				logger.Warn("WebSocket 广播：客户端缓冲区满，断开", "user_id", client.UserID)
				select {
				case h.unregister <- client:
				default:
				}
			}
			return nil
		})
	}
	g.Wait()
}

// broadcastSafe 向所有客户端广播，缓冲区满时静默丢弃，不主动断开连接
// 用于实时统计/在线事件等非关键性广播，避免雪崩式断连
func (h *Hub) broadcastSafe(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.clients {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
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

// buildOnlineUserList 构建当前在线用户列表
func (h *Hub) buildOnlineUserList() []UserBrief {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]UserBrief, 0, len(h.clients))
	for _, clients := range h.clients {
		for client := range clients {
			username := client.Username
			nickname := client.Nickname
			if username == "" {
				username = fmt.Sprintf("user_%d", client.UserID)
			}
			if nickname == "" {
				nickname = fmt.Sprintf("用户%d", client.UserID)
			}
			users = append(users, UserBrief{
				ID:       client.UserID,
				Username: username,
				Nickname: nickname,
			})
			break
		}
	}
	return users
}

// sendUserSync 向指定客户端推送全量用户列表同步（前端全量替换）
func (h *Hub) sendUserSync(client *Client) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn("发送用户列表同步消息到已关闭的连接",
					"user_id", client.UserID,
					"panic", r)
			}
		}()

		users := h.buildOnlineUserList()
		msg := struct {
			Type    string          `json:"type"`
			Payload UserSyncPayload `json:"payload"`
		}{
			Type: "user",
			Payload: UserSyncPayload{
				Action: "sync",
				Users:  users,
			},
		}

		data, err := json.Marshal(msg)
		if err != nil {
			logger.Error("序列化用户列表同步消息失败", "error", err)
			return
		}

		select {
		case client.Send <- data:
		default:
		}
	}()
}

// broadcastUserEvent 广播用户上下线事件
func (h *Hub) broadcastUserEvent(client *Client, action string) {
	username := client.Username
	nickname := client.Nickname
	if username == "" {
		username = fmt.Sprintf("user_%d", client.UserID)
	}
	if nickname == "" {
		nickname = fmt.Sprintf("用户%d", client.UserID)
	}

	msg := &RealtimeMessage{
		Seq:  h.globalSeq.Add(1),
		Type: "user",
		Payload: UserEventPayload{
			Action: action,
			User: UserBrief{
				ID:       client.UserID,
				Username: username,
				Nickname: nickname,
			},
			Timestamp: time.Now().UnixMilli(),
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("序列化用户事件消息失败", "error", err)
		return
	}

	h.broadcastSafe(data)
}

// broadcastStats 广播在线统计
func (h *Hub) broadcastStats() {
	userCount, connCount := h.GetOnlineCount()

	msg := &RealtimeMessage{
		Seq:  h.globalSeq.Add(1),
		Type: "stats",
		Payload: StatsPayload{
			OnlineUsers: userCount,
			Connections: connCount,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("序列化统计消息失败", "error", err)
		return
	}

	h.broadcastSafe(data)
}

// IsOnline 判断用户是否在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

// HandleEvent EventBus 事件处理器，将事件推送给对应用户
func (h *Hub) HandleEvent(event interface{}) {
	userID := extractUserID(event)
	if userID <= 0 {
		return
	}

	seqID, err := h.msgCache.NextSeqID(userID)
	if err != nil {
		logger.Error("获取序列ID失败", "user_id", userID, "error", err)
		return
	}

	msg := NewPushMessage(event, seqID)
	if msg == nil {
		return
	}

	msgJSON, err := msg.Marshal()
	if err != nil {
		logger.Error("WebSocket 消息序列化失败", "error", err)
		return
	}

	if err := h.msgCache.StoreMessage(userID, seqID, msgJSON); err != nil {
		logger.Warn("存储消息缓存失败", "user_id", userID, "error", err)
	}

	h.SendToUser(userID, msgJSON)

	go h.sessionMgr.UpdateLastSeq(userID, seqID)
}

// HandleMessage 处理来自 RabbitMQ 的消息
func (h *Hub) HandleMessage(msg rabbitmq.Message) {
	// 低库存事件广播给所有在线用户
	if msg.RoutingKey == "inventory.low" {
		broadcastMsg := &RealtimeMessage{
				Seq:     h.globalSeq.Add(1),
				Type:    "inventory.low",
				Payload: msg.Payload,
			}
		data, _ := broadcastMsg.Marshal()
		if data != nil {
			h.Broadcast(data)
		}
		return
	}

	userID := extractUserIDFromMessage(msg)
	if userID <= 0 {
		return
	}

	seqID, err := h.msgCache.NextSeqID(userID)
	if err != nil {
		logger.Error("获取序列ID失败", "user_id", userID, "error", err)
		return
	}

	pushMsg := NewPushMessageFromRaw(msg.RoutingKey, msg.Payload, seqID)
	if pushMsg == nil {
		return
	}

	msgJSON, _ := pushMsg.Marshal()

	if err := h.msgCache.StoreMessage(userID, seqID, msgJSON); err != nil {
		logger.Warn("存储消息缓存失败", "user_id", userID, "error", err)
	}

	h.SendToUser(userID, msgJSON)

	go h.sessionMgr.UpdateLastSeq(userID, seqID)
}

// SyncMessages 同步消息（增量补发）
// 客户端上报lastSeq，服务端补发 (lastSeq, currentSeq] 区间的消息
func (h *Hub) SyncMessages(userID int64, lastSeq int64) {
	currentSeq, err := h.msgCache.GetCurrentSeqID(userID)
	if err != nil {
		logger.Error("获取当前序列ID失败", "user_id", userID, "error", err)
		return
	}

	if lastSeq >= currentSeq {
		logger.Info("无需同步，客户端已是最新状态", "user_id", userID, "last_seq", lastSeq, "current_seq", currentSeq)
		return
	}

	cachedMinSeq, cachedMaxSeq, err := h.msgCache.GetCachedSeqRange(userID)
	if err != nil {
		logger.Error("获取缓存序列范围失败", "user_id", userID, "error", err)
		return
	}

	if lastSeq < cachedMinSeq {
		logger.Warn("客户端缺失消息超过缓存窗口，需全量同步",
			"user_id", userID,
			"last_seq", lastSeq,
			"cached_min_seq", cachedMinSeq,
			"cached_max_seq", cachedMaxSeq)

		h.sendFullSyncRequired(userID)
		return
	}

	messages, err := h.msgCache.GetMessages(userID, lastSeq, currentSeq)
	if err != nil {
		logger.Error("获取消息失败", "user_id", userID, "error", err)
		return
	}

	logger.Info("补发消息", "user_id", userID, "count", len(messages), "from", lastSeq+1, "to", currentSeq)

	for _, msg := range messages {
		h.SendToUser(userID, msg)
	}

	go h.sessionMgr.IncrementReconnectCount(userID)
	go h.sessionMgr.UpdateLastSeq(userID, currentSeq)
}

// GetUserLastSeq 获取用户最后收到的消息序列号
func (h *Hub) GetUserLastSeq(userID int64) (int64, error) {
	session, err := h.sessionMgr.GetSession(userID)
	if err != nil {
		return 0, err
	}
	if session == nil {
		return 0, nil
	}
	return session.LastSeq, nil
}

// sendWelcomeMessage 发送欢迎消息
func (h *Hub) sendWelcomeMessage(client *Client) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn("发送欢迎消息到已关闭的连接",
					"user_id", client.UserID,
					"panic", r)
			}
		}()

		currentSeq, _ := h.msgCache.GetCurrentSeqID(client.UserID)
		msg := &SystemMessage{
			Type:            "welcome",
			SequenceID:      currentSeq,
			Message:         "连接成功",
			RequireFullSync: false,
		}
		data, _ := msg.Marshal()

		select {
		case client.Send <- data:
		default:
		}
	}()
}

// sendFullSyncRequired 发送需要全量同步的通知
func (h *Hub) sendFullSyncRequired(userID int64) {
	msg := NewSystemMessage("sync_required", "缺失消息超过缓存窗口，请先获取全量数据", true)
	data, _ := msg.Marshal()
	h.SendToUser(userID, data)
}

// extractSeqID 从消息JSON中提取sequence_id
func extractSeqID(data []byte) int64 {
	var msg PushMessage
	if err := msg.Unmarshal(data); err == nil {
		return msg.SequenceID
	}
	return 0
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
