package storage

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

func NewStorage(path string) Storage {
	if path == "" {
		return NewMapStorage()
	}

	return NewFileStorage(path)
}
