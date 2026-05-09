package eventbus

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"time"
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

// maxRetries 同步 handler 失败重试次数
const maxRetries = 2

// retryDelay 重试间隔（指数退避基值）
const retryDelay = 100 * time.Millisecond

// Publish 发布事件（同步执行，失败自动重试，单个 handler 不影响其他）
func (b *Bus) Publish(event interface{}) {
	eventType := getEventType(event)
	b.mu.RLock()
	handlers, ok := b.handlers[eventType]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		b.publishWithRetry(eventType, handler, event)
	}
}

// publishWithRetry 同步执行 handler，失败时重试
func (b *Bus) publishWithRetry(eventType string, handler EventHandler, event interface{}) {
	for i := 0; ; i++ {
		success := func() bool {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[eventbus] panic in handler for %s (attempt %d): %v\n%s",
						eventType, i+1, r, debug.Stack())
				}
			}()
			handler(event)
			return true
		}()

		if success {
			return
		}

		if i >= maxRetries {
			log.Printf("[eventbus] handler for %s failed after %d attempts, giving up",
				eventType, i+1)
			return
		}

		time.Sleep(retryDelay * (1 << i))
	}
}

// PublishAsync 异步发布事件（goroutine 中执行，含重试）
func (b *Bus) PublishAsync(event interface{}) {
	eventType := getEventType(event)
	b.mu.RLock()
	handlers, ok := b.handlers[eventType]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go b.publishWithRetry(eventType, handler, event)
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
