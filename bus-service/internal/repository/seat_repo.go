package repository

import (
	"github.com/Salman-kp/tripneo/bus-service/db"
	"github.com/Salman-kp/tripneo/bus-service/internal/model"
	"gorm.io/gorm"
)

func FindSeatsByBusID(busID string) ([]model.Seat, error) {
	var seats []model.Seat
	err := db.DB.Where("bus_id = ? AND is_active = true", busID).Find(&seats).Error
	return seats, err
}

// DecrementSeat uses a DB-level atomic update to prevent overselling
func DecrementAvailableSeats(scheduleID string, count int) error {
	return db.DB.Model(&model.Schedule{}).
		Where("id = ? AND available_seats >= ?", scheduleID, count).
		UpdateColumn("available_seats", gorm.Expr("available_seats - ?", count)).
		Error
}

func IncrementAvailableSeats(scheduleID string, count int) error {
	return db.DB.Model(&model.Schedule{}).
		Where("id = ?", scheduleID).
		UpdateColumn("available_seats", gorm.Expr("available_seats + ?", count)).
		Error
}
