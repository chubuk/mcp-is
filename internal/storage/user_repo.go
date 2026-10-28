package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID                    string
	Nickname              string
	Email                 string
	Phone                 string
	FirstName             string
	LastName              string
	PasswordHash          string
	Role                  string
	IsActive              bool
	RequirePasswordChange bool
	EmailVerified         bool
	CreatedAt             string
	UpdatedAt             string
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(nickname, email, phone, firstName, lastName, passwordHash string) (*User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO users (nickname, email, phone, first_name, last_name, password_hash)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, nickname, COALESCE(email,''), COALESCE(phone,''),
		           COALESCE(first_name,''), COALESCE(last_name,''), role, is_active, require_password_change, created_at::text`,
		nickname, nullStr(normalizeEmail(email)), nullStr(phone), nullStr(firstName), nullStr(lastName), passwordHash,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Phone, &u.FirstName, &u.LastName, &u.Role, &u.IsActive, &u.RequirePasswordChange, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByID(id string) (*User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, nickname, COALESCE(email,''), COALESCE(phone,''),
		        COALESCE(first_name,''), COALESCE(last_name,''), password_hash, role, is_active, require_password_change, email_verified, created_at::text, updated_at::text
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Phone, &u.FirstName, &u.LastName, &u.PasswordHash, &u.Role, &u.IsActive, &u.RequirePasswordChange, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByLogin(login string) (*User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, nickname, COALESCE(email,''), COALESCE(phone,''),
		        COALESCE(first_name,''), COALESCE(last_name,''), password_hash, role, is_active, require_password_change, email_verified, created_at::text, updated_at::text
		 FROM users WHERE (LOWER(email) = LOWER($1) OR phone = $1 OR nickname = $1) AND is_active = TRUE`, login,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Phone, &u.FirstName, &u.LastName, &u.PasswordHash, &u.Role, &u.IsActive, &u.RequirePasswordChange, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) ListAll() ([]User, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, nickname, COALESCE(email,''), COALESCE(phone,''),
		        COALESCE(first_name,''), COALESCE(last_name,''), role, is_active, require_password_change, created_at::text
		 FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.Phone, &u.FirstName, &u.LastName, &u.Role, &u.IsActive, &u.RequirePasswordChange, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) ListFiltered(page, perPage int, query string) ([]User, int, error) {
	offset := (page - 1) * perPage
	if offset < 0 {
		offset = 0
	}

	where := ""
	args := []any{}
	if query != "" {
		where = `WHERE nickname ILIKE '%' || $1 || '%' 
		          OR email ILIKE '%' || $1 || '%' 
		          OR phone ILIKE '%' || $1 || '%' 
		          OR first_name ILIKE '%' || $1 || '%' 
		          OR last_name ILIKE '%' || $1 || '%'`
		args = append(args, query)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM users ` + where
	if err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	dataArgs := make([]any, len(args))
	copy(dataArgs, args)
	dataArgs = append(dataArgs, perPage, offset)
	dataQuery := `SELECT id, nickname, COALESCE(email,''), COALESCE(phone,''),
	       COALESCE(first_name,''), COALESCE(last_name,''), role, is_active, require_password_change, created_at::text
		FROM users ` + where + ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.db.Query(context.Background(), dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users filtered: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.Phone, &u.FirstName, &u.LastName, &u.Role, &u.IsActive, &u.RequirePasswordChange, &u.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) UpdateRole(id, role string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`, role, id)
	return err
}

func (r *UserRepository) UpdateStatus(id string, isActive bool) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET is_active=$1, updated_at=NOW() WHERE id=$2`, isActive, id)
	return err
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (r *UserRepository) UpdatePassword(id, passwordHash string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`, passwordHash, id)
	return err
}

func (r *UserRepository) UpdateRequirePasswordChange(id string, value bool) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET require_password_change=$1, updated_at=NOW() WHERE id=$2`, value, id)
	return err
}

func (r *UserRepository) UpdateEmailVerified(id string, verified bool) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET email_verified=$1, updated_at=NOW() WHERE id=$2`, verified, id)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, nickname, COALESCE(email,''), password_hash, role, email_verified
		 FROM users WHERE LOWER(email) = LOWER($1)`, normalizeEmail(email),
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.PasswordHash, &u.Role, &u.EmailVerified)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

// PurgeUnverified deletes never-confirmed accounts older than olderThanSeconds
// so open registration cannot bloat the users table.
func (r *UserRepository) PurgeUnverified(olderThanSeconds int64) (int64, error) {
	tag, err := r.db.Exec(context.Background(),
		`DELETE FROM users WHERE email_verified = FALSE AND created_at < NOW() - ($1 * INTERVAL '1 second')`,
		olderThanSeconds)
	if err != nil {
		return 0, fmt.Errorf("purge unverified users: %w", err)
	}
	return tag.RowsAffected(), nil
}

func (r *UserRepository) CreateResetToken(userID, token string, expiresAt string) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO password_resets (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt)
	return err
}

func (r *UserRepository) GetByResetToken(token string) (*User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`SELECT u.id, u.nickname, COALESCE(u.email,''), u.role
		 FROM users u
		 JOIN password_resets pr ON pr.user_id = u.id
		 WHERE pr.token = $1 AND pr.expires_at > NOW()`, token,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) ClearResetTokens(userID string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM password_resets WHERE user_id = $1`, userID)
	return err
}

func (r *UserRepository) CleanExpiredResetTokens() error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM password_resets WHERE expires_at < NOW()`)
	return err
}

func (r *UserRepository) CreateSuperAdmin(nickname, email, phone, firstName, lastName, passwordHash string) (*User, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	var u User
	err = tx.QueryRow(context.Background(),
		`INSERT INTO users (nickname, email, phone, first_name, last_name, password_hash)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, nickname, COALESCE(email,''), COALESCE(phone,''),
		           COALESCE(first_name,''), COALESCE(last_name,''), role, is_active, require_password_change, created_at::text`,
		nickname, nullStr(normalizeEmail(email)), nullStr(phone), nullStr(firstName), nullStr(lastName), passwordHash,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Phone, &u.FirstName, &u.LastName, &u.Role, &u.IsActive, &u.RequirePasswordChange, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users SET role='super_admin', require_password_change=true, updated_at=NOW() WHERE id=$1`, u.ID)
	if err != nil {
		return nil, fmt.Errorf("set super_admin role: %w", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	u.Role = "super_admin"
	u.RequirePasswordChange = true
	return &u, nil
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// normalizeEmail — единная нормализация e-mail для всех flows (OWASP:
// registration/login/reset/linking обязаны использовать одну нормализацию).
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
