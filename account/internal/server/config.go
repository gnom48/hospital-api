package server

import "account/internal/storage"

type Config struct {
	BindAddress   string `toml:"bind_address"`
	LogLevel      string `toml:"log_level"`
	StorageConfig *storage.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddress:   ":8081",
		LogLevel:      "debug",
		StorageConfig: storage.NewConfig(),
	}
}
