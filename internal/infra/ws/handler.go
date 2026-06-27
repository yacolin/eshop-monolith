package ws

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"eshop-monolith/internal/infra/ws/dto"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// Upgrade 升级 HTTP 连接为 WebSocket
// @Summary 升级WebSocket连接
// @Description 将HTTP连接升级为WebSocket连接，支持断线重连时携带last_seq参数
// @Tags websocket
// @Accept json
// @Produce json
// @Param token query string true "JWT Token"
// @Param last_seq query int false "最后收到的消息序列号（重连时使用）" default(0)
// @Success 101 "WebSocket连接建立成功"
// @Failure 401 "未授权"
// @Failure 400 "请求参数错误"
// @Router /api/v1/ws [get]
func (h *Handler) Upgrade(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	lastSeqStr := c.Query("last_seq")
	var lastSeq int64
	if lastSeqStr != "" {
		lastSeq, err = strconv.ParseInt(lastSeqStr, 10, 64)
		if err != nil {
			c.Error(errcode.ErrInvalidParams)
			return
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(h.hub, conn, userID, lastSeq)

	// 异步加载用户信息，不阻塞 WS 握手路径
	if h.hub.userInfoProvider != nil {
		provider := h.hub.userInfoProvider
		go func() {
			username, nickname, lookupErr := provider(userID)
			if lookupErr == nil {
				client.Username = username
				client.Nickname = nickname
			} else {
				logger.Warn("查询用户信息失败", "user_id", userID, "error", lookupErr)
			}
		}()
	}

	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

// GetOnlineStats 获取在线统计
// @Summary 获取在线统计
// @Description 获取当前WebSocket在线用户数和连接数
// @Tags websocket
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.OnlineStatsResponse}
// @Router /api/v1/ws/stats [get]
func (h *Handler) GetOnlineStats(c *gin.Context) {
	userCount, connCount := h.hub.GetOnlineCount()
	response.Success(c, &dto.OnlineStatsResponse{
		OnlineUsers: userCount,
		Connections: connCount,
	})
}

// Reconnect 重连接口（HTTP）
// @Summary WebSocket重连接口
// @Description 客户端重连时提交lastSeq，服务端返回需要补发的消息
// @Tags websocket
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body dto.ReconnectRequest true "重连请求"
// @Success 200 {object} response.Response{data=dto.ReconnectResponse}
// @Failure 401 "未授权"
// @Failure 400 "请求参数错误"
// @Router /api/v1/ws/reconnect [post]
func (h *Handler) Reconnect(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	var req dto.ReconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	currentSeq, err := h.hub.msgCache.GetCurrentSeqID(userID)
	if err != nil {
		response.SysError(c, err)
		return
	}

	if req.LastSeq >= currentSeq {
		response.Success(c, &dto.ReconnectResponse{
			Status:          "ok",
			Message:         "已是最新状态",
			LastSeq:         req.LastSeq,
			CurrentSeq:      currentSeq,
			NeedFullSync:    false,
			NeedIncremental: false,
			MessageCount:    0,
			Messages:        []interface{}{},
		})
		return
	}

	cachedMinSeq, cachedMaxSeq, err := h.hub.msgCache.GetCachedSeqRange(userID)
	if err != nil {
		response.SysError(c, err)
		return
	}

	if req.LastSeq < cachedMinSeq {
		response.Success(c, &dto.ReconnectResponse{
			Status:          "sync_required",
			Message:         "缺失消息超过缓存窗口，请先获取全量数据",
			LastSeq:         req.LastSeq,
			CurrentSeq:      currentSeq,
			CachedMinSeq:    cachedMinSeq,
			CachedMaxSeq:    cachedMaxSeq,
			NeedFullSync:    true,
			NeedIncremental: false,
			MessageCount:    0,
			Messages:        []interface{}{},
		})
		return
	}

	rawMessages, err := h.hub.msgCache.GetMessages(userID, req.LastSeq, currentSeq)
	if err != nil {
		response.SysError(c, err)
		return
	}

	messages := make([]interface{}, len(rawMessages))
	for i, msg := range rawMessages {
		messages[i] = string(msg)
	}

	response.Success(c, &dto.ReconnectResponse{
		Status:          "ok",
		Message:         "增量同步完成",
		LastSeq:         req.LastSeq,
		CurrentSeq:      currentSeq,
		NeedFullSync:    false,
		NeedIncremental: true,
		MessageCount:    len(messages),
		Messages:        messages,
	})

	go h.hub.sessionMgr.IncrementReconnectCount(userID)
	go h.hub.sessionMgr.UpdateLastSeq(userID, currentSeq)
}

// GetUserSession 获取用户会话信息
// @Summary 获取用户会话信息
// @Description 获取当前用户的WebSocket会话信息，包括最后收到的消息序列号等
// @Tags websocket
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.Response{data=dto.SessionResponse}
// @Failure 401 "未授权"
// @Router /api/v1/ws/session [get]
func (h *Handler) GetUserSession(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	session, err := h.hub.sessionMgr.GetSession(userID)
	if err != nil {
		response.SysError(c, err)
		return
	}

	if session == nil {
		response.Success(c, &dto.SessionResponse{
			Exists:  false,
			LastSeq: 0,
		})
		return
	}

	response.Success(c, &dto.SessionResponse{
		Exists:         true,
		UserID:         session.UserID,
		LastSeq:        session.LastSeq,
		ConnectedAt:    session.ConnectedAt.Format(time.RFC3339),
		LastActiveAt:   session.LastActiveAt.Format(time.RFC3339),
		ReconnectCount: session.ReconnectCount,
	})
}

// TestPushRequest 测试推送请求
type TestPushRequest struct {
	Title   string `json:"title" binding:"required"`   // 通知标题
	Message string `json:"message" binding:"required"` // 通知内容
	Level   string `json:"level"`                      // 级别: info/warning/error
	Target  string `json:"target"`                     // 目标: all(默认) 或 user_id
}

// PushTestMessage 推送测试消息
// @Summary 推送WebSocket测试消息
// @Description 向WebSocket客户端推送一条测试消息，用于调试
// @Tags websocket
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body TestPushRequest true "测试消息"
// @Success 200 {object} response.Response
// @Failure 401 "未授权"
// @Router /api/v1/ws/test/push [post]
func (h *Handler) PushTestMessage(c *gin.Context) {
	var req TestPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	level := req.Level
	if level == "" {
		level = "info"
	}

	payload := struct {
		Title   string `json:"title"`
		Message string `json:"message"`
		Level   string `json:"level"`
	}{
		Title:   req.Title,
		Message: req.Message,
		Level:   level,
	}

	msg := struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}{
		Type:    "notification",
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		response.SysError(c, err)
		return
	}

	if req.Target != "" {
		userID, err := strconv.ParseInt(req.Target, 10, 64)
		if err == nil {
			h.hub.SendToUser(userID, data)
			response.Success(c, gin.H{"target": userID, "message": "测试消息已发送"})
			return
		}
	}

	h.hub.Broadcast(data)
	response.Success(c, gin.H{"target": "all", "message": "测试消息已广播到所有在线用户"})
}

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
