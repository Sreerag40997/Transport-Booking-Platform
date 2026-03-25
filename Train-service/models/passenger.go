package models

import (
	"time"

	"github.com/google/uuid"
)

// Passenger is one traveller within a TrainBooking.
// Adults and children get a berth (SeatID set).
type Passenger struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID      uuid.UUID       `gorm:"type:uuid;not null;index"`
	Booking        TrainBooking    `gorm:"foreignKey:BookingID"`
	SeatID         *uuid.UUID      `gorm:"type:uuid"` // nullable — nil for infants
	Seat           *TrainInventory `gorm:"foreignKey:SeatID"`
	FirstName      string          `gorm:"size:100;not null"`
	LastName       string          `gorm:"size:100;not null"`
	DateOfBirth    time.Time       `gorm:"type:date;not null"`
	Gender         string          `gorm:"size:10;not null"` // male | female | other
	PassengerType  string          `gorm:"size:10;not null"` // adult | child | infant
	IDType         string          `gorm:"size:20;not null"` // AADHAAR | PAN | PASSPORT
	IDNumber       string          `gorm:"size:50;not null"`
	MealPreference string          `gorm:"size:20"`       // VEG | NON_VEG | JAIN | NONE
	IsPrimary      bool            `gorm:"default:false"` // contact passenger for booking
	CreatedAt      time.Time
}
