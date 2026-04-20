package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/junaid9001/tripneo/flight-service/config"
	"github.com/junaid9001/tripneo/flight-service/db"
	"github.com/junaid9001/tripneo/flight-service/handlers"
	"github.com/junaid9001/tripneo/flight-service/jobs"
	"github.com/junaid9001/tripneo/flight-service/kafka"
	"github.com/junaid9001/tripneo/flight-service/redis"
	"github.com/junaid9001/tripneo/flight-service/repository"
	"github.com/junaid9001/tripneo/flight-service/routes"
	"github.com/junaid9001/tripneo/flight-service/rpc"
	"github.com/junaid9001/tripneo/flight-service/seed"
	"github.com/junaid9001/tripneo/flight-service/services"
	"github.com/junaid9001/tripneo/flight-service/ws"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.LoadConfig()

	db.ConnectPostgres(cfg)
	redis.ConnectRedis(cfg)

	go services.StartRedisExpirySubscriber()

	if cfg.RUN_SEED_ON_BOOT == "true" {
		seed.SeedAll(db.DB)
	}

	// initialize gRPC Client
	payClient, err := rpc.NewPaymentClient(cfg.PAYMENT_SERVICE_GRPC_URL)
	if err != nil {
		log.Printf("Warning: Payment gRPC client failed to connect: %v", err)
	} else {
		defer payClient.Close()
	}

	// initialize repos and services for background Workers
	bookingRepo := repository.NewBookingRepository(db.DB)
	bookingService := services.NewBookingService(bookingRepo, payClient, ws.DefaultManager)

	// initialize kafka consumer
	kafkaConsumer := kafka.NewConsumer(cfg.KAFKA_BROKERS, "flight-payment-topic", "flight-service-group")
	if kafkaConsumer != nil {
		defer kafkaConsumer.Close()
		go kafkaConsumer.ConsumePaymentEvents(context.Background(), func(evt kafka.PaymentCompletedEvent) {
			bookingService.ProcessPaymentEvent(evt)
		})
	}

	app := fiber.New()

	app.Get("/api/flights/health", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Use("/api/flights/ws", handlers.WebsocketUpgradeMiddleware)
	app.Get("/api/flights/ws", handlers.HandleWebSocket)

	routes.SetupFlightRoutes(app, db.DB, cfg)
	routes.SetupBookingRoutes(app, db.DB, payClient, ws.DefaultManager)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		jobs.GenerateUpcomingInventory(db.DB)
	})
	c.AddFunc("*/5 * * * *", func() {
		jobs.CleanupExpiredBookings(db.DB)
	})
	c.Start()

	go jobs.GenerateUpcomingInventory(db.DB)

	app.Listen(":" + cfg.APP_PORT)
}
