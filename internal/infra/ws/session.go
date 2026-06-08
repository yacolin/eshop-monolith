package ws

import (
	"context"
	"encoding/json"
	"time"

	"eshop-monolith/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const (
	// 会话存储的Redis键前缀
	sessionKeyPrefix = "ws:session:"
	// 会话过期时间（7天）
	sessionExpire = 7 * 24 * time.Hour
)

// Session 用户会话信息
type Session struct {
	UserID        int64     `json:"user_id"`
	LastSeq       int64     `json:"last_seq"`
	ConnectedAt   time.Time `json:"connected_at"`
	LastActiveAt  time.Time `json:"last_active_at"`
	ClientIP      string    `json:"client_ip"`
	UserAgent     string    `json:"user_agent"`
	ReconnectCount int      `json:"reconnect_count"`
}

// SessionManager 会话管理器
type SessionManager struct {
	redisClient *redis.Client
}

// NewSessionManager 创建会话管理器
func NewSessionManager(redisClient *redis.Client) *SessionManager {
	return &SessionManager{
		redisClient: redisClient,
	}
}

// SaveSession 保存会话
func (sm *SessionManager) SaveSession(session *Session) error {
	key := sessionKeyPrefix + itoa(session.UserID)

	data, err := json.Marshal(session)
	if err != nil {
		logger.Error("序列化会话失败", "user_id", session.UserID, "error", err)
		return err
	}

	err = sm.redisClient.Set(context.Background(), key, data, sessionExpire).Err()
	if err != nil {
		logger.Error("保存会话失败", "user_id", session.UserID, "error", err)
	}
	return err
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(userID int64) (*Session, error) {
	key := sessionKeyPrefix + itoa(userID)

	data, err := sm.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.Error("获取会话失败", "user_id", userID, "error", err)
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		logger.Error("反序列化会话失败", "user_id", userID, "error", err)
		return nil, err
	}

	return &session, nil
}

// UpdateLastSeq 更新最后接收的序列ID
func (sm *SessionManager) UpdateLastSeq(userID int64, lastSeq int64) error {
	key := sessionKeyPrefix + itoa(userID)

	ctx := context.Background()

	// 获取现有会话
	data, err := sm.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// 会话不存在，创建新会话
			session := &Session{
				UserID:       userID,
				LastSeq:      lastSeq,
				ConnectedAt:  time.Now(),
				LastActiveAt: time.Now(),
			}
			return sm.SaveSession(session)
		}
		return err
	}

	// 更新会话
	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return err
	}

	session.LastSeq = lastSeq
	session.LastActiveAt = time.Now()

	return sm.SaveSession(&session)
}

// UpdateLastActive 更新最后活跃时间
func (sm *SessionManager) UpdateLastActive(userID int64) error {
	key := sessionKeyPrefix + itoa(userID)

	ctx := context.Background()

	data, err := sm.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return err
	}

	session.LastActiveAt = time.Now()

	return sm.SaveSession(&session)
}

// IncrementReconnectCount 增加重连次数
func (sm *SessionManager) IncrementReconnectCount(userID int64) error {
	key := sessionKeyPrefix + itoa(userID)

	ctx := context.Background()

	data, err := sm.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return err
	}

	session.ReconnectCount++
	session.LastActiveAt = time.Now()

	return sm.SaveSession(&session)
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(userID int64) error {
	key := sessionKeyPrefix + itoa(userID)
	return sm.redisClient.Del(context.Background(), key).Err()
}

// IsSessionExpired 检查会话是否过期（超过缓存窗口视为长期离线）
// cacheWindowSeconds: 缓存窗口大小（秒），超过此时间视为长期离线
func (sm *SessionManager) IsSessionExpired(userID int64, cacheWindowSeconds int64) (bool, error) {
	session, err := sm.GetSession(userID)
	if err != nil {
		return false, err
	}

	if session == nil {
		return true, nil
	}

	// 计算最后活跃时间距今的秒数
	secondsSinceLastActive := time.Since(session.LastActiveAt).Seconds()
	return secondsSinceLastActive > float64(cacheWindowSeconds), nil
}

// GetReconnectCount 获取重连次数
func (sm *SessionManager) GetReconnectCount(userID int64) (int, error) {
	session, err := sm.GetSession(userID)
	if err != nil {
		return 0, err
	}

	if session == nil {
		return 0, nil
	}

	return session.ReconnectCount, nil
}