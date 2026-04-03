package models

import (
	"time"

	"github.com/google/uuid"
)

type Cancellation struct {
	ID              uuid.UUID  `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	BookingID       uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null"`
	Reason          *string    `gorm:"type:text"`
	RefundAmount    float64    `gorm:"type:decimal(10,2);not null"`
	RefundStatus    string     `gorm:"size:20;not null;default:'PENDING'"`
	PolicyAppliedID *uuid.UUID `gorm:"type:uuid"`
	RequestedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	ProcessedAt     *time.Time
	CreatedAt       time.Time

	Booking       Booking             `gorm:"foreignKey:BookingID"`
	PolicyApplied *CancellationPolicy `gorm:"foreignKey:PolicyAppliedID"`
}
