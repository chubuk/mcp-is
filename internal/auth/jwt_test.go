package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	mgr := NewJWTManager("test-secret-key-32-bytes-long!!", 0)

	token, err := mgr.GenerateToken("user-123", "testuser", "user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := mgr.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got '%s'", claims.UserID)
	}

	if claims.Nickname != "testuser" {
		t.Errorf("expected Nickname 'testuser', got '%s'", claims.Nickname)
	}

	if claims.ExpiresAt == nil {
		t.Error("expected ExpiresAt to be set")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("token already expired")
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	mgr := NewJWTManager("test-secret-key", 0)

	_, err := mgr.ValidateToken("invalid-token-string")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWTManager_WrongSecret(t *testing.T) {
	mgr1 := NewJWTManager("secret-one", 0)
	mgr2 := NewJWTManager("secret-two", 0)

	token, _ := mgr1.GenerateToken("user-1", "user1", "user")

	_, err := mgr2.ValidateToken(token)
	if err == nil {
		t.Error("expected error when validating with wrong secret")
	}
}
