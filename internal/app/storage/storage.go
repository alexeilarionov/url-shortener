package storage

type Storage interface {
	Store(key string, value string) error
	Get(key string) (string, error)
}
