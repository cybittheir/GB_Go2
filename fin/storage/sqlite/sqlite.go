package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY, 
			user_name TEXT, 
			age INTEGER
			);
		CREATE TABLE IF NOT EXISTS friends (
			id1 INTEGER, 
			id2 INTEGER
			);`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database %w", err)
	}
	return &Storage{db: db}, nil
}
