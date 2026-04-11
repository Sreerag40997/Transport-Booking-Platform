package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/tripneo/flight-service/dto"
	"github.com/junaid9001/tripneo/flight-service/models"
	"github.com/junaid9001/tripneo/flight-service/redis"
	"github.com/junaid9001/tripneo/flight-service/repository"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo}
}

func generatePNR() string {
	b := make([]byte, 3)
	rand.Read(b)
	return "TR" + hex.EncodeToString(b)[:4]
}

func (s *BookingService) CreateBooking(userID string, req *dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	if req.FlightInstanceID == "" || req.FareTypeID == "" {
		return nil, errors.New("mandatory flight or fare ID missing")
	}

	flightInstance, err := s.repo.GetFlightInstanceByID(req.FlightInstanceID)
	if err != nil {
		return nil, errors.New("invalid flight instance")
	}

	_, err = s.repo.GetFareTypeByID(req.FareTypeID)
	if err != nil {
		return nil, errors.New("invalid fare type")
	}

	if len(req.Passengers) > 9 {
		return nil, errors.New("maximum 9 passengers allowed per booking")
	}

	bdUserId, _ := uuid.Parse(userID)
	fiId, _ := uuid.Parse(req.FlightInstanceID)
	ftId, _ := uuid.Parse(req.FareTypeID)

	var baseFare float64 = 0
	if req.SeatClass == "BUSINESS" {
		baseFare = flightInstance.CurrentPriceBusiness
	} else {
		baseFare = flightInstance.CurrentPriceEconomy
	}

	totalBase := baseFare * float64(len(req.Passengers))

	var ancillaries []models.Ancillary
	var ancTotal float64 = 0
	for _, a := range req.Ancillaries {
		ancTotal += a.Price * float64(a.Quantity)
		ancillaries = append(ancillaries, models.Ancillary{
			Type:        a.Type,
			Description: a.Description,
			Price:       a.Price,
			Quantity:    a.Quantity,
		})
	}

	taxes := totalBase * 0.18 // 18% tax
	serviceFee := 500.0
	totalAmount := totalBase + taxes + serviceFee + ancTotal

	var passengers []models.Passenger
	var lockedSeats []string
	ctx := context.Background()

	// safe cleanup for failed transactions
	lockSucceeded := false
	defer func() {
		if !lockSucceeded {
			for _, ls := range lockedSeats {
				_ = redis.ReleaseSeatLock(ctx, ls)
			}
		}
	}()

	for _, p := range req.Passengers {
		dob, _ := time.Parse("2006-01-02", p.DateOfBirth)
		var sId *uuid.UUID = nil

		// infants don't occupy seats
		if p.PassengerType == "infant" {
			sId = nil
		} else if p.SeatID != "" {
			seat, err := s.repo.GetSeatByID(p.SeatID)
			if err != nil || !seat.IsAvailable {
				return nil, errors.New("seat already permanently booked or currently held")
			}

			// 10 minute hold
			acquired, err := redis.AcquireSeatLock(ctx, p.SeatID, userID, 10*time.Minute)
			if err != nil || !acquired {
				return nil, errors.New("seat is currently held by another user")
			}
			lockedSeats = append(lockedSeats, p.SeatID)

			sidVal, _ := uuid.Parse(p.SeatID)
			sId = &sidVal
		}

		var mp *string = nil
		if p.MealPreference != "" {
			mp = &p.MealPreference
		}

		passengers = append(passengers, models.Passenger{
			FirstName:      p.FirstName,
			LastName:       p.LastName,
			DateOfBirth:    dob,
			Gender:         p.Gender,
			PassengerType:  p.PassengerType,
			IDType:         p.IDType,
			IDNumber:       p.IDNumber,
			SeatID:         sId,
			MealPreference: mp,
		})
	}

	pnr := generatePNR()
	expiresAt := time.Now().Add(10 * time.Minute)
	var gstin *string = nil
	if req.GSTIN != "" {
		gstin = &req.GSTIN
	}

	booking := &models.Booking{
		PNR:              pnr,
		UserID:           bdUserId,
		FlightInstanceID: fiId,
		FareTypeID:       ftId,
		Source:           "live", // default to live for now
		TripType:         req.TripType,
		SeatClass:        req.SeatClass,
		Status:           "PENDING_PAYMENT",
		BaseFare:         totalBase,
		Taxes:            taxes,
		ServiceFee:       serviceFee,
		AncillariesTotal: ancTotal,
		TotalAmount:      totalAmount,
		Currency:         "INR",
		GSTIN:            gstin,
		ExpiresAt:        &expiresAt,
		Passengers:       passengers,
		Ancillaries:      ancillaries,
	}

	if err := s.repo.CreateBooking(booking); err != nil {
		return nil, err
	}

	lockSucceeded = true

	return mapBookingToDTO(booking), nil
}

func mapBookingToDTO(booking *models.Booking) *dto.BookingResponse {
	var passengers []dto.PassengerDto
	for _, p := range booking.Passengers {
		var sId string
		if p.SeatID != nil {
			sId = p.SeatID.String()
		}
		var mp string
		if p.MealPreference != nil {
			mp = *p.MealPreference
		}
		passengers = append(passengers, dto.PassengerDto{
			FirstName:      p.FirstName,
			LastName:       p.LastName,
			DateOfBirth:    p.DateOfBirth.Format("2006-01-02"),
			Gender:         p.Gender,
			PassengerType:  p.PassengerType,
			IDType:         p.IDType,
			IDNumber:       p.IDNumber,
			SeatID:         sId,
			MealPreference: mp,
		})
	}

	var ancillaries []dto.AncillaryBookDto
	for _, a := range booking.Ancillaries {
		ancillaries = append(ancillaries, dto.AncillaryBookDto{
			Type:        a.Type,
			Description: a.Description,
			Price:       a.Price,
			Quantity:    a.Quantity,
		})
	}

	var checkedExpiresAt string
	if booking.ExpiresAt != nil {
		checkedExpiresAt = booking.ExpiresAt.Format(time.RFC3339)
	}

	return &dto.BookingResponse{
		ID:               booking.ID.String(),
		PNR:              booking.PNR,
		FlightInstanceID: booking.FlightInstanceID.String(),
		Source:           booking.Source,
		Status:           booking.Status,
		SeatClass:        booking.SeatClass,
		TripType:         booking.TripType,
		BaseFare:         booking.BaseFare,
		Taxes:            booking.Taxes,
		ServiceFee:       booking.ServiceFee,
		AncillariesTotal: booking.AncillariesTotal,
		TotalAmount:      booking.TotalAmount,
		Currency:         booking.Currency,
		BookedAt:         booking.BookedAt.Format(time.RFC3339),
		ExpiresAt:        checkedExpiresAt,
		Passengers:       passengers,
		Ancillaries:      ancillaries,
	}
}

func (s *BookingService) GetBookingByID(id string) (*dto.BookingResponse, error) {
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return nil, err
	}
	return mapBookingToDTO(booking), nil
}

func (s *BookingService) GetBookingByPNR(pnr string) (*dto.BookingResponse, error) {
	booking, err := s.repo.GetBookingByPNR(pnr)
	if err != nil {
		return nil, err
	}
	return mapBookingToDTO(booking), nil
}

func (s *BookingService) GetBookingsByUserID(userID string) ([]dto.BookingResponse, error) {
	bookings, err := s.repo.GetBookingsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.BookingResponse
	for _, b := range bookings {
		responses = append(responses, *mapBookingToDTO(&b))
	}
	return responses, nil
}

func (s *BookingService) ConfirmBooking(id string) error {
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return err
	}

	if booking.Status != "PENDING_PAYMENT" {
		return errors.New("booking is not pending payment")
	}

	booking.Status = "CONFIRMED"
	now := time.Now()
	booking.ConfirmedAt = &now

	qrData := "MOCK_SIGNED_QR_DATA_" + booking.PNR
	eTicket := &models.ETicket{
		BookingID:    booking.ID,
		TicketNumber: "TKT-" + booking.PNR,
		QRCodeURL:    "https://storage.tripneo.com/qr/" + booking.PNR + ".png",
		QRData:       qrData,
	}

	if err := s.repo.ConfirmBookingAndSeats(booking, eTicket); err != nil {
		return err
	}

	ctx := context.Background()
	for _, p := range booking.Passengers {
		if p.SeatID != nil {
			_ = redis.ReleaseSeatLock(ctx, p.SeatID.String())
		}
	}

	log.Println("[KAFKA MOCK] Published event: flight.booking.confirmed for PNR:", booking.PNR)
	return nil
}

func (s *BookingService) CancelBooking(id string, req *dto.CancelBookingRequest) error {
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return err
	}

	if booking.Status != "CONFIRMED" {
		return errors.New("cannot cancel a non-confirmed booking")
	}

	refundAmount := booking.TotalAmount * 0.9 // 90% refund
	reason := "User requested"
	if req.Reason != "" {
		reason = req.Reason
	}

	cancelData := &models.Cancellation{
		BookingID:    booking.ID,
		Reason:       &reason,
		RefundAmount: refundAmount,
		RefundStatus: "PROCESSING",
	}

	if err := s.repo.CreateCancellation(cancelData, booking); err != nil {
		return err
	}

	log.Println("[KAFKA MOCK] Published event: flight.booking.cancelled for PNR:", booking.PNR)
	return nil
}

func (s *BookingService) GetTicket(bookingID string) (*dto.TicketResponse, error) {
	t, err := s.repo.GetETicketByBookingID(bookingID)
	if err != nil {
		return nil, err
	}

	return &dto.TicketResponse{
		BookingID:    t.BookingID.String(),
		TicketNumber: t.TicketNumber,
		QRCodeURL:    t.QRCodeURL,
	}, nil
}
