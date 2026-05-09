package eventbus

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
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

// Publish 发布事件（同步执行，每个 handler 独立 recover，单个失败不影响其他 handler）
func (b *Bus) Publish(event interface{}) {
	eventType := getEventType(event)
	b.mu.RLock()
	handlers, ok := b.handlers[eventType]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[eventbus] panic in handler for %s: %v\n%s",
						eventType, r, debug.Stack())
				}
			}()
			h(event)
		}(handler)
	}
}

// PublishAsync 异步发布事件（fire-and-forget，不保证投递）
func (b *Bus) PublishAsync(event interface{}) {
	eventType := getEventType(event)
	b.mu.RLock()
	handlers, ok := b.handlers[eventType]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[eventbus] async panic in handler for %s: %v\n%s",
						eventType, r, debug.Stack())
				}
			}()
			h(event)
		}(handler)
	}
}

// getEventType 获取事件类型
func getEventType(event interface{}) string {
	t := reflect.TypeOf(event)
	if t == nil {
		return fmt.Sprintf("%T", event)
	}
	return t.String()
}
