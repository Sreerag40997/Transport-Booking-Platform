package models

import (
	"time"

	"github.com/google/uuid"
)

type CancellationPolicy struct {
	ID                   uuid.UUID `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Name                 string    `gorm:"size:100;not null"`
	HoursBeforeDeparture int       `gorm:"not null"`
	RefundPercentage     float64   `gorm:"type:decimal(5,2);not null"`
	CancellationFee      float64   `gorm:"type:decimal(10,2);not null;default:0"`
	IsActive             bool      `gorm:"default:true"`
	CreatedAt            time.Time
}
