package services

import (
	"context"
	"fmt"
	"log"

	"github.com/junaid9001/tripneo/flight-service/redis"
	"github.com/junaid9001/tripneo/flight-service/ws"
)

// StartRedisExpirySubscriber listens for Redis key expiration events and notifies the corresponding user via WebSocket.
func StartRedisExpirySubscriber() {
	if redis.Client == nil {
		log.Fatal("Cannot start expiry subscriber: Redis client is not initialized")
	}

	ctx := context.Background()

	// 1. Enable Keyspace events. 'Ex' stands for Expired events.
	// This usually should be configured in redis.conf, but we can do it here to ensure it's on.
	// We want events for keys starting with "seat_lock:"
	err := redis.Client.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	if err != nil {
		log.Printf("Warning: Failed to set notify-keyspace-events config. Make sure it's enabled in Redis server: %v", err)
	}

	pubsub := redis.Client.Subscribe(ctx, "__keyevent@0__:expired")
	defer pubsub.Close()
	
	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Warning: Failed to subscribe to keyspace events: %v", err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	log.Println("Started Redis Keyspace Expiry Subscriber")

	for msg := range ch {
		// msg.Payload will contain the name of the expired key, e.g., "shadow:seat_lock:<userID>:<seatID>"
		key := msg.Payload
		
		var userID, seatID string
		_, err := fmt.Sscanf(key, "shadow:seat_lock:%[^:]:%s", &userID, &seatID)
		if err == nil && userID != "" && seatID != "" {
			log.Printf("EXPIRED event received for seat %s locked by %s", seatID, userID)
			
			// Notify the user via WebSocket
			message := map[string]interface{}{
				"type": "SESSION_EXPIRED",
				"message": "Your hold on the selected seat has expired.",
				"seat_id": seatID,
			}
			ws.DefaultManager.SendToUser(userID, message)
		}
	}
}
