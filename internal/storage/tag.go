package storage

import (
	"context"
	"fmt"
)

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

func (s *Storage) GetTagForPoint(ctx context.Context, pid int, tags *[]Tag) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT type_id, tag_types.name FROM tags LEFT JOIN tag_types on tags.type_id = tag_types.id WHERE point_id = $1`

	rows, err := s.pool.Query(ctx, query, pid)
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

func (s *Storage) AddTagsToPoint(ctx context.Context, pid int, tags *[]int) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `INSERT INTO tags (point_id, type_id) VALUES `

	addTags := *tags

	var values string

	for i := range addTags {
		if i == 0 {
			values = fmt.Sprintf("(%v, %v)", pid, addTags[i])

			continue
		}

		values = fmt.Sprintf("%v, (%v, %v)", values, pid, addTags[i])
	}

	query += values

	if _, err := s.pool.Query(ctx, query); err != nil {
		return err
	}

	return nil
}
