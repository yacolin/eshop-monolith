package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerConfig struct {
	Queue      string
	Bindings   []string // routing key 绑定列表
	Prefetch   int
	RetryLimit int
}

type Consumer struct {
	client   *Client
	ch       *amqp.Channel
	cfg      ConsumerConfig
	handlers map[string]func(Message) error
	mu       sync.RWMutex
}

func NewConsumer(client *Client, cfg ConsumerConfig) *Consumer {
	// 创建独立 channel，避免多消费者共用导致消息分发异常
	if cfg.Prefetch <= 0 {
		cfg.Prefetch = client.cfg.PrefetchCount
	}
	if cfg.RetryLimit <= 0 {
		cfg.RetryLimit = client.cfg.RetryLimit
	}
		ch, err := client.NewChannel()
	if err != nil {
		return nil
	}
	return &Consumer{
		client:   client,
		ch:       ch,
		cfg:      cfg,
		handlers: make(map[string]func(Message) error),
	}
}

func (c *Consumer) HandleFunc(key string, fn func(Message) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[key] = fn
}

func (c *Consumer) Start(ctx context.Context) error {
	ch := c.ch
	if ch == nil {
		return fmt.Errorf("rabbitmq: channel 不可用")
	}

	// 1. 声明死信队列
	dlqName := c.cfg.Queue + ".dlq"
	_, err := ch.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare dlq %s: %w", dlqName, err)
	}
	// 绑定 dlq 到 exchange（便于查看原始消息）
	if err := ch.QueueBind(dlqName, "#", c.client.cfg.Exchange, false, nil); err != nil {
		return fmt.Errorf("bind dlq %s: %w", dlqName, err)
	}

	// 2. 声明主队列，死信指向 dlq
	mainArgs := amqp.Table{
		"x-dead-letter-exchange":    c.client.cfg.Exchange,
		"x-dead-letter-routing-key": dlqName,
	}
	_, err = ch.QueueDeclare(c.cfg.Queue, true, false, false, false, mainArgs)
	if err != nil {
		return fmt.Errorf("declare queue %s: %w", c.cfg.Queue, err)
	}

	// 3. 绑定 routing keys
	for _, key := range c.cfg.Bindings {
		if err := ch.QueueBind(c.cfg.Queue, key, c.client.cfg.Exchange, false, nil); err != nil {
			return fmt.Errorf("bind %s -> %s: %w", key, c.cfg.Queue, err)
		}
	}

	// 4. 声明重试队列（TTL 后回到主 exchange）
	retryQueue := c.cfg.Queue + ".retry"
	retryArgs := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": c.cfg.Queue,
		"x-message-ttl":             int32(c.client.cfg.RetryDelayMs),
	}
	_, err = ch.QueueDeclare(retryQueue, true, false, false, false, retryArgs)
	if err != nil {
		return fmt.Errorf("declare retry queue %s: %w", retryQueue, err)
	}
	// 绑定重试队列到 exchange
	if err := ch.QueueBind(retryQueue, retryQueue, c.client.cfg.Exchange, false, nil); err != nil {
		return fmt.Errorf("bind retry queue %s: %w", retryQueue, err)
	}

	// 5. 开始消费
	deliveries, err := ch.Consume(c.cfg.Queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume %s: %w", c.cfg.Queue, err)
	}

	go c.consumeLoop(ctx, deliveries)
	log.Printf("rabbitmq: 消费者 [%s] 已启动, 绑定: %v", c.cfg.Queue, c.cfg.Bindings)
	return nil
}

func (c *Consumer) consumeLoop(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				return
			}
			c.processDelivery(d)
		}
	}
}

func (c *Consumer) processDelivery(d amqp.Delivery) {
	var msg Message
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		log.Printf("rabbitmq [%s]: 消息解析失败, 转入 dlq: %v", c.cfg.Queue, err)
		d.Nack(false, false)
		return
	}

	c.mu.RLock()
	handler, ok := c.handlers[msg.RoutingKey]
	if !ok {
		handler, ok = c.handlers["*"]
	}
	c.mu.RUnlock()

	if !ok {
		d.Ack(false)
		return
	}

	if err := handler(msg); err != nil {
		retryInfo := extractRetryInfo(d)
		if retryInfo.RetryCount >= c.cfg.RetryLimit {
			log.Printf("rabbitmq [%s]: 消息 %s 重试 %d 次仍失败, 转入 dlq: %v",
				c.cfg.Queue, msg.ID, retryInfo.RetryCount, err)
			d.Nack(false, false)
			return
		}
		c.sendToRetry(msg, retryInfo.RetryCount+1, err.Error())
		d.Ack(false)
		return
	}

	d.Ack(false)
}

type retryHeader struct {
	RetryCount int
	LastError  string
}

func extractRetryInfo(d amqp.Delivery) retryHeader {
	h := retryHeader{}
	if d.Headers != nil {
		if count, ok := d.Headers["x-retry-count"].(int32); ok {
			h.RetryCount = int(count)
		}
		if errStr, ok := d.Headers["x-retry-error"].(string); ok {
			h.LastError = errStr
		}
	}
	return h
}

func (c *Consumer) sendToRetry(msg Message, retryCount int, lastErr string) {
	retryQueue := c.cfg.Queue + ".retry"
	body, _ := json.Marshal(msg)

	headers := amqp.Table{
		"x-retry-count": int32(retryCount),
		"x-retry-error": lastErr,
	}

	if err := c.ch.Publish(
		c.client.cfg.Exchange,
		retryQueue,
		false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
			Timestamp:    time.Now(),
		},
	); err != nil {
		log.Printf("rabbitmq [%s]: 发送到重试队列失败: %v", c.cfg.Queue, err)
	}
}
