package mapstorage

import (
	"context"
	"shortener/internal/models"
	"shortener/internal/short"
)

type storage struct {
	records map[string]string
}

func (s *storage) Get(ctx context.Context, key string) (string, error) {
	v, ok := s.records[key]
	if !ok {
		return "", nil
	}

	return v, nil
}

func (s *storage) Put(ctx context.Context, key, value string) error {
	s.records[key] = value

	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	return nil
}

func (s *storage) Batch(ctx context.Context, urls models.BatchRequest) error {
	for _, url := range urls {
		if err := s.Put(ctx, short.URL([]byte(url.OriginalURL)), url.OriginalURL); err != nil {
			return err
		}
	}
	return nil
}

func NewStorage() (*storage, error) {
	return &storage{records: map[string]string{}}, nil
}
