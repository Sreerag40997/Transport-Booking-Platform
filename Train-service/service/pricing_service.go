package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nabeel-mp/tripneo/train-service/config"
	"github.com/nabeel-mp/tripneo/train-service/models"
	"gorm.io/gorm"
)

// RunPricingEngine runs every N minutes (set by PRICING_ENGINE_INTERVAL_MINUTES).
// Reads active pricing_rules and recalculates train_inventory.price on upcoming schedules.
func RunPricingEngine(db *gorm.DB, cfg *config.Config) {
	interval := time.Duration(cfg.PricingEngineIntervalMins) * time.Minute
	log.Printf("[pricing-engine] Started — running every %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on startup
	runPricingPass(db)

	for range ticker.C {
		runPricingPass(db)
	}
}

func runPricingPass(db *gorm.DB) {
	log.Println("[pricing-engine] Running pricing recalculation...")

	// Fetch all active rules ordered by priority
	var rules []models.PricingRule
	if err := db.Where("is_active = true").Order("priority ASC").Find(&rules).Error; err != nil {
		log.Printf("[pricing-engine] Failed to fetch rules: %v", err)
		return
	}
	if len(rules) == 0 {
		return
	}

	// Fetch upcoming schedules (next 30 days)
	cutoff := time.Now().AddDate(0, 0, 30)
	var schedules []models.TrainSchedule
	if err := db.Where("departure_at > ? AND departure_at < ? AND status != 'CANCELLED'",
		time.Now(), cutoff).Find(&schedules).Error; err != nil {
		log.Printf("[pricing-engine] Failed to fetch schedules: %v", err)
		return
	}

	updated := 0
	for _, schedule := range schedules {
		// Fetch all inventory for this schedule
		var inventory []models.TrainInventory
		if err := db.Where("train_schedule_id = ? AND status = 'AVAILABLE'", schedule.ID).
			Find(&inventory).Error; err != nil {
			continue
		}
		if len(inventory) == 0 {
			continue
		}

		// Compute total capacity per class for fill-rate calculation
		totalByClass := map[string]int{}
		soldByClass := map[string]int{}
		for _, item := range inventory {
			totalByClass[item.Class]++
			if item.Status == "BOOKED" {
				soldByClass[item.Class]++
			}
		}

		for _, item := range inventory {
			multiplier := applyRules(rules, item, schedule, totalByClass, soldByClass)
			newPrice := item.WholesalePrice * 1.15 * multiplier // base margin 15%
			// Round to 2 decimal places
			newPrice = float64(int(newPrice*100+0.5)) / 100

			if err := db.Model(&models.TrainInventory{}).
				Where("id = ?", item.ID).
				Update("price", newPrice).Error; err != nil {
				log.Printf("[pricing-engine] Failed to update price for seat %s: %v", item.ID, err)
			} else {
				updated++
			}
		}
	}
	log.Printf("[pricing-engine] Updated %d seat prices", updated)
}

// applyRules applies all matching pricing rules and stacks their multipliers.
func applyRules(rules []models.PricingRule, item models.TrainInventory, schedule models.TrainSchedule,
	totalByClass, soldByClass map[string]int) float64 {

	multiplier := 1.0
	daysBeforeDep := time.Until(schedule.DepartureAt).Hours() / 24

	for _, rule := range rules {
		var conditions map[string]interface{}
		if err := json.Unmarshal(rule.Conditions, &conditions); err != nil {
			continue
		}

		switch rule.RuleType {
		case "DEMAND":
			total := totalByClass[item.Class]
			sold := soldByClass[item.Class]
			if total > 0 {
				fillRate := float64(sold) / float64(total)
				if threshold, ok := conditions["fill_rate_above"].(float64); ok {
					if fillRate > threshold {
						multiplier *= float64(rule.Multiplier)
					}
				}
			}
		case "TIME_TO_DEPARTURE":
			if above, ok := conditions["days_before_above"].(float64); ok {
				if daysBeforeDep > above {
					multiplier *= float64(rule.Multiplier)
				}
			}
			if below, ok := conditions["days_before_below"].(float64); ok {
				if daysBeforeDep < below {
					multiplier *= float64(rule.Multiplier)
				}
			}
		case "SEASONAL":
			if months, ok := conditions["months"].([]interface{}); ok {
				currentMonth := int(time.Now().Month())
				for _, m := range months {
					if int(m.(float64)) == currentMonth {
						multiplier *= float64(rule.Multiplier)
						break
					}
				}
			}
		}
	}
	return multiplier
}
