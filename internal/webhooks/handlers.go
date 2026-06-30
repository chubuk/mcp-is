package webhooks

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	event := r.PathValue("event")
	slog.Info("webhook received", "event", event)

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("webhook decode error", "error", err)
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	slog.Info("webhook payload", "event", event, "keys", len(payload))
	w.WriteHeader(http.StatusOK)
}
