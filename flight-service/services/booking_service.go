package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/tripneo/flight-service/dto"
	"github.com/junaid9001/tripneo/flight-service/models"
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
	flightInstance, err := s.repo.GetFlightInstanceByID(req.FlightInstanceID)
	if err != nil {
		return nil, errors.New("invalid flight instance")
	}

	_, err = s.repo.GetFareTypeByID(req.FareTypeID)
	if err != nil {
		return nil, errors.New("invalid fare type")
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

	taxes := totalBase * 0.18 // fake 18% tax
	serviceFee := 500.0
	totalAmount := totalBase + taxes + serviceFee + ancTotal

	var passengers []models.Passenger
	for _, p := range req.Passengers {
		dob, _ := time.Parse("2006-01-02", p.DateOfBirth)
		var sId *uuid.UUID = nil
		if p.SeatID != "" {
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
	expiresAt := time.Now().Add(15 * time.Minute)
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

	return &dto.BookingResponse{
		ID:               booking.ID.String(),
		PNR:              booking.PNR,
		FlightInstanceID: booking.FlightInstanceID.String(),
		Source:           booking.Source,
		Status:           booking.Status,
		TotalAmount:      booking.TotalAmount,
		Currency:         booking.Currency,
		BookedAt:         booking.CreatedAt,
		ExpiresAt:        booking.ExpiresAt,
	}, nil
}

func (s *BookingService) GetBookingByID(id string) (*models.Booking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *BookingService) GetBookingByPNR(pnr string) (*models.Booking, error) {
	return s.repo.GetBookingByPNR(pnr)
}

func (s *BookingService) GetBookingsByUserID(userID string) ([]models.Booking, error) {
	return s.repo.GetBookingsByUserID(userID)
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

	if err := s.repo.UpdateBooking(booking); err != nil {
		return err
	}

	qrData := "MOCK_SIGNED_QR_DATA_" + booking.PNR
	eTicket := &models.ETicket{
		BookingID:    booking.ID,
		TicketNumber: "TKT-" + booking.PNR,
		QRCodeURL:    "https://storage.tripneo.com/qr/" + booking.PNR + ".png",
		QRData:       qrData,
	}
	s.repo.SaveETicket(eTicket)

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

	refundAmount := booking.TotalAmount * 0.9 // flat 90% mock refund
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

	if err := s.repo.CreateCancellation(cancelData); err != nil {
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
