package storage

type UserRepositoryInterface interface {
	Create(nickname, email, phone, firstName, lastName, passwordHash string) (*User, error)
	GetByID(id string) (*User, error)
	GetByLogin(login string) (*User, error)
	GetByEmail(email string) (*User, error)
	ListAll() ([]User, error)
	ListFiltered(page, perPage int, query string) ([]User, int, error)
	UpdateRole(id, role string) error
	UpdateStatus(id string, isActive bool) error
	UpdatePassword(id, passwordHash string) error
	UpdateRequirePasswordChange(id string, value bool) error
	UpdateEmail(id, email string) error
	Delete(id string) error
	Count() (int, error)
	CreateResetToken(userID, token string, expiresAt string) error
	GetByResetToken(token string) (*User, error)
	ClearResetTokens(userID string) error
}

type ShopRepositoryInterface interface {
	Create(userID, name, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault, isReadOnly bool) (*Shop, error)
	GetByID(id string) (*Shop, error)
	ListByUserID(userID string) ([]Shop, error)
	ListByUserIDFiltered(userID string, page, perPage int) ([]Shop, int, error)
	GetDefault(userID string) (*Shop, error)
	Update(id, name, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault, isReadOnly bool) error
	Delete(id string) error
}

type AIKeyRepositoryInterface interface {
	Create(userID, keyHash, name string) (*AIKey, error)
	GetByHash(keyHash string) (*AIKey, error)
	GetByID(id string) (*AIKey, error)
	ListByUserID(userID string) ([]AIKey, error)
	ListByUserIDFiltered(userID string, page, perPage int) ([]AIKey, int, error)
	UpdateLastUsed(id string) error
	UpdateCommands(id, commands string) error
	UpdateCommandsByUserID(userID, commands string) error
	UpdatePDNEnabled(id string, enabled bool) error
	Delete(id string) error
}
