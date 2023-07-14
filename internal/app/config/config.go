package config

import (
	"flag"
	"os"
)

type Config struct {
	StartAddr       string
	ShortURLAddr    string
	StorageType     string
	LogLevel        string
	FileStoragePath string
}

func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.StartAddr, "a", "localhost:8080", "The starting address (format: host:port)")
	flag.StringVar(&cfg.ShortURLAddr, "b", "http://localhost:8080", "The short URL server address string")
	flag.StringVar(&cfg.StorageType, "s", "file", "type of storage to use (memory/file)")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db.json", "Path to file storage")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.StartAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.ShortURLAddr = envBaseAddr
	}
	if envStorageType := os.Getenv("STORAGE_TYPE"); envStorageType != "" {
		cfg.StorageType = envStorageType
	}
	if cfg.StorageType == "file" {
		if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
			cfg.FileStoragePath = envFileStoragePath
		}
	}
	if envLogLevel := os.Getenv("LOGLEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	}

	return cfg
}
