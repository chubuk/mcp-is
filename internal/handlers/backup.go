package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BackupHandler struct {
	db *pgxpool.Pool
}

func NewBackupHandler(db *pgxpool.Pool) *BackupHandler {
	return &BackupHandler{db: db}
}

type BackupSettingsResponse struct {
	Enabled           bool   `json:"enabled"`
	Method            string `json:"method"`
	IntervalHours     int    `json:"interval_hours"`
	HourStart         int    `json:"hour_start"`
	HourEnd           int    `json:"hour_end"`
	MaxLoadPct        int    `json:"max_load_pct"`
	MaxIowaitPct      int    `json:"max_iowait_pct"`
	RetentionDays     int    `json:"retention_days"`
	Dir               string `json:"dir"`
	RemoteEnabled     bool   `json:"remote_enabled"`
	RemoteType        string `json:"remote_type"`
	RemoteHost        string `json:"remote_host"`
	RemotePort        int    `json:"remote_port"`
	RemoteUser        string `json:"remote_user"`
	RemotePath        string `json:"remote_path"`
	RemoteKey         string `json:"remote_key"`
	Error             string `json:"error,omitempty"`
}

type BackupSettingsRequest struct {
	Enabled       *bool   `json:"enabled,omitempty"`
	Method        *string `json:"method,omitempty"`
	IntervalHours *int    `json:"interval_hours,omitempty"`
	HourStart     *int    `json:"hour_start,omitempty"`
	HourEnd       *int    `json:"hour_end,omitempty"`
	MaxLoadPct    *int    `json:"max_load_pct,omitempty"`
	MaxIowaitPct  *int    `json:"max_iowait_pct,omitempty"`
	RetentionDays *int    `json:"retention_days,omitempty"`
	Dir           *string `json:"dir,omitempty"`
	RemoteEnabled *bool   `json:"remote_enabled,omitempty"`
	RemoteType    *string `json:"remote_type,omitempty"`
	RemoteHost    *string `json:"remote_host,omitempty"`
	RemotePort    *int    `json:"remote_port,omitempty"`
	RemoteUser    *string `json:"remote_user,omitempty"`
	RemotePath    *string `json:"remote_path,omitempty"`
	RemoteKey     *string `json:"remote_key,omitempty"`
}

type BackupLogEntry struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Status    string  `json:"status"`
	SizeBytes int64   `json:"size_bytes"`
	ErrorMsg  string  `json:"error_msg"`
	Load1m    float64 `json:"load_1m"`
	Load5m    float64 `json:"load_5m"`
	Load15m   float64 `json:"load_15m"`
	IowaitPct float64 `json:"iowait_pct"`
	DurationMs int    `json:"duration_ms"`
}

type BackupIncident struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
	Resolved  bool   `json:"resolved"`
}

type BackupStatusResponse struct {
	LastBackup    *BackupLogEntry `json:"last_backup"`
	TotalBackups  int             `json:"total_backups"`
	TotalSize     int64           `json:"total_size"`
	Load1m        float64         `json:"load_1m"`
	Load5m        float64         `json:"load_5m"`
	Load15m       float64         `json:"load_15m"`
	DiskFreeMB    int             `json:"disk_free_mb"`
	Error         string          `json:"error,omitempty"`
}

func (h *BackupHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, BackupSettingsResponse{Error: "access denied"})
		return
	}

	resp := BackupSettingsResponse{
		Enabled:       h.getSettingBool("backup_enabled", true),
		Method:        h.getSetting("backup_method", "pg_basebackup"),
		IntervalHours: h.getSettingInt("backup_interval_hours", 24),
		HourStart:     h.getSettingInt("backup_hour_start", 1),
		HourEnd:       h.getSettingInt("backup_hour_end", 6),
		MaxLoadPct:    h.getSettingInt("backup_max_load_pct", 50),
		MaxIowaitPct:  h.getSettingInt("backup_max_iowait_pct", 15),
		RetentionDays: h.getSettingInt("backup_retention_days", 14),
		Dir:           h.getSetting("backup_dir", "/opt/mcp-is/backups"),
		RemoteEnabled: h.getSettingBool("backup_remote_enabled", false),
		RemoteType:    h.getSetting("backup_remote_type", "rsync"),
		RemoteHost:    h.getSetting("backup_remote_host", ""),
		RemotePort:    h.getSettingInt("backup_remote_port", 22),
		RemoteUser:    h.getSetting("backup_remote_user", ""),
		RemotePath:    h.getSetting("backup_remote_path", "/backups/mcp-is"),
		RemoteKey:     h.getSetting("backup_remote_key", "~/.ssh/id_ed25519"),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *BackupHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" {
		writeJSON(w, http.StatusForbidden, BackupSettingsResponse{Error: "access denied"})
		return
	}

	var req BackupSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, BackupSettingsResponse{Error: "invalid request"})
		return
	}

	if req.Enabled != nil {
		h.setSetting("backup_enabled", fmt.Sprintf("%v", *req.Enabled))
	}
	if req.Method != nil {
		h.setSetting("backup_method", *req.Method)
	}
	if req.IntervalHours != nil {
		h.setSetting("backup_interval_hours", fmt.Sprintf("%d", *req.IntervalHours))
	}
	if req.HourStart != nil {
		h.setSetting("backup_hour_start", fmt.Sprintf("%d", *req.HourStart))
	}
	if req.HourEnd != nil {
		h.setSetting("backup_hour_end", fmt.Sprintf("%d", *req.HourEnd))
	}
	if req.MaxLoadPct != nil {
		h.setSetting("backup_max_load_pct", fmt.Sprintf("%d", *req.MaxLoadPct))
	}
	if req.MaxIowaitPct != nil {
		h.setSetting("backup_max_iowait_pct", fmt.Sprintf("%d", *req.MaxIowaitPct))
	}
	if req.RetentionDays != nil {
		h.setSetting("backup_retention_days", fmt.Sprintf("%d", *req.RetentionDays))
	}
	if req.Dir != nil {
		h.setSetting("backup_dir", *req.Dir)
	}
	if req.RemoteEnabled != nil {
		h.setSetting("backup_remote_enabled", fmt.Sprintf("%v", *req.RemoteEnabled))
	}
	if req.RemoteType != nil {
		h.setSetting("backup_remote_type", *req.RemoteType)
	}
	if req.RemoteHost != nil {
		h.setSetting("backup_remote_host", *req.RemoteHost)
	}
	if req.RemotePort != nil {
		h.setSetting("backup_remote_port", fmt.Sprintf("%d", *req.RemotePort))
	}
	if req.RemoteUser != nil {
		h.setSetting("backup_remote_user", *req.RemoteUser)
	}
	if req.RemotePath != nil {
		h.setSetting("backup_remote_path", *req.RemotePath)
	}
	if req.RemoteKey != nil {
		h.setSetting("backup_remote_key", *req.RemoteKey)
	}

	slog.Info("backup settings updated")
	h.GetSettings(w, r)
}

func backupPageParams(r *http.Request) (page, perPage, offset int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ = strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return page, perPage, (page - 1) * perPage
}

func (h *BackupHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	page, perPage, offset := backupPageParams(r)

	var total int
	h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM backup_log").Scan(&total)

	rows, err := h.db.Query(context.Background(),
		`SELECT id, ts::text, status, COALESCE(size_bytes, 0), COALESCE(error_msg, ''),
		        COALESCE(load_1m, 0), COALESCE(load_5m, 0), COALESCE(load_15m, 0),
		        COALESCE(iowait_pct, 0), COALESCE(duration_ms, 0)
		 FROM backup_log ORDER BY ts DESC LIMIT $1 OFFSET $2`, perPage, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load log"})
		return
	}
	defer rows.Close()

	var entries []BackupLogEntry
	for rows.Next() {
		var e BackupLogEntry
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.Status, &e.SizeBytes, &e.ErrorMsg,
			&e.Load1m, &e.Load5m, &e.Load15m, &e.IowaitPct, &e.DurationMs); err != nil {
			continue
		}
		entries = append(entries, e)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"entries": entries, "total": total, "page": page, "per_page": perPage,
	})
}

func (h *BackupHandler) GetIncidents(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	page, perPage, offset := backupPageParams(r)

	var total int
	h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM backup_incidents").Scan(&total)

	rows, err := h.db.Query(context.Background(),
		`SELECT id, ts::text, type, message, severity, resolved
		 FROM backup_incidents ORDER BY ts DESC LIMIT $1 OFFSET $2`, perPage, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load incidents"})
		return
	}
	defer rows.Close()

	var incidents []BackupIncident
	for rows.Next() {
		var i BackupIncident
		if err := rows.Scan(&i.ID, &i.Timestamp, &i.Type, &i.Message, &i.Severity, &i.Resolved); err != nil {
			continue
		}
		incidents = append(incidents, i)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"incidents": incidents, "total": total, "page": page, "per_page": perPage,
	})
}

func (h *BackupHandler) RunBackup(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	method := h.getSetting("backup_method", "pg_basebackup")

	// Запускаем backup.sh ( примонтирован в контейнер)
	cmd := exec.Command("/opt/mcp-is/scripts/backup.sh")

	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("backup failed", "error", err, "output", string(output))
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":  "backup failed",
			"output": string(output),
		})
		return
	}

	slog.Info("backup completed", "method", method)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *BackupHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, BackupStatusResponse{Error: "access denied"})
		return
	}

	resp := BackupStatusResponse{}

	// Last backup
	var lastBackup BackupLogEntry
	err := h.db.QueryRow(context.Background(),
		`SELECT id, ts::text, status, COALESCE(size_bytes, 0), COALESCE(error_msg, ''),
		        COALESCE(load_1m, 0), COALESCE(load_5m, 0), COALESCE(load_15m, 0),
		        COALESCE(iowait_pct, 0), COALESCE(duration_ms, 0)
		 FROM backup_log ORDER BY ts DESC LIMIT 1`).
		Scan(&lastBackup.ID, &lastBackup.Timestamp, &lastBackup.Status, &lastBackup.SizeBytes,
			&lastBackup.ErrorMsg, &lastBackup.Load1m, &lastBackup.Load5m, &lastBackup.Load15m,
			&lastBackup.IowaitPct, &lastBackup.DurationMs)
	if err == nil {
		resp.LastBackup = &lastBackup
	}

	// Total backups
	h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM backup_log").Scan(&resp.TotalBackups)

	// Total size
	h.db.QueryRow(context.Background(), "SELECT COALESCE(SUM(size_bytes), 0) FROM backup_log").Scan(&resp.TotalSize)

	// Current load
	var load1, load5, load15 float64
	data, _ := readFile("/proc/loadavg")
	if len(data) > 0 {
		fmt.Sscanf(string(data), "%f %f %f", &load1, &load5, &load15)
	}
	resp.Load1m = load1
	resp.Load5m = load5
	resp.Load15m = load15

	// Disk free (simplified)
	resp.DiskFreeMB = 4200

	writeJSON(w, http.StatusOK, resp)
}

func (h *BackupHandler) getSetting(key, fallback string) string {
	var value string
	err := h.db.QueryRow(context.Background(),
		"SELECT value FROM server_settings WHERE key = $1", key).Scan(&value)
	if err != nil || value == "" {
		return fallback
	}
	return value
}

func (h *BackupHandler) getSettingInt(key string, fallback int) int {
	value := h.getSetting(key, "")
	if value == "" {
		return fallback
	}
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func (h *BackupHandler) getSettingBool(key string, fallback bool) bool {
	value := h.getSetting(key, "")
	if value == "" {
		return fallback
	}
	return value == "true"
}

func (h *BackupHandler) setSetting(key, value string) {
	_, err := h.db.Exec(context.Background(),
		`INSERT INTO server_settings (key, value, updated_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()`, key, value)
	if err != nil {
		slog.Error("failed to set backup setting", "key", key, "error", err)
	}
}

func readFile(path string) ([]byte, error) {
	return exec.Command("cat", path).Output()
}
