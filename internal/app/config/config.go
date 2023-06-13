package config

import (
	"flag"
)

type Config struct {
	StartAddr    string
	ShortUrlAddr string
	Storage      string
}

func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.StartAddr, "a", "localhost:8080", "The starting address (format: host:port)")
	flag.StringVar(&cfg.ShortUrlAddr, "b", "http://localhost:8080", "The short URL server address string")
	flag.StringVar(&cfg.Storage, "s", "memory", "type of storage to use (memory)")
	flag.Parse()
	return cfg
}
