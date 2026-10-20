package handlers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type parsedTool struct {
	Name     string
	Desc     string
	ReadOnly bool
	HasRO    bool
}

// parseToolsGo extracts tool registrations from internal/mcp/tools.go via go/ast
// (no import of the mcp package — avoids cycles).
func parseToolsGo(t *testing.T, rel string) []parsedTool {
	t.Helper()
	path := filepath.Join("..", "..", "internal", "mcp", rel)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	var tools []parsedTool
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		fun, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if fun.Sel == nil || fun.Sel.Name != "NewTool" || len(call.Args) == 0 {
			return true
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		name, err := strconv.Unquote(lit.Value)
		if err != nil {
			return true
		}
		tool := parsedTool{Name: name}
		for _, arg := range call.Args[1:] {
			opt, ok := arg.(*ast.CallExpr)
			if !ok || len(opt.Args) == 0 {
				continue
			}
			sel, ok := opt.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel == nil {
				continue
			}
			switch sel.Sel.Name {
			case "WithDescription":
				if dl, ok := opt.Args[0].(*ast.BasicLit); ok && dl.Kind == token.STRING {
					tool.Desc, _ = strconv.Unquote(dl.Value)
				}
			case "WithReadOnlyHintAnnotation":
				if bl, ok := opt.Args[0].(*ast.Ident); ok {
					tool.HasRO = true
					tool.ReadOnly = bl.Name == "true"
				}
			}
		}
		tools = append(tools, tool)
		return true
	})
	return tools
}

func TestCommandListsSyncWithTools(t *testing.T) {
	tools := parseToolsGo(t, "tools.go")
	disc := parseToolsGo(t, "discovery.go")
	tools = append(tools, disc...)
	if len(tools) < 200 {
		t.Fatalf("parsed too few tools from tools.go: %d", len(tools))
	}

	wantAll := map[string]parsedTool{}
	for _, tool := range tools {
		if _, dup := wantAll[tool.Name]; dup {
			t.Errorf("duplicate tool name %s", tool.Name)
			continue
		}
		wantAll[tool.Name] = tool
	}

	gotAll := map[string]bool{}
	for _, name := range allCommandNames {
		if gotAll[name] {
			t.Errorf("duplicate in allCommandNames: %s", name)
		}
		gotAll[name] = true
	}
	gotRO := map[string]bool{}
	for _, name := range readCommands {
		gotRO[name] = true
	}

	for name, tool := range wantAll {
		if !gotAll[name] {
			t.Errorf("tool %s missing from allCommandNames", name)
		}
		if _, ok := commandDescriptions[name]; !ok {
			t.Errorf("tool %s missing from commandDescriptions", name)
		} else if strings.TrimSpace(commandDescriptions[name]) == "" {
			t.Errorf("tool %s has empty description", name)
		}
		if tool.HasRO && tool.ReadOnly && !gotRO[name] {
			t.Errorf("read-only tool %s missing from readCommands", name)
		}
		if tool.HasRO && !tool.ReadOnly && gotRO[name] {
			t.Errorf("write tool %s must not be in readCommands", name)
		}
		if !tool.HasRO {
			t.Errorf("tool %s has no explicit ReadOnlyHint", name)
		}
	}
	for name := range gotAll {
		if _, ok := wantAll[name]; !ok {
			t.Errorf("allCommandNames has unknown tool %s", name)
		}
	}
	for name := range gotRO {
		if _, ok := wantAll[name]; !ok {
			t.Errorf("readCommands has unknown tool %s", name)
		}
		if tool, ok := wantAll[name]; ok && tool.HasRO && !tool.ReadOnly {
			t.Errorf("readCommands contains write tool %s", name)
		}
	}

	if len(allCommandNames) != len(wantAll) {
		t.Errorf("allCommandNames len=%d, tools.go len=%d", len(allCommandNames), len(wantAll))
	}
	if len(commandDescriptions) != len(wantAll) {
		t.Errorf("commandDescriptions len=%d, tools.go len=%d", len(commandDescriptions), len(wantAll))
	}
}

func TestToolsGoExists(t *testing.T) {
	if _, err := os.Stat(filepath.Join("..", "..", "internal", "mcp", "tools.go")); err != nil {
		t.Fatal(err)
	}
}
