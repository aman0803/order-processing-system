package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name        string    `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
}

type Inventory struct {
	gorm.Model
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Quantity  int       `gorm:"not null"`
	Reserved  int       `gorm:"default:0"`
}
