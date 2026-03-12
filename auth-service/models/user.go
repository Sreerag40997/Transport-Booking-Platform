package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"size:254;uniqueIndex;not null"`
	PasswordHash string
	Role         string `gorm:"default:'user'"`
	IsVerified   bool   `gorm:"default:false"`
	CreatedAt    time.Time
}
