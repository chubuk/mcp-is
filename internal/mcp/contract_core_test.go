package mcp

import (
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// TestContractCoreTools — контракт core-каталога по реестру: схемы, аннотации,
// метаданные. Форма запросов к InSales покрыта мок-тестами client-слоя
// (internal/integration/insales/client_test.go); здесь — контракт MCP-поверхности.
func TestContractCoreTools(t *testing.T) {
	s := newDiscoveryTestServer()
	RegisterTools(s, nil)

	for name := range coreToolNames {
		st := s.GetTool(name)
		if st == nil {
			t.Errorf("core tool %q not registered", name)
			continue
		}
		tool := st.Tool
		if tool.Description == "" {
			t.Errorf("tool %q: empty description", name)
		}
		if tool.Title == "" {
			t.Errorf("tool %q: empty title", name)
		}
		if tool.InputSchema.Type != "object" {
			t.Errorf("tool %q: input schema type %q, want object", name, tool.InputSchema.Type)
		}
		for prop, def := range tool.InputSchema.Properties {
			if _, isBool := def.(bool); isBool {
				t.Errorf("tool %q: property %q is a bool schema", name, prop)
			}
		}
		// Аннотация безопасности обязательна явно (nil ReadOnlyHint = утечка write в RO filter).
		if tool.Annotations.ReadOnlyHint == nil {
			t.Errorf("tool %q: ReadOnlyHint annotation missing", name)
		}
		meta := ToolMetaFor(name)
		if meta.Class == "" {
			t.Errorf("tool %q: no toolMeta entry", name)
		}
	}
}

// TestContractReadOnlyAnnotationMatchesClass — read-класс = readOnlyHint,
// кроме call_insales_tool (диспетчер, вызывает и мутации). Фиксирует расхождение
// 122 vs 123 как осознанный инвариант, а не дрейф.
func TestContractReadOnlyAnnotationMatchesClass(t *testing.T) {
	s := server.NewMCPServer("test", "1.0")
	RegisterTools(s, nil)
	RegisterDiscoveryTools(s, nil, nil)

	checked := 0
	for name := range toolMeta {
		if ToolMetaFor(name).Class != ClassRead {
			continue
		}
		checked++
		st := s.GetTool(name)
		if st == nil {
			t.Fatalf("read-class tool %q not registered", name)
		}
		ro := st.Tool.Annotations.ReadOnlyHint != nil && *st.Tool.Annotations.ReadOnlyHint
		if name == "call_insales_tool" {
			if ro {
				t.Error("call_insales_tool dispatches mutations — must not claim ReadOnlyHint")
			}
			continue
		}
		if !ro {
			t.Errorf("read-class tool %q lacks ReadOnlyHint(true)", name)
		}
	}
	if checked != 123 {
		t.Errorf("read-class tools = %d, want 123 (docs are synced to this number)", checked)
	}
}

// TestContractCreateToolsHaveIdempotencyKey — все create_* принимают idempotency_key.
func TestContractCreateToolsHaveIdempotencyKey(t *testing.T) {
	s := server.NewMCPServer("test", "1.0")
	RegisterTools(s, nil)
	RegisterDiscoveryTools(s, nil, nil)

	count := 0
	for name := range toolMeta {
		if !strings.HasPrefix(name, "create_") {
			continue
		}
		count++
		st := s.GetTool(name)
		if st == nil {
			t.Fatalf("create tool %q not registered", name)
		}
		if _, ok := st.Tool.InputSchema.Properties["idempotency_key"]; !ok {
			t.Errorf("create tool %q lacks idempotency_key", name)
		}
	}
	if count != 52 {
		t.Errorf("create_* tools = %d, want 52 (docs/CHANGELOG)", count)
	}
	// Обратная проверка: ключ только у create_* — add_*/update_*/delete_* без него.
	for name := range toolMeta {
		if strings.HasPrefix(name, "create_") {
			continue
		}
		st := s.GetTool(name)
		if st == nil {
			continue
		}
		if _, ok := st.Tool.InputSchema.Properties["idempotency_key"]; ok {
			t.Errorf("non-create tool %q must not expose idempotency_key", name)
		}
	}
}
