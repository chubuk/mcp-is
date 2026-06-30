package storage

type UserRepositoryInterface interface {
	Create(nickname, email, phone, firstName, lastName, passwordHash string) (*User, error)
	GetByID(id string) (*User, error)
	GetByLogin(login string) (*User, error)
	GetByEmail(email string) (*User, error)
	ListAll() ([]User, error)
	UpdateRole(id, role string) error
	UpdateStatus(id string, isActive bool) error
	UpdatePassword(id, passwordHash string) error
	UpdateEmail(id, email string) error
	Delete(id string) error
	Count() (int, error)
	CreateResetToken(userID, token string, expiresAt string) error
	GetByResetToken(token string) (*User, error)
	ClearResetTokens(userID string) error
}

type ShopRepositoryInterface interface {
	Create(userID, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault bool) (*Shop, error)
	GetByID(id string) (*Shop, error)
	ListByUserID(userID string) ([]Shop, error)
	GetDefault(userID string) (*Shop, error)
	Update(id, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault bool) error
	Delete(id string) error
}

type AIKeyRepositoryInterface interface {
	Create(userID, keyHash, name string) (*AIKey, error)
	GetByHash(keyHash string) (*AIKey, error)
	ListByUserID(userID string) ([]AIKey, error)
	UpdateLastUsed(id string) error
	Delete(id string) error
}
