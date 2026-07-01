//go:build integration

package storage

import (
	"fmt"
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

func TestUserRepo_ListFiltered_NoQuery(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("alpha", "a@test.com", "", "", "", "h")
	repo.Create("beta", "b@test.com", "", "", "", "h")
	repo.Create("gamma", "c@test.com", "", "", "", "h")

	users, total, err := repo.ListFiltered(1, 10, "")
	if err != nil {
		t.Fatalf("ListFiltered failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestUserRepo_ListFiltered_WithQuery(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("ivan", "ivan@test.com", "+79991112233", "Ivan", "Petrov", "h")
	repo.Create("petr", "petr@test.com", "+79994445566", "Petr", "Ivanov", "h")
	repo.Create("sidor", "sidor@test.com", "", "Sidor", "Sidorov", "h")

	users, total, err := repo.ListFiltered(1, 10, "ivan")
	if err != nil {
		t.Fatalf("ListFiltered failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2 (ivan + petr with Ivanov), got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserRepo_ListFiltered_Pagination(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	for i := 0; i < 5; i++ {
		repo.Create(fmt.Sprintf("page%d", i), fmt.Sprintf("p%d@test.com", i), "", "", "", "h")
	}

	page1, total, _ := repo.ListFiltered(1, 2, "")
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(page1) != 2 {
		t.Errorf("expected 2 users on page 1, got %d", len(page1))
	}

	page2, _, _ := repo.ListFiltered(2, 2, "")
	if len(page2) != 2 {
		t.Errorf("expected 2 users on page 2, got %d", len(page2))
	}

	page3, _, _ := repo.ListFiltered(3, 2, "")
	if len(page3) != 1 {
		t.Errorf("expected 1 user on page 3, got %d", len(page3))
	}
}

func TestUserRepo_ListFiltered_SearchByPhone(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("user1", "", "+79991234567", "", "", "h")
	repo.Create("user2", "", "+79998765432", "", "", "h")

	users, total, _ := repo.ListFiltered(1, 10, "123")
	if total != 1 {
		t.Errorf("expected 1 user by phone search, got %d", total)
	}
	if len(users) == 1 && users[0].Phone != "+79991234567" {
		t.Errorf("expected phone '+79991234567', got '%s'", users[0].Phone)
	}
}

func TestUserRepo_Count(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("cnt1", "", "", "", "", "h")
	repo.Create("cnt2", "", "", "", "", "h")

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 users, got %d", count)
	}
}

func TestUserRepo_UpdateRole(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	user, _ := repo.Create("roleuser", "", "", "", "", "h")

	err := repo.UpdateRole(user.ID, "admin")
	if err != nil {
		t.Fatalf("UpdateRole failed: %v", err)
	}

	found, _ := repo.GetByID(user.ID)
	if found.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", found.Role)
	}
}

func TestUserRepo_UpdateStatus(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	user, _ := repo.Create("statususer", "", "", "", "", "h")

	err := repo.UpdateStatus(user.ID, false)
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	found, _ := repo.GetByID(user.ID)
	if found.IsActive {
		t.Error("expected is_active=false")
	}
}

func TestUserRepo_UpdatePassword(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	user, _ := repo.Create("pwduser", "", "", "", "", "old_hash")

	err := repo.UpdatePassword(user.ID, "new_hash")
	if err != nil {
		t.Fatalf("UpdatePassword failed: %v", err)
	}

	found, _ := repo.GetByID(user.ID)
	if found.PasswordHash != "new_hash" {
		t.Errorf("expected password_hash 'new_hash', got '%s'", found.PasswordHash)
	}
}

func TestUserRepo_UpdateEmail(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	user, _ := repo.Create("emailupd", "", "", "", "", "h")

	err := repo.UpdateEmail(user.ID, "newemail@test.com")
	if err != nil {
		t.Fatalf("UpdateEmail failed: %v", err)
	}

	found, _ := repo.GetByID(user.ID)
	if found.Email != "newemail@test.com" {
		t.Errorf("expected email 'newemail@test.com', got '%s'", found.Email)
	}
}

func TestUserRepo_PasswordResetFlow(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	user, _ := repo.Create("resetuser", "reset@test.com", "", "", "", "h")

	err := repo.CreateResetToken(user.ID, "token123", "2030-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("CreateResetToken failed: %v", err)
	}

	found, err := repo.GetByResetToken("token123")
	if err != nil {
		t.Fatalf("GetByResetToken failed: %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("expected user ID '%s', got '%s'", user.ID, found.ID)
	}

	err = repo.ClearResetTokens(user.ID)
	if err != nil {
		t.Fatalf("ClearResetTokens failed: %v", err)
	}

	_, err = repo.GetByResetToken("token123")
	if err == nil {
		t.Error("expected error after clearing tokens")
	}
}

func TestUserRepo_GetByEmail(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	repo.Create("emailget", "findme@test.com", "", "", "", "h")

	user, err := repo.GetByEmail("findme@test.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if user.Nickname != "emailget" {
		t.Errorf("expected nickname 'emailget', got '%s'", user.Nickname)
	}
}
