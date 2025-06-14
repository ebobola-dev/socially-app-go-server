package config

import "os"

type SMTPConfig struct {
	EMAIL_ADDRESS  string
	EMAIL_PASSWORD string
}

func LoadSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		EMAIL_ADDRESS:  os.Getenv("APP_EMAIL_ADDRESS"),
		EMAIL_PASSWORD: os.Getenv("APP_EMAIL_PASSWORD"),
	}
}
