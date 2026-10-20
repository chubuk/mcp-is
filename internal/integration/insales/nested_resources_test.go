package insales

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/ratelimit"
	"github.com/chubuk/mcp-is/internal/resilience"
)

func testClientWithCapture(t *testing.T, handler http.HandlerFunc) (*InSalesClient, *[]capturedReq) {
	t.Helper()
	var reqs []capturedReq
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var parsed any
		if len(body) > 0 {
			_ = json.Unmarshal(body, &parsed)
		}
		reqs = append(reqs, capturedReq{Method: r.Method, Path: r.URL.Path, Body: parsed})
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	c := &InSalesClient{
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
		CircuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		RateTracker:    ratelimit.NewRateTracker(500, 5*time.Minute),
		BreakerPool:    ratelimit.NewPerKeyBreaker(5, 30*time.Second),
	}
	return c, &reqs
}

type capturedReq struct {
	Method string
	Path   string
	Body   any
}

func jsonResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func TestCreateProductFieldValue_Payload(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 201, `{"id": 7, "product_field_id": 2, "value": "100 mg"}`)
	})
	_, err := c.CreateProductFieldValue(context.Background(), "12", map[string]any{
		"product_field_handle": "dosage",
		"value":                "100 mg",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(*reqs) != 1 {
		t.Fatalf("requests: %#v", *reqs)
	}
	r := (*reqs)[0]
	if r.Method != http.MethodPost || r.Path != "/products/12/product_field_values.json" {
		t.Errorf("got %s %s", r.Method, r.Path)
	}
	body := r.Body.(map[string]any)
	pfv := body["product_field_value"].(map[string]any)
	if pfv["product_field_handle"] != "dosage" || pfv["value"] != "100 mg" {
		t.Errorf("body=%#v", body)
	}
}

func TestUpdateProductFieldValue_ByHandle_Payload(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, `{"id": 1, "value": "new"}`)
	})
	_, err := c.UpdateProductFieldValue(context.Background(), "12", "dosage", map[string]any{"value": "new"})
	if err != nil {
		t.Fatal(err)
	}
	r := (*reqs)[0]
	if r.Method != http.MethodPut || r.Path != "/products/12/product_field_values/dosage.json" {
		t.Errorf("got %s %s", r.Method, r.Path)
	}
	body := r.Body.(map[string]any)
	if body["product_field_value"].(map[string]any)["value"] != "new" {
		t.Errorf("body=%#v", body)
	}
}

func TestListMenuItems_Path(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, `[{"id": 1, "title": "About", "related_object_id": "Page_2"}]`)
	})
	items, err := c.ListMenuItems(context.Background(), "3")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("items=%#v", items)
	}
	if (*reqs)[0].Path != "/menus/3/items.json" {
		t.Errorf("path=%s", (*reqs)[0].Path)
	}
}

func TestCreateMenuItem_RelatedObjectComposite(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 201, `{"id": 9, "related_object_id": "Blog_1"}`)
	})
	_, err := c.CreateMenuItem(context.Background(), "3", map[string]any{
		"title":             "News",
		"related_object_id": "Blog_1",
	})
	if err != nil {
		t.Fatal(err)
	}
	r := (*reqs)[0]
	if r.Path != "/menus/3/items.json" {
		t.Errorf("path=%s", r.Path)
	}
	mi := r.Body.(map[string]any)["menu_item"].(map[string]any)
	if mi["related_object_id"] != "Blog_1" {
		t.Errorf("body=%#v", mi)
	}
}

func TestCreateDeliveryVariant_TypePayload(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 201, `{"id": 4, "type": "DeliveryVariant::Fixed"}`)
	})
	_, err := c.CreateDeliveryVariant(context.Background(), map[string]any{
		"title": "Courier",
		"type":  "DeliveryVariant::Fixed",
	})
	if err != nil {
		t.Fatal(err)
	}
	r := (*reqs)[0]
	if r.Method != http.MethodPost || r.Path != "/delivery_variants.json" {
		t.Errorf("got %s %s", r.Method, r.Path)
	}
	dv := r.Body.(map[string]any)["delivery_variant"].(map[string]any)
	if dv["type"] != "DeliveryVariant::Fixed" {
		t.Errorf("body=%#v", r.Body)
	}
}

func TestCreatePaymentGateway_TypePayload(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 201, `{"id": 2, "type": "PaymentGateway::Cod"}`)
	})
	_, err := c.CreatePaymentGateway(context.Background(), map[string]any{
		"title": "Cash",
		"type":  "PaymentGateway::Cod",
	})
	if err != nil {
		t.Fatal(err)
	}
	r := (*reqs)[0]
	if r.Method != http.MethodPost || r.Path != "/payment_gateways.json" {
		t.Errorf("got %s %s", r.Method, r.Path)
	}
	pg := r.Body.(map[string]any)["payment_gateway"].(map[string]any)
	if pg["type"] != "PaymentGateway::Cod" {
		t.Errorf("body=%#v", r.Body)
	}
}

func TestListProductAccessoryLinks_Path(t *testing.T) {
	c, reqs := testClientWithCapture(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, `[]`)
	})
	_, err := c.ListProductAccessoryLinks(context.Background(), "55")
	if err != nil {
		t.Fatal(err)
	}
	if (*reqs)[0].Path != "/products/55/product_accessory_links.json" {
		t.Errorf("path=%s", (*reqs)[0].Path)
	}
}
