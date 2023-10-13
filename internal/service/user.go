package service

import (
	"context"
	"errors"
	"github.com/google/uuid"

	"kekaton/back/internal/storage"
)

var (
	ErrEmptyUserEmail = errors.New("empty user email")
	ErrEmptyUserID    = errors.New("empty user id")
)

func (s *Service) RegisterUser(ctx context.Context, user *storage.User) error {
	hased, err := s.MakePassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hased.Hash
	user.Salt = hased.Salt

	return s.storage.CreateUser(ctx, user)
}

func (s *Service) GetUserByID(ctx context.Context, user *storage.User) error {
	if user.ID == uuid.Nil {
		return ErrEmptyUserID
	}

	return s.storage.GetUserByID(ctx, user)
}

func (s *Service) GetUserByEmail(ctx context.Context, user *storage.User) error {
	if user.Email == "" {
		return ErrEmptyUserEmail
	}

	return s.storage.GetUserByEmail(ctx, user)
}
