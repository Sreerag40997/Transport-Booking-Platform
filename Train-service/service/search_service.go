package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nabeel-mp/tripneo/train-service/repository"
	goredis "github.com/redis/go-redis/v9"
)

const searchCacheTTL = 2 * time.Minute

func SearchTrains(ctx context.Context, rdb *goredis.Client, fromCode, toCode, dateStr, class string) ([]repository.SearchResult, error) {
	// Parse the incoming date string (e.g., "2026-03-28") into a time.Time object
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	// Default to "SL" if no class is provided so the repository query doesn't fail
	if class == "" {
		class = "SL"
	}

	// Call the repository function
	results, err := repository.SearchTrains(fromCode, toCode, class, parsedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to search trains: %w", err)
	}

	return results, nil
}

// GetScheduleDetail returns a single schedule with its train details.
func GetScheduleDetail(scheduleID string) (interface{}, error) {
	schedule, err := repository.GetScheduleByID(scheduleID)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func GetSeatMap(
	ctx context.Context,
	rdb *goredis.Client,
	scheduleID, class string,
) (interface{}, error) {
	seats, err := repository.GetSeatsByScheduleAndClass(scheduleID, class)
	if err != nil {
		return nil, err
	}

	type SeatWithLock struct {
		ID         string  `json:"id"`
		SeatNumber string  `json:"seat_number"`
		Coach      string  `json:"coach"`
		Class      string  `json:"class"`
		BerthType  string  `json:"berth_type"`
		Status     string  `json:"status"`
		Price      float64 `json:"price"`
		IsLocked   bool    `json:"is_locked"`
	}

	result := make([]SeatWithLock, len(seats))
	for i, s := range seats {
		isLocked := false
		if s.Status == "AVAILABLE" {
			// Check Redis lock — display only, not for booking gate
			locked, _ := checkLockStatus(ctx, rdb, scheduleID, s.ID.String())
			isLocked = locked
		}
		result[i] = SeatWithLock{
			ID:         s.ID.String(),
			SeatNumber: s.SeatNumber,
			Coach:      s.Coach,
			Class:      s.Class,
			BerthType:  s.BerthType,
			Status:     s.Status,
			Price:      s.Price,
			IsLocked:   isLocked,
		}
	}
	return result, nil
}

// checkLockStatus is an internal helper that reads the Redis lock key.
func checkLockStatus(ctx context.Context, rdb *goredis.Client, scheduleID, seatID string) (bool, error) {
	key := fmt.Sprintf("seat:lock:train:%s:%s", scheduleID, seatID)
	_, err := rdb.Get(ctx, key).Result()
	if err == goredis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
