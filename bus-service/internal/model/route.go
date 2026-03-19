package model

import (
	"github.com/google/uuid"
	"time"
)

type Route struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Source        string    `gorm:"size:150;not null"`
	Destination   string    `gorm:"size:150;not null"`
	DistanceKm    int       `gorm:"not null"`
	EstimatedTime int       `gorm:"not null"` // in minutes
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
