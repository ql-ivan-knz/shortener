package storage

import (
	"context"
	"shortener/config"
	"shortener/internal/models"
	"shortener/internal/storage/db"
	"shortener/internal/storage/file"
	mapStorage "shortener/internal/storage/map"
)

type Storage interface {
	Get(ctx context.Context, key string) (models.URLItem, error)
	Put(ctx context.Context, key, value string, userID string) error
	Batch(ctx context.Context, urls []models.URLItem, userID string) error
	GetAllURLs(ctx context.Context, userID string) ([]models.URLItem, error)
	DeleteURLs(ctx context.Context, urls []string, userID string) error
	Ping(ctx context.Context) error
}

func NewStorage(config config.Config) (Storage, error) {
	if config.DatabaseDSN != "" {
		return db.NewStorage(config.DatabaseDSN)
	}

	if config.FileStoragePath != "" {
		return file.NewStorage(config.FileStoragePath)
	}

	return mapStorage.NewStorage()
}
