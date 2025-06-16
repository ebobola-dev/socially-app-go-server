package config

import (
	"os"

	env "github.com/ebobola-dev/socially-app-go-server/internal/util/env"
)

type JWTConfig struct {
	ENCODE_ALGORITNM        string
	ACCESS_SERCER_KEY       []byte
	REFRESH_SERCER_KEY      []byte
	ACCESS_DURABILITY_MIN   int
	REFRESH_DURABILITY_DAYS int
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		ENCODE_ALGORITNM:        os.Getenv("JWT_ENCODE_ALGORITNM"),
		ACCESS_SERCER_KEY:       []byte(os.Getenv("JWT_ACCESS_SERCER_KEY")),
		REFRESH_SERCER_KEY:      []byte(os.Getenv("JWT_REFRESH_SERCER_KEY")),
		ACCESS_DURABILITY_MIN:   env.GetInt("JWT_ACCESS_DURABILITY_MIN"),
		REFRESH_DURABILITY_DAYS: env.GetInt("JWT_REFRESH_DURABILITY_DAYS"),
	}
}
