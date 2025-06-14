package config

import (
	"log"

	"github.com/ebobola-dev/socially-app-go-server/internal/util/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	BuildType BuildType
	OwnerKey  string

	Database *DatabaseConfig
	JWT      *JWTConfig
	SMTP     *SMTPConfig
	Minio    *MinioConfig
}

func Initialize() *Config {
	_ = godotenv.Load()
	btStr := env.GetString("BUILD_TYPE")
	buildType, err := ParseBuildType(btStr)

	if err != nil {
		log.Printf("Invalid BUILD_TYPE '%s', falling back to DEV\n", btStr)
		buildType = Development
	}

	return &Config{
		Port:      env.GetString("INTERNAL_PORT"),
		BuildType: buildType,
		OwnerKey:  env.GetString("OWNER_KEY"),
		Database:  LoadDatabaseConfig(),
		JWT:       LoadJWTConfig(),
		SMTP:      LoadSMTPConfig(),
		Minio:     LoadMinioConfig(),
	}
}
