package service

import (
	"context"

	"kekaton/back/internal/storage"
)

func (s *Service) GetTagByID(ctx context.Context, tag *storage.Tag) error {
	return s.storage.GetTagByID(ctx, tag)
}

func (s *Service) GetTags(ctx context.Context, tags *[]storage.Tag) error {
	return s.storage.GetTags(ctx, tags)
}
