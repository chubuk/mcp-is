package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chubuk/mcp-is/internal/storage"
)

type mockMCPClientTypeRepo struct {
	types []storage.MCPClientType
	err   error
}

func (m *mockMCPClientTypeRepo) ListActive() ([]storage.MCPClientType, error) {
	return m.types, m.err
}

func TestMCPClientTypeHandler_List_Success(t *testing.T) {
	repo := &mockMCPClientTypeRepo{
		types: []storage.MCPClientType{
			{ID: "1", Name: "Claude Desktop", ConfigKey: "mcpServers", ConfigFormat: "json",
				ConfigTemplate: "{}", FileName: "claude_desktop_config.json", FileDescription: "~/Library/Application Support/Claude/", SortOrder: 1, IsActive: true},
			{ID: "2", Name: "Cursor", ConfigKey: "mcpServers", ConfigFormat: "json",
				ConfigTemplate: "{}", FileName: "mcp.json", FileDescription: "~/.cursor/mcp.json", SortOrder: 2, IsActive: true},
		},
	}
	h := NewMCPClientTypeHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/mcp-client-types", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result []storage.MCPClientType
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 2 {
		t.Errorf("expected 2 types, got %d", len(result))
	}
	if result[0].Name != "Claude Desktop" {
		t.Errorf("expected 'Claude Desktop', got '%s'", result[0].Name)
	}
}

func TestMCPClientTypeHandler_List_Empty(t *testing.T) {
	repo := &mockMCPClientTypeRepo{types: nil}
	h := NewMCPClientTypeHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/mcp-client-types", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result []storage.MCPClientType
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 0 {
		t.Errorf("expected 0 types, got %d", len(result))
	}
}

func TestMCPClientTypeHandler_List_Fields(t *testing.T) {
	repo := &mockMCPClientTypeRepo{
		types: []storage.MCPClientType{
			{
				ID: "1", Name: "VS Code (Copilot)", ConfigKey: "servers", ConfigFormat: "json",
				ConfigTemplate: `{"servers":{"mcp-is":{"type":"http"}}}`,
				FileName: "mcp.json", FileDescription: ".vscode/mcp.json", SortOrder: 3, IsActive: true,
			},
		},
	}
	h := NewMCPClientTypeHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/mcp-client-types", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	var result []storage.MCPClientType
	json.NewDecoder(w.Body).Decode(&result)

	if len(result) != 1 {
		t.Fatalf("expected 1 type, got %d", len(result))
	}
	item := result[0]
	if item.ConfigKey != "servers" {
		t.Errorf("expected config_key 'servers', got '%s'", item.ConfigKey)
	}
	if item.FileName != "mcp.json" {
		t.Errorf("expected file_name 'mcp.json', got '%s'", item.FileName)
	}
	if item.FileDescription != ".vscode/mcp.json" {
		t.Errorf("expected file_description '.vscode/mcp.json', got '%s'", item.FileDescription)
	}
}

func TestMCPClientTypeHandler_List_Error(t *testing.T) {
	repo := &mockMCPClientTypeRepo{err: &mockError{"db error"}}
	h := NewMCPClientTypeHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/mcp-client-types", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }
