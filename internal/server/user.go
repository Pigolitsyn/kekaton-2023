package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"kekaton/back/internal/storage"
)

type RequestSignUp struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) handleUserSignUp(fcx *fiber.Ctx) error {
	req := RequestSignUp{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return ErrData
	}

	user := storage.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	if err := s.service.RegisterUser(fcx.UserContext(), &user); err != nil {
		return ErrInternal
	}

	token, err := s.MakeJWT(&user)
	if err != nil {
		return ErrInternal
	}

	expires := time.Now().Add(time.Hour * 72).UTC()

	fcx.Cookie(&fiber.Cookie{
		Name:     s.config.TokenName,
		Value:    token,
		Expires:  expires,
		HTTPOnly: false,
		Secure:   true,
	})

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful sign up",
		"jwt":     token,
		"exp":     expires,
	})
}

type RequestSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleUserSignIn(fcx *fiber.Ctx) error {
	req := RequestSignIn{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	if req.Email == "" || req.Password == "" {
		return ErrData
	}

	user := storage.User{
		Email: req.Email,
	}

	if err := s.service.GetUserByEmail(fcx.UserContext(), &user); err != nil {
		return ErrInternal
	}

	if s.service.VerifyPassword(user.Password, user.Salt, req.Password) {
		return ErrInternal
	}

	token, err := s.MakeJWT(&user)
	if err != nil {
		return ErrInternal
	}

	expires := time.Now().Add(time.Hour * 72).UTC()

	fcx.Cookie(&fiber.Cookie{
		Name:     s.config.TokenName,
		Value:    token,
		Expires:  expires,
		HTTPOnly: false,
		Secure:   true,
	})

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful sign in",
		"jwt":     token,
		"exp":     expires,
	})
}

func (s *Server) handleUserSignOut(fcx *fiber.Ctx) error {
	fcx.Cookie(&fiber.Cookie{
		Name:    s.config.TokenName,
		Value:   "",
		Expires: time.Now().Add(-time.Hour * 24),
	})

	return fiber.NewError(fiber.StatusOK, "successful sign out")
}

func (s *Server) handleGetUser(fcx *fiber.Ctx) error {
	id := fcx.Query("id")

	if id == "" {
		return ErrRequest
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrData
	}

	user := storage.User{
		ID: uid,
	}

	if err = s.service.GetUserByID(fcx.UserContext(), &user); err != nil {
		return ErrInternal
	}

	return fcx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successful",
		"user":    user,
	})
}

func (s *Server) handleUpdateUser(fcx *fiber.Ctx) error {
	req := RequestSignUp{}

	if err := fcx.BodyParser(&req); err != nil {
		return ErrRequest
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return ErrData
	}

	usr, ok := fcx.Locals("user").(storage.User)
	if !ok {
		return ErrRequest
	}

	user := storage.User{
		ID:       usr.ID,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	if err := s.service.UpdateUser(fcx.UserContext(), &user); err != nil {
		return ErrInternal
	}

	return ErrSuccess
}
