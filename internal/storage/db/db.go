package db

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type storage struct {
	pool *pgxpool.Pool
}

func NewStorage(dsn string) (*storage, error) {
	pool, err := initDB(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &storage{
		pool: pool,
	}, nil
}

func initDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS urls (short VARCHAR(8) PRIMARY KEY , original TEXT)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return pool, nil
}

func (s *storage) Get(ctx context.Context, key string) (string, error) {
	row := s.pool.QueryRow(ctx, `SELECT original FROM urls WHERE short = $1`, key)

	var url string
	if err := row.Scan(&url); err != nil {
		return "", fmt.Errorf("failed to read row: %v", err)
	}
	return url, nil
}

func (s *storage) Put(ctx context.Context, key string, value string) error {
	_, err := s.pool.Exec(
		ctx,
		`INSERT INTO urls (short, original) VALUES ($1, $2)`, key, value,
	)
	if err != nil {
		return fmt.Errorf("failed insert url: %w", err)
	}

	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
