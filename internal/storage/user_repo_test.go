//go:build integration

package storage

import (
	"testing"
)

func TestUserRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	user, err := repo.Create("testuser", "test@example.com", "+79991234567", "Test", "User", "hash123")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == "" {
		t.Error("expected non-empty ID")
	}
	if user.Nickname != "testuser" {
		t.Errorf("expected nickname 'testuser', got '%s'", user.Nickname)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUserRepo_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	created, _ := repo.Create("getbyid", "getbyid@test.com", "", "", "", "hash")

	found, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.Nickname != "getbyid" {
		t.Errorf("expected nickname 'getbyid', got '%s'", found.Nickname)
	}
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	_, err := repo.GetByID("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestUserRepo_GetByLogin_Nickname(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("loginuser", "", "", "", "", "hash")

	user, err := repo.GetByLogin("loginuser")
	if err != nil {
		t.Fatalf("GetByLogin failed: %v", err)
	}
	if user.Nickname != "loginuser" {
		t.Errorf("expected 'loginuser', got '%s'", user.Nickname)
	}
}

func TestUserRepo_GetByLogin_Email(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("emailuser", "email@test.com", "", "", "", "hash")

	user, err := repo.GetByLogin("email@test.com")
	if err != nil {
		t.Fatalf("GetByLogin failed: %v", err)
	}
	if user.Email != "email@test.com" {
		t.Errorf("expected 'email@test.com', got '%s'", user.Email)
	}
}

func TestUserRepo_GetByLogin_Phone(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("phoneuser", "", "+79990001122", "", "", "hash")

	user, err := repo.GetByLogin("+79990001122")
	if err != nil {
		t.Fatalf("GetByLogin failed: %v", err)
	}
	if user.Phone != "+79990001122" {
		t.Errorf("expected '+79990001122', got '%s'", user.Phone)
	}
}

func TestUserRepo_Create_DuplicateNickname(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("dup", "", "", "", "", "hash")

	_, err := repo.Create("dup", "", "", "", "", "hash")
	if err == nil {
		t.Error("expected error for duplicate nickname")
	}
}
