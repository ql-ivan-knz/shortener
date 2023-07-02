package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"shortener/internal/models"
	"time"
)

type storage struct {
	pool *pgxpool.Pool
}

var ErrorConflict = errors.New("url is already saved")

func NewStorage(dsn string) (*storage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}
	return &storage{
		pool: pool,
	}, nil
}

//go:embed migrations/*.sql
var fs embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iosf driver: %v", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %v", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %v", err)
		}
	}

	return nil
}

func (s *storage) Get(ctx context.Context, key string) (string, error) {
	row := s.pool.QueryRow(
		ctx,
		`SELECT original_url FROM links WHERE hash_url = $1`,
		key,
	)

	var url string
	if err := row.Scan(&url); err != nil {
		return "", fmt.Errorf("failed to read row: %v", err)
	}
	return url, nil
}

func (s *storage) Put(ctx context.Context, key string, value string, userID string) error {
	tag, err := s.pool.Exec(
		ctx,
		`INSERT INTO links (hash_url, original_url, user_id) 
		 VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING`,
		key, value, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert row: %v", err)
	}

	count := tag.RowsAffected()

	if count == 0 {
		return ErrorConflict
	}

	return nil
}

func (s *storage) Batch(ctx context.Context, rows []models.URLItem, userID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			log.Println("transaction rolled back")
		}
	}()

	batch := &pgx.Batch{}
	for _, r := range rows {
		batch.Queue(
			`INSERT INTO links (hash_url, original_url, user_id) VALUES ($1, $2, $3)`,
			r.ShortURL, r.OriginalURL, userID,
		)
	}

	results := s.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("error executing statement %v", err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *storage) GetAllURLs(ctx context.Context, userID string) ([]models.URLItem, error) {
	rows, err := s.pool.Query(ctx, "SELECT hash_url, original_url FROM links WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed urls belong to userID=%s: %v", userID, err)
	}
	defer rows.Close()

	var urls []models.URLItem
	for rows.Next() {
		var row models.URLItem

		err = rows.Scan(&row.ShortURL, &row.OriginalURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan urls %v", err)
		}

		urls = append(urls, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed getting urls %v", err)
	}

	return urls, nil
}

func (s *storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
