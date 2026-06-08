package ws

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"eshop-monolith/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const (
	// 消息缓存的Redis键前缀
	msgCacheKeyPrefix = "ws:msg:"
	// 序列ID的Redis键前缀
	seqIDKeyPrefix = "ws:seq:"
	// 默认缓存消息数量
	defaultCacheSize = 1000
	// 缓存过期时间（24小时）
	cacheExpire = 24 * time.Hour
)

// MessageCache 消息缓存管理器
type MessageCache struct {
	redisClient *redis.Client
	cacheSize   int
}

// NewMessageCache 创建消息缓存管理器
func NewMessageCache(redisClient *redis.Client) *MessageCache {
	return &MessageCache{
		redisClient: redisClient,
		cacheSize:   defaultCacheSize,
	}
}

// NextSeqID 获取下一个序列ID
func (mc *MessageCache) NextSeqID(userID int64) (int64, error) {
	key := seqIDKeyPrefix + itoa(userID)
	return mc.redisClient.Incr(context.Background(), key).Result()
}

// GetCurrentSeqID 获取当前序列ID（最新的）
func (mc *MessageCache) GetCurrentSeqID(userID int64) (int64, error) {
	key := seqIDKeyPrefix + itoa(userID)
	val, err := mc.redisClient.Get(context.Background(), key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// StoreMessage 存储消息到缓存
// 使用ZSet存储，score为sequenceId，value为消息JSON
func (mc *MessageCache) StoreMessage(userID int64, seqID int64, message []byte) error {
	key := msgCacheKeyPrefix + itoa(userID)

	ctx := context.Background()
	pipe := mc.redisClient.TxPipeline()

	// 存储消息
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(seqID),
		Member: message,
	})

	// 设置过期时间
	pipe.Expire(ctx, key, cacheExpire)

	// 限制缓存大小（保留最近N条）
	pipe.ZRemRangeByRank(ctx, key, 0, -int64(mc.cacheSize+1))

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error("存储消息失败", "user_id", userID, "seq_id", seqID, "error", err)
	}
	return err
}

// GetMessages 获取指定序列ID范围的消息
// 返回 (lastSeq, currentSeq] 区间的消息
func (mc *MessageCache) GetMessages(userID int64, lastSeq int64, currentSeq int64) ([][]byte, error) {
	key := msgCacheKeyPrefix + itoa(userID)

	// ZRANGEBYSCORE 获取 (lastSeq, currentSeq] 区间的消息
	// min: (lastSeq (exclusive)
	// max: currentSeq (inclusive)
	messages, err := mc.redisClient.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min:    "(" + itoa(lastSeq), // 开区间
		Max:    itoa(currentSeq),     // 闭区间
		Offset: 0,
		Count:  -1, // 返回所有
	}).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.Error("获取消息失败", "user_id", userID, "last_seq", lastSeq, "current_seq", currentSeq, "error", err)
		return nil, err
	}

	result := make([][]byte, 0, len(messages))
	for _, msg := range messages {
		result = append(result, []byte(msg))
	}
	return result, nil
}

// GetCachedSeqRange 获取缓存中消息的序列ID范围
func (mc *MessageCache) GetCachedSeqRange(userID int64) (minSeq, maxSeq int64, err error) {
	key := msgCacheKeyPrefix + itoa(userID)

	ctx := context.Background()

	// 获取最小score
	min, err := mc.redisClient.ZRange(ctx, key, 0, 0).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	// 获取最大score
	max, err := mc.redisClient.ZRange(ctx, key, -1, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	if len(min) == 0 || len(max) == 0 {
		return 0, 0, nil
	}

	// 解析消息获取序列ID
	var minMsg PushMessage
	if err := json.Unmarshal([]byte(min[0]), &minMsg); err != nil {
		return 0, 0, err
	}

	var maxMsg PushMessage
	if err := json.Unmarshal([]byte(max[0]), &maxMsg); err != nil {
		return 0, 0, err
	}

	return minMsg.SequenceID, maxMsg.SequenceID, nil
}

// GetCacheSize 获取缓存大小配置
func (mc *MessageCache) GetCacheSize() int {
	return mc.cacheSize
}

// SetCacheSize 设置缓存大小
func (mc *MessageCache) SetCacheSize(size int) {
	mc.cacheSize = size
}

// itoa 辅助函数：int64转字符串
func itoa(n int64) string {
	return json.Number(strconv.FormatInt(n, 10)).String()
}