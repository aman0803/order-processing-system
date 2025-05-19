package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusProcessed OrderStatus = "PROCESSED"
	OrderStatusFailed    OrderStatus = "FAILED"
)

type Order struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
	TotalPrice float64     `gorm:"not null"`
	Status     OrderStatus `gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`
}
