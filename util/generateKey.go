package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 24) // 24 random bytes → ~32 chars after base64
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// URL-safe base64 (no "+" or "/" or "=")
	key := base64.RawURLEncoding.EncodeToString(bytes)

	// Add brand + type prefix
	return "cgnz_sk_" + key, nil
}
