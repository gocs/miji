package postgres

import (
	"fmt"

	"github.com/gocs/miji"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) Comment(id uuid.UUID) (miji.Comment, error) {
	var c miji.Comment
	if err := s.Get(&c, `SELECT * FROM comments WHERE id = $1`, id); err != nil {
		return miji.Comment{}, fmt.Errorf("error getting comment: %w", err)
	}

	return c, nil
}

func (s *CommentStore) CommentsByPost(postID uuid.UUID) ([]miji.Comment, error) {
	var cc []miji.Comment

	if err := s.Select(&cc, `SELECT * FROM comments WHERE post_id = $1`, postID); err != nil {
		return []miji.Comment{}, fmt.Errorf("error getting comments: %w", err)
	}

	return cc, nil
}

func (s *CommentStore) CreateComment(c *miji.Comment) error {
	if err := s.Get(c, `INSERT INTO comments VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		c.ID, c.PostID, c.Content, c.Votes); err != nil {
		return fmt.Errorf("error creating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(c *miji.Comment) error {
	if err := s.Get(c, `UPDATE comments SET thread_id = $1, content = $2, votes = $3 WHERE id = $4 RETURNING *`,
		c.PostID, c.Content, c.Votes, c.ID); err != nil {
		return fmt.Errorf("error updating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM comments WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}
	return nil
}