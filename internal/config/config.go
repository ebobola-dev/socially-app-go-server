package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	BuildType BuildType
} //

func Initialize() *Config {
	_ = godotenv.Load()
	btStr := _getEnv("BUILD_TYPE", "DEV")
	buildType, err := ParseBuildType(btStr)

	if err != nil {
		log.Printf("Invalid BUILD_TYPE '%s', falling back to DEV\n", btStr)
		buildType = Development
	}

	return &Config{
		Port:      _getEnv("INTERNAL_PORT", "8080"),
		BuildType: buildType,
	}
}

func _getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
