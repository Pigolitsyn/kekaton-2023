package server

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

type RequestNewPoint struct {
	Coordinates storage.Coordinates `json:"coordinates"`
	Description string              `json:"description"`
	OpenTime    time.Duration       `json:"open_time"`
	CloseTime   time.Duration       `json:"close_time"`
	Tags        []int               `json:"tags"`
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
	ID int `json:"id"`
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
		"type":     "FeatureCollection",
		"features": pointsToGeo(points),
	})
}

type RequestUpdatePoint struct {
	ID          int                 `json:"id"`
	Coordinates storage.Coordinates `json:"coordinates"`
	Description string              `json:"description"`
	OpenTime    time.Duration       `json:"open_time"`
	CloseTime   time.Duration       `json:"close_time"`
	Tags        []int               `json:"tags"`
}

func (s *Server) handleUpdatePoint(fcx *fiber.Ctx) error {
	req := RequestUpdatePoint{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	if req.Coordinates == [2]float64{0.0, 0.0} || req.Description == "" {
		return ErrData
	}

	point := storage.Point{
		ID:          req.ID,
		Coordinates: req.Coordinates,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
	}

	if err := s.service.UpdatePoint(fcx.UserContext(), &point, &req.Tags); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}
