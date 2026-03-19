package model

import (
	"github.com/google/uuid"
	"time"
)

type BusTracking struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusID     uuid.UUID `gorm:"type:uuid;not null"`
	Latitude  float64   `gorm:"type:decimal(10,6);not null"`
	Longitude float64   `gorm:"type:decimal(10,6);not null"`
	Speed     float64
	UpdatedAt time.Time
	Bus       Bus
}

type Inventory struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ScheduleID     uuid.UUID `gorm:"type:uuid;not null"`
	TotalSeats     int       `gorm:"not null"`
	AvailableSeats int       `gorm:"not null"`
	PurchaseDate   time.Time `gorm:"not null"`
	ValidDate      time.Time `gorm:"not null"`
	CreatedAt      time.Time
	Schedule       Schedule
}

type QrScan struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID uuid.UUID `gorm:"type:uuid;not null"`
	ScannedBy uuid.UUID `gorm:"type:uuid;not null"`
	ScanTime  time.Time `gorm:"autoCreateTime"`
	Status    string    `gorm:"size:30;not null"` // valid, invalid
	Booking   Booking
}
