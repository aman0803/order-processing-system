package kafka

import (
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) HandleEvent(msg *kafka.Message) error {
	var event struct {
		OrderID string `json:"order_id"`
		UserID  string `json:"user_id"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal: %v", err)
		return err
	}

	log.Printf("Notifying user %s about order %s", event.UserID, event.OrderID)
	return nil
}

type OrderEvent struct {
	OrderID   string      `json:"order_id"`
	UserID    string      `json:"user_id"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	Items     []OrderItem `json:"items"`
	CreatedAt string      `json:"created_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
