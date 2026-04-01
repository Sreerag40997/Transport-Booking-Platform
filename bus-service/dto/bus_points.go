package dto

import "time"

type BoardingPointResponse struct {
	ID            string    `json:"id"`
	StopID        string    `json:"stop_id"`
	StopName      string    `json:"stop_name"`
	City          string    `json:"city"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	PickupTime    time.Time `json:"pickup_time"`
	SequenceOrder int       `json:"sequence_order"`
	Landmark      string    `json:"landmark"`
}

type DroppingPointResponse struct {
	ID            string    `json:"id"`
	StopID        string    `json:"stop_id"`
	StopName      string    `json:"stop_name"`
	City          string    `json:"city"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	DropTime      time.Time `json:"drop_time"`
	SequenceOrder int       `json:"sequence_order"`
	Landmark      string    `json:"landmark"`
}

type RouteStop struct {
	StopName  string    `json:"stop_name"`
	City      string    `json:"city"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Time      time.Time `json:"time"`
	Type      string    `json:"type"`
	Sequence  int       `json:"sequence"`
}
