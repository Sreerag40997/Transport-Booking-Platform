package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/junaid9001/tripneo/auth-service/service"
)

var Validate = validator.New()

func Register() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=4"`
		}

		if err := c.Bind().Body(&req); err != nil {
			log.Print(err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON body",
			})
		}

		if err := Validate.Struct(req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON body",
			})
		}

		err := service.CreateUser(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, service.EmailAlreadyTaken) {
				return c.Status(http.StatusConflict).JSON(fiber.Map{
					"error": "email already registered",
				})
			}

			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}

		return nil

	}

}
