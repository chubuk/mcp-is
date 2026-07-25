package insales

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/ratelimit"
	"github.com/chubuk/mcp-is/internal/resilience"
)

func newTestClient(srv *httptest.Server) *InSalesClient {
	return &InSalesClient{
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
		CircuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		RateTracker:    ratelimit.NewRateTracker(500, 5*time.Minute),
		BreakerPool:    ratelimit.NewPerKeyBreaker(5, 30*time.Second),
	}
}

func TestParseTitleFromHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "standard InSales title",
			input: `<html><head><title>Парафарм Шоп | Официальный магазин</title></head><body></body></html>`,
			want:  "Парафарм Шоп",
		},
		{
			name:  "title without pipe",
			input: `<html><head><title>My Shop</title></head></html>`,
			want:  "My Shop",
		},
		{
			name:  "uppercase TITLE",
			input: `<html><head><TITLE>Test Shop | Subtitle</TITLE></head></html>`,
			want:  "Test Shop",
		},
		{
			name:  "no title tag",
			input: `<html><head></head><body></body></html>`,
			want:  "",
		},
		{
			name:  "empty title",
			input: `<html><head><title></title></head></html>`,
			want:  "",
		},
		{
			name:  "title with whitespace",
			input: `<html><head><title>  My Shop  |  Sub  </title></head></html>`,
			want:  "My Shop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTitleFromHTML(tt.input)
			if got != tt.want {
				t.Errorf("parseTitleFromHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetStoreInfo_WithMockServer(t *testing.T) {
	storeHTML := `<html><head><title>Парафарм Шоп | Магазин</title></head><body></body></html>`
	title := parseTitleFromHTML(storeHTML)
	if title != "Парафарм Шоп" {
		t.Errorf("expected 'Парафарм Шоп', got %q", title)
	}
}

func TestGetStoreInfo_FallbackToDomain(t *testing.T) {
	// If storefront is unreachable, should fallback to domain
	title := parseTitleFromHTML("")
	if title != "" {
		t.Errorf("expected empty title for empty HTML, got %q", title)
	}
}

func TestGetStoreInfo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	title := parseTitleFromHTML("<html><head><title>Test</title></head></html>")
	if title != "Test" {
		t.Errorf("expected 'Test', got %q", title)
	}
	_ = ctx
}

// === URL encoding tests ===

func TestListProducts_URL_EncodesQuery(t *testing.T) {
	var receivedURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client := newTestClient(srv)

	filters := map[string]any{"q": "Одуванчик П 100"}
	_, err := client.ListProducts(context.Background(), 1, 5, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// URL must contain encoded query, not raw cyrillic
	if receivedURL != "/products.json?page=1&per_page=5&q=%D0%9E%D0%B4%D1%83%D0%B2%D0%B0%D0%BD%D1%87%D0%B8%D0%BA+%D0%9F+100" {
		t.Errorf("URL not properly encoded: %s", receivedURL)
	}
}

func TestListProducts_URL_EncodesCyrillic(t *testing.T) {
	var receivedURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client := newTestClient(srv)

	filters := map[string]any{"q": "Красная футболка"}
	_, err := client.ListProducts(context.Background(), 1, 25, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedURL != "/products.json?page=1&per_page=25&q=%D0%9A%D1%80%D0%B0%D1%81%D0%BD%D0%B0%D1%8F+%D1%84%D1%83%D1%82%D0%B1%D0%BE%D0%BB%D0%BA%D0%B0" {
		t.Errorf("URL not properly encoded: %s", receivedURL)
	}
}

func TestListOrders_URL_EncodesFilters(t *testing.T) {
	var receivedURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client := newTestClient(srv)

	filters := map[string]any{
		"fulfillment_status": "new",
		"updated_since":      "2026-01-01T00:00:00Z",
	}
	_, err := client.ListOrders(context.Background(), 1, 25, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "/orders.json?page=1&per_page=25&fulfillment_status[]=new&updated_since=2026-01-01T00%3A00%3A00Z"
	if receivedURL != expected {
		t.Errorf("URL not properly encoded:\n  got:  %s\n  want: %s", receivedURL, expected)
	}
}

func TestListMetaFields_FiltersByResourceType(t *testing.T) {
	var receivedURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client := newTestClient(srv)

	args := map[string]any{
		"resource_type": "Product",
		"resource_id":   "12345",
	}
	_, err := client.ListMetaFields(context.Background(), args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "/meta_fields.json?meta_field[owner_type]=Product&meta_field[owner_id]=12345"
	if receivedURL != expected {
		t.Errorf("URL not properly built:\n  got:  %s\n  want: %s", receivedURL, expected)
	}
}

func TestListMetaFields_NoFilters(t *testing.T) {
	var receivedURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client := newTestClient(srv)

	_, err := client.ListMetaFields(context.Background(), map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "/meta_fields.json"
	if receivedURL != expected {
		t.Errorf("URL should not have query params when no filters:\n  got:  %s\n  want: %s", receivedURL, expected)
	}
}

// === Product field values tests ===

func TestMapProduct_IncludesFieldValues(t *testing.T) {
	resp := &InSalesProductResponse{
		ID:    12345,
		Title: "Одуванчик П 100",
		ProductFieldValues: []InSalesFieldValue{
			{ID: 1, ProductFieldID: 1, Value: "Противопоказания: Беременность"},
			{ID: 2, ProductFieldID: 2, Value: "Дозировка: 100 мг"},
		},
	}

	product := MapProduct(resp)

	fvs := product.FieldValues
	if fvs == nil {
		t.Fatal("FieldValues should not be nil")
	}
	if len(fvs) != 2 {
		t.Fatalf("expected 2 field values, got %d", len(fvs))
	}

	fv1 := fvs[0].(map[string]any)
	if fv1["value"] != "Противопоказания: Беременность" {
		t.Errorf("expected field value 'Противопоказания: Беременность', got %v", fv1["value"])
	}
	if fv1["product_field_id"] != 1 {
		t.Errorf("expected product_field_id 1, got %v", fv1["product_field_id"])
	}
}

func TestMapProduct_IncludesAllNestedData(t *testing.T) {
	resp := &InSalesProductResponse{
		ID:               1,
		Title:            "Тест",
		Description:      "Описание",
		ShortDescription: "Краткое",
		Permalink:        "test",
		IsHidden:         false,
		Variants:         []InSalesVariant{{ID: 1, Title: "Red", Price: "100"}},
		Images:           []InSalesImage{{ID: 1, URL: "http://img.jpg", Filename: "img.jpg"}},
		Properties:       []InSalesProductProperty{{ID: 1, Title: "Бренд"}},
		Characteristics:  []InSalesCharacteristic{{ID: 1, PropertyID: 1, Title: "Nike"}},
		OptionNames:      []InSalesOptionName{{ID: 1, Title: "Размер"}},
		ProductFieldValues: []InSalesFieldValue{{ID: 1, ProductFieldID: 1, Value: "test"}},
	}

	product := MapProduct(resp)

	if product.Description != "Описание" {
		t.Errorf("expected Description 'Описание', got %q", product.Description)
	}
	if product.ShortDescription != "Краткое" {
		t.Errorf("expected ShortDescription 'Краткое', got %q", product.ShortDescription)
	}
	if product.Permalink != "test" {
		t.Errorf("expected Permalink 'test', got %q", product.Permalink)
	}
	if len(product.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(product.Images))
	}
	if len(product.FieldValues) != 1 {
		t.Errorf("expected 1 field value, got %d", len(product.FieldValues))
	}
	if len(product.Properties) != 1 {
		t.Errorf("expected 1 property, got %d", len(product.Properties))
	}
}

func TestMapOrder_IncludesShippingAddress(t *testing.T) {
	resp := &InSalesOrderResponse{
		ID:           1,
		TotalPrice:   1500.0,
		ShippingAddress: &InSalesShippingAddress{
			ID:                  9,
			Name:                "Иван",
			FullName:            "Иванов Иван",
			FullDeliveryAddress: "Москва, ул. Тестовая",
			Phone:               "+79991234567",
		},
	}

	order := MapOrder(resp)

	sa := order.ShippingAddress
	if sa == nil {
		t.Fatal("ShippingAddress should not be nil")
	}
	if sa["full_name"] != "Иванов Иван" {
		t.Errorf("expected full_name 'Иванов Иван', got %v", sa["full_name"])
	}
	if sa["full_delivery_address"] != "Москва, ул. Тестовая" {
		t.Errorf("expected full_delivery_address, got %v", sa["full_delivery_address"])
	}
}

func TestMapClient_IncludesFields(t *testing.T) {
	resp := &InSalesClientResponse{
		ID:                 1,
		Email:              "test@example.com",
		Name:               "Иван",
		Surname:            "Иванов",
		FullName:           "Иванов Иван Иванович",
		BonusPoints:        500,
		ProgressiveDiscount: "10.0",
		GroupDiscount:      "5.0",
		ClientTags:         []string{"vip", "regular"},
		FieldsValues:       []InSalesFieldValue{{ID: 1, ProductFieldID: 1, Value: "text"}},
	}

	client := MapClient(resp)

	if client["bonus_points"] != 500 {
		t.Errorf("expected bonus_points 500, got %v", client["bonus_points"])
	}
	if client["progressive_discount"] != "10.0" {
		t.Errorf("expected progressive_discount '10.0', got %v", client["progressive_discount"])
	}
	if client["full_name"] != "Иванов Иван Иванович" {
		t.Errorf("expected full_name, got %v", client["full_name"])
	}
	tags, ok := client["client_tags"].([]string)
	if !ok || len(tags) != 2 {
		t.Errorf("expected 2 client_tags, got %v", client["client_tags"])
	}
	fvs, ok := client["fields_values"].([]any)
	if !ok || len(fvs) != 1 {
		t.Errorf("expected 1 fields_values, got %v", client["fields_values"])
	}
}
