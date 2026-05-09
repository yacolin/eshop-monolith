package eventbus

import (
	"reflect"
	"sync"
)

// Bus 事件总线
type Bus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// EventHandler 事件处理器
type EventHandler func(event interface{})

// NewBus 创建事件总线
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (b *Bus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish 发布事件
func (b *Bus) Publish(event interface{}) {
	eventType := getEventType(event)
	b.mu.RLock()
	handlers, ok := b.handlers[eventType]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(event)
	}
}

// getEventType 获取事件类型
func getEventType(event interface{}) string {
	return reflect.TypeOf(event).String()
}
