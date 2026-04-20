package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// acquireSeatLock attempts to lock a seat for a specific user and duration.
// returns true if the lock was acquired successfully.
func AcquireSeatLock(ctx context.Context, seatID string, userID string, ttl time.Duration) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}

	key := fmt.Sprintf("seat_lock:%s", seatID)
	shadowKey := fmt.Sprintf("shadow:seat_lock:%s:%s", userID, seatID)

	isSet, err := Client.SetNX(ctx, key, userID, ttl).Result()
	if err != nil {
		return false, err
	}

	if isSet {
		// set the shadow key for the notification to trigger correctly

		Client.Set(ctx, shadowKey, "1", ttl)
	}

	return isSet, nil
}

// releaseseatLock unlocks a seat.
// returns true if the lock was successfully deleted.
func ReleaseSeatLock(ctx context.Context, seatID string) error {
	if Client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	key := fmt.Sprintf("seat_lock:%s", seatID)

	err := Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// getseatLockOwner returns the current owner (user ID) of the lock if any.
func GetSeatLockOwner(ctx context.Context, seatID string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}

	key := fmt.Sprintf("seat_lock:%s", seatID)

	val, err := Client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", nil // lock does not exist
	} else if err != nil {
		return "", err
	}

	return val, nil
}
