package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PDNLogEntry — запись доступа к ПДн (статистика в админке, раздел «Персональные данные»).
type PDNLogEntry struct {
	UserID       string
	ShopID       string
	Tool         string
	DataTypes    []string
	Jurisdiction string
	LMProvider   string
	Masked       bool
}

type PDNLogRepository struct {
	db *pgxpool.Pool
}

func NewPDNLogRepository(db *pgxpool.Pool) *PDNLogRepository {
	return &PDNLogRepository{db: db}
}

// Insert is best-effort for the caller: сбой лога не должен ломать вызов инструмента.
func (r *PDNLogRepository) Insert(ctx context.Context, e PDNLogEntry) error {
	var userID, shopID any // пустой UUID не принимается — nil вместо ""
	if e.UserID != "" {
		userID = e.UserID
	}
	if e.ShopID != "" {
		shopID = e.ShopID
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO pdn_access_log (user_id, shop_id, tool, data_types, jurisdiction, llm_provider, masked)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		userID, shopID, e.Tool, e.DataTypes, e.Jurisdiction, e.LMProvider, e.Masked)
	if err != nil {
		return fmt.Errorf("insert pdn_access_log: %w", err)
	}
	return nil
}
