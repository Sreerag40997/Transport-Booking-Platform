package models

import (
	"time"

	"github.com/google/uuid"
)

type AgentInventory struct {
	ID                 uuid.UUID `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	AgentID            uuid.UUID `gorm:"type:uuid;not null"`
	FlightInstanceID   uuid.UUID `gorm:"type:uuid;not null"`
	FareTypeID         uuid.UUID `gorm:"type:uuid;not null"`
	SeatClass          string    `gorm:"size:20;not null"`
	QuantityPurchased  int       `gorm:"not null"`
	QuantitySold       int       `gorm:"not null;default:0"`
	WholesalePrice     float64   `gorm:"type:decimal(10,2);not null"`
	SellingPrice       float64   `gorm:"type:decimal(10,2);not null"`
	Status             string    `gorm:"size:20;not null;default:'ACTIVE'"`
	PurchasedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ExpiresAt          time.Time `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Agent          Agent          `gorm:"foreignKey:AgentID"`
	FlightInstance FlightInstance `gorm:"foreignKey:FlightInstanceID"`
	FareType       FareType       `gorm:"foreignKey:FareTypeID"`
}
