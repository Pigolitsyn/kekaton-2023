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

func (s *Storage) GetTagsForPoint(ctx context.Context, pid int, tags *[]Tag) error {
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

	newTags := *tags

	var values string

	for i := range newTags {
		if i == 0 {
			values = fmt.Sprintf("(%v, %v)", pid, newTags[i])

			continue
		}

		values = fmt.Sprintf("%v, (%v, %v)", values, pid, newTags[i])
	}

	query += values

	if _, err := s.pool.Query(ctx, query); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdatePointTags(ctx context.Context, pid int, tags *[]int) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	ognTags := make([]Tag, 0)

	if err := s.GetTagsForPoint(ctx, pid, &ognTags); err != nil {
		return err
	}

	newTags := *tags

	oldTags := make([]int, len(ognTags))

	for i := range oldTags {
		oldTags[i] = ognTags[i].Type
	}

	add, rem := tagDifference(oldTags, newTags)

	var (
		valueAdd string
		valueRem string
	)

	for i := range add {
		if i == 0 {
			valueAdd = fmt.Sprintf("(%v, %v)", pid, add[i])

			continue
		}

		valueAdd = fmt.Sprintf("%v, (%v, %v)", valueAdd, pid, add[i])
	}

	for i := range rem {
		if i == 0 {
			valueRem = fmt.Sprintf("%v", rem[i])

			continue
		}

		valueRem = fmt.Sprintf("%v, %v", valueRem, rem[i])
	}

	queryRem := `DELETE FROM tags WHERE point_id = $1 AND type_id = ANY($2::INT[])`
	queryAdd := `INSERT INTO tags (point_id, type_id) VALUES ` + valueAdd

	if _, err := s.pool.Query(ctx, queryRem, pid, "{"+valueRem+"}"); err != nil {
		return err
	}

	if _, err := s.pool.Query(ctx, queryAdd); err != nil {
		return err
	}

	return nil
}

func tagDifference(a, b []int) ([]int, []int) {
	var (
		ina = make(map[int]struct{}, len(a))
		inb = make(map[int]struct{}, len(b))
	)

	for _, x := range a {
		ina[x] = struct{}{}
	}

	for _, x := range b {
		inb[x] = struct{}{}
	}

	var (
		add = make([]int, 0)
		rem = make([]int, 0)
	)

	for _, x := range a {
		if _, ok := inb[x]; !ok {
			rem = append(rem, x)
		}
	}

	for _, x := range b {
		if _, ok := ina[x]; !ok {
			add = append(add, x)
		}
	}

	return add, rem
}
