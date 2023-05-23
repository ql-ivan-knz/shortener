package storage

type Storage map[string]string

var storage Storage = make(map[string]string)

func (s Storage) Set(key, value string) {
	s[key] = value
}

func (s Storage) Get(key string) (string, bool) {
	v, ok := s[key]

	return v, ok
}

func NewStorage() *Storage {
	return &storage
}
