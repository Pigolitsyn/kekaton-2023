package service

import (
	"context"

	"kekaton/back/internal/storage"
)

func (s *Service) CreateComment(ctx context.Context, comment *storage.Comment) error {
	return s.storage.CreateComment(ctx, comment)
}

func (s *Service) UpdateComment(ctx context.Context, comment *storage.Comment) error {
	return s.storage.UpdateComment(ctx, comment)
}

func (s *Service) DeleteComment(ctx context.Context, comment *storage.Comment) error {
	return s.storage.DeleteComment(ctx, comment)
}

func (s *Service) GetComment(ctx context.Context, comment *storage.Comment) error {
	return s.storage.GetComment(ctx, comment)
}

func (s *Service) GetCommentsForPoint(ctx context.Context, pointID int, comments *[]storage.Comment) error {
	return s.storage.GetCommentsForPoint(ctx, pointID, comments)
}
