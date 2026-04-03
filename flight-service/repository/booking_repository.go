package repository

import (
	"github.com/junaid9001/tripneo/flight-service/models"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateBooking(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) GetBookingByID(id string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Passengers").Preload("Ancillaries").Preload("FlightInstance").Preload("FareType").First(&booking, "id = ?", id).Error
	return &booking, err
}

func (r *BookingRepository) GetBookingByPNR(pnr string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Passengers").Preload("Ancillaries").Preload("FlightInstance").Preload("FareType").First(&booking, "pnr = ?", pnr).Error
	return &booking, err
}

func (r *BookingRepository) GetBookingsByUserID(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("FlightInstance").Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) UpdateBooking(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) CreateCancellation(cancellation *models.Cancellation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cancellation).Error; err != nil {
			return err
		}
		return tx.Model(&models.Booking{}).Where("id = ?", cancellation.BookingID).Update("status", "CANCELLED").Error
	})
}

func (r *BookingRepository) SaveETicket(ticket *models.ETicket) error {
	return r.db.Create(ticket).Error
}

func (r *BookingRepository) GetETicketByBookingID(bookingID string) (*models.ETicket, error) {
	var ticket models.ETicket
	err := r.db.First(&ticket, "booking_id = ?", bookingID).Error
	return &ticket, err
}

func (r *BookingRepository) GetFlightInstanceByID(id string) (*models.FlightInstance, error) {
	var instance models.FlightInstance
	err := r.db.First(&instance, "id = ?", id).Error
	return &instance, err
}

func (r *BookingRepository) GetFareTypeByID(id string) (*models.FareType, error) {
	var fare models.FareType
	err := r.db.First(&fare, "id = ?", id).Error
	return &fare, err
}
