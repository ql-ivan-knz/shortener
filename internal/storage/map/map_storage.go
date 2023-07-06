package mapstorage

import (
	"context"
	"shortener/internal/models"
)

type storage struct {
	records map[string]storageItem
}

type storageItem struct {
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

func (s *storage) Get(ctx context.Context, key string) (models.URLItem, error) {
	v := s.records[key]

	return models.URLItem{
		OriginalURL: v.OriginalURL,
		ShortURL:    key,
		IsDeleted:   v.IsDeleted,
	}, nil
}

func (s *storage) Put(ctx context.Context, key, value, userID string) error {
	s.records[key] = storageItem{OriginalURL: value, UserID: userID, IsDeleted: false}

	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	return nil
}

func (s *storage) Batch(ctx context.Context, urls []models.URLItem, userID string) error {
	for _, url := range urls {
		if err := s.Put(ctx, url.ShortURL, url.OriginalURL, userID); err != nil {
			return err
		}
	}
	return nil
}

func (s *storage) DeleteURLs(ctx context.Context, shortURLs []string, userID string) error {
	for _, shortURL := range shortURLs {
		item, ok := s.records[shortURL]
		if ok && item.UserID == userID {
			s.records[shortURL] = storageItem{OriginalURL: item.OriginalURL, UserID: userID, IsDeleted: true}
		}
	}
	return nil
}

func (s *storage) GetAllURLs(ctx context.Context, userID string) ([]models.URLItem, error) {
	return nil, nil
}

func NewStorage() (*storage, error) {
	return &storage{records: map[string]storageItem{}}, nil
}
