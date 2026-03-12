package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/junaid9001/tripneo/auth-service/config"
	"github.com/junaid9001/tripneo/auth-service/handlers"
)

func Register(app *fiber.App, cfg *config.Config) {
	auth := app.Group("/auth")

	auth.Post("/register", handlers.Register())
}
