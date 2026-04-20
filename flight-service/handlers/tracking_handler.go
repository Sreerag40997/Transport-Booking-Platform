package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/junaid9001/tripneo/flight-service/services"
)

type TrackingHandler struct {
	trackingService *services.TrackingService
}

func NewTrackingHandler(service *services.TrackingService) *TrackingHandler {
	return &TrackingHandler{trackingService: service}
}

func (h *TrackingHandler) GetFlightStatus(c fiber.Ctx) error {
	pnr := c.Params("pnr")
	if pnr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "PNR is required"})
	}

	data, err := h.trackingService.GetLiveRadar(pnr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
