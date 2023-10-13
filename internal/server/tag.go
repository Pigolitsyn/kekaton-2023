package server

import (
	"github.com/gofiber/fiber/v2"

	"kekaton/back/internal/storage"
)

func (s *Server) handleGetTag(fcx *fiber.Ctx) error {
	ttype := fcx.QueryInt("type")

	tag := storage.Tag{
		Type: ttype,
	}

	if err := s.service.GetTagByID(fcx.UserContext(), &tag); err != nil {
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
