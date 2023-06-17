package mapstorage

import (
	"context"
)

type storage struct {
	records map[string]string
}

func (s *storage) Get(key string) (string, error) {
	v, ok := s.records[key]
	if !ok {
		return "", nil
	}

	return v, nil
}

func (s *storage) Set(key, value string) error {
	s.records[key] = value

	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	return nil
}

func NewStorage() (*storage, error) {
	return &storage{records: map[string]string{}}, nil
}
