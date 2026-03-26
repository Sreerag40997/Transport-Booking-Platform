package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/junaid9001/tripneo/api-gateway/config"
	"github.com/junaid9001/tripneo/api-gateway/proxy"
	"github.com/redis/go-redis/v9"
)

func RegisterTrainRoutes(app *fiber.App, cfg *config.Config, rdb *redis.Client) {
	api := app.Group("/api/train")

	api.Get("/health", proxy.To(cfg.TRAIN_SERVICE_URL))

	// ------ Public ------
	api.Get("/search", proxy.To(cfg.TRAIN_SERVICE_URL))
	api.Get("/:id", proxy.To(cfg.TRAIN_SERVICE_URL))
	api.Get("/:id/live-status", proxy.To(cfg.TRAIN_SERVICE_URL))
	api.Get("/:id/seats", proxy.To(cfg.TRAIN_SERVICE_URL))

	//----- Protected -----
	api.Post("/book", proxy.To(cfg.TRAIN_SERVICE_URL))
	api.Get("/bookings/:id", proxy.To(cfg.TRAIN_SERVICE_URL))
	api.Get("/bookings/user/history", proxy.To(cfg.TRAIN_SERVICE_URL))
	api.Post("/bookings/:id/cancel", proxy.To(cfg.TRAIN_SERVICE_URL))

}
