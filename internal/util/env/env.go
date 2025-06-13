package env

import (
	"log"
	"os"
	"strconv"
)

func GetIntEnv(key string) int {
	val := os.Getenv(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid integer in %s: %v", key, err)
	}
	return i
}
