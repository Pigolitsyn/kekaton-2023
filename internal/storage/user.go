package storage

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	Salt     string    `json:"-"`
}

func (s *Storage) CreateUser(ctx context.Context, user *User) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

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
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT id, username, password, salt FROM users WHERE email = $1`

	if err := s.pool.QueryRow(ctx, query, user.Email).Scan(&user.ID, &user.Username, &user.Password, &user.Salt); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserByID(ctx context.Context, user *User) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT email, username, password, salt FROM users WHERE id = $1`

	if err := s.pool.QueryRow(ctx, query, user.ID).Scan(&user.Email, &user.Username, &user.Password, &user.Salt); err != nil {
		return err
	}

	return nil
}
