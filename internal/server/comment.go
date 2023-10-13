package server

import "github.com/gofiber/fiber/v2"

func (s *Server) handleGetComment(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}

func (s *Server) handleGetComments(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}

func (s *Server) handleAddComment(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}

func (s *Server) handleUpdateComment(fcx *fiber.Ctx) error {
	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successful"})
}
