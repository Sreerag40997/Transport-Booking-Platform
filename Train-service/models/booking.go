package models

import (
	"time"

	"github.com/google/uuid"
)

// TrainBooking is the core booking record.
// One booking = one user + one schedule + one class + one or more berths.
type TrainBooking struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PNR             string        `gorm:"size:6;uniqueIndex;not null"` // ABC123
	UserID          string        `gorm:"size:36;not null;index"`      // from X-User-ID header
	TrainScheduleID uuid.UUID     `gorm:"type:uuid;not null;index"`
	TrainSchedule   TrainSchedule `gorm:"foreignKey:TrainScheduleID"`
	SeatClass       string        `gorm:"size:5;not null"`                            // SL | 3AC | 2AC | 1AC
	Status          string        `gorm:"size:30;not null;default:'PENDING_PAYMENT'"` // PENDING_PAYMENT|CONFIRMED|CANCELLED|REFUNDED|EXPIRED|FAILED
	BaseFare        float64       `gorm:"type:decimal(10,2);not null"`
	Taxes           float64       `gorm:"type:decimal(10,2);not null;default:0"`
	ServiceFee      float64       `gorm:"type:decimal(10,2);not null;default:0"`
	TotalAmount     float64       `gorm:"type:decimal(10,2);not null"`
	Currency        string        `gorm:"size:3;not null;default:'INR'"`
	PaymentRef      string        `gorm:"size:100"` // Razorpay payment intent ID
	BookedAt        time.Time     `gorm:"default:now()"`
	ConfirmedAt     *time.Time    // pointer = nullable timestamp
	CancelledAt     *time.Time
	ExpiresAt       *time.Time // 15 minutes from creation
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (TrainBooking) TableName() string {
	return "train_bookings"
}

// BookingSeat links one booking to one specific inventory berth.
// Multiple rows per booking for multi-passenger bookings.
type BookingSeat struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Booking   TrainBooking   `gorm:"foreignKey:BookingID"`
	SeatID    uuid.UUID      `gorm:"type:uuid;not null"`
	Seat      TrainInventory `gorm:"foreignKey:SeatID"`
}

func (BookingSeat) TableName() string {
	return "booking_seats"
}
