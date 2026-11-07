package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrIdempotencyInProgress — ключ уже исполняется параллельным запросом.
var ErrIdempotencyInProgress = errors.New("idempotency request in progress")

// ErrIdempotencyKeyReuse — ключ переиспользован с другими данными.
var ErrIdempotencyKeyReuse = errors.New("idempotency key reused with different payload")

// IdempotencyRecord — сохранённый результат create-операции.
type IdempotencyRecord struct {
	Tool        string
	RequestHash string
	Status      string // pending | done
	Response    []byte // JSON: {"text": "...", "structured": ...}
}

// IdempotencyRepository — дедупликация create-вызовов MCP по клиентскому ключу.
type IdempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

// Begin резервирует ключ (insert-before-execute, закрывает race параллельных retry).
// created=true — строка новая, нужно выполнить операцию и вызвать Finish/Cancel.
// created=false — ключ уже есть: запись возвращается как есть (replay при status=done).
func (r *IdempotencyRepository) Begin(ctx context.Context, aiKeyID, key, tool, requestHash string) (rec *IdempotencyRecord, created bool, err error) {
	if r == nil || r.db == nil {
		return nil, false, errors.New("idempotency repo not initialized")
	}
	// Ленивая очистка TTL 24ч — без отдельного job'а.
	if _, derr := r.db.Exec(ctx, `DELETE FROM idempotency_keys WHERE created_at < NOW() - INTERVAL '24 hours'`); derr != nil {
		// best-effort: чистка не должна блокировать вызов
		_ = derr
	}
	tag, err := r.db.Exec(ctx,
		`INSERT INTO idempotency_keys (ai_key_id, key, tool, request_hash, status)
		 VALUES ($1, $2, $3, $4, 'pending')
		 ON CONFLICT (ai_key_id, key) DO NOTHING`,
		aiKeyID, key, tool, requestHash,
	)
	if err != nil {
		return nil, false, fmt.Errorf("idempotency begin: %w", err)
	}
	if tag.RowsAffected() == 1 {
		return &IdempotencyRecord{Tool: tool, RequestHash: requestHash, Status: "pending"}, true, nil
	}
	var got IdempotencyRecord
	var resp []byte
	err = r.db.QueryRow(ctx,
		`SELECT tool, request_hash, status, response FROM idempotency_keys WHERE ai_key_id=$1 AND key=$2`,
		aiKeyID, key,
	).Scan(&got.Tool, &got.RequestHash, &got.Status, &resp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// строка исчезла между INSERT и SELECT (TTL-чистка) — считаем ключ свободным
			return &IdempotencyRecord{Tool: tool, RequestHash: requestHash, Status: "pending"}, true, nil
		}
		return nil, false, fmt.Errorf("idempotency lookup: %w", err)
	}
	got.Response = resp
	return &got, false, nil
}

// Finish сохраняет успешный ответ и переводит ключ в status=done.
func (r *IdempotencyRepository) Finish(ctx context.Context, aiKeyID, key string, response []byte) error {
	if r == nil || r.db == nil {
		return errors.New("idempotency repo not initialized")
	}
	_, err := r.db.Exec(ctx,
		`UPDATE idempotency_keys SET status='done', response=$3 WHERE ai_key_id=$1 AND key=$2`,
		aiKeyID, key, response,
	)
	if err != nil {
		return fmt.Errorf("idempotency finish: %w", err)
	}
	return nil
}

// Cancel удаляет зарезервированный ключ после ошибки — ретрай разрешён.
func (r *IdempotencyRepository) Cancel(ctx context.Context, aiKeyID, key string) error {
	if r == nil || r.db == nil {
		return errors.New("idempotency repo not initialized")
	}
	_, err := r.db.Exec(ctx,
		`DELETE FROM idempotency_keys WHERE ai_key_id=$1 AND key=$2 AND status='pending'`,
		aiKeyID, key,
	)
	if err != nil {
		return fmt.Errorf("idempotency cancel: %w", err)
	}
	return nil
}
