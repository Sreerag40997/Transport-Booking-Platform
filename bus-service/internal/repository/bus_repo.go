package repository

import (
	"github.com/Salman-kp/tripneo/bus-service/db"
	"github.com/Salman-kp/tripneo/bus-service/internal/model"
)

func SearchBuses(src, dest, date, busType string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	q := db.DB.Preload("Bus").Preload("Route").
		Joins("JOIN buses ON buses.id = schedules.bus_id").
		Joins("JOIN routes ON routes.id = schedules.route_id").
		Where("routes.source = ? AND routes.destination = ?", src, dest).
		Where("DATE(schedules.departure_time) = ?", date).
		Where("schedules.available_seats > 0").
		Where("schedules.status = 'active'")
	if busType != "" {
		q = q.Where("buses.type = ?", busType)
	}
	return schedules, q.Find(&schedules).Error
}

func FindBusByID(id string) (*model.Bus, error) {
	var bus model.Bus
	err := db.DB.Where("id = ?", id).First(&bus).Error
	return &bus, err
}
