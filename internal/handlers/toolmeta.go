package handlers

import (
	"net/http"
	"sort"

	"github.com/chubuk/mcp-is/internal/mcp"
)

// ToolMeta отдаёт метаданные инструментов (name/group/class) для UI-редактора
// политики ИИ-ключа. Источник — единый маппинг tool_domains.go (без дублей в JS).
func ToolMeta(w http.ResponseWriter, _ *http.Request) {
	names := mcp.AllToolNames()
	sort.Strings(names)
	type item struct {
		Name   string `json:"name"`
		Group  string `json:"group"`
		Class  string `json:"class"`
		Domain string `json:"domain"`
	}
	items := make([]item, 0, len(names))
	for _, n := range names {
		m := mcp.ToolMetaFor(n)
		items = append(items, item{Name: n, Group: m.Group, Class: m.Class, Domain: m.Domain})
	}
	writeJSON(w, http.StatusOK, map[string]any{"tools": items})
}
