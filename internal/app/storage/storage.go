package storage

import "github.com/alexeilarionov/url-shortener/internal/app/config"

type ShortenedData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	Store(data ShortenedData) error
	Get(key string) (ShortenedData, error)
}

func NewStorage(cfg config.Config) Storage {
	switch cfg.StorageType {
	case "memory":
		return NewInMemoryStorage()
	case "file":
		return NewFileStorage(cfg.FileStoragePath)
	default:
		return NewInMemoryStorage()
	}
}
