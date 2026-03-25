package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PricingRule defines a dynamic pricing adjustment rule.
type PricingRule struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name       string         `gorm:"size:100;not null"`
	RuleType   string         `gorm:"size:30;not null"`           // DEMAND | TIME_TO_DEPARTURE | SEASONAL
	Conditions datatypes.JSON `gorm:"type:jsonb;not null"`        // {"fill_rate_above": 0.70}
	Multiplier float64        `gorm:"type:decimal(5,3);not null"` // 1.25 = +25% | 0.90 = -10%
	Priority   int            `gorm:"not null;default:0"`
	IsActive   bool           `gorm:"default:true"`
	CreatedAt  time.Time
}

func (PricingRule) TableName() string { return "pricing_rules" }
