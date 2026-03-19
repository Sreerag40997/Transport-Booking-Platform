package model

import (
	"github.com/google/uuid"
	"time"
)

type Bus struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name       string    `gorm:"size:150;not null"`
	Number     string    `gorm:"size:50;uniqueIndex;not null"`
	Type       string    `gorm:"size:50;not null"` // AC, Non-AC, Sleeper
	Operator   string    `gorm:"size:150;not null"`
	TotalSeats int       `gorm:"not null"`
	Amenities  string
	Status     string `gorm:"size:30;default:'active'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}
