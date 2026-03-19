package main

import (
	//"github.com/Salman-kp/tripneo/bus-service/cache"
	"github.com/Salman-kp/tripneo/bus-service/config"
	"github.com/Salman-kp/tripneo/bus-service/db"

	//	"github.com/Salman-kp/tripneo/bus-service/kafka"
	//"github.com/Salman-kp/tripneo/bus-service/routes"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()
	db.ConnectPostgres(cfg)
	//	rdb := cache.NewRedisClient(cfg)
	//	producer := kafka.NewProducer(cfg)
	//	defer producer.Close()

	app := fiber.New()
	app.Get("/api/bus/health", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})
	//	routes.Register(app, cfg, rdb, producer)

	if err := app.Listen(":" + cfg.APP_PORT); err != nil {
		log.Fatal(err)
	}
}
