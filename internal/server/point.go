package server

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

type RequestNewPoint struct {
	Coordinates storage.Coordinates
	Description string
	OpenTime    time.Duration
	CloseTime   time.Duration
	Tags        []int
}

func (s *Server) handleAddPoint(fcx *fiber.Ctx) error {
	req := RequestNewPoint{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	if req.Coordinates == [2]float64{0.0, 0.0} || req.Description == "" {
		return ErrData
	}

	usr, ok := fcx.Locals("user").(storage.User)
	if !ok {
		return ErrRequest
	}

	point := storage.Point{
		Coordinates: req.Coordinates,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		Creator:     usr,
	}

	if err := s.service.RegisterPoint(fcx.UserContext(), &point, &req.Tags); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"point":   point.ID,
	})
}

type RequestPoint struct {
	ID int
}

func (s *Server) handleGetPoint(fcx *fiber.Ctx) error {
	pid, err := strconv.ParseInt(fcx.Query("id"), 10, 0)
	if err != nil || pid == 0 {
		return ErrRequest
	}

	point := storage.Point{
		ID: int(pid),
	}

	if err = s.service.GetPointByID(fcx.UserContext(), &point); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"point":   point,
	})
}

func (s *Server) handleGetPoints(fcx *fiber.Ctx) error {
	points := make([]storage.Point, 0)

	if err := s.service.GetPoints(fcx.UserContext(), &points); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"points":  points,
	})
}

func (s *Server) handleUpdatePoint(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}
