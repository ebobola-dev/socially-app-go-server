package config

import (
	"os"

	"github.com/ebobola-dev/socially-app-go-server/internal/util/env"
)

type SMTPConfig struct {
	ADDRESS  string
	PASSWORD string
	HOST     string
	PORT     int
}

func LoadSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		ADDRESS:  os.Getenv("APP_EMAIL_ADDRESS"),
		PASSWORD: os.Getenv("APP_EMAIL_PASSWORD"),
		HOST:     os.Getenv("APP_EMAIL_HOST"),
		PORT:     env.GetInt("APP_EMAIL_PORT"),
	}
}
