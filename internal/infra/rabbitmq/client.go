package rabbitmq

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	cfg    Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	notify chan *amqp.Error
	mu     sync.Mutex
	closed bool
}

func NewClient(cfg Config) *Client {
	c := &Client{cfg: cfg, closed: false}
	if err := c.connect(); err != nil {
		log.Fatalf("rabbitmq: 连接失败: %v", err)
	}
	c.notify = make(chan *amqp.Error, 1)
	c.conn.NotifyClose(c.notify)
	go c.reconnectLoop()
	return c
}

func (c *Client) connect() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		c.cfg.Username, c.cfg.Password, c.cfg.Host, c.cfg.Port, c.cfg.VHost)

	conn, err := amqp.DialTLS(url, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		// fallback to plain connection
		conn, err = amqp.Dial(url)
		if err != nil {
			return fmt.Errorf("dial rabbitmq: %w", err)
		}
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		c.cfg.Exchange,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // args
	); err != nil {
		conn.Close()
		return fmt.Errorf("declare exchange: %w", err)
	}

	if err := ch.Qos(c.cfg.PrefetchCount, 0, false); err != nil {
		conn.Close()
		return fmt.Errorf("set qos: %w", err)
	}

	c.mu.Lock()
	oldConn := c.conn
	oldCh := c.ch
	c.conn = conn
	c.ch = ch
	c.mu.Unlock()

	// 关闭旧连接（新连接已就绪，不会影响 notify）
	if oldCh != nil {
		oldCh.Close()
	}
	if oldConn != nil {
		oldConn.Close()
	}

	return nil
}

func (c *Client) reconnectLoop() {
	for {
		err := <-c.notify
		c.mu.Lock()
		closed := c.closed
		c.mu.Unlock()
		if closed {
			return
		}

		if err != nil {
			log.Printf("rabbitmq: 连接断开: %v, 准备重连...", err)
		}

		for {
			time.Sleep(3 * time.Second)

			c.mu.Lock()
			closed = c.closed
			c.mu.Unlock()
			if closed {
				return
			}

			if err := c.connect(); err != nil {
				log.Printf("rabbitmq: 重连失败: %v", err)
				continue
			}

			c.mu.Lock()
			c.notify = make(chan *amqp.Error, 1)
			c.conn.NotifyClose(c.notify)
			c.mu.Unlock()

			log.Printf("rabbitmq: 重连成功")
			break
		}
	}
}

func (c *Client) Publish(ctx context.Context, event RoutingKeyer) error {
	msg, err := NewMessage(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	return c.publishMessage(ctx, msg)
}

func (c *Client) PublishWithKey(ctx context.Context, routingKey string, event RoutingKeyer) error {
	msg, err := NewMessage(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	msg.RoutingKey = routingKey
	return c.publishMessage(ctx, msg)
}

func (c *Client) publishMessage(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	c.mu.Lock()
	ch := c.ch
	c.mu.Unlock()

	if ch == nil {
		return fmt.Errorf("rabbitmq: channel 不可用")
	}

	return ch.PublishWithContext(ctx,
		c.cfg.Exchange,
		msg.RoutingKey,
		false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

func (c *Client) Channel() *amqp.Channel {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ch
}

// NewChannel 为消费者创建独立 channel
func (c *Client) NewChannel() (*amqp.Channel, error) {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return nil, fmt.Errorf("rabbitmq: 连接不可用")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("打开 channel 失败: %w", err)
	}
	if err := ch.Qos(c.cfg.PrefetchCount, 0, false); err != nil {
		ch.Close()
		return nil, fmt.Errorf("设置 Qos 失败: %w", err)
	}
	return ch, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	conn := c.conn
	ch := c.ch
	c.conn = nil
	c.ch = nil
	c.mu.Unlock()

	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		// 异步关闭避免等待 broker 回应阻塞退出
		done := make(chan error, 1)
		go func() { done <- conn.Close() }()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
		}
	}
	return nil
}
