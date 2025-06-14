package config

import (
	"os"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Name:     os.Getenv("MYSQL_NAME"),
	}
}
