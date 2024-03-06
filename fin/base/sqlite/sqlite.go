package sqlite

import (
	"GB_Study_02/fin/base"
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

// func (s *Storage) Query(ctx context.Context)(v any...,error){

// }

func (s *Storage) RemoveUser(ctx context.Context, user *base.User) error {
	q := `DELETE FROM users WHERE id = ?`
	if _, err := s.db.ExecContext(ctx, q, user.Id); err != nil {
		return fmt.Errorf("can't delete User: %s", err)
	}

	return nil
}

func (s *Storage) AddUser(ctx context.Context, u *base.User) error {
	q := `INSERT INTO users (name,age) VALUES (?,?)`

	if _, err := s.db.ExecContext(ctx, q, u.Name, u.Age); err != nil {
		return fmt.Errorf("can't add User: %w", err)
	}

	return nil
}

func (s *Storage) FriendsList(ctx context.Context, u *base.User) (*base.Stora, error) {
	q := `SELECT id1 AS l,users.user_name,users.age FROM(
		SELECT friends.id1 FROM friends
		WHERE id2=?
		UNION
		SELECT friends.id2 FROM friends
		WHERE id1=?
		)s
		LEFT JOIN users ON l=id`

	var id, age int
	var name string
	if err := s.db.QueryRowContext(ctx, q, u.Id, u.Id).Scan(&id, &name, &age); err != nil {
		return nil, fmt.Errorf("can't read Users: %w", err)
	}

	return nil, nil
}
