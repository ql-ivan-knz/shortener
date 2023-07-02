package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"shortener/internal/models"
	"strconv"
)

type storage struct {
	filePath string
	numLines int
}

type ShortenURL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

var increment = 0

func countLines(path string) int {
	fmt.Println("init increment")
	count := 0

	file, _ := os.OpenFile(path, os.O_RDONLY, 0666)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		count++
	}

	return count
}

func (s *storage) Put(ctx context.Context, key, value, userID string) error {
	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		su := ShortenURL{}

		if err := json.Unmarshal(scanner.Bytes(), &su); err != nil {
			return err
		}

		if su.ShortURL == key {
			return nil
		}
	}

	increment++
	su := ShortenURL{UUID: strconv.Itoa(increment), ShortURL: key, OriginalURL: value, UserID: userID}
	data, err := json.Marshal(&su)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) Get(ctx context.Context, key string) (string, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("in get")
		su := ShortenURL{}

		if err := json.Unmarshal(scanner.Bytes(), &su); err != nil {
			return "", err
		}

		if su.ShortURL == key {
			return su.OriginalURL, nil
		}
	}

	return "", nil
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

func (s *storage) GetAllURLs(ctx context.Context, userID string) ([]models.URLItem, error) {
	return nil, nil
}

func NewStorage(path string) (*storage, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create file if no exists
		err = os.MkdirAll(filepath.Dir(path), 0666)
		if err != nil {
			return nil, err
		}

		_, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}

	return &storage{filePath: path, numLines: countLines(path)}, nil
}
