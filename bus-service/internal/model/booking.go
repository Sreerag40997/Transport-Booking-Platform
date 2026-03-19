package model

import (
	"github.com/google/uuid"
	"time"
)

type Booking struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	ScheduleID    uuid.UUID `gorm:"type:uuid;not null"`
	TotalAmount   int64     `gorm:"not null"`
	Status        string    `gorm:"size:30;not null"` // confirmed, cancelled
	PaymentStatus string    `gorm:"size:30;not null"` // pending, success, failed
	PaymentRefID  string    `gorm:"size:150"`
	QRCode        string    `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Schedule      Schedule
	BookingSeats  []BookingSeat
}

type BookingSeat struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID       uuid.UUID `gorm:"type:uuid;not null"`
	SeatID          uuid.UUID `gorm:"type:uuid;not null"`
	PassengerName   string    `gorm:"size:150;not null"`
	PassengerAge    int       `gorm:"not null"`
	PassengerGender string    `gorm:"size:20;not null"`
	CreatedAt       time.Time
	Seat            Seat
}
