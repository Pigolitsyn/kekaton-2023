package server

import (
	"math"
	"time"

	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

type RequestNewPoint struct {
	Coordinates storage.Coordinates `json:"coordinates"`
	Address     string              `json:"address"`
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
		Address:	 req.Address,
		Creator:     usr,
	}

	if err := s.service.RegisterPoint(fcx.UserContext(), &point, &req.Tags); err != nil {
		return ErrInternal
	}

	return ErrSuccess
}

func (s *Server) handleGetPoint(fcx *fiber.Ctx) error {
	pid := fcx.QueryInt("id")

	point := storage.Point{
		ID: int(pid),
	}

	if err := s.service.GetPointByID(fcx.UserContext(), &point); err != nil {
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
	Address     string              `json:"address"`
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
		Address:     req.Address,
		Description: req.Description,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
	}

	if err := s.service.UpdatePoint(fcx.UserContext(), &point, &req.Tags); err != nil {
		return ErrInternal
	}

	return ErrSuccess
}

func (s *Server) handleGetClosestPoint(fcx *fiber.Ctx) error {
	lon := fcx.QueryFloat("lon")
	lat := fcx.QueryFloat("lat")
	rad := fcx.QueryFloat("rad")

	points := make([]storage.Point, 0)

	if err := s.service.GetPoints(fcx.UserContext(), &points); err != nil {
		return ErrInternal
	}

	center := storage.Point{
		Coordinates: [2]float64{lon, lat},
	}

	point := storage.Point{}

	for i := range points {
		if p, ok := pointInRadius(center, points[i], point, rad); ok {
			point = p
		}
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "successful",
		"coordinates": point.Coordinates,
	})
}

func pointInRadius(c, p, o storage.Point, r float64) (storage.Point, bool) {
	x1 := math.Abs(c.Coordinates[0] - p.Coordinates[0])
	y1 := math.Abs(c.Coordinates[1] - p.Coordinates[1])

	x2 := math.Abs(c.Coordinates[0] - o.Coordinates[0])
	y2 := math.Abs(c.Coordinates[1] - o.Coordinates[1])

	cl := storage.Point{}

	if (x1 + y1) <= (x2 + y2) {
		cl = p
	} else {
		cl = o
	}

	return cl, x1+y1 <= r
}
