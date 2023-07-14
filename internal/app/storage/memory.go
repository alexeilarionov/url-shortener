package storage

import "errors"

type InMemoryStorage struct {
	data map[string]ShortenedData
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]ShortenedData),
	}
}

func (s *InMemoryStorage) Store(data ShortenedData) error {
	s.data[data.ShortURL] = data
	return nil
}

func (s *InMemoryStorage) Get(key string) (ShortenedData, error) {
	value, exists := s.data[key]
	if !exists {
		return ShortenedData{}, errors.New("key not found: " + key)
	}
	return value, nil
}
