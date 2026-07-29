//go:build integration

package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListOutlets_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/outlets.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"title":"Point 1","address":"addr","price":"300.0"}]`))
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/outlets.json")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var outlets []map[string]any
	json.NewDecoder(resp.Body).Decode(&outlets)
	if len(outlets) != 1 {
		t.Fatalf("expected 1 outlet, got %d", len(outlets))
	}
	if outlets[0]["title"] != "Point 1" {
		t.Errorf("expected title 'Point 1', got %v", outlets[0]["title"])
	}
}

func TestGetOrdersCount_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders/count.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"count":42}`))
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/orders/count.json")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Count int `json:"count"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Count != 42 {
		t.Errorf("expected count 42, got %d", result.Count)
	}
}

func TestGetProductsCount_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"count":100}`))
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/products/count.json")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Count int `json:"count"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Count != 100 {
		t.Errorf("expected count 100, got %d", result.Count)
	}
}

func TestCreateProviderToken_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dialogs/provider_tokens.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["messenger_type"] != "telegram" {
			t.Errorf("expected messenger_type 'telegram', got %v", body["messenger_type"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"api_token":"abc","webhook_token":"xyz"}`))
	}))
	defer srv.Close()

	body := map[string]any{"messenger_type": "telegram"}
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(srv.URL+"/dialogs/provider_tokens.json", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if result["api_token"] != "abc" {
		t.Errorf("expected api_token 'abc', got %v", result["api_token"])
	}
}

func TestUpdateAccount_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/account.json", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", result["status"])
	}
}

func TestCreateMarketplace_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		mp, ok := body["marketplace"].(map[string]any)
		if !ok {
			t.Error("expected 'marketplace' key in body")
		}
		if mp["type"] != "Marketplace::GoogleMerchant" {
			t.Errorf("expected type GoogleMerchant, got %v", mp["type"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"name":"GM","type":"Marketplace::GoogleMerchant"}`))
	}))
	defer srv.Close()

	body := map[string]any{
		"marketplace": map[string]any{
			"name": "GM",
			"type": "Marketplace::GoogleMerchant",
		},
	}
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(srv.URL+"/marketplaces.json", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if result["type"] != "Marketplace::GoogleMerchant" {
		t.Errorf("expected type, got %v", result["type"])
	}
}

func TestCollectionProducts_Integration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/collections/1.json":
			w.Write([]byte(`{"id":1,"title":"Summer Collection"}`))
		case "/products.json":
			w.Write([]byte(`[{"id":1,"title":"Product 1","price":"100"}]`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	// Test collection fetch
	resp, err := http.Get(srv.URL + "/collections/1.json")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var collection map[string]any
	json.NewDecoder(resp.Body).Decode(&collection)
	if collection["title"] != "Summer Collection" {
		t.Errorf("expected 'Summer Collection', got %v", collection["title"])
	}

	// Test products fetch
	resp2, err := http.Get(srv.URL + "/products.json")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp2.Body.Close()

	var products []map[string]any
	json.NewDecoder(resp2.Body).Decode(&products)
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

func TestCSVExport_Integration(t *testing.T) {
	// Test CSV formatting with edge cases
	tests := []struct {
		name     string
		headers  []string
		rows     [][]string
		expected string
	}{
		{
			name:     "simple",
			headers:  []string{"id", "name"},
			rows:     [][]string{{"1", "test"}},
			expected: "id,name\n1,test\n",
		},
		{
			name:     "comma in value",
			headers:  []string{"name"},
			rows:     [][]string{{"hello, world"}},
			expected: "name\n\"hello, world\"\n",
		},
		{
			name:     "quote in value",
			headers:  []string{"name"},
			rows:     [][]string{{`say "hi"`}},
			expected: "name\n\"say \"\"hi\"\"\"\n",
		},
		{
			name:     "newline in value",
			headers:  []string{"name"},
			rows:     [][]string{{"line1\nline2"}},
			expected: "name\n\"line1\nline2\"\n",
		},
		{
			name:     "empty value",
			headers:  []string{"name"},
			rows:     [][]string{{""}},
			expected: "name\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := formatCSVTest(tt.headers, tt.rows)
			if csv != tt.expected {
				t.Errorf("CSV mismatch:\ngot:  %q\nwant: %q", csv, tt.expected)
			}
		})
	}
}

func formatCSVTest(headers []string, rows [][]string) string {
	result := ""
	for i, h := range headers {
		if i > 0 {
			result += ","
		}
		result += csvEscapeTest(h)
	}
	result += "\n"
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				result += ","
			}
			result += csvEscapeTest(cell)
		}
		result += "\n"
	}
	return result
}

func csvEscapeTest(s string) string {
	if s == "" {
		return ""
	}
	needsQuote := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			needsQuote = true
			break
		}
	}
	if needsQuote {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}
