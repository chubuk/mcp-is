package handlers

import (
	"fmt"
	"sync"

	"github.com/your-project/mcp-insales/internal/storage"
)

type mockUserRepo struct {
	mu      sync.RWMutex
	users   map[string]*storage.User
	counter int
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*storage.User)}
}

func (m *mockUserRepo) Create(nickname, email, phone, firstName, lastName, passwordHash string) (*storage.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.Nickname == nickname || (email != "" && u.Email == email) || (phone != "" && u.Phone == phone) {
			return nil, fmt.Errorf("duplicate")
		}
	}

	m.counter++
	u := &storage.User{
		ID:           fmt.Sprintf("user-%d", m.counter),
		Nickname:     nickname,
		Email:        email,
		Phone:        phone,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		CreatedAt:    "2026-01-01T00:00:00Z",
	}
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserRepo) GetByID(id string) (*storage.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockUserRepo) GetByLogin(login string) (*storage.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Nickname == login || u.Email == login || u.Phone == login {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockUserRepo) GetByEmail(email string) (*storage.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockUserRepo) ListAll() ([]storage.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var users []storage.User
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, nil
}

func (m *mockUserRepo) UpdateRole(id, role string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, ok := m.users[id]; ok {
		u.Role = role
		return nil
	}
	return fmt.Errorf("not found")
}

func (m *mockUserRepo) UpdateStatus(id string, isActive bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, ok := m.users[id]; ok {
		u.IsActive = isActive
		return nil
	}
	return fmt.Errorf("not found")
}

func (m *mockUserRepo) UpdatePassword(id, passwordHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, ok := m.users[id]; ok {
		u.PasswordHash = passwordHash
		return nil
	}
	return fmt.Errorf("not found")
}

func (m *mockUserRepo) UpdateEmail(id, email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, ok := m.users[id]; ok {
		u.Email = email
		return nil
	}
	return fmt.Errorf("not found")
}

func (m *mockUserRepo) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) Count() (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.users), nil
}

func (m *mockUserRepo) CreateResetToken(userID, token string, expiresAt string) error {
	return nil
}

func (m *mockUserRepo) GetByResetToken(token string) (*storage.User, error) {
	return nil, fmt.Errorf("not found")
}

func (m *mockUserRepo) ClearResetTokens(userID string) error {
	return nil
}

type mockShopRepo struct {
	mu      sync.RWMutex
	shops   map[string]*storage.Shop
	counter int
}

func newMockShopRepo() *mockShopRepo {
	return &mockShopRepo{shops: make(map[string]*storage.Shop)}
}

func (m *mockShopRepo) Create(userID, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault bool) (*storage.Shop, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	s := &storage.Shop{
		ID:                   fmt.Sprintf("shop-%d", m.counter),
		UserID:               userID,
		ShopURL:              shopURL,
		APIKeyID:             apiKeyID,
		APIPasswordEncrypted: apiPasswordEncrypted,
		IsDefault:            isDefault,
		CreatedAt:            "2026-01-01T00:00:00Z",
	}
	m.shops[s.ID] = s
	return s, nil
}

func (m *mockShopRepo) GetByID(id string) (*storage.Shop, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if s, ok := m.shops[id]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockShopRepo) ListByUserID(userID string) ([]storage.Shop, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []storage.Shop
	for _, s := range m.shops {
		if s.UserID == userID {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (m *mockShopRepo) GetDefault(userID string) (*storage.Shop, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, s := range m.shops {
		if s.UserID == userID && s.IsDefault {
			return s, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockShopRepo) Update(id, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.shops[id]; ok {
		s.ShopURL = shopURL
		s.APIKeyID = apiKeyID
		s.APIPasswordEncrypted = apiPasswordEncrypted
		s.IsDefault = isDefault
		return nil
	}
	return fmt.Errorf("not found")
}

func (m *mockShopRepo) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.shops, id)
	return nil
}

type mockAIKeyRepo struct {
	mu      sync.RWMutex
	keys    map[string]*storage.AIKey
	counter int
}

func newMockAIKeyRepo() *mockAIKeyRepo {
	return &mockAIKeyRepo{keys: make(map[string]*storage.AIKey)}
}

func (m *mockAIKeyRepo) Create(userID, keyHash, name string) (*storage.AIKey, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	k := &storage.AIKey{
		ID:        fmt.Sprintf("key-%d", m.counter),
		UserID:    userID,
		KeyHash:   keyHash,
		Name:      name,
		CreatedAt: "2026-01-01T00:00:00Z",
	}
	m.keys[k.ID] = k
	return k, nil
}

func (m *mockAIKeyRepo) GetByHash(keyHash string) (*storage.AIKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, k := range m.keys {
		if k.KeyHash == keyHash {
			return k, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockAIKeyRepo) ListByUserID(userID string) ([]storage.AIKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []storage.AIKey
	for _, k := range m.keys {
		if k.UserID == userID {
			result = append(result, *k)
		}
	}
	return result, nil
}

func (m *mockAIKeyRepo) UpdateLastUsed(id string) error {
	return nil
}

func (m *mockAIKeyRepo) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.keys, id)
	return nil
}
