package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"kekaton/back/internal/storage"
)

var (
	ErrInvalidJWTMetadata = errors.New("can't parse jwt metadata")
	ErrInvalidJWT         = errors.New("invalid jwt")
)

type TokenMetadata struct {
	UserID uuid.UUID
}

func (s *Server) MakeJWT(user *storage.User) (string, error) {
	claims := jwt.MapClaims{
		"uid": user.ID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.config.Secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Server) ValidateJWT(fcx *fiber.Ctx) (*jwt.Token, error) {
	var tokenString string

	if cookie := fcx.Cookies(s.config.TokenName); cookie != "" {
		tokenString = cookie
	}

	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) { return s.config.Secret, nil })
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Server) ExtractJWTMetadata(token *jwt.Token) (*TokenMetadata, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidJWTMetadata
	}

	if !token.Valid {
		return nil, ErrInvalidJWT
	}

	uid, err := uuid.Parse(claims["uid"].(string))
	if err != nil {
		return nil, err
	}

	return &TokenMetadata{
		UserID: uid,
	}, nil
}
