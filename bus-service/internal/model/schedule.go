package model

import (
	"github.com/google/uuid"
	"time"
)

type Schedule struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusID          uuid.UUID `gorm:"type:uuid;not null"`
	RouteID        uuid.UUID `gorm:"type:uuid;not null"`
	DepartureTime  time.Time `gorm:"not null"`
	ArrivalTime    time.Time `gorm:"not null"`
	Price          int64     `gorm:"not null"`
	AvailableSeats int       `gorm:"not null"`
	Status         string    `gorm:"size:30;default:'active'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Bus            Bus
	Route          Route
}
