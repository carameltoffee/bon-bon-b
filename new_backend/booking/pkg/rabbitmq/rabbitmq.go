package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQPublisher interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}

type DefaultRMQPublisher struct {
	url     string
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex
}

func NewDefaultRMQPublisher(uri string) *DefaultRMQPublisher {
	return &DefaultRMQPublisher{url: uri}
}

func (p *DefaultRMQPublisher) ensureConnection() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil && !p.conn.IsClosed() && p.channel != nil && !p.channel.IsClosed() {
		return nil
	}

	if p.conn != nil && !p.conn.IsClosed() {
		_ = p.conn.Close()
	}
	var err error
	p.conn, err = amqp.Dial(p.url)
	if err != nil {
		return err
	}
	p.channel, err = p.conn.Channel()
	return err
}

func (p *DefaultRMQPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	if err := p.ensureConnection(); err != nil {
		return err
	}

	return p.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (p *DefaultRMQPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

type MQConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func ConnectRabbitMQ(uri string) (*MQConnection, error) {
	conn, err := amqp.DialConfig(uri, amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &MQConnection{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (mq *MQConnection) Close() {
	if mq.Channel != nil {
		if err := mq.Channel.Close(); err != nil {
			log.Printf("failed to close channel: %v", err)
		}
	}
	if mq.Conn != nil {
		if err := mq.Conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}
}
