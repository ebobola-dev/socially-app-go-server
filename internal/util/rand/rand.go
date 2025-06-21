package rand_util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateHEX(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even, got %d", length)
	}
	byteLen := length / 2
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
