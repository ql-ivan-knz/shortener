package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type storage struct {
	db *sql.DB
}

func (s *storage) Get(key string) (string, error) {
	return "", nil
}

func (s *storage) Put() {

}

func (s *storage) Set(key string, value string) error {
	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func NewStorage(databaseDSN string) (*storage, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	return &storage{db: db}, nil
}
