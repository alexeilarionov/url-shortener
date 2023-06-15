package config

import (
	"flag"
	"os"
)

type Config struct {
	StartAddr    string
	ShortURLAddr string
	StorageType  string
}

func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.StartAddr, "a", "localhost:8080", "The starting address (format: host:port)")
	flag.StringVar(&cfg.ShortURLAddr, "b", "http://localhost:8080", "The short URL server address string")
	flag.StringVar(&cfg.StorageType, "s", "memory", "type of storage to use (memory)")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.StartAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.ShortURLAddr = envBaseAddr
	}

	return cfg
}
