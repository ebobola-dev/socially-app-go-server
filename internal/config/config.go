package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func Initialize() *Config {
	_ = godotenv.Load()
	return &Config{
		Port: _getEnv("PORT", "8080"),
	}
}

func _getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
