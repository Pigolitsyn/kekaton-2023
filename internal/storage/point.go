package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Point struct {
	ID          int           `json:"id"`
	Coordinates Coordinates   `json:"coordinates"`
	Address     string        `json:"address"`
	Description string        `json:"description"`
	OpenTime    time.Duration `json:"open"`
	CloseTime   time.Duration `json:"close"`
	Creator     User          `json:"creator"`
	Tags        []Tag         `json:"tags"`
}

type Coordinates [2]float64

func (s *Storage) CreatePoint(ctx context.Context, point *Point) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `INSERT INTO points (coordinates, address, description, open_time, close_time, created_by) VALUES (POINT($1, $2), $3, $4, $5, $6, $7) RETURNING id`

	row := s.pool.QueryRow(
		ctx,
		query,
		point.Coordinates[0],
		point.Coordinates[1],
		point.Address,
		point.Description,
		pgtype.Time{
			Microseconds: point.OpenTime.Microseconds(),
			Valid:        true,
		},
		pgtype.Time{
			Microseconds: point.CloseTime.Microseconds(),
			Valid:        true,
		},
		point.Creator.ID,
	)

	if err := row.Scan(&point.ID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetPointByID(ctx context.Context, point *Point) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT coordinates, address, description, open_time, close_time, created_by FROM points WHERE id = $1`

	var (
		coords    = pgtype.Point{}
		openTime  = pgtype.Time{}
		closeTime = pgtype.Time{}
	)

	if err := s.pool.QueryRow(ctx, query, point.ID).Scan(&coords, &point.Address, &point.Description, &openTime, &closeTime, &point.Creator.ID); err != nil {
		return err
	}

	point.Coordinates = [2]float64{coords.P.X, coords.P.Y}
	point.OpenTime = time.Duration(openTime.Microseconds)
	point.CloseTime = time.Duration(closeTime.Microseconds)

	return nil
}

func (s *Storage) GetPoints(ctx context.Context, points *[]Point) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT id, coordinates, address, description, open_time, close_time, created_by FROM points`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	newPoints := make([]Point, 0)

	for rows.Next() {
		point := Point{}

		coords := pgtype.Point{}
		openTime := pgtype.Time{}
		closeTime := pgtype.Time{}

		if err = rows.Scan(&point.ID, &coords, &point.Address, &point.Description, &openTime, &closeTime, &point.Creator.ID); err != nil {
			return err
		}

		point.Coordinates = [2]float64{coords.P.X, coords.P.Y}
		point.OpenTime = time.Duration(openTime.Microseconds)
		point.CloseTime = time.Duration(closeTime.Microseconds)

		newPoints = append(newPoints, point)
	}

	*points = newPoints

	return nil
}

func (s *Storage) UpdatePoint(ctx context.Context, point *Point) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `UPDATE points SET coordinates = POINT($2, $3), address = $4, description = $5, open_time = $6, close_time = $7 WHERE id = $1`

	var (
		openTime = pgtype.Time{
			Microseconds: point.OpenTime.Microseconds(),
			Valid:        true,
		}
		closeTime = pgtype.Time{
			Microseconds: point.CloseTime.Microseconds(),
			Valid:        true,
		}
	)

	if _, err := s.pool.Query(ctx, query, point.ID, point.Coordinates[0], point.Coordinates[1], point.Address, point.Description, openTime, closeTime); err != nil {
		return err
	}

	return nil
}
