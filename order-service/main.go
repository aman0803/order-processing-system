package main

import (
	"log"
	"os"

	"github.com/aman0803/order-processing-system/order-service/handlers"
	"github.com/aman0803/order-processing-system/order-service/kafka"
	"github.com/aman0803/order-processing-system/order-service/models"

	"github.com/gin-gonic/gin"
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
	if err := db.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize kafka producer
	kafkaProducer, err := kafka.NewProducer(os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		log.Fatalf("Failed to create kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Create router
	router := gin.Default()

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(db, kafkaProducer)

	// Define routes
	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders/:id", orderHandler.GetOrder)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
