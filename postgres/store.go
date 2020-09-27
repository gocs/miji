package postgres

import (
	"fmt"

	"github.com/gocs/miji"
	"github.com/jmoiron/sqlx"

	// go inject the pq driver
	_ "github.com/lib/pq"
)

type Store struct {
	miji.ThreadStore
	miji.PostStore
	miji.CommentStore
}

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Store{
		ThreadStore:  NewThreadStore(db),
		PostStore:    NewPostStore(db),
		CommentStore: NewCommentStore(db),
	}, nil
}
