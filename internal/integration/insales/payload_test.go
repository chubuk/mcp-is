package insales

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCleanArgs_DropsNilAndEmptyStrings(t *testing.T) {
	got := cleanArgs(map[string]any{
		"title":    "x",
		"sku":      "",
		"vendor":   nil,
		"quantity": 0,
		"price":    "0",
		"hidden":   false,
	})
	want := map[string]any{
		"title":    "x",
		"quantity": 0,
		"price":    "0",
		"hidden":   false,
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("key %q: got %v, want %v", k, got[k], v)
		}
	}
}

func TestApplyClearFields_WritesEmptyAndValidates(t *testing.T) {
	payload := map[string]any{"title": "x"}
	if err := applyClearFields(payload, productUpdateFields, []any{"html_title", "meta_keywords"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload["html_title"] != "" || payload["meta_keywords"] != "" {
		t.Errorf("clear fields not applied: %v", payload)
	}
	if err := applyClearFields(payload, productUpdateFields, []any{"id"}); err == nil {
		t.Error("expected error for non-whitelisted clear field")
	}
	if err := applyClearFields(payload, productUpdateFields, []any{"category_id"}); err == nil {
		t.Error("expected error for clear_fields category_id")
	}
	if err := applyClearFields(payload, productUpdateFields, "html_title"); err == nil {
		t.Error("expected error for non-array clear_fields")
	}
}

// captured records the last request seen by an httptest server.
type captured struct {
	req  *http.Request
	body []byte
}

// captureRequest is an httptest server recording the last JSON body and headers.
func captureRequest(t *testing.T) (*httptest.Server, *captured) {
	t.Helper()
	c := &captured{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.req = r
		buf := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buf)
		c.body = buf
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": 1, "title": "t"}`))
	}))
	t.Cleanup(srv.Close)
	return srv, c
}

func TestUpdateProduct_WhitelistAndClearFields(t *testing.T) {
	srv, cap := captureRequest(t)
	client := newTestClient(srv)

	_, err := client.UpdateProduct(context.Background(), "42", map[string]any{
		"product_id":    "42",
		"title":         "New",
		"price":         "777",
		"sku":           "",
		"quantity":      0,
		"is_hidden":     true,
		"archived":      true,
		"variants":      []any{},
		"zzz":           "junk",
		"category_id":   float64(7),
		"clear_fields":  []any{"meta_keywords"},
		"meta_keywords": "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]map[string]any
	if err := json.Unmarshal(cap.body, &payload); err != nil {
		t.Fatalf("bad body %q: %v", cap.body, err)
	}
	p := payload["product"]
	// whitelist: only declared fields; empties dropped; clear_fields wins over drop
	// price is not whitelisted: it silently no-ops on InSales (variants own it)
	for _, k := range []string{"product_id", "price", "quantity", "is_hidden", "archived", "variants", "zzz"} {
		if _, ok := p[k]; ok {
			t.Errorf("field %q must not be sent, got %v", k, p)
		}
	}
	if p["title"] != "New" {
		t.Errorf("title = %v", p["title"])
	}
	if _, ok := p["sku"]; ok {
		t.Errorf("empty sku must be dropped, got %v", p)
	}
	if p["category_id"] != float64(7) {
		t.Errorf("category_id = %v, want 7", p["category_id"])
	}
	if v, ok := p["meta_keywords"]; !ok || v != "" {
		t.Errorf("clear_fields meta_keywords must be sent as empty string, got %v", p)
	}
	if cap.req.Header.Get("Cache-Control") != "no-cache, no-store" {
		t.Errorf("missing no-cache headers: %v", cap.req.Header)
	}
}

func TestCreateProduct_NoExplicitNulls(t *testing.T) {
	srv, cap := captureRequest(t)
	client := newTestClient(srv)

	_, err := client.CreateProduct(context.Background(), "T", "10", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]map[string]any
	if err := json.Unmarshal(cap.body, &payload); err != nil {
		t.Fatalf("bad body %q: %v", cap.body, err)
	}
	p := payload["product"]
	for _, k := range []string{"description", "vendor", "sku", "quantity", "price"} {
		if v, ok := p[k]; ok {
			t.Errorf("absent field %q must not be sent, got %v", k, v)
		}
	}
	if p["title"] != "T" {
		t.Errorf("required fields wrong: %v", p)
	}
	// InSales rejects creates without a variant carrying price/sku/stock.
	vs, ok := p["variants_attributes"].([]any)
	if !ok || len(vs) != 1 {
		t.Fatalf("variants_attributes missing: %v", p)
	}
	v := vs[0].(map[string]any)
	if v["price"] != "10" {
		t.Errorf("variant price wrong: %v", v)
	}
	for _, k := range []string{"sku", "quantity"} {
		if _, ok := v[k]; ok {
			t.Errorf("absent variant field %q must not be sent, got %v", k, v)
		}
	}
}

func TestAddImage_UsesSrcField(t *testing.T) {
	srv, cap := captureRequest(t)
	client := newTestClient(srv)

	_, err := client.AddImage(context.Background(), "123", map[string]any{"url": "https://x/img.png"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]map[string]any
	if err := json.Unmarshal(cap.body, &payload); err != nil {
		t.Fatalf("bad body %q: %v", cap.body, err)
	}
	img := payload["image"]
	if img["src"] != "https://x/img.png" {
		t.Errorf("src missing/wrong: %v", img)
	}
	if _, ok := img["url"]; ok {
		t.Errorf("url must not be sent (InSales wants src): %v", img)
	}
}

func TestFetchPaged_SmallPageWindow(t *testing.T) {
	// 15 items catalog; per_page=3 page=5 must return items 13-15 despite
	// InSales flooring per_page to 10 (we emulate via 250-window fetch).
	var urls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urls = append(urls, r.URL.RequestURI())
		items := []map[string]any{}
		for i := 1; i <= 15; i++ {
			items = append(items, map[string]any{"id": i})
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()
	client := newTestClient(srv)

	got, err := fetchPaged(context.Background(), client, "/products.json", 5, 3, func(b []byte) ([]InSalesProductResponse, error) {
		var v []InSalesProductResponse
		return v, json.Unmarshal(b, &v)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(urls) != 1 || urls[0] != "/products.json?page=1&per_page=250" {
		t.Errorf("window fetch URLs = %v", urls)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 items, got %d", len(got))
	}
	if got[0].ID != 13 || got[2].ID != 15 {
		t.Errorf("window slice wrong: got ids %d..%d, want 13..15", got[0].ID, got[2].ID)
	}
}

func TestFetchPaged_BoundarySpansTwoRequests(t *testing.T) {
	// offset 248 with per_page=3 spans two 250-pages (248..250).
	var urls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urls = append(urls, r.URL.RequestURI())
		start := 0
		if r.URL.Query().Get("page") == "2" {
			start = 250
		}
		items := []map[string]any{}
		for i := start + 1; i <= start+250; i++ {
			items = append(items, map[string]any{"id": i})
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()
	client := newTestClient(srv)

	// page=84, perPage=3 → offset 249 → items 250-252 → spans 250-pages 1 and 2
	got, err := fetchPaged(context.Background(), client, "/products.json?q=x", 84, 3, func(b []byte) ([]InSalesProductResponse, error) {
		var v []InSalesProductResponse
		return v, json.Unmarshal(b, &v)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(urls) != 2 {
		t.Fatalf("want 2 window requests, got %v", urls)
	}
	for _, u := range urls {
		if !strings.Contains(u, "q=x") {
			t.Errorf("filter missing in window request %s", u)
		}
	}
	if len(got) != 3 || got[0].ID != 250 || got[2].ID != 252 {
		t.Errorf("slice wrong: %d..%d", got[0].ID, got[len(got)-1].ID)
	}
}

func TestFetchPaged_NormalPagePassesThrough(t *testing.T) {
	var urls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urls = append(urls, r.URL.RequestURI())
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	client := newTestClient(srv)

	_, err := fetchPaged(context.Background(), client, "/products.json", 2, 25, func(b []byte) ([]InSalesProductResponse, error) {
		var v []InSalesProductResponse
		return v, json.Unmarshal(b, &v)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(urls) != 1 || urls[0] != "/products.json?page=2&per_page=25" {
		t.Errorf("URLs = %v", urls)
	}
}

func TestAddImage_UsesAttachment(t *testing.T) {
	srv, cap := captureRequest(t)
	client := newTestClient(srv)

	_, err := client.AddImage(context.Background(), "123", map[string]any{
		"image_data": "aGVsbG8=",
		"filename":   "cover.png",
		"position":   float64(2),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]map[string]any
	if err := json.Unmarshal(cap.body, &payload); err != nil {
		t.Fatalf("bad body %q: %v", cap.body, err)
	}
	img := payload["image"]
	if img["attachment"] != "aGVsbG8=" {
		t.Errorf("attachment missing/wrong: %v", img)
	}
	if img["filename"] != "cover.png" {
		t.Errorf("filename missing/wrong: %v", img)
	}
	if _, ok := img["src"]; ok {
		t.Errorf("src must not be sent in attachment mode: %v", img)
	}
	if _, ok := img["image_data"]; ok {
		t.Errorf("image_data must not leak into payload: %v", img)
	}
}
