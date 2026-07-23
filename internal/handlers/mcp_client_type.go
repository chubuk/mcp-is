package handlers

import (
	"net/http"

	"github.com/chubuk/mcp-is/internal/storage"
)

type MCPClientTypeLister interface {
	ListActive() ([]storage.MCPClientType, error)
}

type MCPClientTypeHandler struct {
	repo MCPClientTypeLister
}

func NewMCPClientTypeHandler(repo MCPClientTypeLister) *MCPClientTypeHandler {
	return &MCPClientTypeHandler{repo: repo}
}

func (h *MCPClientTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	types, err := h.repo.ListActive()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list MCP client types"})
		return
	}
	if types == nil {
		types = []storage.MCPClientType{}
	}
	writeJSON(w, http.StatusOK, types)
}
