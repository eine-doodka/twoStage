package config

import "time"

type Config struct {
	RedisAddr       string        `toml:"redis_addr"`
	StorageDuration time.Duration `toml:"storage_duration"`
	CodeLength      int           `toml:"code_length"`
}

func Default() *Config {
	return &Config{
		StorageDuration: 30 * time.Second,
	}
}
