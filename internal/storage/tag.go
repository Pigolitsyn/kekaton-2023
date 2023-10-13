package storage

import "context"

type Tag struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

func (s *Storage) GetTagByID(ctx context.Context, tag *Tag) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT name FROM tag_types WHERE id = $1`

	if err := s.pool.QueryRow(ctx, query, tag.Type).Scan(&tag.Name); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTags(ctx context.Context, tags *[]Tag) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT id, name FROM tag_types`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	newTags := make([]Tag, 0)

	for rows.Next() {
		tag := Tag{}

		if err = rows.Scan(&tag.Type, &tag.Name); err != nil {
			return err
		}

		newTags = append(newTags, tag)
	}

	*tags = newTags

	return nil
}
