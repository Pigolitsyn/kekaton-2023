package storage

import (
	"context"

	"github.com/google/uuid"
)

type Comment struct {
	ID      int       `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	PointID int       `json:"point_id"`
	Text    string    `json:"text"`
	Rating  int8      `json:"rating"`
}

func (s *Storage) CreateComment(ctx context.Context, comment *Comment) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `INSERT INTO comments (user_id, point_id, comment_text, rating) VALUES ($1, $2, $3, $4) RETURNING id`

	row := s.pool.QueryRow(
		ctx,
		query,
		comment.UserID,
		comment.PointID,
		comment.Text,
		comment.Rating,
	)

	if err := row.Scan(&comment.ID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetComment(ctx context.Context, comment *Comment) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := `SELECT user_id, point_id, comment_text, rating FROM comments WHERE comment_id = $1`

	if err := s.pool.QueryRow(ctx, query, comment.ID).Scan(
		&comment.UserID,
		&comment.PointID,
		&comment.Text,
		&comment.Rating,
	); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateComment(ctx context.Context, comment *Comment) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := `UPDATE comments SET comment_text = $1, rating = $2 WHERE comment_id = $3 and user_id = $4`

	if _, err := s.pool.Query(ctx, query, comment.Text, comment.Rating, comment.ID, comment.UserID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteComment(ctx context.Context, comment *Comment) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := `DELETE FROM comments WHERE comment_id = $1`

	if _, err := s.pool.Query(ctx, query, comment.ID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetCommentsForPoint(ctx context.Context, pointID int, comments *[]Comment) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := `SELECT user_id, point_id, comment_text, rating FROM comments WHERE point_id = $1`

	rows, err := s.pool.Query(ctx, query, pointID)
	if err != nil {
		return err
	}

	newComment := make([]Comment, 0)

	for rows.Next() {
		var comment Comment

		if err = rows.Scan(
			&comment.UserID,
			&comment.PointID,
			&comment.Text,
			&comment.Rating,
		); err != nil {
			return err
		}

		newComment = append(newComment, comment)
	}

	*comments = newComment

	return nil
}
