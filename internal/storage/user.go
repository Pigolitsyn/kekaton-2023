package storage

import (
	"context"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Username string
	Password string
	Salt     string
}

func (s *Storage) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, username, password, salt) VALUES ($1, $2, $3, $4) RETURNING id`

	row := s.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.Username,
		user.Password,
		user.Salt,
	)

	if err := row.Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, user *User) error {
	query := `SELECT id, username, password, salt FROM users WHERE email = $1`

	if err := s.pool.QueryRow(ctx, query, user.Email).Scan(&user.ID, &user.Username, &user.Password, &user.Salt); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserByID(ctx context.Context, user *User) error {
	query := `SELECT email, username, password, salt FROM users WHERE id = $1`

	if err := s.pool.QueryRow(ctx, query, user.ID).Scan(&user.Email, &user.Username, &user.Password, &user.Salt); err != nil {
		return err
	}

	return nil
}
