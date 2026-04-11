package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// AcquireSeatLock attempts to lock a seat for a specific user and duration.
// Returns true if the lock was acquired successfully.
func AcquireSeatLock(ctx context.Context, seatID string, userID string, ttl time.Duration) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}
	
	// We use two keys: 
	// 1. The actual lock that prevents others: "seat_lock:<seatID>" -> value "userID" with TTL
	// 2. The shadow key used for expiry notifications: "shadow:seat_lock:<userID>:<seatID>" with same TTL
	
	key := fmt.Sprintf("seat_lock:%s", seatID)
	shadowKey := fmt.Sprintf("shadow:seat_lock:%s:%s", userID, seatID)
	
	// SET key value NX EX ttl
	isSet, err := Client.SetNX(ctx, key, userID, ttl).Result()
	if err != nil {
		return false, err
	}
	
	if isSet {
		// Set the shadow key for the notification to trigger correctly
		// We can safely ignore errors here as it doesn't break the actual lock
		Client.Set(ctx, shadowKey, "1", ttl)
	}
	
	return isSet, nil
}

// ReleaseSeatLock unlocks a seat.
// Returns true if the lock was successfully deleted.
func ReleaseSeatLock(ctx context.Context, seatID string) error {
	if Client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	
	key := fmt.Sprintf("seat_lock:%s", seatID)
	
	err := Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	
	// Also attempt to delete the shadow key, but don't strictly care if it fails
	// because we're just releasing early, so no timeout notification is needed.
	// NOTE: To delete it we'd need to know the userID. We can fetch it first if we want, or just ignore.
	// We will just let the shadow key expire on its own. The WS notification will be sent,
	// but we can have the client ignore it if the booking was already finalized.
	
	return nil
}

// GetSeatLockOwner returns the current owner (user ID) of the lock if any.
func GetSeatLockOwner(ctx context.Context, seatID string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}

	key := fmt.Sprintf("seat_lock:%s", seatID)

	val, err := Client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", nil // Lock does not exist
	} else if err != nil {
		return "", err
	}

	return val, nil
}
