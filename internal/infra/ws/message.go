package ws

import (
	"encoding/json"
	"strconv"
	"time"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/trade"
	flashEvents "eshop-monolith/internal/flashsale/events"
)

// PushMessage WebSocket 推送消息
type PushMessage struct {
	Type        string      `json:"type"`
	SequenceID  int64       `json:"sequence_id"`
	Timestamp   int64       `json:"timestamp"`
	Data        interface{} `json:"data"`
}

// NewPushMessage 从事件创建推送消息
func NewPushMessage(event interface{}, sequenceID int64) *PushMessage {
	now := time.Now().UnixMilli()

	switch e := event.(type) {
	case trade.OrderPaidEvent:
		return &PushMessage{Type: "order.paid", SequenceID: sequenceID, Timestamp: now, Data: e}
	case trade.OrderShippedEvent:
		return &PushMessage{Type: "order.shipped", SequenceID: sequenceID, Timestamp: now, Data: e}
	case trade.OrderDeliveredEvent:
		return &PushMessage{Type: "order.delivered", SequenceID: sequenceID, Timestamp: now, Data: e}
	case trade.OrderCancelledEvent:
		return &PushMessage{Type: "order.cancelled", SequenceID: sequenceID, Timestamp: now, Data: e}

	case trade.PaymentSuccessEvent:
		return &PushMessage{Type: "payment.success", SequenceID: sequenceID, Timestamp: now, Data: e}
	case trade.PaymentFailedEvent:
		return &PushMessage{Type: "payment.failed", SequenceID: sequenceID, Timestamp: now, Data: e}
	case trade.RefundCreatedEvent:
		return &PushMessage{Type: "refund.created", SequenceID: sequenceID, Timestamp: now, Data: e}

	case flashEvents.FlashOrderCreatedEvent:
		return &PushMessage{Type: "flash.created", SequenceID: sequenceID, Timestamp: now, Data: e}
	case flashEvents.FlashOrderPaidEvent:
		return &PushMessage{Type: "flash.paid", SequenceID: sequenceID, Timestamp: now, Data: e}
	case flashEvents.FlashOrderCancelledEvent:
		return &PushMessage{Type: "flash.cancelled", SequenceID: sequenceID, Timestamp: now, Data: e}

	case inventory.InventoryLowEvent:
		return &PushMessage{Type: "inventory.low", SequenceID: sequenceID, Timestamp: now, Data: e}
	}

	return nil
}

// Marshal 序列化为 JSON
func (m *PushMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal 从JSON反序列化
func (m *PushMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

// extractUserID 从事件中提取 user_id
func extractUserID(event interface{}) int64 {
	switch e := event.(type) {
	case trade.OrderPaidEvent:
		return parseCustomerID(e.CustomerID)
	case trade.OrderShippedEvent:
		return parseCustomerID(e.CustomerID)
	case trade.OrderDeliveredEvent:
		return parseCustomerID(e.CustomerID)
	case trade.OrderCancelledEvent:
		return parseCustomerID(e.CustomerID)

	case flashEvents.FlashOrderCreatedEvent:
		return e.UserID
	case flashEvents.FlashOrderPaidEvent:
		return e.UserID
	case flashEvents.FlashOrderCancelledEvent:
		return e.UserID

	case trade.PaymentSuccessEvent:
		return 0
	case trade.PaymentFailedEvent:
		return 0
	case trade.RefundCreatedEvent:
		return 0
	}

	return 0
}

// parseCustomerID 将 CustomerID (string) 解析为 int64
func parseCustomerID(customerID string) int64 {
	if customerID == "" {
		return 0
	}
	id, err := strconv.ParseInt(customerID, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// extractUserIDFromMessage 从 RabbitMQ 消息中提取用户ID
func extractUserIDFromMessage(msg rabbitmq.Message) int64 {
	var data struct {
		CustomerID string `json:"customer_id"`
		UserID     int64  `json:"user_id"`
	}
	json.Unmarshal(msg.Payload, &data)
	if data.UserID > 0 {
		return data.UserID
	}
	id, _ := strconv.ParseInt(data.CustomerID, 10, 64)
	return id
}

// NewPushMessageFromRaw 从原始消息创建推送消息
func NewPushMessageFromRaw(routingKey string, payload json.RawMessage, seqID int64) *PushMessage {
	return &PushMessage{
		Type:       routingKey,
		SequenceID: seqID,
		Timestamp:  time.Now().UnixMilli(),
		Data:       payload,
	}
}

// RealtimeMessage 实时推送消息（全局广播，不按用户缓存）
type RealtimeMessage struct {
	Seq     int64       `json:"seq"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// NewRealtimeMessage 创建实时推送消息
func NewRealtimeMessage(msgType string, payload interface{}) *RealtimeMessage {
	return &RealtimeMessage{
		Type:    msgType,
		Payload: payload,
	}
}

// Marshal 序列化为 JSON
func (m *RealtimeMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UserBrief 用户简要信息
type UserBrief struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// UserEventPayload 用户上下线事件负载
type UserEventPayload struct {
	Action    string    `json:"action"` // "online" / "offline"
	User      UserBrief `json:"user"`
	Timestamp int64     `json:"timestamp"`
}

// StatsPayload 在线统计数据负载
type StatsPayload struct {
	OnlineUsers int `json:"online_users"`
	Connections int `json:"connections"`
}

// UserSyncPayload 全量用户列表同步负载（连接时推送，前端全量替换）
type UserSyncPayload struct {
	Action string      `json:"action"` // "sync"
	Users  []UserBrief `json:"users"`
}

// SystemMessage 系统消息类型
type SystemMessage struct {
	Type        string `json:"type"`
	SequenceID  int64  `json:"sequence_id"`
	Message     string `json:"message"`
	RequireFullSync bool `json:"require_full_sync"`
}

// NewSystemMessage 创建系统消息
func NewSystemMessage(msgType string, message string, requireFullSync bool) *SystemMessage {
	return &SystemMessage{
		Type:        msgType,
		SequenceID:  0,
		Message:     message,
		RequireFullSync: requireFullSync,
	}
}

// Marshal 序列化为JSON
func (m *SystemMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}