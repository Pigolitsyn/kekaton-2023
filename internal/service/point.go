package service

import (
	"context"

	"kekaton/back/internal/storage"
)

func (s *Service) RegisterPoint(ctx context.Context, point *storage.Point, tags *[]int) error {
	if err := s.storage.CreatePoint(ctx, point); err != nil {
		return err
	}

	if len(*tags) > 0 {
		if err := s.storage.AddTagsToPoint(ctx, point.ID, tags); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetPointByID(ctx context.Context, point *storage.Point) error {
	if err := s.storage.GetPointByID(ctx, point); err != nil {
		return err
	}

	if err := s.storage.GetUserByID(ctx, &point.Creator); err != nil {
		return err
	}

	if err := s.storage.GetTagForPoint(ctx, point.ID, &point.Tags); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetPoints(ctx context.Context, points *[]storage.Point) error {
	if err := s.storage.GetPoints(ctx, points); err != nil {
		return err
	}

	newPoints := *points

	for i := range newPoints {
		if err := s.storage.GetUserByID(ctx, &newPoints[i].Creator); err != nil {
			return err
		}

		if err := s.storage.GetTagForPoint(ctx, newPoints[i].ID, &newPoints[i].Tags); err != nil {
			return err
		}
	}

	return nil
}
