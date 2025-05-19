package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aman0803/order-processing-system/order-service/kafka"
	"github.com/aman0803/order-processing-system/order-service/models"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db            *gorm.DB
	kafkaProducer *kafka.Producer
}

func NewOrderHandler(db *gorm.DB, kafkaProducer *kafka.Producer) *OrderHandler {
	return &OrderHandler{
		db:            db,
		kafkaProducer: kafkaProducer,
	}
}

type CreateOrderRequest struct {
	UserID uuid.UUID          `json:"user_id" binding:"required"`
	Items  []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create order
	order := models.Order{
		UserID: req.UserID,
		Status: models.OrderStatusCreated,
	}

	// Calculate total price and create order items
	var totalPrice float64
	var orderItems []models.OrderItem
	for _, item := range req.Items {
		// In a real application, you would fetch product price from database
		productPrice := 10.0 // Mock price, replace with actual lookup

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productPrice,
		})
		totalPrice += productPrice * float64(item.Quantity)
	}

	order.TotalPrice = totalPrice
	order.Items = orderItems

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse("failed to create order"))
		return
	}

	// Prepare order event for Kafka
	orderEvent := OrderEvent{
		OrderID:   order.ID.String(),
		UserID:    order.UserID.String(),
		Total:     order.TotalPrice,
		Status:    string(order.Status),
		Items:     req.Items,
		CreatedAt: order.CreatedAt,
	}

	eventBytes, err := json.Marshal(orderEvent)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse("failed to prepare order event"))
		return
	}

	// Publish to Kafka using our updated kafka.Producer
	if err := h.kafkaProducer.ProduceOrderEvent(
		order.ID.String(),
		order.UserID.String(),
		eventBytes,
	); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse("failed to publish order event"))
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("failed to commit transaction"))
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse(order))
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("order ID is required"))
		return
	}

	var order models.Order
	if err := h.db.Preload("Items").First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse("order not found"))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(order))
}

type OrderEvent struct {
	OrderID   string             `json:"order_id"`
	UserID    string             `json:"user_id"`
	Total     float64            `json:"total"`
	Status    string             `json:"status"`
	Items     []OrderItemRequest `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
}
