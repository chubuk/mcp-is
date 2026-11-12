package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, "null"},
		{"string", "hello", `"hello"`},
		{"number", 42, "42"},
		{"bool", true, "true"},
		{"map", map[string]any{"key": "value"}, "{\n  \"key\": \"value\"\n}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatResult(tt.input)
			if result != tt.expected {
				t.Errorf("formatResult(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// structuredContent по MCP-спеке обязан сериализоваться в JSON-объект.
func TestWrapEnvelope(t *testing.T) {
	tests := []struct {
		name string
		in   any
	}{
		{"slice", []any{map[string]any{"id": "1"}}},
		{"dynamic slice in any", any([]any{map[string]any{"id": "1"}})},
		{"typed nil slice", []map[string]any(nil)},
		{"string", "hello"},
		{"int", 42},
		{"byte slice", []byte("raw")},
		{"raw message array", json.RawMessage(`[1,2]`)},
		{"raw message object", json.RawMessage(`{"a":1}`)},
		{"map", map[string]any{"a": 1}},
		{"ptr to slice", &[]int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := wrapEnvelope(tt.in)
			b, err := json.Marshal(out)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if b[0] != '{' {
				t.Errorf("marshals to %s, want JSON object", b)
			}
		})
	}
	// nil остаётся nil — поле omitempty, спека разрешает отсутствие.
	if got := wrapEnvelope(nil); got != nil {
		t.Errorf("wrapEnvelope(nil) = %v, want nil", got)
	}
	// слайс → домашний конверт {items,total}
	out := wrapEnvelope([]any{"a", "b"}).(map[string]any)
	if out["total"] != 2 {
		t.Errorf("total = %v, want 2", out["total"])
	}
	if items, ok := out["items"].([]any); !ok || len(items) != 2 {
		t.Errorf("items = %v, want 2-element slice", out["items"])
	}
	// объект не модифицируется
	obj := map[string]any{"a": 1}
	if wrapEnvelope(obj).(map[string]any)["a"] != 1 {
		t.Error("map payload modified")
	}
}

func TestToolResult_ObjectStructuredContent(t *testing.T) {
	res := toolResult([]any{map[string]any{"id": "1"}})
	if res.StructuredContent == nil {
		t.Fatal("structuredContent missing")
	}
	b, err := json.Marshal(res.StructuredContent)
	if err != nil || b[0] != '{' {
		t.Errorf("structuredContent = %s (%v), want JSON object", b, err)
	}
}

func TestGetFloat(t *testing.T) {
	args := map[string]any{
		"page":     float64(5),
		"per_page": float64(50),
	}

	if got := getFloat(args, "page", 1); got != 5 {
		t.Errorf("getFloat(page) = %v, want 5", got)
	}
	if got := getFloat(args, "missing", 10); got != 10 {
		t.Errorf("getFloat(missing) = %v, want 10", got)
	}
	if got := getFloat(nil, "key", 7); got != 7 {
		t.Errorf("getFloat(nil) = %v, want 7", got)
	}
}

func TestHashArgsStable(t *testing.T) {
	a := hashArgs(map[string]any{"b": 2, "a": 1})
	b := hashArgs(map[string]any{"a": 1, "b": 2})
	if a == "" || a != b {
		t.Errorf("hash not stable: %q vs %q", a, b)
	}
	if hashArgs(nil) != "" {
		t.Error("empty args should hash empty")
	}
}

func TestAuditFromContext(t *testing.T) {
	e := auditFromContext(context.Background(), "create_product", "call", map[string]any{"title": "x"}, true)
	if e.ToolName != "create_product" || e.Action != "call" || !e.ResultOK {
		t.Errorf("entry = %+v", e)
	}
	if e.ArgsHash == "" {
		t.Error("expected args hash")
	}
}

func TestServerJSONVersionMatches(t *testing.T) {
	raw, err := os.ReadFile("../../server.json")
	if err != nil {
		t.Skip("server.json not found")
	}
	var doc struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	// keep in sync with server.go NewMCPServer version literal
	src, err := os.ReadFile("server.go")
	if err != nil {
		t.Fatal(err)
	}
	want := regexp.MustCompile(`"(\d+\.\d+\.\d+)"`).FindStringSubmatch(string(src))
	if want == nil {
		t.Fatal("version literal not found in server.go")
	}
	if doc.Version != want[1] {
		t.Errorf("server.json version %q != server.go %q", doc.Version, want[1])
	}
}

// --- ПДн: applyPDNToResult ---

func pdnCtx() context.Context {
	return context.WithValue(context.Background(), mcpShopPDNKey, true)
}

func TestApplyPDNToResult_TypedOrderRoundTrip(t *testing.T) {
	orders := []insales.Order{
		{ID: "1", ClientName: "Иванов Иван", ShippingAddress: map[string]any{
			"phone": "+79991234567", "full_name": "Иванов Иван",
		}},
	}
	out := applyPDNToResult(pdnCtx(), orders, "list_orders")
	typed, ok := out.([]insales.Order)
	if !ok {
		t.Fatalf("result type %T, want []insales.Order", out)
	}
	if len(typed) != 1 || typed[0].ID != "1" {
		t.Fatalf("structure broken: %+v", typed)
	}
	if typed[0].ClientName != "И***" {
		t.Errorf("client_name = %q, want masked", typed[0].ClientName)
	}
	addr := typed[0].ShippingAddress
	if addr["phone"] != "+7 (***) ***-**-**" || addr["full_name"] != "И***" {
		t.Errorf("shipping_address not masked: %v", addr)
	}
}

func TestApplyPDNToResult_MapTree(t *testing.T) {
	in := map[string]any{
		"client": map[string]any{"email": "ivan@example.com"},
	}
	out := applyPDNToResult(pdnCtx(), in, "client_full_profile").(map[string]any)
	if out["client"].(map[string]any)["email"] != "i***@example.com" {
		t.Errorf("nested map not masked: %v", out)
	}
}

func TestApplyPDNToResult_DisabledKeepsData(t *testing.T) {
	orders := []insales.Order{{ID: "1", ClientName: "Иванов"}}
	out := applyPDNToResult(context.Background(), orders, "list_orders").([]insales.Order)
	if out[0].ClientName != "Иванов" {
		t.Errorf("masked without PDN flag: %v", out[0].ClientName)
	}
}

// --- outputSchema: boolean-узлы ---

func TestFixSchemaNodes(t *testing.T) {
	in := map[string]any{
		"items": map[string]any{"items": true, "type": "array"},
		"ok":    false,
		"list":  []any{true},
	}
	out := fixSchemaNodes(in).(map[string]any)
	if _, isObj := out["items"].(map[string]any)["items"].(map[string]any); !isObj {
		t.Errorf("true not replaced: %v", out["items"])
	}
	if out["ok"] != false {
		t.Errorf("false must stay: %v", out["ok"])
	}
	if _, isObj := out["list"].([]any)[0].(map[string]any); !isObj {
		t.Errorf("true in slice not replaced: %v", out["list"])
	}
}

// --- outputSchema: валидация фикстурами (null + отсутствующие ключи) ---

func validateAgainstSchema(t *testing.T, tool mcp.Tool, payload any) {
	t.Helper()
	sb, err := json.Marshal(tool.OutputSchema)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	c := jsonschema.NewCompiler()
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(sb))
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	if err := c.AddResource("s.json", doc); err != nil {
		t.Fatalf("add schema: %v", err)
	}
	sch, err := c.Compile("s.json")
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	pb, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	var plain any
	if err := json.Unmarshal(pb, &plain); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if err := sch.Validate(plain); err != nil {
		t.Errorf("validation failed: %v\npayload: %s", err, pb)
	}
}

func TestOutputSchema_FixturesNullAndMissing(t *testing.T) {
	// null-поля: CategoryID/Variants без omitempty отдают null
	prod := insales.Product{ID: "1", Title: "Т"}
	validateAgainstSchema(t,
		mcp.NewTool("t", fixSchemaForTest(mcp.WithOutputSchema[insales.Product]())), prod)

	// отсутствующие ключи (omitempty): client=null / поля не обязательны
	details := OrderFullDetailsOut{Order: map[string]any{"id": "1"}, Client: nil}
	validateAgainstSchema(t,
		mcp.NewTool("t", fixSchemaForTest(mcp.WithOutputSchema[OrderFullDetailsOut]())), details)

	// map-payload без части ключей (как service.go pipeline)
	profile := map[string]any{"client": map[string]any{"id": "1"}}
	validateAgainstSchema(t,
		mcp.NewTool("t", fixSchemaForTest(mcp.WithOutputSchema[ClientFullProfileOut]())), profile)
}

// fixSchemaForTest применяет тот же фикс boolean-схем, что и registerTool.
func fixSchemaForTest(opt mcp.ToolOption) mcp.ToolOption {
	return func(t *mcp.Tool) {
		opt(t)
		if t.OutputSchema.Properties != nil {
			t.OutputSchema.Properties = fixSchemaNodes(t.OutputSchema.Properties).(map[string]any)
		}
		if t.OutputSchema.Defs != nil {
			t.OutputSchema.Defs = fixSchemaNodes(t.OutputSchema.Defs).(map[string]any)
		}
	}
}

func TestValidateImageArgs(t *testing.T) {
	const b64 = "aGVsbG8=" // "hello"
	tests := []struct {
		name    string
		args    map[string]any
		exactly bool
		wantErr string
		wantRaw string // ожидаемое содержимое спул-файла после успеха ("-" = без картинки)
	}{
		{
			name:    "url only ok",
			args:    map[string]any{"image_url": "https://x/i.png"},
			wantRaw: "-",
		},
		{
			name:    "data only ok",
			args:    map[string]any{"image_data": b64, "image_filename": "a.png"},
			wantRaw: "hello",
		},
		{
			name:    "both rejected",
			args:    map[string]any{"image_url": "https://x/i.png", "image_data": b64, "image_filename": "a.png"},
			wantErr: "взаимоисключимы",
		},
		{
			name:    "exactly one required",
			args:    map[string]any{},
			exactly: true,
			wantErr: "укажите",
		},
		{
			name:    "filename required with data",
			args:    map[string]any{"image_data": b64},
			wantErr: "обязателен",
		},
		{
			name:    "bad base64",
			args:    map[string]any{"image_data": "не base64!", "image_filename": "a.png"},
			wantErr: "base64",
		},
		{
			name:    "oversize",
			args:    map[string]any{"image_data": strings.Repeat("A", 28<<20), "image_filename": "a.png"},
			wantErr: "слишком большое",
		},
		{
			name:    "data uri prefix stripped",
			args:    map[string]any{"image_data": "data:image/png;base64," + b64, "image_filename": "a.png"},
			wantRaw: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup, err := validateImageArgs(context.Background(), tt.args, "image_url", "image_filename", tt.exactly)
			defer cleanup()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want substring %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			switch tt.wantRaw {
			case "-":
				if path, _ := tt.args["image_file"].(string); path != "" {
					t.Errorf("image_file must be empty without image_data, got %q", path)
				}
			default:
				path, _ := tt.args["image_file"].(string)
				if path == "" {
					t.Fatal("image_file not set")
				}
				got, rerr := os.ReadFile(path)
				if rerr != nil {
					t.Fatalf("read spool: %v", rerr)
				}
				if string(got) != tt.wantRaw {
					t.Errorf("spool bytes = %q, want %q", got, tt.wantRaw)
				}
				if mask, _ := tt.args["image_data"].(string); !strings.HasPrefix(mask, "<base64,") {
					t.Errorf("image_data must be masked, got %q", mask)
				}
			}
		})
	}
}

func TestSpoolImageData_CleanupRemovesFile(t *testing.T) {
	args := map[string]any{"image_data": "aGVsbG8=", "image_filename": "a.png"}
	cleanup, err := spoolImageData(context.Background(), args)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	path := args["image_file"].(string)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("spool file missing: %v", err)
	}
	cleanup()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("spool file not removed: %v", err)
	}
}

func TestSpoolImageData_MaskContentBound(t *testing.T) {
	// Разные картинки одного размера → разные маски (иначе коллизия ArgsHash).
	a := map[string]any{"image_data": "aGVsbG8=", "image_filename": "a.png"}
	b := map[string]any{"image_data": "d29ybGQ=", "image_filename": "b.png"} // "world"
	ca, err := spoolImageData(context.Background(), a)
	defer ca()
	if err != nil {
		t.Fatal(err)
	}
	cb, err := spoolImageData(context.Background(), b)
	defer cb()
	if err != nil {
		t.Fatal(err)
	}
	if a["image_data"] == b["image_data"] {
		t.Errorf("same mask for different content: %v", a["image_data"])
	}
	if !strings.Contains(a["image_data"].(string), "sha256:") {
		t.Errorf("mask must carry content hash: %v", a["image_data"])
	}
}

func TestSpoolImageData_EmptyPayloadRejected(t *testing.T) {
	noop, err := spoolImageData(context.Background(), map[string]any{"image_data": "data:image/png;base64,", "image_filename": "a.png"})
	defer noop()
	if err == nil || !strings.Contains(err.Error(), "пустое") {
		t.Errorf("err = %v, want пустое изображение", err)
	}
}

func TestHashArgs_BigStringMasked(t *testing.T) {
	big := strings.Repeat("A", 200<<10)
	h1 := hashArgs(map[string]any{"image_data": big, "title": "x"})
	h2 := hashArgs(map[string]any{"image_data": big, "title": "x"})
	h3 := hashArgs(map[string]any{"image_data": big + "B", "title": "x"})
	if h1 == "" || h1 != h2 {
		t.Errorf("hash not stable: %q vs %q", h1, h2)
	}
	if h1 == h3 {
		t.Error("different content must differ")
	}
	if strings.Contains(h1, strings.Repeat("A", 64)) {
		t.Error("hash leaked raw string")
	}
	// image_file — служебный путь, в хэш не входит
	if hashArgs(map[string]any{"image_file": "/tmp/x", "title": "x"}) != hashArgs(map[string]any{"title": "x"}) {
		t.Error("image_file must not affect hash")
	}
}

func TestDiffFromArgs_MasksImageData(t *testing.T) {
	diff := diffFromArgs(nil, map[string]any{
		"title":      "t",
		"image_data": strings.Repeat("A", 5000),
	}, nil)
	if diff["title"].New != "t" {
		t.Errorf("title masked unexpectedly: %#v", diff["title"])
	}
	newVal, _ := diff["image_data"].New.(string)
	if newVal == "" || newVal == strings.Repeat("A", 5000) {
		t.Errorf("image_data not masked: %#v", diff["image_data"])
	}
	if !strings.Contains(newVal, "base64") {
		t.Errorf("mask should say base64, got %q", newVal)
	}
}

func TestIdString(t *testing.T) {
	if got := idString(float64(3418753025)); got != "3418753025" {
		t.Errorf("idString(float64) = %q, want 3418753025", got)
	}
	if got := idString("16460017"); got != "16460017" {
		t.Errorf("idString(string) = %q, want 16460017", got)
	}
}
