package models

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID             uuid.UUID `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	UserID         uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	CompanyName    string    `gorm:"size:200;not null"`
	LicenseNumber  string    `gorm:"size:100"`
	CommissionRate float64   `gorm:"type:decimal(5,2);not null;default:5.00"`
	CreditLimit    float64   `gorm:"type:decimal(12,2);not null;default:0"`
	CreditUsed     float64   `gorm:"type:decimal(12,2);not null;default:0"`
	Status         string    `gorm:"size:20;not null;default:'PENDING'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
