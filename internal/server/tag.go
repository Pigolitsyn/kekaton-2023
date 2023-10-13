package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

type RequestTag struct {
	Type int
}

func (s *Server) handleGetTag(fcx *fiber.Ctx) error {
	ttype, err := strconv.ParseInt(fcx.Query("type"), 10, 0)
	if err != nil || ttype == 0 {
		return ErrRequest
	}

	tag := storage.Tag{
		Type: int(ttype),
	}

	if err = s.service.GetTagByID(fcx.UserContext(), &tag); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"tag":     tag,
	})
}

func (s *Server) handleGetTags(fcx *fiber.Ctx) error {
	tags := make([]storage.Tag, 0)

	if err := s.service.GetTags(fcx.UserContext(), &tags); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"tag":     tags,
	})
}
