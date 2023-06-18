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
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key, value string) error
	Batch(ctx context.Context, urls models.BatchRequest) error
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
