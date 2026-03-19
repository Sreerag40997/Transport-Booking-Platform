package model

import (
	"github.com/google/uuid"
	"time"
)

type Seat struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusID      uuid.UUID `gorm:"type:uuid;not null"`
	SeatNumber string    `gorm:"size:20;not null"`
	SeatType   string    `gorm:"size:50;not null"` // window, aisle, sleeper
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
	Bus        Bus
}
