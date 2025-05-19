package main

import (
	"log"
	"os"

	"github.com/aman0803/order-processing-system/inventory-service/kafka"
	"github.com/aman0803/order-processing-system/inventory-service/models"
	"github.com/aman0803/order-processing-system/inventory-service/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize database
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate models
	if err := db.AutoMigrate(&models.Product{}, &models.Inventory{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize inventory service
	inventoryService := services.NewInventoryService(db)

	// Initialize Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(
		os.Getenv("KAFKA_BROKERS"),
		"inventory-group",
		[]string{"orders"},
		inventoryService.ProcessOrder,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	// Start consuming messages
	log.Println("Inventory service started. Waiting for messages...")
	kafkaConsumer.Consume()
}
