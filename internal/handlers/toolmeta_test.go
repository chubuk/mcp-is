package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chubuk/mcp-is/internal/mcp"
)

func TestToolMeta_ReturnsAllToolsWithGroupAndClass(t *testing.T) {
	rec := httptest.NewRecorder()
	ToolMeta(rec, httptest.NewRequest(http.MethodGet, "/api/tool-meta", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		Tools []struct {
			Name   string `json:"name"`
			Group  string `json:"group"`
			Class  string `json:"class"`
			Domain string `json:"domain"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	want := len(mcp.AllToolNames())
	if len(resp.Tools) != want {
		t.Fatalf("len(tools) = %d, want %d", len(resp.Tools), want)
	}
	seen := make(map[string]bool, len(resp.Tools))
	for _, item := range resp.Tools {
		if item.Name == "" || item.Group == "" || item.Class == "" {
			t.Fatalf("empty field in %+v", item)
		}
		if seen[item.Name] {
			t.Fatalf("duplicate tool %q", item.Name)
		}
		seen[item.Name] = true
		if m := mcp.ToolMetaFor(item.Name); m.Group != item.Group || m.Class != item.Class || m.Domain != item.Domain {
			t.Fatalf("tool %q: got (%s,%s,%s), want (%s,%s,%s)", item.Name, item.Group, item.Class, item.Domain, m.Group, m.Class, m.Domain)
		}
	}
}

func TestToolMeta_SortedByName(t *testing.T) {
	rec := httptest.NewRecorder()
	ToolMeta(rec, httptest.NewRequest(http.MethodGet, "/api/tool-meta", nil))
	var resp struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for i := 1; i < len(resp.Tools); i++ {
		if resp.Tools[i-1].Name > resp.Tools[i].Name {
			t.Fatalf("tools not sorted: %q > %q", resp.Tools[i-1].Name, resp.Tools[i].Name)
		}
	}
}
