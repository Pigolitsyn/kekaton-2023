package service

import (
	"context"
	"kekaton/back/internal/storage"
)

func (s *Service) RegisterPoint(ctx context.Context, point *storage.Point) error {
	return s.storage.CreatePoint(ctx, point)
}
