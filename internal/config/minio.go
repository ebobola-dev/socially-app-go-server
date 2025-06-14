package config

import "os"

type MinioConfig struct {
	USER     string
	PASSWORD string
}

func LoadMinioConfig() *MinioConfig {
	return &MinioConfig{
		USER:     os.Getenv("MINIO_ROOT_USER"),
		PASSWORD: os.Getenv("MINIO_ROOT_PASSWORD"),
	}
}
