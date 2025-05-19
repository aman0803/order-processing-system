package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        "orders",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return &Producer{
		writer: writer,
		topic:  "orders",
	}, nil
}

func (p *Producer) ProduceOrderEvent(orderID string, userID string, payload []byte) error {
	headers := []kafka.Header{
		{Key: "user_id", Value: []byte(userID)},
	}

	msg := kafka.Message{
		Key:     []byte(orderID),
		Value:   payload,
		Headers: headers,
		Time:    time.Now(),
	}

	ctx := context.Background()
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	}
}
