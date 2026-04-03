package dto

import "time"

type PassengerDto struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"` // YYYY-MM-DD
	Gender         string `json:"gender"`
	PassengerType  string `json:"passenger_type"` // adult | child | infant
	IDType         string `json:"id_type"`        // PASSPORT | AADHAAR | PAN
	IDNumber       string `json:"id_number"`
	SeatID         string `json:"seat_id,omitempty"`
	MealPreference string `json:"meal_preference,omitempty"`
}

type AncillaryBookDto struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	PassengerID string  `json:"passenger_id,omitempty"`
}

type CreateBookingRequest struct {
	FlightInstanceID string             `json:"flight_instance_id"`
	FareTypeID       string             `json:"fare_type_id"`
	TripType         string             `json:"trip_type"`
	SeatClass        string             `json:"seat_class"`
	Passengers       []PassengerDto     `json:"passengers"`
	Ancillaries      []AncillaryBookDto `json:"ancillaries"`
	GSTIN            string             `json:"gstin,omitempty"`
}

type BookingResponse struct {
	ID               string    `json:"id"`
	PNR              string    `json:"pnr"`
	FlightInstanceID string    `json:"flight_instance_id"`
	Source           string    `json:"source"`
	Status           string    `json:"status"`
	TotalAmount      float64   `json:"total_amount"`
	Currency         string    `json:"currency"`
	BookedAt         time.Time `json:"booked_at"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason,omitempty"`
}

type TicketResponse struct {
	BookingID    string `json:"booking_id"`
	TicketNumber string `json:"ticket_number"`
	QRCodeURL    string `json:"qr_code_url"`
}
