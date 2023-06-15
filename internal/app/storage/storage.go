package storage

type Storage interface {
	Store(key string, value string) error
	Get(key string) (string, error)
}

func NewStorage(storageType string) Storage {
	switch storageType {
	case "memory":
		return NewInMemoryStorage()
	default:
		return NewInMemoryStorage()
	}
}
