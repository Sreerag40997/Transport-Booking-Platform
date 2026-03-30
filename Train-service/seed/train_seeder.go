package seed

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/nabeel-mp/tripneo/train-service/config"
	"github.com/nabeel-mp/tripneo/train-service/models"
	"github.com/nabeel-mp/tripneo/train-service/utils"
	"gorm.io/gorm"
)

func SeedStations(tx *gorm.DB) error {
	bytes, err := os.ReadFile("data/stations.json")
	if err != nil {
		return err
	}

	var stations []models.Station
	if err := json.Unmarshal(bytes, &stations); err != nil {
		return err
	}

	for _, s := range stations {
		if err := tx.Where("code = ?", s.Code).FirstOrCreate(&s).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedTrains(tx *gorm.DB) error {
	bytes, err := os.ReadFile("data/trains.json")
	if err != nil {
		return err
	}

	var rawTrains []struct {
		TrainNumber string  `json:"train_number"`
		TrainName   string  `json:"train_name"`
		DaysOfWeek  []int32 `json:"days_of_week"`
		IsActive    bool    `json:"is_active"`
		Stops       []struct {
			StationCode string `json:"station_code"`
			Sequence    int    `json:"sequence"`
			Arrival     string `json:"arrival"`
			Departure   string `json:"departure"`
			DayOffset   int    `json:"day_offset"`
			Distance    int    `json:"distance"`
		} `json:"stops"`
	}

	if err := json.Unmarshal(bytes, &rawTrains); err != nil {
		return err
	}

	for _, r := range rawTrains {

		var origin, destination, depTime, arrTime string
		var totalDuration int
		if len(r.Stops) > 0 {
			firstStop := r.Stops[0]
			lastStop := r.Stops[len(r.Stops)-1]

			origin = firstStop.StationCode
			depTime = firstStop.Departure

			destination = lastStop.StationCode
			arrTime = lastStop.Arrival
			totalDuration = lastStop.Distance
		}
		train := models.Train{
			TrainNumber:        r.TrainNumber,
			TrainName:          r.TrainName,
			OriginStation:      origin,
			DestinationStation: destination,
			DepartureTime:      depTime,
			ArrivalTime:        arrTime,
			DurationMinutes:    totalDuration,
			DaysOfWeek:         pq.Int32Array(r.DaysOfWeek),
			IsActive:           r.IsActive,
		}
		if err := tx.Where("train_number = ?", train.TrainNumber).
			Assign(models.Train{OriginStation: origin, DestinationStation: destination, DepartureTime: depTime, ArrivalTime: arrTime, DurationMinutes: totalDuration}).
			FirstOrCreate(&train).Error; err != nil {
			return err
		}

		for _, stop := range r.Stops {
			var station models.Station
			if err := tx.Where("code = ?", stop.StationCode).First(&station).Error; err != nil {
				log.Printf("Warning: Station %s not found, skipping stop", stop.StationCode)
				continue
			}
			trainStop := models.TrainStop{
				TrainID:       train.ID,
				StationID:     station.ID,
				StopSequence:  stop.Sequence,
				ArrivalTime:   stop.Arrival,
				DepartureTime: stop.Departure,
				DayOffset:     stop.DayOffset,
				DistanceKm:    stop.Distance,
			}
			tx.Where("train_id = ? AND stop_sequence = ?", train.ID, stop.Sequence).FirstOrCreate(&trainStop)
		}
	}
	return nil
}

func SeedTrainsFromAPI(tx *gorm.DB, cfg *config.Config) error {
	trainNumbersToSeed := []string{"12004", "12951", "12431"}

	for _, trainNo := range trainNumbersToSeed {
		log.Printf("Fetching details for train %s from API...", trainNo)

		apiData, err := utils.FetchTrainDetails(trainNo, cfg.RAPID_API_KEY, cfg.RAPID_API_HOST)
		if err != nil {
			log.Printf("Error fetching train %s: %v", trainNo, err)

			// 2. Add a pause even if there's an error, before trying the next train
			time.Sleep(2 * time.Second)
			continue
		}

		if len(apiData.Data.Route) == 0 {
			log.Printf("Warning: Train %s returned empty route, skipping", trainNo)
			time.Sleep(2 * time.Second)
			continue
		}

		firstStop := apiData.Data.Route[0]
		lastStop := apiData.Data.Route[len(apiData.Data.Route)-1]

		origin := firstStop.StationCode
		depTime := firstStop.DepartureTime
		destination := lastStop.StationCode
		arrTime := lastStop.ArrivalTime

		totalDuration, _ := strconv.Atoi(lastStop.Distance)

		train := models.Train{
			TrainNumber:        apiData.Data.TrainNo,
			TrainName:          apiData.Data.TrainName,
			OriginStation:      origin,
			DestinationStation: destination,
			DepartureTime:      depTime,
			ArrivalTime:        arrTime,
			DurationMinutes:    totalDuration,
			DaysOfWeek:         pq.Int32Array{1, 2, 3, 4, 5, 6, 7},
			IsActive:           true,
		}

		if err := tx.Where("train_number = ?", train.TrainNumber).
			Assign(train).
			FirstOrCreate(&train).Error; err != nil {
			return err
		}

		for _, stop := range apiData.Data.Route {
			station := models.Station{
				Code: stop.StationCode,
				Name: stop.StationName,
			}
			tx.Where("code = ?", station.Code).FirstOrCreate(&station)

			dist, _ := strconv.Atoi(stop.Distance)

			trainStop := models.TrainStop{
				TrainID:       train.ID,
				StationID:     station.ID,
				StopSequence:  stop.StnSerialNumber,
				ArrivalTime:   stop.ArrivalTime,
				DepartureTime: stop.DepartureTime,
				DayOffset:     stop.Day,
				DistanceKm:    dist,
			}

			tx.Where("train_id = ? AND stop_sequence = ?", train.ID, stop.StnSerialNumber).
				Assign(trainStop).
				FirstOrCreate(&trainStop)
		}

		log.Printf("Successfully seeded train %s", trainNo)

		// 3. Add a 2-second pause before fetching the next train
		time.Sleep(2 * time.Second)
	}
	return nil
}
