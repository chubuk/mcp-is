package insales

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/ratelimit"
	"github.com/chubuk/mcp-is/internal/resilience"
)

func TestBuildOrderBody(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
		want map[string]any
	}{
		{
			name: "empty args",
			args: nil,
			want: map[string]any{},
		},
		{
			name: "status alias maps to financial_status",
			args: map[string]any{"status": "paid", "comment": "hi"},
			want: map[string]any{"financial_status": "paid", "comment": "hi"},
		},
		{
			name: "financial_status wins over status alias",
			args: map[string]any{"status": "paid", "financial_status": "pending"},
			want: map[string]any{"financial_status": "pending"},
		},
		{
			name: "order_id never enters body",
			args: map[string]any{"order_id": "42", "comment": "x"},
			want: map[string]any{"comment": "x"},
		},
		{
			name: "order line create with product_id",
			args: map[string]any{
				"order_lines_attributes": []any{
					map[string]any{"product_id": 1, "quantity": 2, "vat": 10},
				},
			},
			want: map[string]any{
				"order_lines_attributes": []any{
					map[string]any{"product_id": 1, "quantity": 2, "vat": 10},
				},
			},
		},
		{
			name: "order line destroy",
			args: map[string]any{
				"order_lines_attributes": []any{
					map[string]any{"id": 7, "_destroy": true},
				},
			},
			want: map[string]any{
				"order_lines_attributes": []any{
					map[string]any{"id": 7, "_destroy": true},
				},
			},
		},
		{
			name: "discounts_attributes",
			args: map[string]any{
				"discounts_attributes": []any{
					map[string]any{"discount": 100, "type_id": 2, "description": "-100"},
					map[string]any{"discount": 10, "discount_order_lines_ids": []any{25.0}},
				},
			},
			want: map[string]any{
				"discounts_attributes": []any{
					map[string]any{"discount": 100, "type_id": 2, "description": "-100"},
					map[string]any{"discount": 10, "discount_order_lines_ids": []any{25.0}},
				},
			},
		},
		{
			name: "shipping_address_attributes",
			args: map[string]any{
				"shipping_address_attributes": map[string]any{
					"address":            "Moscow, Znamenka 12",
					"full_locality_name": "Moscow",
				},
			},
			want: map[string]any{
				"shipping_address_attributes": map[string]any{
					"address":            "Moscow, Znamenka 12",
					"full_locality_name": "Moscow",
				},
			},
		},
		{
			name: "custom_status_permalink",
			args: map[string]any{"custom_status_permalink": "v-obrabotke"},
			want: map[string]any{"custom_status_permalink": "v-obrabotke"},
		},
		{
			name: "unknown keys dropped",
			args: map[string]any{"order_id": "1", "random": "nope", "comment": "ok"},
			want: map[string]any{"comment": "ok"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildOrderBody(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildOrderBody() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func captureOrderRequest(t *testing.T, method string, fn func(c *InSalesClient) (*Order, error)) map[string]any {
	t.Helper()
	var captured map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 {
			_ = json.Unmarshal(body, &captured)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1, "order_number": "A1"}`))
	}))
	t.Cleanup(srv.Close)

	c := &InSalesClient{
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
		CircuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		RateTracker:    ratelimit.NewRateTracker(500, 5*time.Minute),
		BreakerPool:    ratelimit.NewPerKeyBreaker(5, 30*time.Second),
	}
	if _, err := fn(c); err != nil {
		t.Fatalf("call failed: %v", err)
	}
	return captured
}

func TestCreateOrder_PayloadShape(t *testing.T) {
	args := map[string]any{
		"comment": "from tool",
		"status":  "paid",
		"order_lines_attributes": []any{
			map[string]any{"product_id": 5, "quantity": 2},
		},
		"shipping_address_attributes": map[string]any{
			"address":            "Moscow",
			"full_locality_name": "Moscow",
		},
	}
	got := captureOrderRequest(t, http.MethodPost, func(c *InSalesClient) (*Order, error) {
		return c.CreateOrder(context.Background(), 99, args)
	})
	order, ok := got["order"].(map[string]any)
	if !ok {
		t.Fatalf("payload missing order: %#v", got)
	}
	if order["client_id"] != float64(99) && order["client_id"] != 99 {
		t.Errorf("client_id = %v, want 99", order["client_id"])
	}
	if order["comment"] != "from tool" {
		t.Errorf("comment = %v", order["comment"])
	}
	if order["financial_status"] != "paid" {
		t.Errorf("financial_status = %v, want paid (from status alias)", order["financial_status"])
	}
	if _, ok := order["status"]; ok {
		t.Error("status alias leaked into body")
	}
	if _, ok := order["order_id"]; ok {
		t.Error("order_id leaked into body")
	}
	lines, ok := order["order_lines_attributes"].([]any)
	if !ok || len(lines) != 1 {
		t.Fatalf("order_lines_attributes = %#v", order["order_lines_attributes"])
	}
	ship, ok := order["shipping_address_attributes"].(map[string]any)
	if !ok || ship["address"] != "Moscow" {
		t.Errorf("shipping_address_attributes = %#v", order["shipping_address_attributes"])
	}
}

func TestUpdateOrder_PayloadShape(t *testing.T) {
	args := map[string]any{
		"order_id": "42",
		"status":   "paid",
		"order_lines_attributes": []any{
			map[string]any{"id": 1, "quantity": 3},
			map[string]any{"id": 2, "_destroy": true},
		},
		"discounts_attributes": []any{
			map[string]any{"discount": 10, "type_id": 1},
		},
	}
	got := captureOrderRequest(t, http.MethodPut, func(c *InSalesClient) (*Order, error) {
		return c.UpdateOrder(context.Background(), "42", args)
	})
	order, ok := got["order"].(map[string]any)
	if !ok {
		t.Fatalf("payload missing order: %#v", got)
	}
	if _, ok := order["order_id"]; ok {
		t.Error("order_id must not be inside order body")
	}
	if order["financial_status"] != "paid" {
		t.Errorf("financial_status = %v", order["financial_status"])
	}
	lines, ok := order["order_lines_attributes"].([]any)
	if !ok || len(lines) != 2 {
		t.Fatalf("lines = %#v", order["order_lines_attributes"])
	}
	second := lines[1].(map[string]any)
	if second["_destroy"] != true {
		t.Errorf("second line _destroy = %#v, want true", second["_destroy"])
	}
	disc, ok := order["discounts_attributes"].([]any)
	if !ok || len(disc) != 1 {
		t.Fatalf("discounts = %#v", order["discounts_attributes"])
	}
}

func TestEnrichFieldValues(t *testing.T) {
	values := []any{
		map[string]any{"id": 9, "product_field_id": 1, "value": "123"},
		map[string]any{"id": 10, "product_field_id": 2, "value": "456", "handle": "keep"},
	}
	catalog := []any{
		map[string]any{"id": 1, "handle": "dosage", "title": "Дозировка"},
		map[string]any{"id": 2, "handle": "sku_ext", "title": "Внешний артикул"},
	}
	got := EnrichFieldValues(values, catalog)
	m0 := got[0].(map[string]any)
	if m0["handle"] != "dosage" || m0["title"] != "Дозировка" {
		t.Errorf("first value not enriched: %#v", m0)
	}
	m1 := got[1].(map[string]any)
	if m1["handle"] != "keep" {
		t.Errorf("existing handle must be kept: %#v", m1)
	}
}
