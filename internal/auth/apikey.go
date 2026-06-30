package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func GenerateAPIKey() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw := hex.EncodeToString(b)
	hash := HashAPIKey(raw)
	return raw, hash, nil
}
