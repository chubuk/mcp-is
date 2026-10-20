package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditLogEntry records one MCP write-tool call. Args are stored as hash only (PDN).
type AuditLogEntry struct {
	RequestID string
	UserID    string
	APIKeyID  string
	ShopID    string
	ToolName  string
	Action    string // call | denied | error
	ArgsHash  string
	ResultOK  bool
}

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

// Insert is best-effort for the caller: audit failures must not break the tool call.
func (r *AuditRepository) Insert(ctx context.Context, e AuditLogEntry) error {
	if r == nil || r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO audit_log (request_id, user_id, api_key_id, shop_id, tool_name, action, args_hash, result_ok)
		 VALUES (NULLIF($1,''), NULLIF($2,''), NULLIF($3,''), NULLIF($4,''), $5, $6, NULLIF($7,''), $8)`,
		e.RequestID, e.UserID, e.APIKeyID, e.ShopID, e.ToolName, e.Action, e.ArgsHash, e.ResultOK,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

// ListRecent returns newest entries first (admin UI).
func (r *AuditRepository) ListRecent(ctx context.Context, limit int) ([]AuditLogEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(ctx,
		`SELECT COALESCE(request_id,''), COALESCE(user_id,''), COALESCE(api_key_id,''), COALESCE(shop_id,''),
		        tool_name, action, COALESCE(args_hash,''), COALESCE(result_ok,false)
		 FROM audit_log ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("list audit log: %w", err)
	}
	defer rows.Close()
	var out []AuditLogEntry
	for rows.Next() {
		var e AuditLogEntry
		if err := rows.Scan(&e.RequestID, &e.UserID, &e.APIKeyID, &e.ShopID,
			&e.ToolName, &e.Action, &e.ArgsHash, &e.ResultOK); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		out = append(out, e)
	}
	return out, nil
}
