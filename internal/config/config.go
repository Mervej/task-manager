package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DatabasePath string `json:"database_path"`
	Port         string `json:"port"`
	LogLevel     string `json:"log_level"`
}

func LoadConfig() (*Config, error) {
	// set the default db paht
	defaultDBPath := filepath.Join("data", "tasks.db")

	config := &Config{
		DatabasePath: getEnv("DATABASE_PATH", defaultDBPath),
		Port:         getEnv("PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
