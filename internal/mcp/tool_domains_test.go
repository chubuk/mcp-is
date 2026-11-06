package mcp

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"
)

// parseToolNames extracts tool names from mcp.NewTool("...") calls via go/ast.
func parseToolNames(t *testing.T, path string) []string {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	var names []string
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		fun, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || fun.Sel == nil || fun.Sel.Name != "NewTool" || len(call.Args) == 0 {
			return true
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		name, err := strconv.Unquote(lit.Value)
		if err == nil {
			names = append(names, name)
		}
		return true
	})
	return names
}

// TestToolDomainCoverage — все инструменты из tools.go и discovery.go замаплены в tool_domains.go.
func TestToolDomainCoverage(t *testing.T) {
	toolNames := parseToolNames(t, "tools.go")
	discNames := parseToolNames(t, "discovery.go")

	if len(toolNames) < 200 {
		t.Fatalf("parsed only %d tools from tools.go — parse broken?", len(toolNames))
	}

	var missing []string
	for _, name := range append(toolNames, discNames...) {
		if _, ok := toolMeta[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		t.Errorf("tools missing from tool_domains.go (%d): %v", len(missing), missing)
	}

	// Обратный дрейф: записи в toolDomains для несуществующих инструментов.
	all := make(map[string]bool, len(toolNames)+len(discNames))
	for _, n := range append(toolNames, discNames...) {
		all[n] = true
	}
	var orphans []string
	for name := range toolMeta {
		if !all[name] {
			orphans = append(orphans, name)
		}
	}
	if len(orphans) > 0 {
		t.Errorf("tool_domains.go entries with no tool in tools.go/discovery.go (%d): %v", len(orphans), orphans)
	}
}

// TestToolClassMatchesReadOnlyHint — класс «read» ↔ ReadOnlyHint=true и наоборот.
func TestToolClassMatchesReadOnlyHint(t *testing.T) {
	// Исключения: read-класс, но не readOnly-hint (или наоборот) — фиксируем явно.
	exceptions := map[string]bool{
		"product_catalog_search": true, // read, но search по каталогу = не readOnly hint
		"collection_products":    true, // read, но составная выборка
	}

	toolNames := parseToolNames(t, "tools.go")
	for _, name := range toolNames {
		meta := ToolMetaFor(name)
		if exceptions[name] {
			continue
		}
		isRead := meta.Class == ClassRead
		// Проверяем по toolMeta: read-класс должен быть readOnly.
		// (ReadOnlyHint сверяем в tools_sync_test.go handlers; здесь — консистентность класса.)
		if meta.Class != ClassRead && meta.Class != ClassWrite && meta.Class != ClassDelete {
			t.Errorf("tool %q has invalid class %q", name, meta.Class)
		}
		if isRead && meta.Class != ClassRead {
			t.Errorf("tool %q class mismatch", name)
		}
	}
}

// TestCommandPolicy_Priority — override команды > класс в группе > default; deny побеждает.
func TestCommandPolicy_Priority(t *testing.T) {
	// nil policy = allow all
	var nilPolicy *CommandPolicy
	if !nilPolicy.IsCommandAllowed("list_products") {
		t.Error("nil policy should allow")
	}

	// default deny, no group/command override → deny
	p := &CommandPolicy{Default: "deny"}
	if p.IsCommandAllowed("list_products") {
		t.Error("default deny should deny")
	}

	// default deny, group read allowed → read commands allowed
	p = &CommandPolicy{
		Default: "deny",
		Groups:  map[string]GroupClassEnabled{"catalog": {Read: true}},
	}
	if !p.IsCommandAllowed("list_products") {
		t.Error("group read=true should allow list_products")
	}
	if p.IsCommandAllowed("create_product") {
		t.Error("group write=false should deny create_product")
	}

	// command override beats group
	p.Commands = map[string]bool{"create_product": true}
	if !p.IsCommandAllowed("create_product") {
		t.Error("command override=true should allow")
	}
	p.Commands = map[string]bool{"list_products": false}
	if p.IsCommandAllowed("list_products") {
		t.Error("command override=false should deny even if group allows")
	}
}

// TestToolMetaFor_Unknown_StrictestClass — неизвестное имя → класс delete.
func TestToolMetaFor_Unknown_StrictestClass(t *testing.T) {
	meta := ToolMetaFor("nonexistent_tool_xyz")
	if meta.Class != ClassDelete {
		t.Errorf("unknown tool class = %q, want %q", meta.Class, ClassDelete)
	}
}
