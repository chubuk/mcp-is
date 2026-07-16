package pdn

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func LogAccess(db *pgxpool.Pool, userID, shopID, tool string, dataTypes []string, jurisdiction, llmProvider string, masked bool) {
	go func() {
		_, err := db.Exec(context.Background(),
			`INSERT INTO pdn_access_log (user_id, shop_id, tool, data_types, jurisdiction, llm_provider, masked)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			userID, shopID, tool, dataTypes, jurisdiction, llmProvider, masked,
		)
		if err != nil {
			slog.Error("failed to log PDN access", "error", err)
		}
	}()
}
