package server

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

type RequestNewPoint struct {
	Coordinates storage.Coordinates
	Description string
	OpenTime    time.Duration
	CloseTime   time.Duration
	Tags        []storage.Tag
}

func (s *Server) handleAddPoint(fcx *fiber.Ctx) error {
	req := RequestNewPoint{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	// TODO validate

	point := storage.Point{
		Coordinates: req.Coordinates,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		CreatedBy:   fcx.Locals("user").(storage.User).ID,
	}

	if err := s.service.RegisterPoint(fcx.UserContext(), &point); err != nil {
		return ErrInternal
	}

	// TODO register tags

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"point":   point.ID,
	})
}

func (s *Server) handleGetPoints(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}
