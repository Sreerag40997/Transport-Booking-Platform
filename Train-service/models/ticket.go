package models

import (
	"time"

	"github.com/google/uuid"
)

// TrainTicket stores the generated e-ticket after booking confirmation.
type TrainTicket struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID    uuid.UUID    `gorm:"type:uuid;uniqueIndex;not null"` // one ticket per booking
	Booking      TrainBooking `gorm:"foreignKey:BookingID"`
	TicketNumber string       `gorm:"size:20;uniqueIndex;not null"` // TRN-20260414-ABC123
	QRCodeURL    string       `gorm:"not null"`                     // CDN URL to QR image
	QRData       string       `gorm:"type:text;not null"`           // HMAC-signed base64 payload
	IssuedAt     time.Time    `gorm:"default:now()"`
}

func (TrainTicket) TableName() string {
	return "train_tickets"
}
