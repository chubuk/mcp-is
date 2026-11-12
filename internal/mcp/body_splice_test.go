package mcp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// spliceReq прогоняет body через spliceImageBody как это делает ServeHTTP.
func spliceReq(t *testing.T, body string) (*http.Request, string) {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	path, cleanup := spliceImageBody(r)
	t.Cleanup(cleanup)
	return r, path
}

func bodyString(t *testing.T, r *http.Request) string {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

const (
	rpcAddImage = `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"add_image","arguments":{"product_id":"5","filename":"a.png","image_data":"%s"}}}`
)

func TestSpliceImageBody_Basic(t *testing.T) {
	// "hello" → base64 aGVsbG8=; значение уходит в спул, в теле остаётся маска.
	r, path := spliceReq(t, strings.Replace(rpcAddImage, "%s", "aGVsbG8=", 1))
	if path == "" {
		t.Fatal("spool path not returned")
	}
	got := bodyString(t, r)
	if strings.Contains(got, "aGVsbG8=") {
		t.Errorf("raw base64 leaked into body: %s", got)
	}
	if strings.Contains(got, "image_file") {
		t.Errorf("image_file must not be injected into body: %s", got)
	}
	var msg struct {
		Params struct {
			Arguments map[string]any `json:"arguments"`
		} `json:"params"`
	}
	if err := json.Unmarshal([]byte(got), &msg); err != nil {
		t.Fatalf("spliced body is not valid JSON: %v\n%s", err, got)
	}
	mask, _ := msg.Params.Arguments["image_data"].(string)
	if !strings.HasPrefix(mask, "<base64,") {
		t.Errorf("image_data must be masked, got %q", mask)
	}
	// Маска обязана совпадать с маской handler-пути (diff/хэш не разъезжаются).
	sum := sha256.Sum256([]byte("hello"))
	want := imageMask(5, sum[:])
	if mask != want {
		t.Errorf("mask = %q, want %q (spoolImageData contract)", mask, want)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read spool: %v", err)
	}
	if string(raw) != "hello" {
		t.Errorf("spool bytes = %q, want %q", raw, "hello")
	}
}

func TestSpliceImageBody_Escapes(t *testing.T) {
	tests := []struct {
		name  string
		value string // как лежит в JSON (уже экранированное)
	}{
		{"newline in base64", `aGVs\nbG8=`}, // PEM-подобный ввод
		{"escaped slash in data uri", `data:image\/png;base64,aGVsbG8=`},
		{"unicode escape", `aGVs\u0062G8=`}, // \u0062 = 'b'
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, path := spliceReq(t, strings.Replace(rpcAddImage, "%s", tt.value, 1))
			if path == "" {
				t.Fatal("spool path not returned")
			}
			got := bodyString(t, r)
			if strings.Contains(got, "aGVsbG8=") || strings.Contains(got, "hello") {
				t.Errorf("raw value leaked into body: %s", got)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read spool: %v", err)
			}
			if string(raw) != "hello" {
				t.Errorf("spool bytes = %q, want %q", raw, "hello")
			}
		})
	}
}

func TestSpliceImageBody_FalsePositiveInString(t *testing.T) {
	// Текст "image_data":"…" внутри строки (body_html) — не ключ, тело не трогается.
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_blog_post","arguments":{"blog_id":"1","title":"t","body_html":"bla \"image_data\":\"aGVsbG8=\" bla","published_at":"2026-01-01","author":"a"}}}`
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("must not splice inside string values, got spool %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body must be byte-identical:\n got %s\nwant %s", got, body)
	}
	// Настоящий ключ при этом режется, текст в body_html остаётся целым.
	bodyBoth := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_blog_post","arguments":{"blog_id":"1","body_html":"bla \"image_data\":\"aGVsbG8=\" bla","image_data":"aGVsbG8=","image_filename":"a.png","published_at":"2026-01-01","author":"a"}}}`
	r2, path2 := spliceReq(t, bodyBoth)
	if path2 == "" {
		t.Fatal("real image_data must be spliced")
	}
	got := bodyString(t, r2)
	if !strings.Contains(got, `"bla \"image_data\":\"aGVsbG8=\" bla"`) {
		t.Errorf("body_html text mangled: %s", got)
	}
	if strings.Contains(got, `"image_data":"aGVsbG8="`) {
		t.Errorf("raw base64 leaked into body: %s", got)
	}
}

func TestSpliceImageBody_WrongPath(t *testing.T) {
	// image_data не на пути params.arguments.image_data — не трогаем.
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"add_image","arguments":{"nested":{"image_data":"aGVsbG8="},"filename":"a.png"}}}`
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("must not splice nested image_data, got spool %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body changed:\n got %s\nwant %s", got, body)
	}
}

func TestSpliceImageBody_WrongTool(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"name before args", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_products","arguments":{"image_data":"aGVsbG8="}}}`},
		{"name after args", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"arguments":{"image_data":"aGVsbG8="},"name":"list_products"}}`},
		{"no name", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"arguments":{"image_data":"aGVsbG8="}}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, path := spliceReq(t, tt.body)
			if path != "" {
				t.Errorf("must not splice for foreign tool, got spool %q", path)
			}
			if got := bodyString(t, r); got != tt.body {
				t.Errorf("body changed:\n got %s\nwant %s", got, tt.body)
			}
		})
	}
}

func TestSpliceImageBody_BadBase64RestoresBody(t *testing.T) {
	// Ошибка decode → тело восстанавливается байт-в-байт, tool-ошибку
	// вернёт штатный handler.
	body := strings.Replace(rpcAddImage, "%s", "не-base64!!!", 1)
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("bad base64 must not produce spool, got %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body must be byte-identical:\n got %s\nwant %s", got, body)
	}
}

func TestSpliceImageBody_EmptyValueRestoresBody(t *testing.T) {
	body := strings.Replace(rpcAddImage, "%s", "", 1)
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("empty value must not produce spool, got %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body must be byte-identical:\n got %s\nwant %s", got, body)
	}
}

func TestSpliceImageBody_BatchUntouched(t *testing.T) {
	body := `[{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"add_image","arguments":{"image_data":"aGVsbG8="}}}]`
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("batch must not be spliced, got spool %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body changed:\n got %s\nwant %s", got, body)
	}
}

func TestSpliceImageBody_DataURI(t *testing.T) {
	r, path := spliceReq(t, strings.Replace(rpcAddImage, "%s", "data:image/png;base64,aGVsbG8=", 1))
	if path == "" {
		t.Fatal("spool path not returned")
	}
	if got := bodyString(t, r); strings.Contains(got, "aGVsbG8=") {
		t.Errorf("raw base64 leaked into body: %s", got)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read spool: %v", err)
	}
	if string(raw) != "hello" {
		t.Errorf("spool bytes = %q, want %q", raw, "hello")
	}
}

func TestSpliceImageBody_DuplicateKeyUntouched(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"add_image","arguments":{"image_data":"aGVsbG8=","image_data":"d29ybGQ="}}}`
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("duplicate image_data must not be spliced, got spool %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body must be byte-identical:\n got %s\nwant %s", got, body)
	}
}

func TestSpliceImageBody_NoImageUntouched(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_products","arguments":{"per_page":10}}}`
	r, path := spliceReq(t, body)
	if path != "" {
		t.Errorf("no image_data — no spool, got %q", path)
	}
	if got := bodyString(t, r); got != body {
		t.Errorf("body must be byte-identical:\n got %s\nwant %s", got, body)
	}
}

func TestSpliceImageBody_OtherLiteralsUntouched(t *testing.T) {
	// Числа/булевы/null и вложенные массивы не должны ломать сканер.
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"add_image","arguments":{"position":3,"alt":null,"tags":["a","b",{"image_data":"x"}],"flag":true,"image_data":"aGVsbG8=","filename":"a.png"}}}`
	r, path := spliceReq(t, body)
	if path == "" {
		t.Fatal("spool path not returned")
	}
	got := bodyString(t, r)
	if strings.Contains(got, "aGVsbG8=") {
		t.Errorf("raw base64 leaked into body: %s", got)
	}
	if !strings.Contains(got, `"tags":["a","b",{"image_data":"x"}]`) {
		t.Errorf("nested literals mangled: %s", got)
	}
}

func TestSpoolImageData_CtxPathSkipsDecode(t *testing.T) {
	// HTTP-уровень уже распаковал тело: handler только подставляет путь.
	r, path := spliceReq(t, strings.Replace(rpcAddImage, "%s", "aGVsbG8=", 1))
	var msg struct {
		Params struct {
			Arguments map[string]any `json:"arguments"`
		} `json:"params"`
	}
	if err := json.Unmarshal([]byte(bodyString(t, r)), &msg); err != nil {
		t.Fatal(err)
	}
	args := msg.Params.Arguments
	ctx := context.WithValue(context.Background(), mcpSpoolPathKey, path)
	cleanup, err := spoolImageData(ctx, args)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	cleanup()
	if args["image_file"] != path {
		t.Errorf("image_file = %v, want %q", args["image_file"], path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("spool must outlive handler cleanup (ServeHTTP owns it): %v", err)
	}
}

// benchImageBody собирает JSON-RPC тело с image_data = base64 из mb МБ.
func benchImageBody(mb int) []byte {
	raw := make([]byte, mb<<20)
	b64 := base64.StdEncoding.EncodeToString(raw)
	return []byte(`{"jsonrpc":"2.0","id":"m1","method":"tools/call","params":{"name":"add_image","arguments":{"product_id":"5","filename":"a.png","image_data":"` + b64 + `"}}}`)
}

// BenchmarkSpliceImageBody — CPU/диск полного ingress-пути (сплайс + decode в
// спул). Сравнить с BenchmarkDecodeImageStream: разница = цена body-splice.
func BenchmarkSpliceImageBody(b *testing.B) {
	body := benchImageBody(20)
	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
		path, cleanup := spliceImageBody(r)
		if path == "" {
			b.Fatal("no splice")
		}
		cleanup()
	}
}

// BenchmarkDecodeImageStream — прежний handler-путь: decode base64 из строки в
// памяти (как spoolImageData до body-splice), без сканирования тела.
func BenchmarkDecodeImageStream(b *testing.B) {
	raw := make([]byte, 20<<20)
	b64 := base64.StdEncoding.EncodeToString(raw)
	b.SetBytes(int64(len(b64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path, _, _, err := decodeImageStream(strings.NewReader(b64))
		if err != nil {
			b.Fatal(err)
		}
		_ = os.Remove(path)
	}
}
