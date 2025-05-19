package services

import (
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type InventoryService struct {
	db *gorm.DB
}

func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{db: db}
}

func (s *InventoryService) ProcessOrder(msg *kafka.Message) error {
	var event struct {
		Items []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range event.Items {
		var inv struct {
			Quantity int
			Reserved int
		}
		if err := tx.Table("inventories").Where("product_id = ?", item.ProductID).First(&inv).Error; err != nil {
			tx.Rollback()
			return errors.New("product not found")
		}

		if inv.Quantity-inv.Reserved < item.Quantity {
			tx.Rollback()
			return errors.New("insufficient stock")
		}

		if err := tx.Exec("UPDATE inventories SET reserved = reserved + ? WHERE product_id = ?",
			item.Quantity, item.ProductID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
