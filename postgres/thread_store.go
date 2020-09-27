package postgres

import (
	"fmt"

	"github.com/gocs/miji"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ThreadStore struct {
	*sqlx.DB
}

func NewThreadStore(db *sqlx.DB) *ThreadStore {
	return &ThreadStore{DB: db}
}

func (s *ThreadStore) Thread(id uuid.UUID) (miji.Thread, error) {
	var t miji.Thread
	if err := s.Get(&t, `SELECT * FROM threads WHERE id = $1`, id); err != nil {
		return miji.Thread{}, fmt.Errorf("error getting thread: %W", err)
	}

	return t, nil
}

func (s *ThreadStore) Threads() ([]miji.Thread, error) {
	var tt []miji.Thread

	if err := s.Get(&tt, `SELECT * FROM threads`); err != nil {
		return []miji.Thread{}, fmt.Errorf("error getting threads: %W", err)
	}

	return tt, nil
}

func (s *ThreadStore) CreateThread(t *miji.Thread) error {
	if err := s.Get(t, `INSERT INTO threads VALUES ($1, $2, $3) RETURNING *`,
		t.ID, t.Title, t.Description); err != nil {
		return fmt.Errorf("error creating thread: %W", err)
	}
	return nil
}

func (s *ThreadStore) UpdateThread(t *miji.Thread) error {
	if err := s.Get(t, `UPDATE threads SET title = $1, description = $2 WHERE id = $3 RETURNING *`,
		t.Title, t.Description, t.ID); err != nil {
		return fmt.Errorf("error updating thread: %W", err)
	}
	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM threads WHERE id = $1`, id); err != nil {
	return fmt.Errorf("error deleting thread: %W", err)
	}
	return nil
}
