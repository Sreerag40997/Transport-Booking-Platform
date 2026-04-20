package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/junaid9001/tripneo/flight-service/config"
	"github.com/junaid9001/tripneo/flight-service/repository"
)

type TrackingService struct {
	repo   *repository.BookingRepository
	config *config.Config
}

func NewTrackingService(repo *repository.BookingRepository, cfg *config.Config) *TrackingService {
	return &TrackingService{
		repo:   repo,
		config: cfg,
	}
}

type AviationStackFlight struct {
	FlightStatus string `json:"flight_status"`
	Live         *struct {
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
		Altitude        float64 `json:"altitude"`
		Direction       float64 `json:"direction"`
		SpeedHorizontal float64 `json:"speed_horizontal"`
		IsGround        bool    `json:"is_ground"`
	} `json:"live"`
	Departure struct {
		Iata string `json:"iata"`
	} `json:"departure"`
	Arrival struct {
		Iata string `json:"iata"`
	} `json:"arrival"`
}

type AviationStackResponse struct {
	Data []AviationStackFlight `json:"data"`
}

func (s *TrackingService) GetLiveRadar(pnr string) (map[string]interface{}, error) {
	booking, err := s.repo.GetBookingByPNR(pnr)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	apiKey := s.config.AVIATIONSTACK_API_KEY
	if apiKey == "" {
		return nil, fmt.Errorf("tracking API key not configured")
	}

	flightNumber := booking.FlightInstance.Flight.FlightNumber

	url := fmt.Sprintf("http://api.aviationstack.com/v1/flights?access_key=%s&flight_status=active&limit=100", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to reach tracking API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tracking API response: %v", err)
	}

	var avResp AviationStackResponse
	if err := json.Unmarshal(body, &avResp); err != nil {
		log.Printf("[Tracking API Raw Response] %s", string(body))
		return nil, fmt.Errorf("failed to parse tracking data: %v", err)
	}

	var activeFlight *AviationStackFlight
	for _, f := range avResp.Data {
		if f.Live != nil && f.Live.Latitude != 0 {
			activeFlight = &f
			break
		}
	}

	if activeFlight == nil {
		return nil, fmt.Errorf("no active radar data found currently")
	}

	return map[string]interface{}{
		"pnr":           pnr,
		"flight_number": flightNumber,
		"origin":        booking.FlightInstance.Flight.OriginAirport.IataCode,
		"destination":   booking.FlightInstance.Flight.DestinationAirport.IataCode,
		"status":        "EN_ROUTE",
		"live": map[string]interface{}{
			"latitude":  activeFlight.Live.Latitude,
			"longitude": activeFlight.Live.Longitude,
			"altitude":  activeFlight.Live.Altitude,
			"speed_mph": activeFlight.Live.SpeedHorizontal,
			"heading":   activeFlight.Live.Direction,
			"is_ground": activeFlight.Live.IsGround,
		},
	}, nil
}
