package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/chubuk/mcp-is/internal/ratelimit"
	"github.com/chubuk/mcp-is/internal/resilience"
)

type stubProvider struct{ c *insales.InSalesClient }

func (p stubProvider) GetClient(ctx context.Context) (*insales.InSalesClient, context.Context, error) {
	return p.c, ctx, nil
}

func newFieldTestService(srv *httptest.Server) *Service {
	client := &insales.InSalesClient{
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
		CircuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		RateTracker:    ratelimit.NewRateTracker(500, 5*time.Minute),
		BreakerPool:    ratelimit.NewPerKeyBreaker(5, 30*time.Second),
	}
	return NewService(stubProvider{c: client})
}

func TestPopFieldValues(t *testing.T) {
	args := map[string]any{"title": "t", "field_values": map[string]any{"sostav": "x"}}
	fv := popFieldValues(args)
	if !reflect.DeepEqual(fv, map[string]any{"sostav": "x"}) {
		t.Errorf("fv = %#v", fv)
	}
	if _, ok := args["field_values"]; ok {
		t.Errorf("field_values must be removed from args: %#v", args)
	}
	if popFieldValues(args) != nil {
		t.Error("second pop must be nil")
	}
}

func TestInlineFieldAttrs(t *testing.T) {
	got := inlineFieldAttrs(map[string]any{"a": "1", "b": 2})
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	seen := map[string]any{}
	for _, m := range got {
		seen[m["handle"].(string)] = m["value"]
	}
	if seen["a"] != "1" || seen["b"] != 2 {
		t.Errorf("attrs = %#v", got)
	}
}

// InSales POST /orders.json игнорирует fields_values_attributes — CreateOrder
// обязан дописать их отдельным PUT (иначе значения молча теряются).
func TestCreateOrder_WritesFieldValuesFollowUp(t *testing.T) {
	var putBody map[string]any
	var putPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/fields.json":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": 777, "handle": "note", "destiny": 3}})
		case r.Method == http.MethodPost && r.URL.Path == "/orders.json":
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 123})
		case r.Method == http.MethodPut && r.URL.Path == "/orders/123.json":
			putPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&putBody)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 123})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	svc := newFieldTestService(srv)
	args := map[string]any{"comment": "x", "field_values": map[string]any{"note": "v1"}}
	order, err := svc.CreateOrder(context.Background(), 0, args)
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if order.ID != "123" {
		t.Fatalf("order id = %q", order.ID)
	}
	if putPath != "/orders/123.json" {
		t.Fatalf("follow-up PUT not sent, path = %q", putPath)
	}
	orderBody, _ := putBody["order"].(map[string]any)
	attrs, _ := orderBody["fields_values_attributes"].([]any)
	if len(attrs) != 1 {
		t.Fatalf("attrs = %#v", putBody)
	}
	attr, _ := attrs[0].(map[string]any)
	if anyToID(attr["field_id"]) != "777" || attr["value"] != "v1" {
		t.Errorf("attr = %#v", attr)
	}
}

func TestCollectionFieldAttrs_UnknownHandle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/collection_fields.json" {
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": 1, "handle": "size"}})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	svc := newFieldTestService(srv)
	_, err := svc.collectionFieldAttrs(context.Background(), map[string]any{"nope": "x"})
	if err == nil || !strings.Contains(err.Error(), "неизвестные handle") {
		t.Fatalf("want unknown-handle error, got %v", err)
	}
	if !strings.Contains(err.Error(), "size") {
		t.Errorf("error must list available handles: %v", err)
	}
	attrs, err := svc.collectionFieldAttrs(context.Background(), map[string]any{"size": "L"})
	if err != nil {
		t.Fatalf("valid handle: %v", err)
	}
	if len(attrs) != 1 || attrs[0]["handle"] != "size" || attrs[0]["value"] != "L" {
		t.Errorf("attrs = %#v", attrs)
	}
}
