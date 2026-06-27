package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message 通用消息包装
type Message struct {
	ID         string          `json:"id"`
	RoutingKey string          `json:"routing_key"`
	Timestamp  int64           `json:"timestamp"`
	Payload    json.RawMessage `json:"payload"`
}

// RoutingKeyer 事件通过此接口获取 routing key
type RoutingKeyer interface {
	RoutingKey() string
}

// NewMessage 从事件创建消息
func NewMessage(event RoutingKeyer) (Message, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return Message{}, err
	}
	return Message{
		ID:         uuid.New().String(),
		RoutingKey: event.RoutingKey(),
		Timestamp:  time.Now().UnixMilli(),
		Payload:    payload,
	}, nil
}
