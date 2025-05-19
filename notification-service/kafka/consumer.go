package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// MessageHandler defines the processing function signature
type MessageHandler func(msg *kafka.Message) error

// Consumer wraps the Kafka consumer
type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewConsumer(brokers string, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	// Create a new reader with the given configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		GroupID:     groupID,
		GroupTopics: topics,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     time.Second,
		StartOffset: kafka.FirstOffset, // equivalent to "earliest"
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		reader:  reader,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (c *Consumer) Consume() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(c.ctx)
			if err != nil {
				log.Printf("Consumer error: %v\n", err)
				continue
			}

			if err := c.handler(&msg); err != nil {
				log.Printf("Handler error: %v\n", err)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.cancel()
	if err := c.reader.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}
