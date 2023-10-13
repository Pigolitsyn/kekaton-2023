package server

import (
	"math"

	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

type RequestNewComment struct {
	PointID int    `json:"point_id"`
	Text    string `json:"text"`
	Rating  int8   `json:"rating"`
}

func (s *Server) handleGetComment(fcx *fiber.Ctx) error {
	req := RequestNewComment{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	if req.PointID == 0 || req.Text == "" || req.Rating == 0 {
		return ErrData
	}

	usr, ok := fcx.Locals("user").(storage.User)
	if !ok {
		return ErrRequest
	}

	comment := storage.Comment{
		PointID: req.PointID,
		Text:    req.Text,
		Rating:  int8(math.Min(math.Min(float64(req.Rating), 5), 1)),
		UserID:  usr.ID,
	}

	if err := s.service.GetComment(fcx.UserContext(), &comment); err != nil {
		return err
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"comment": comment,
	})
}

func (s *Server) handleGetCommentsForPoint(fcx *fiber.Ctx) error {
	oid, err := fcx.ParamsInt("ID")
	if err != nil {
		return ErrRequest
	}

	comments := make([]storage.Comment, 0)

	if err = s.service.GetCommentsForPoint(fcx.UserContext(), oid, &comments); err != nil {
		return err
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "successful",
		"comments": comments,
	})
}

func (s *Server) handleAddComment(fcx *fiber.Ctx) error {
	return ErrSuccess
}

func (s *Server) handleUpdateComment(fcx *fiber.Ctx) error {
	return ErrSuccess
}
