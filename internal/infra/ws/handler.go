package ws

import (
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

	// 查询用户信息用于在线事件广播
	if h.hub.userInfoProvider != nil {
		username, nickname, lookupErr := h.hub.userInfoProvider(userID)
		if lookupErr == nil {
			client.Username = username
			client.Nickname = nickname
		} else {
			logger.Warn("查询用户信息失败", "user_id", userID, "error", lookupErr)
		}
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
