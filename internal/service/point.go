package service

import (
	"context"

	"kekaton/back/internal/storage"
)

func (s *Service) RegisterPoint(ctx context.Context, point *storage.Point) error {
	return s.storage.CreatePoint(ctx, point)
}

func (s *Service) GetPointByID(ctx context.Context, point *storage.Point) error {
	if err := s.storage.GetPointByID(ctx, point); err != nil {
		return err
	}

	if err := s.storage.GetUserByID(ctx, &point.Creator); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetPoints(ctx context.Context, points *[]storage.Point) error {
	if err := s.storage.GetPoints(ctx, points); err != nil {
		return err
	}

	newPoints := *points

	// Currently no need to optimize.
	// Keep rolling.
	for i := range newPoints {
		if err := s.storage.GetUserByID(ctx, &newPoints[i].Creator); err != nil {
			return err
		}
	}

	return nil
}
