package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PDNHandler struct {
	db            *pgxpool.Pool
	settingsCache *pdn.PDNSettingsCache
}

func NewPDNHandler(db *pgxpool.Pool) *PDNHandler {
	return &PDNHandler{
		db:            db,
		settingsCache: pdn.NewPDNSettingsCache(),
	}
}

type PDNSettingsResponse struct {
	Enabled     bool   `json:"enabled"`
	Log         bool   `json:"log"`
	MaskPII     bool   `json:"mask_pii"`
	AutoProtect bool   `json:"auto_protect"`
	Error       string `json:"error,omitempty"`
}

type PDNSettingsRequest struct {
	Enabled     *bool `json:"enabled,omitempty"`
	Log         *bool `json:"log,omitempty"`
	MaskPII     *bool `json:"mask_pii,omitempty"`
	AutoProtect *bool `json:"auto_protect,omitempty"`
}

type PDNLogEntry struct {
	ID           string   `json:"id"`
	Timestamp    string   `json:"timestamp"`
	Tool         string   `json:"tool"`
	DataTypes    []string `json:"data_types"`
	Jurisdiction string   `json:"jurisdiction"`
	LLMProvider  string   `json:"llm_provider"`
	Masked       bool     `json:"masked"`
}

func (h *PDNHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, PDNSettingsResponse{Error: "access denied"})
		return
	}

	settings := h.settingsCache.Get(h.db)

	writeJSON(w, http.StatusOK, PDNSettingsResponse{
		Enabled:     settings.Enabled,
		Log:         settings.Log,
		MaskPII:     settings.MaskPII,
		AutoProtect: settings.AutoProtect,
	})
}

func (h *PDNHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" {
		writeJSON(w, http.StatusForbidden, PDNSettingsResponse{Error: "only super_admin can change PDN settings"})
		return
	}

	var req PDNSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, PDNSettingsResponse{Error: "invalid request body"})
		return
	}

	if req.Enabled != nil {
		h.settingsCache.Invalidate()
	}
	if req.Log != nil {
		h.settingsCache.Invalidate()
	}
	if req.MaskPII != nil {
		h.settingsCache.Invalidate()
	}
	if req.AutoProtect != nil {
		h.settingsCache.Invalidate()
	}

	h.GetSettings(w, r)
}

func (h *PDNHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rows, err := h.db.Query(context.Background(),
		`SELECT id, ts::text, tool, COALESCE(data_types, ARRAY[]::TEXT[]), COALESCE(jurisdiction,''), COALESCE(llm_provider,''), masked
		 FROM pdn_access_log ORDER BY ts DESC LIMIT $1`, limit,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load log"})
		return
	}
	defer rows.Close()

	var entries []PDNLogEntry
	for rows.Next() {
		var e PDNLogEntry
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.Tool, &e.DataTypes, &e.Jurisdiction, &e.LLMProvider, &e.Masked); err != nil {
			continue
		}
		entries = append(entries, e)
	}

	writeJSON(w, http.StatusOK, map[string]any{"entries": entries, "total": len(entries)})
}

func (h *PDNHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	var totalToday, totalMasked int
	h.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM pdn_access_log WHERE ts >= CURRENT_DATE`).Scan(&totalToday)
	h.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM pdn_access_log WHERE ts >= CURRENT_DATE AND masked = TRUE`).Scan(&totalMasked)

	writeJSON(w, http.StatusOK, map[string]any{
		"total_today":    totalToday,
		"masked_today":   totalMasked,
		"unmasked_today": totalToday - totalMasked,
	})
}

func (h *PDNHandler) GetJurisdictions(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	// Count shops by jurisdiction
	jurisdictions := map[string]int{"RU": 0, "KZ": 0, "BY": 0, "INTL": 0}
	rows, err := h.db.Query(context.Background(),
		`SELECT COALESCE(jurisdiction, 'AUTO') FROM shops`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var j string
			if err := rows.Scan(&j); err == nil {
				jurisdictions[j]++
			}
		}
	}

	// Count auto-protected vs manual today
	var autoProtected, manualProtected int
	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM pdn_access_log WHERE ts >= CURRENT_DATE AND masked = TRUE AND llm_provider != 'RU' AND llm_provider != ''`,
	).Scan(&autoProtected)
	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM pdn_access_log WHERE ts >= CURRENT_DATE AND masked = TRUE AND llm_provider = 'RU'`,
	).Scan(&manualProtected)

	writeJSON(w, http.StatusOK, map[string]any{
		"jurisdictions":    jurisdictions,
		"auto_protected":   autoProtected,
		"manual_protected": manualProtected,
	})
}
