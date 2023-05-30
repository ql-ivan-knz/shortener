package storage

type MapStorage map[string]string

var storage MapStorage = make(map[string]string)

func (s MapStorage) Get(key string) (string, error) {
	v, ok := s[key]
	if !ok {
		return "", nil
	}

	return v, nil
}

func (s MapStorage) Set(key, value string) error {
	s[key] = value

	return nil
}

func NewMapStorage() Storage {
	return storage
}
