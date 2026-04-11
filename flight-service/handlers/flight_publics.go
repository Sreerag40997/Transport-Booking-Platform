package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (h *FlightHandler) SearchAirport(c fiber.Ctx) error {

	search := c.Query("search", "del")
	search = strings.TrimSpace(search)
	if search == "" || len(search) < 3 || len(search) > 20 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid search"})
	}

	airports, err := h.flightService.SearchAirport(strings.ToLower(search))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "ok",
		"data":    airports,
	})
}

func (h *FlightHandler) GetAirlines(c fiber.Ctx) error {
	airlines, err := h.flightService.GetAirlines()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "ok",
		"data":    airlines,
	})
}
