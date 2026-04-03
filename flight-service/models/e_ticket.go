package models

import (
	"time"

	"github.com/google/uuid"
)

type ETicket struct {
	ID           uuid.UUID `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	BookingID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	TicketNumber string    `gorm:"size:20;uniqueIndex;not null"`
	QRCodeURL    string    `gorm:"type:text;not null"`
	QRData       string    `gorm:"type:text;not null"`
	IssuedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Booking Booking `gorm:"foreignKey:BookingID"`
}
