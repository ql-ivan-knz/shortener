package db

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/models"
	"shortener/internal/short"
	"strings"
	"time"
)

type storage struct {
	pool *pgxpool.Pool
}

var ErrorConflict = errors.New("url is already saved")

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

	_, err = pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS id_url ON urls (original)`)
	if err != nil {
		return nil, fmt.Errorf("failed to set index: %v", err)
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
	tag, err := s.pool.Exec(ctx, "INSERT INTO urls (short, original) VALUES ($1, $2) ON CONFLICT (original) DO NOTHING", key, value)
	if err != nil {
		return fmt.Errorf("failed to insert row: %v", err)
	}

	count := tag.RowsAffected()

	if count == 0 {
		return ErrorConflict
	}

	return nil
}

func (s *storage) Batch(ctx context.Context, urls models.BatchRequest) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	values := make([]string, len(urls))

	for index, url := range urls {
		values[index] = "('" + short.URL([]byte(url.OriginalURL)) + "', '" + url.OriginalURL + "')"
	}

	query := "INSERT INTO urls (short, original) VALUES " + strings.Join(values, ", ")

	_, err = tx.Exec(ctx, query)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("failed transacton rollback %v", err)
		}
		return fmt.Errorf("failed to insert line in table with %v", err)
	}

	return tx.Commit(ctx)
}

func (s *storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
