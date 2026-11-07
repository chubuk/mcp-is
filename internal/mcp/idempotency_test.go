package mcp

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestWithIdempotencyStripsKeyWithoutRepo(t *testing.T) {
	var got map[string]any
	inner := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		got = req.GetArguments()
		return mcp.NewToolResultText("ok"), nil
	}
	h := withIdempotency("create_product", inner)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "create_product",
			Arguments: map[string]any{"title": "X", "idempotency_key": "k1"},
		},
	}
	res, err := h(context.Background(), req)
	if err != nil || res == nil || res.IsError {
		t.Fatalf("unexpected result: %v %v", res, err)
	}
	if _, ok := got["idempotency_key"]; ok {
		t.Errorf("idempotency_key leaked to handler: %v", got)
	}
	if got["title"] != "X" {
		t.Errorf("real args lost: %v", got)
	}
}

func TestWithIdempotencyRejectsLongKey(t *testing.T) {
	inner := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		t.Fatal("handler must not run on invalid key")
		return nil, nil
	}
	h := withIdempotency("create_product", inner)
	long := make([]byte, 129)
	for i := range long {
		long[i] = 'a'
	}
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "create_product",
			Arguments: map[string]any{"idempotency_key": string(long)},
		},
	}
	res, err := h(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Errorf("expected error for long key, got %v", res)
	}
}

func TestCreateToolsExposeIdempotencyKey(t *testing.T) {
	s := server.NewMCPServer("test", "1.0")
	found := 0
	for _, name := range []string{"create_product", "create_order", "create_client"} {
		tool := mcp.NewTool(name,
			mcp.WithDescription("t"),
			mcp.WithString("title", mcp.Description("x")),
			mcp.WithReadOnlyHintAnnotation(false),
		)
		addTool(s, tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText("ok"), nil
		})
		st := s.GetTool(name)
		if st == nil {
			t.Fatalf("tool %s not registered", name)
		}
		if _, ok := st.Tool.InputSchema.Properties["idempotency_key"]; !ok {
			t.Errorf("tool %s lacks idempotency_key property", name)
		}
		found++
	}
	if found != 3 {
		t.Fatalf("expected 3 tools, got %d", found)
	}
}

// fakeStore — in-memory idempotencyStore для тестов обёртки без БД.
// Ветки состояний (replay / in-progress / reuse / cancel) — детерминированные;
// сама гонка INSERT доказывается на реальной БД (TestIdempotencyRepoConcurrentBegin).
type fakeStore struct {
	mu       sync.Mutex
	records  map[string]*storage.IdempotencyRecord
	beginErr error
}

func newFakeStore() *fakeStore {
	return &fakeStore{records: map[string]*storage.IdempotencyRecord{}}
}

func (f *fakeStore) Begin(ctx context.Context, aiKeyID, key, tool, requestHash string) (*storage.IdempotencyRecord, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.beginErr != nil {
		return nil, false, f.beginErr
	}
	id := aiKeyID + "|" + key
	if rec, ok := f.records[id]; ok {
		cp := *rec
		return &cp, false, nil
	}
	rec := &storage.IdempotencyRecord{Tool: tool, RequestHash: requestHash, Status: "pending"}
	f.records[id] = rec
	cp := *rec
	return &cp, true, nil
}

func (f *fakeStore) Finish(ctx context.Context, aiKeyID, key string, response []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if rec, ok := f.records[aiKeyID+"|"+key]; ok {
		rec.Status = "done"
		rec.Response = response
	}
	return nil
}

func (f *fakeStore) Cancel(ctx context.Context, aiKeyID, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	id := aiKeyID + "|" + key
	if rec, ok := f.records[id]; ok && rec.Status == "pending" {
		delete(f.records, id)
	}
	return nil
}

func idemCtx() context.Context {
	return context.WithValue(context.Background(), auth.KeyIDKey, "ai-1")
}

func idemReq(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: "create_thing", Arguments: args},
	}
}

func okInner(calls *int) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		(*calls)++
		res := mcp.NewToolResultStructured(map[string]any{"id": 42}, `{"id":42}`)
		return res, nil
	}
}

// TestWithIdempotencyReplayFlag — повтор с тем же payload после успеха:
// handler не вызывается второй раз, ответ сохранён + idempotent_replay.
func TestWithIdempotencyReplayFlag(t *testing.T) {
	store := newFakeStore()
	idemRepo = store
	t.Cleanup(func() { idemRepo = nil })

	calls := 0
	h := withIdempotency("create_thing", okInner(&calls))

	res, err := h(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"}))
	if err != nil || res == nil || res.IsError {
		t.Fatalf("first call: %v %v", res, err)
	}
	res, err = h(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"}))
	if err != nil || res == nil || res.IsError {
		t.Fatalf("replay call: %v %v", res, err)
	}
	if calls != 1 {
		t.Errorf("handler calls = %d, want 1", calls)
	}
	sc, _ := res.StructuredContent.(map[string]any)
	if sc["idempotent_replay"] != true {
		t.Errorf("idempotent_replay missing: %v", sc)
	}
}

// TestWithIdempotencyInProgress — повтор, пока первый ещё выполняется (pending):
// осмысленная ошибка, handler не вызывается (не дубль!).
func TestWithIdempotencyInProgress(t *testing.T) {
	store := newFakeStore()
	idemRepo = store
	t.Cleanup(func() { idemRepo = nil })

	if _, created, err := store.Begin(idemCtx(), "ai-1", "k1", "create_thing", hashArgs(map[string]any{"title": "X"})); err != nil || !created {
		t.Fatalf("seed pending: created=%v err=%v", created, err)
	}

	calls := 0
	h := withIdempotency("create_thing", okInner(&calls))
	res, err := h(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"}))
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if res == nil || !res.IsError {
		t.Fatalf("want error result, got %+v", res)
	}
	if calls != 0 {
		t.Errorf("handler must not run on pending, calls=%d", calls)
	}
	text, _ := res.Content[0].(mcp.TextContent)
	if !strings.Contains(text.Text, "выполняется") {
		t.Errorf("error text = %q, want in-progress message", text.Text)
	}
}

// TestWithIdempotencyKeyReuse — тот же ключ с другим payload: ошибка,
// handler не вызывается.
func TestWithIdempotencyKeyReuse(t *testing.T) {
	store := newFakeStore()
	idemRepo = store
	t.Cleanup(func() { idemRepo = nil })

	calls := 0
	h := withIdempotency("create_thing", okInner(&calls))
	if res, err := h(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"})); err != nil || res.IsError {
		t.Fatalf("first call: %v %v", res, err)
	}
	res, err := h(idemCtx(), idemReq(map[string]any{"title": "OTHER", "idempotency_key": "k1"}))
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if res == nil || !res.IsError {
		t.Fatalf("want reuse error, got %+v", res)
	}
	if calls != 1 {
		t.Errorf("handler calls = %d, want 1 (second must be rejected)", calls)
	}
	text, _ := res.Content[0].(mcp.TextContent)
	if !strings.Contains(text.Text, "другими данными") {
		t.Errorf("error text = %q, want reuse message", text.Text)
	}
}

// TestWithIdempotencyCancelOnFailure — упавший handler освобождает ключ:
// повтор с тем же ключом выполняется заново.
func TestWithIdempotencyCancelOnFailure(t *testing.T) {
	store := newFakeStore()
	idemRepo = store
	t.Cleanup(func() { idemRepo = nil })

	calls := 0
	fail := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calls++
		return mcp.NewToolResultError("boom"), nil
	}
	h := withIdempotency("create_thing", fail)
	res, err := h(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"}))
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if res == nil || !res.IsError {
		t.Fatalf("want error result, got %+v", res)
	}

	h2 := withIdempotency("create_thing", okInner(&calls))
	res, err = h2(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"}))
	if err != nil || res == nil || res.IsError {
		t.Fatalf("retry after cancel: %v %v", res, err)
	}
	if calls != 2 {
		t.Errorf("handler calls = %d, want 2 (retry allowed after failure)", calls)
	}
}

// TestWithIdempotencyFailOpen — ошибка БД на Begin: дедупликация молча
// пропускается, handler выполняется как обычно.
func TestWithIdempotencyFailOpen(t *testing.T) {
	store := newFakeStore()
	store.beginErr = context.DeadlineExceeded
	idemRepo = store
	t.Cleanup(func() { idemRepo = nil })

	calls := 0
	h := withIdempotency("create_thing", okInner(&calls))
	res, err := h(idemCtx(), idemReq(map[string]any{"title": "X", "idempotency_key": "k1"}))
	if err != nil || res == nil || res.IsError {
		t.Fatalf("fail-open call: %v %v", res, err)
	}
	if calls != 1 {
		t.Errorf("handler calls = %d, want 1", calls)
	}
	sc, _ := res.StructuredContent.(map[string]any)
	if sc["idempotent_replay"] == true {
		t.Errorf("fail-open must not mark replay: %v", sc)
	}
}
