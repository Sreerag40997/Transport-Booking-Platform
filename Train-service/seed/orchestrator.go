package seed

import (
	"log"

	"github.com/nabeel-mp/tripneo/train-service/config"
	"gorm.io/gorm"
)

// SeedAll now accepts the application configuration to pass RapidAPI credentials to the seeder
func SeedAll(db *gorm.DB, cfg *config.Config) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Seed standard stations (if you still have local data for this)
		if err := SeedStations(tx); err != nil {
			log.Println("Error seeding stations:", err)
			return err
		}

		// 2. Seed trains directly from the RapidAPI using our new function
		if err := SeedTrainsFromAPI(tx, cfg); err != nil {
			log.Println("Error seeding trains from RapidAPI:", err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Println("Train Seeding completed successfully")
	return nil
}
