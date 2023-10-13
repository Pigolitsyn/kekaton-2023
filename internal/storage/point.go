package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Point struct {
	ID          int           `json:"id"`
	Coordinates Coordinates   `json:"coordinates"`
	Description string        `json:"description"`
	OpenTime    time.Duration `json:"open"`
	CloseTime   time.Duration `json:"close"`
	CreatedBy   uuid.UUID     `json:"creator"`
}

type Coordinates [2]float64

func (s *Storage) CreatePoint(ctx context.Context, point *Point) error {
	query := `INSERT INTO points (coordinates, description, open_time, close_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	row := s.pool.QueryRow(
		ctx,
		query,
		pgtype.Point{
			P: pgtype.Vec2{
				X: point.Coordinates[0],
				Y: point.Coordinates[1],
			},
			Valid: true,
		},
		point.Description,
		pgtype.Time{
			Microseconds: point.OpenTime.Microseconds(),
			Valid:        true,
		},
		pgtype.Time{
			Microseconds: point.CloseTime.Microseconds(),
			Valid:        true,
		},
		point.CreatedBy,
	)

	if err := row.Scan(&point.ID); err != nil {
		return err
	}

	return nil
}
