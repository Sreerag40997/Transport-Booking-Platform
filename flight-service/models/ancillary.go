package models

import (
	"time"

	"github.com/google/uuid"
)

type Ancillary struct {
	ID          uuid.UUID  `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	BookingID   uuid.UUID  `gorm:"type:uuid;not null"`
	PassengerID *uuid.UUID `gorm:"type:uuid"`
	Type        string     `gorm:"size:30;not null"`
	Description string     `gorm:"size:200;not null"`
	Price       float64    `gorm:"type:decimal(10,2);not null"`
	Quantity    int        `gorm:"not null;default:1"`
	CreatedAt   time.Time
}
