package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/nabeel-mp/tripneo/train-service/config"
	"github.com/nabeel-mp/tripneo/train-service/models"
	"gorm.io/gorm"
)

// Adjust these structs based on your specific RapidAPI JSON response
type IRCTCTrainScheduleResponse struct {
	Data struct {
		TrainNumber string             `json:"train_number"`
		TrainName   string             `json:"train_name"`
		Route       []IRCTCStationStop `json:"route"`
		Days        []int32            `json:"run_days"` // e.g., [1,0,1,0,0,0,0]
	} `json:"data"`
}

type IRCTCStationStop struct {
	StationCode   string `json:"station_code"`
	StationName   string `json:"station_name"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
	Distance      int    `json:"distance"`
	Day           int    `json:"day"` // 1 for Day 1, 2 for Day 2, etc.
}

func FetchTrainDataFromAPI(trainNumber string, cfg *config.Config) (*IRCTCTrainScheduleResponse, error) {
	// Example URL - replace with your specific RapidAPI endpoint
	url := fmt.Sprintf("https://%s/v1/trains/schedule?train_number=%s", cfg.RAPID_API_HOST, trainNumber)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-key", cfg.RAPID_API_KEY)
	req.Header.Add("x-rapidapi-host", cfg.RAPID_API_HOST)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch data, status code: %d", res.StatusCode)
	}

	var response IRCTCTrainScheduleResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func SeedTrainToDB(tx *gorm.DB, trainData *IRCTCTrainScheduleResponse) error {
	data := trainData.Data
	if len(data.Route) == 0 {
		return fmt.Errorf("no route data found for train %s", data.TrainNumber)
	}

	firstStop := data.Route[0]
	lastStop := data.Route[len(data.Route)-1]

	// 1. Create Train Record
	train := models.Train{
		TrainNumber:        data.TrainNumber,
		TrainName:          data.TrainName,
		OriginStation:      firstStop.StationCode,
		DestinationStation: lastStop.StationCode,
		DepartureTime:      firstStop.DepartureTime,
		ArrivalTime:        lastStop.ArrivalTime,
		DurationMinutes:    lastStop.Distance, // You might need to calculate actual minutes here
		DaysOfWeek:         pq.Int32Array(data.Days),
		IsActive:           true,
	}

	if err := tx.Where("train_number = ?", train.TrainNumber).FirstOrCreate(&train).Error; err != nil {
		return err
	}

	// 2. Loop through stops and create Stations and TrainStops
	for i, stop := range data.Route {
		// Ensure the station exists in the DB
		var station models.Station
		if err := tx.Where("code = ?", stop.StationCode).Attrs(models.Station{
			Name: stop.StationName,
		}).FirstOrCreate(&station).Error; err != nil {
			log.Printf("Failed to save station %s: %v", stop.StationCode, err)
			continue
		}

		// Calculate Day Offset (IRCTC usually gives Day 1, Day 2. We need 0, 1, etc.)
		dayOffset := stop.Day - 1
		if dayOffset < 0 {
			dayOffset = 0
		}

		// Save the Train Stop
		trainStop := models.TrainStop{
			TrainID:       train.ID,
			StationID:     station.ID,
			StopSequence:  i + 1,
			ArrivalTime:   stop.ArrivalTime,
			DepartureTime: stop.DepartureTime,
			DayOffset:     dayOffset,
			DistanceKm:    stop.Distance,
		}

		tx.Where("train_id = ? AND stop_sequence = ?", train.ID, trainStop.StopSequence).
			Assign(trainStop).
			FirstOrCreate(&trainStop)
	}

	return nil
}
