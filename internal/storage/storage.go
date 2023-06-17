package storage

import (
	"context"
	"shortener/config"
	"shortener/internal/storage/db"
	"shortener/internal/storage/file"
	mapStorage "shortener/internal/storage/map"
)

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
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
