package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IRCTCTrainResponse perfectly matches the RapidAPI /getTrainSchedule response
type IRCTCTrainResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TrainNo            string                 `json:"trainNo"`
		TrainName          string                 `json:"trainName"`
		SourceStation      string                 `json:"sourceStation"`
		DestinationStation string                 `json:"destinationStation"`
		RunDays            map[string]interface{} `json:"runDays"`
		Route              []struct {
			StationCode     string `json:"stationCode"`
			StationName     string `json:"stationName"`
			StnSerialNumber int    `json:"stnSerialNumber"`
			ArrivalTime     string `json:"arrivalTime"`
			DepartureTime   string `json:"departureTime"`
			HaltTime        string `json:"haltTime"`
			Distance        string `json:"distance"` // API returns this as a string!
			Day             int    `json:"day"`
		} `json:"route"`
	} `json:"data"`
}

// FetchTrainDetails calls the RapidAPI endpoint for a specific train number.
func FetchTrainDetails(trainNumber string, apiKey string, apiHost string) (*IRCTCTrainResponse, error) {
	// Updated to the exact URL from your screenshot
	url := fmt.Sprintf("https://%s/api/v1/getTrainSchedule?trainNo=%s", apiHost, trainNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", apiHost)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch train %s: status %d", trainNumber, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var trainData IRCTCTrainResponse
	if err := json.Unmarshal(body, &trainData); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	// Check if the API itself returned an error message inside the successful HTTP response
	if !trainData.Status {
		return nil, fmt.Errorf("API returned false status: %s", trainData.Message)
	}

	return &trainData, nil
}
