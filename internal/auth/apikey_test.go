package auth

import (
	"testing"
)

func TestGenerateAPIKey(t *testing.T) {
	raw, hash, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey failed: %v", err)
	}

	if len(raw) != 64 {
		t.Errorf("expected raw key length 64, got %d", len(raw))
	}

	if len(hash) != 64 {
		t.Errorf("expected hash length 64 (SHA-256 hex), got %d", len(hash))
	}
}

func TestHashAPIKey_Deterministic(t *testing.T) {
	key := "test-api-key-12345"

	hash1 := HashAPIKey(key)
	hash2 := HashAPIKey(key)

	if hash1 != hash2 {
		t.Error("same input should produce same hash")
	}
}

func TestHashAPIKey_DifferentInputs(t *testing.T) {
	hash1 := HashAPIKey("key-one")
	hash2 := HashAPIKey("key-two")

	if hash1 == hash2 {
		t.Error("different inputs should produce different hashes")
	}
}

func TestGenerateAPIKey_Unique(t *testing.T) {
	_, hash1, _ := GenerateAPIKey()
	_, hash2, _ := GenerateAPIKey()

	if hash1 == hash2 {
		t.Error("two generated keys should have different hashes")
	}
}
