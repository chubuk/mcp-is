package auth

import (
	"testing"
)

func TestAESCipher_EncryptDecrypt(t *testing.T) {
	cipher, err := NewAESCipher("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("NewAESCipher failed: %v", err)
	}

	plaintext := "my-secret-password-123"

	encrypted, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == plaintext {
		t.Error("encrypted text should differ from plaintext")
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestAESCipher_WrongKeyLength(t *testing.T) {
	_, err := NewAESCipher("short-key")
	if err == nil {
		t.Error("expected error for wrong key length")
	}
}

func TestAESCipher_DecryptInvalid(t *testing.T) {
	cipher, _ := NewAESCipher("12345678901234567890123456789012")

	_, err := cipher.Decrypt("not-valid-base64-or-ciphertext")
	if err == nil {
		t.Error("expected error for invalid ciphertext")
	}
}

func TestAESCipher_DifferentEncryptions(t *testing.T) {
	cipher, _ := NewAESCipher("12345678901234567890123456789012")

	enc1, _ := cipher.Encrypt("same-text")
	enc2, _ := cipher.Encrypt("same-text")

	if enc1 == enc2 {
		t.Error("same plaintext should produce different ciphertext (random nonce)")
	}
}
