package server

import "github.com/gofiber/fiber/v2"

func (s *Server) handleGetTag(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}

func (s *Server) handleGetTags(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}
