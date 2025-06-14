package env

import (
	"log"
	"os"
	"strconv"
)

func GetInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Unable to read env '%s'", key)
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Unable to convert env to int '%s' -> %s,  %v", key, val, err)
	}
	return i
}

func GetString(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Unable to read env '%s'", key)
	}
	return val
}
