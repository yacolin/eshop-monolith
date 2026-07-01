package dto

// ReconnectRequest 重连接口请求
type ReconnectRequest struct {
	LastSeq int64 `json:"last_seq" binding:"required"` // 客户端最后收到的消息序列号
}

// ReconnectResponse 重连接口响应
type ReconnectResponse struct {
	Status          string        `json:"status"`                   // 状态：ok, sync_required
	Message         string        `json:"message"`                  // 提示信息
	LastSeq         int64         `json:"last_seq"`                 // 客户端上报的序列号
	CurrentSeq      int64         `json:"current_seq"`              // 服务端当前序列号
	NeedFullSync    bool          `json:"need_full_sync"`           // 是否需要全量同步
	NeedIncremental bool          `json:"need_incremental"`         // 是否需要增量同步
	CachedMinSeq    int64         `json:"cached_min_seq,omitempty"` // 缓存最小序列号
	CachedMaxSeq    int64         `json:"cached_max_seq,omitempty"` // 缓存最大序列号
	MessageCount    int           `json:"message_count,omitempty"`  // 补发消息数量
	Messages        []interface{} `json:"messages,omitempty"`       // 需要补发的消息列表
}

// SessionResponse 用户会话信息响应
type SessionResponse struct {
	Exists         bool   `json:"exists"`          // 会话是否存在
	UserID         int64  `json:"user_id"`         // 用户ID
	LastSeq        int64  `json:"last_seq"`        // 最后收到的消息序列号
	ConnectedAt    string `json:"connected_at"`    // 首次连接时间
	LastActiveAt   string `json:"last_active_at"`  // 最后活跃时间
	ReconnectCount int    `json:"reconnect_count"` // 重连次数
}

// OnlineStatsResponse 在线统计响应
type OnlineStatsResponse struct {
	OnlineUsers int `json:"online_users"` // 在线用户数
	Connections int `json:"connections"`  // 连接数
}
