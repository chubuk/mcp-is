package insales

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chubuk/mcp-is/internal/ratelimit"
	"github.com/chubuk/mcp-is/internal/resilience"
)

type InSalesClient struct {
	BaseURL        string
	ShopDomain     string
	APIKey         string
	Password       string
	HTTPClient     *http.Client
	CircuitBreaker *resilience.CircuitBreaker
	RateTracker    *ratelimit.RateTracker
	BreakerPool    *ratelimit.PerKeyBreaker
}

func NewClient(shopDomain, apiKey, password string) (*InSalesClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key cannot be empty")
	}
	if shopDomain == "" {
		return nil, fmt.Errorf("shop domain cannot be empty")
	}
	return &InSalesClient{
		BaseURL:    fmt.Sprintf("https://%s/admin", shopDomain),
		ShopDomain: shopDomain,
		APIKey:     apiKey,
		Password:   password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		CircuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		RateTracker:    ratelimit.NewRateTracker(500, 5*time.Minute),
		BreakerPool:    ratelimit.NewPerKeyBreaker(5, 30*time.Second),
	}, nil
}

func (c *InSalesClient) doRequest(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	breaker := c.BreakerPool.GetBreaker(c.APIKey)
	maxRetries := 3

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Throttle if approaching rate limit
		delay := c.RateTracker.GetThrottleDelay(c.APIKey)
		if delay > 0 {
			slog.Debug("throttling request", "delay", delay, "shop", c.ShopDomain, "attempt", attempt)
			time.Sleep(delay)
		}

		// Check circuit breaker
		if !breaker.Allow() {
			slog.Warn("circuit breaker open", "shop", c.ShopDomain, "path", path)
			return nil, 0, fmt.Errorf("service temporarily unavailable (circuit breaker open)")
		}

		result, statusCode, err := c.doRequestRaw(ctx, method, path, body)

		if err == nil {
			breaker.RecordSuccess()
			return result, statusCode, nil
		}

		// 429: retry with Retry-After
		if statusCode == 429 {
			retryAfter := parseRetryAfter(err.Error())
			slog.Warn("rate limited by InSales", "retry_after", retryAfter, "attempt", attempt+1, "shop", c.ShopDomain)
			breaker.RecordFailure()
			time.Sleep(retryAfter)
			continue
		}

		// Other errors: don't retry
		breaker.RecordFailure()
		return result, statusCode, err
	}

	return nil, 0, fmt.Errorf("rate limited after %d retries", maxRetries)
}

func parseRetryAfter(errMsg string) time.Duration {
	// Extract seconds from "rate limited (429), retry after 60"
	for i := len(errMsg) - 1; i >= 0; i-- {
		if errMsg[i] == ' ' {
			numStr := errMsg[i+1:]
			if n, err := strconv.Atoi(numStr); err == nil && n > 0 {
				return time.Duration(n) * time.Second
			}
			break
		}
	}
	return 60 * time.Second // default
}

func (c *InSalesClient) doRequestRaw(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.APIKey, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	// Record request for rate tracking
	c.RateTracker.RecordRequest(c.APIKey)

	slog.Debug("insales api response", "method", method, "path", path, "status", resp.StatusCode)

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		return nil, resp.StatusCode, fmt.Errorf("rate limited (429), retry after %s", retryAfter)
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}

type AccountInfo struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Subdomain    string `json:"subdomain"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	ContactPhone string `json:"contact_phone"`
}

func (c *InSalesClient) GetStoreInfo(ctx context.Context) (*StoreInfo, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/account.json", nil)
	if err != nil {
		// Fallback: try validating credentials via products endpoint
		_, status, err2 := c.doRequest(ctx, http.MethodGet, "/products.json?per_page=1", nil)
		if err2 == nil || status == 404 || status == 403 {
			return c.buildStoreInfoFallback(ctx), nil
		}
		return nil, fmt.Errorf("connection failed — check domain and credentials")
	}

	var account AccountInfo
	if err := json.Unmarshal(body, &account); err != nil {
		return c.buildStoreInfoFallback(ctx), nil
	}

	name := c.ShopDomain
	if account.Title != "" {
		name = account.Title
	} else if title := c.fetchStoreTitle(ctx); title != "" {
		name = title
	}

	return &StoreInfo{
		Name:     name,
		Domain:   c.ShopDomain,
		Email:    account.Email,
		Currency: "RUB",
	}, nil
}

func (c *InSalesClient) buildStoreInfoFallback(ctx context.Context) *StoreInfo {
	name := c.ShopDomain
	if title := c.fetchStoreTitle(ctx); title != "" {
		name = title
	}
	return &StoreInfo{
		Name:     name,
		Domain:   c.ShopDomain,
		Email:    "",
		Currency: "RUB",
	}
}

func (c *InSalesClient) fetchStoreTitle(ctx context.Context) string {
	url := fmt.Sprintf("https://%s/", c.ShopDomain)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "MCP-IS/1.0")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return ""
	}

	return parseTitleFromHTML(string(body))
}

func parseTitleFromHTML(html string) string {
	// Find <title>...</title>
	start := strings.Index(html, "<title>")
	if start == -1 {
		start = strings.Index(html, "<TITLE>")
	}
	if start == -1 {
		return ""
	}
	start += 7 // len("<title>")
	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		end = strings.Index(html[start:], "</TITLE>")
	}
	if end == -1 {
		return ""
	}
	title := strings.TrimSpace(html[start : start+end])
	// Take part before " | " if present (InSales format: "Name | Subtitle")
	if idx := strings.Index(title, " | "); idx > 0 {
		title = strings.TrimSpace(title[:idx])
	}
	return title
}

func (c *InSalesClient) ListProducts(ctx context.Context, page, perPage int, filters map[string]any) ([]Product, error) {
	path := fmt.Sprintf("/products.json?page=%d&per_page=%d", page, perPage)
	if v, ok := filters["category_id"].(string); ok && v != "" {
		path += "&category_id=" + url.QueryEscape(v)
	}
	if v, ok := filters["collection_id"].(string); ok && v != "" {
		path += "&collection_id=" + url.QueryEscape(v)
	}
	if v, ok := filters["q"].(string); ok && v != "" {
		path += "&q=" + url.QueryEscape(v)
	}
	if v, ok := filters["updated_since"].(string); ok && v != "" {
		path += "&updated_since=" + url.QueryEscape(v)
	}
	if v, ok := filters["from_id"].(float64); ok && v > 0 {
		path += fmt.Sprintf("&from_id=%d", int(v))
	}
	if v, ok := filters["with_deleted"].(bool); ok && v {
		path += "&with_deleted=true"
	}
	if v, ok := filters["deleted"].(bool); ok && v {
		path += "&deleted=true"
	}
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesProductResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode products: %w", err)
	}
	result := make([]Product, len(items))
	for i, item := range items {
		result[i] = *MapProduct(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetProduct(ctx context.Context, productID string) (*Product, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/products/%s.json", productID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesProductResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode product: %w", err)
	}
	return MapProduct(&resp), nil
}

func (c *InSalesClient) CreateProduct(ctx context.Context, title, price string, args map[string]any) (*Product, error) {
	payload := map[string]any{
		"product": map[string]any{
			"title":       title,
			"price":       price,
			"description": args["description"],
			"vendor":      args["vendor"],
			"sku":         args["sku"],
			"quantity":    args["quantity"],
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/products.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesProductResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode product: %w", err)
	}
	return MapProduct(&resp), nil
}

func (c *InSalesClient) UpdateProduct(ctx context.Context, productID string, args map[string]any) (*Product, error) {
	payload := map[string]any{
		"product": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/products/%s.json", productID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesProductResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode product: %w", err)
	}
	return MapProduct(&resp), nil
}

func (c *InSalesClient) DeleteProduct(ctx context.Context, productID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/products/%s.json", productID), nil)
	return err
}

func (c *InSalesClient) ListOrders(ctx context.Context, page, perPage int, filters map[string]any) ([]Order, error) {
	path := fmt.Sprintf("/orders.json?page=%d&per_page=%d", page, perPage)
	if v, ok := filters["fulfillment_status"].(string); ok && v != "" {
		path += "&fulfillment_status[]=" + url.QueryEscape(v)
	}
	if v, ok := filters["updated_since"].(string); ok && v != "" {
		path += "&updated_since=" + url.QueryEscape(v)
	}
	if v, ok := filters["from_id"].(int); ok && v > 0 {
		path += fmt.Sprintf("&from_id=%d", v)
	}
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesOrderResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode orders: %w", err)
	}
	result := make([]Order, len(items))
	for i, item := range items {
		result[i] = *MapOrder(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/orders/%s.json", orderID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode order: %w", err)
	}
	return MapOrder(&resp), nil
}

func (c *InSalesClient) UpdateOrder(ctx context.Context, orderID string, args map[string]any) (*Order, error) {
	payload := map[string]any{
		"order": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/orders/%s.json", orderID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode order: %w", err)
	}
	return MapOrder(&resp), nil
}

func (c *InSalesClient) CreateOrder(ctx context.Context, clientID int, args map[string]any) (*Order, error) {
	payload := map[string]any{
		"order": map[string]any{
			"client_id": clientID,
			"comment":   args["comment"],
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/orders.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode order: %w", err)
	}
	return MapOrder(&resp), nil
}

func (c *InSalesClient) ListClients(ctx context.Context, page, perPage int) ([]any, error) {
	path := fmt.Sprintf("/clients.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesClientResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode clients: %w", err)
	}
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = MapClient(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetClient(ctx context.Context, clientID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/clients/%s.json", clientID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesClientResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode client: %w", err)
	}
	return MapClient(&resp), nil
}

func (c *InSalesClient) CreateClient(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{
		"client": map[string]any{
			"name":       args["name"],
			"surname":    args["surname"],
			"middlename": args["middlename"],
			"email":      args["email"],
			"phone":      args["phone"],
			"password":   args["password"],
			"registered": args["registered"],
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/clients.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesClientResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode client: %w", err)
	}
	return MapClient(&resp), nil
}

func (c *InSalesClient) UpdateClient(ctx context.Context, clientID string, args map[string]any) (any, error) {
	payload := map[string]any{
		"client": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/clients/%s.json", clientID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesClientResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode client: %w", err)
	}
	return MapClient(&resp), nil
}

func (c *InSalesClient) DeleteClient(ctx context.Context, clientID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/clients/%s.json", clientID), nil)
	return err
}

func (c *InSalesClient) ListCategories(ctx context.Context, page, perPage int) ([]any, error) {
	path := fmt.Sprintf("/categories.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCategoryResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode categories: %w", err)
	}
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = MapCategory(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetCategory(ctx context.Context, categoryID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/categories/%s.json", categoryID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesCategoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode category: %w", err)
	}
	return MapCategory(&resp), nil
}

func (c *InSalesClient) CreateCategory(ctx context.Context, title string, args map[string]any) (any, error) {
	payload := map[string]any{
		"category": map[string]any{
			"title":     title,
			"parent_id": args["parent_id"],
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/categories.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCategoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode category: %w", err)
	}
	return MapCategory(&resp), nil
}

func (c *InSalesClient) UpdateCategory(ctx context.Context, categoryID string, args map[string]any) (any, error) {
	payload := map[string]any{
		"category": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/categories/%s.json", categoryID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCategoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode category: %w", err)
	}
	return MapCategory(&resp), nil
}

func (c *InSalesClient) DeleteCategory(ctx context.Context, categoryID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/categories/%s.json", categoryID), nil)
	return err
}

func (c *InSalesClient) ListCollections(ctx context.Context, page, perPage int) ([]any, error) {
	path := fmt.Sprintf("/collections.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCollectionResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode collections: %w", err)
	}
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = MapCollection(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetCollection(ctx context.Context, collectionID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/collections/%s.json", collectionID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesCollectionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode collection: %w", err)
	}
	return MapCollection(&resp), nil
}

func (c *InSalesClient) CreateCollection(ctx context.Context, title string) (any, error) {
	payload := map[string]any{
		"collection": map[string]any{
			"title": title,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/collections.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCollectionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode collection: %w", err)
	}
	return MapCollection(&resp), nil
}

func (c *InSalesClient) UpdateCollection(ctx context.Context, collectionID, title string) (any, error) {
	payload := map[string]any{
		"collection": map[string]any{
			"title": title,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/collections/%s.json", collectionID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCollectionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode collection: %w", err)
	}
	return MapCollection(&resp), nil
}

func (c *InSalesClient) DeleteCollection(ctx context.Context, collectionID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/collections/%s.json", collectionID), nil)
	return err
}

func (c *InSalesClient) ListImages(ctx context.Context, productID string) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/products/%s/images.json", productID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesImage
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode images: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = items[i]
	}
	return result, nil
}

func (c *InSalesClient) AddImage(ctx context.Context, productID, url string, args map[string]any) (any, error) {
	payload := map[string]any{
		"image": map[string]any{
			"url":      url,
			"filename": args["filename"],
			"position": args["position"],
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/products/%s/images.json", productID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesImage
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	return resp, nil
}

func (c *InSalesClient) DeleteImage(ctx context.Context, productID, imageID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/products/%s/images/%s.json", productID, imageID), nil)
	return err
}

func (c *InSalesClient) ListMetaFields(ctx context.Context, args map[string]any) ([]any, error) {
	path := "/meta_fields.json"
	params := []string{}
	if v, ok := args["resource_type"].(string); ok && v != "" {
		params = append(params, "meta_field[owner_type]="+url.QueryEscape(v))
	}
	if v, ok := args["resource_id"].(string); ok && v != "" {
		params = append(params, "meta_field[owner_id]="+url.QueryEscape(v))
	} else if v, ok := args["resource_id"].(float64); ok && v > 0 {
		params = append(params, fmt.Sprintf("meta_field[owner_id]=%d", int(v)))
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesMetaFieldResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode meta fields: %w", err)
	}
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = MapMetaField(&item)
	}
	return result, nil
}

func (c *InSalesClient) CreateMetaField(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{
		"meta_field": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/meta_fields.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesMetaFieldResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode meta field: %w", err)
	}
	return MapMetaField(&resp), nil
}

func (c *InSalesClient) UpdateMetaField(ctx context.Context, metaFieldID string, args map[string]any) (any, error) {
	payload := map[string]any{
		"meta_field": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/meta_fields/%s.json", metaFieldID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesMetaFieldResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode meta field: %w", err)
	}
	return MapMetaField(&resp), nil
}

func (c *InSalesClient) DeleteMetaField(ctx context.Context, metaFieldID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/meta_fields/%s.json", metaFieldID), nil)
	return err
}

func (c *InSalesClient) ListWebhooks(ctx context.Context) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/webhooks.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesWebhookResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode webhooks: %w", err)
	}
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = MapWebhook(&item)
	}
	return result, nil
}

func (c *InSalesClient) CreateWebhook(ctx context.Context, event, url string, args map[string]any) (any, error) {
	format, _ := args["format"].(string)
	if format == "" {
		format = "json"
	}
	payload := map[string]any{
		"webhook": map[string]any{
			"event":  event,
			"url":    url,
			"format": format,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/webhooks.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesWebhookResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode webhook: %w", err)
	}
	return MapWebhook(&resp), nil
}

func (c *InSalesClient) UpdateWebhook(ctx context.Context, webhookID string, args map[string]any) (any, error) {
	payload := map[string]any{
		"webhook": args,
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/webhooks/%s.json", webhookID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesWebhookResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode webhook: %w", err)
	}
	return MapWebhook(&resp), nil
}

func (c *InSalesClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/webhooks/%s.json", webhookID), nil)
	return err
}

func (c *InSalesClient) ListCurrencies(ctx context.Context) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/currencies.json", nil)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode currencies: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) ListDeliveryVariants(ctx context.Context) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/delivery_variants.json", nil)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode delivery variants: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) ListPaymentGateways(ctx context.Context) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/payment_gateways.json", nil)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode payment gateways: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) ListPages(ctx context.Context, page, perPage int) ([]Page, error) {
	path := fmt.Sprintf("/pages.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var raw []InSalesPageResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode pages: %w", err)
	}
	result := make([]Page, len(raw))
	for i, p := range raw {
		result[i] = *MapPage(&p)
	}
	return result, nil
}

func (c *InSalesClient) GetPage(ctx context.Context, pageID string) (*Page, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/pages/%s.json", pageID), nil)
	if err != nil {
		return nil, err
	}
	var raw InSalesPageResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}
	return MapPage(&raw), nil
}

func (c *InSalesClient) CreatePage(ctx context.Context, title, content, handle string) (*Page, error) {
	payload := map[string]any{
		"page": map[string]any{
			"title":   title,
			"content": content,
			"handle":  handle,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/pages.json", payload)
	if err != nil {
		return nil, err
	}
	var raw InSalesPageResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}
	return MapPage(&raw), nil
}

func (c *InSalesClient) UpdatePage(ctx context.Context, pageID string, args map[string]any) (*Page, error) {
	payload := map[string]any{"page": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/pages/%s.json", pageID), payload)
	if err != nil {
		return nil, err
	}
	var raw InSalesPageResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}
	return MapPage(&raw), nil
}

func (c *InSalesClient) DeletePage(ctx context.Context, pageID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/pages/%s.json", pageID), nil)
	return err
}

func (c *InSalesClient) ListAccessories(ctx context.Context, page, perPage int) ([]Accessory, error) {
	path := fmt.Sprintf("/accessories.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var raw []InSalesAccessoryResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode accessories: %w", err)
	}
	result := make([]Accessory, len(raw))
	for i, a := range raw {
		result[i] = *MapAccessory(&a)
	}
	return result, nil
}

func (c *InSalesClient) GetAccessory(ctx context.Context, accessoryID string) (*Accessory, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/accessories/%s.json", accessoryID), nil)
	if err != nil {
		return nil, err
	}
	var raw InSalesAccessoryResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode accessory: %w", err)
	}
	return MapAccessory(&raw), nil
}

func (c *InSalesClient) CreateAccessory(ctx context.Context, title string, args map[string]any) (*Accessory, error) {
	payload := map[string]any{
		"accessory": map[string]any{
			"name":      title,
			"max_count": args["max_count"],
			"min_count": args["min_count"],
			"permalink": args["permalink"],
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/accessories.json", payload)
	if err != nil {
		return nil, err
	}
	var raw InSalesAccessoryResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode accessory: %w", err)
	}
	return MapAccessory(&raw), nil
}

func (c *InSalesClient) UpdateAccessory(ctx context.Context, accessoryID string, args map[string]any) (*Accessory, error) {
	payload := map[string]any{"accessory": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/accessories/%s.json", accessoryID), payload)
	if err != nil {
		return nil, err
	}
	var raw InSalesAccessoryResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode accessory: %w", err)
	}
	return MapAccessory(&raw), nil
}

func (c *InSalesClient) DeleteAccessory(ctx context.Context, accessoryID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/accessories/%s.json", accessoryID), nil)
	return err
}

// === ВАРИАНТЫ ===

func (c *InSalesClient) ListVariants(ctx context.Context, productID string) ([]Variant, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/products/%s/variants.json", productID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesVariant
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode variants: %w", err)
	}
	return mapVariants(items), nil
}

func (c *InSalesClient) GetVariant(ctx context.Context, productID, variantID string) (*Variant, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/products/%s/variants/%s.json", productID, variantID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesVariant
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode variant: %w", err)
	}
	result := mapVariants([]InSalesVariant{resp})
	return &result[0], nil
}

func (c *InSalesClient) UpdateVariant(ctx context.Context, productID, variantID string, args map[string]any) (*Variant, error) {
	payload := map[string]any{"variant": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/products/%s/variants/%s.json", productID, variantID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesVariant
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode variant: %w", err)
	}
	result := mapVariants([]InSalesVariant{resp})
	return &result[0], nil
}

func (c *InSalesClient) DeleteVariant(ctx context.Context, productID, variantID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/products/%s/variants/%s.json", productID, variantID), nil)
	return err
}

// === ЗНАЧЕНИЯ СВОЙСТВ ===

func (c *InSalesClient) ListPropertyValues(ctx context.Context, accessoryID string) ([]PropertyValue, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/accessories/%s/values.json", accessoryID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesPropertyValueResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode property values: %w", err)
	}
	result := make([]PropertyValue, len(items))
	for i, item := range items {
		result[i] = *MapPropertyValue(&item)
	}
	return result, nil
}

func (c *InSalesClient) CreatePropertyValue(ctx context.Context, accessoryID, title string) (*PropertyValue, error) {
	payload := map[string]any{
		"accessory_value": map[string]any{
			"name":  title,
			"price": 0,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/accessories/%s/values.json", accessoryID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesPropertyValueResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode property value: %w", err)
	}
	return MapPropertyValue(&resp), nil
}

func (c *InSalesClient) UpdatePropertyValue(ctx context.Context, accessoryID, propertyValueID string, args map[string]any) (*PropertyValue, error) {
	payload := map[string]any{"accessory_value": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/accessories/%s/values/%s.json", accessoryID, propertyValueID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesPropertyValueResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode property value: %w", err)
	}
	return MapPropertyValue(&resp), nil
}

func (c *InSalesClient) DeletePropertyValue(ctx context.Context, accessoryID, propertyValueID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/accessories/%s/values/%s.json", accessoryID, propertyValueID), nil)
	return err
}

// === МЕНЮ ===

func (c *InSalesClient) ListMenus(ctx context.Context) ([]Menu, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/menus.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesMenuResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode menus: %w", err)
	}
	result := make([]Menu, len(items))
	for i, item := range items {
		result[i] = *MapMenu(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetMenu(ctx context.Context, menuID string) (*Menu, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/menus/%s.json", menuID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesMenuResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode menu: %w", err)
	}
	return MapMenu(&resp), nil
}

func (c *InSalesClient) CreateMenu(ctx context.Context, title, handle string) (*Menu, error) {
	payload := map[string]any{
		"menu": map[string]any{
			"title":  title,
			"handle": handle,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/menus.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesMenuResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode menu: %w", err)
	}
	return MapMenu(&resp), nil
}

func (c *InSalesClient) UpdateMenu(ctx context.Context, menuID string, args map[string]any) (*Menu, error) {
	payload := map[string]any{"menu": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/menus/%s.json", menuID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesMenuResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode menu: %w", err)
	}
	return MapMenu(&resp), nil
}

func (c *InSalesClient) DeleteMenu(ctx context.Context, menuID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/menus/%s.json", menuID), nil)
	return err
}

// === РЕДИРЕКТЫ ===

func (c *InSalesClient) ListRedirects(ctx context.Context) ([]Redirect, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/redirects.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesRedirectResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode redirects: %w", err)
	}
	result := make([]Redirect, len(items))
	for i, item := range items {
		result[i] = *MapRedirect(&item)
	}
	return result, nil
}

func (c *InSalesClient) CreateRedirect(ctx context.Context, url, target string) (*Redirect, error) {
	payload := map[string]any{
		"redirect": map[string]any{
			"url":    url,
			"target": target,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/redirects.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesRedirectResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode redirect: %w", err)
	}
	return MapRedirect(&resp), nil
}

func (c *InSalesClient) DeleteRedirect(ctx context.Context, redirectID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/redirects/%s.json", redirectID), nil)
	return err
}

// === БЛОГИ ===

func (c *InSalesClient) ListBlogs(ctx context.Context) ([]Blog, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/blogs.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesBlogResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode blogs: %w", err)
	}
	result := make([]Blog, len(items))
	for i, item := range items {
		result[i] = *MapBlog(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetBlog(ctx context.Context, blogID string) (*Blog, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/blogs/%s.json", blogID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesBlogResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode blog: %w", err)
	}
	return MapBlog(&resp), nil
}

func (c *InSalesClient) CreateBlog(ctx context.Context, title, handle string) (*Blog, error) {
	payload := map[string]any{
		"blog": map[string]any{
			"title":  title,
			"handle": handle,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/blogs.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesBlogResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode blog: %w", err)
	}
	return MapBlog(&resp), nil
}

func (c *InSalesClient) DeleteBlog(ctx context.Context, blogID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/blogs/%s.json", blogID), nil)
	return err
}

func (c *InSalesClient) UpdateBlog(ctx context.Context, blogID string, args map[string]any) (*Blog, error) {
	payload := map[string]any{"blog": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/blogs/%s.json", blogID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesBlogResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode blog: %w", err)
	}
	return MapBlog(&resp), nil
}

func (c *InSalesClient) ListBlogPosts(ctx context.Context, blogID string) ([]BlogPost, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/blogs/%s/articles.json", blogID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesBlogPostResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode blog posts: %w", err)
	}
	result := make([]BlogPost, len(items))
	for i, item := range items {
		result[i] = *MapBlogPost(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetBlogPost(ctx context.Context, blogID, postID string) (*BlogPost, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/blogs/%s/articles/%s.json", blogID, postID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesBlogPostResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode blog post: %w", err)
	}
	return MapBlogPost(&resp), nil
}

func (c *InSalesClient) CreateBlogPost(ctx context.Context, blogID, title, bodyHTML string) (*BlogPost, error) {
	payload := map[string]any{
		"article": map[string]any{
			"title":        title,
			"content":      bodyHTML,
			"published_at": "2025-01-01",
			"author":       "Admin",
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/blogs/%s/articles.json", blogID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesBlogPostResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode blog post: %w", err)
	}
	return MapBlogPost(&resp), nil
}

func (c *InSalesClient) UpdateBlogPost(ctx context.Context, blogID, postID string, args map[string]any) (*BlogPost, error) {
	payload := map[string]any{"article": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/blogs/%s/articles/%s.json", blogID, postID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesBlogPostResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode blog post: %w", err)
	}
	return MapBlogPost(&resp), nil
}

func (c *InSalesClient) DeleteBlogPost(ctx context.Context, blogID, postID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/blogs/%s/articles/%s.json", blogID, postID), nil)
	return err
}

// === DISCOUNT CODES ===

func (c *InSalesClient) ListDiscountCodes(ctx context.Context, page, perPage int) ([]any, error) {
	path := fmt.Sprintf("/discount_codes.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesDiscountCodeResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode discount codes: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapDiscountCode(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetDiscountCode(ctx context.Context, codeID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/discount_codes/%s.json", codeID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesDiscountCodeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode discount code: %w", err)
	}
	return MapDiscountCode(&resp), nil
}

func (c *InSalesClient) CreateDiscountCode(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"discount_code": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/discount_codes.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesDiscountCodeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode discount code: %w", err)
	}
	return MapDiscountCode(&resp), nil
}

func (c *InSalesClient) UpdateDiscountCode(ctx context.Context, codeID string, args map[string]any) (any, error) {
	payload := map[string]any{"discount_code": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/discount_codes/%s.json", codeID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesDiscountCodeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode discount code: %w", err)
	}
	return MapDiscountCode(&resp), nil
}

func (c *InSalesClient) DeleteDiscountCode(ctx context.Context, codeID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/discount_codes/%s.json", codeID), nil)
	return err
}

// === COLLECTS ===

func (c *InSalesClient) ListCollects(ctx context.Context, page, perPage int) ([]any, error) {
	path := fmt.Sprintf("/collects.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCollectResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode collects: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapCollect(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetCollectsByProduct(ctx context.Context, productID string) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/collects.json?product_id=%s", productID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCollectResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode collects: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapCollect(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) CreateCollect(ctx context.Context, collectionID, productID int) (any, error) {
	payload := map[string]any{
		"collect": map[string]any{
			"collection_id": collectionID,
			"product_id":    productID,
		},
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/collects.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCollectResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode collect: %w", err)
	}
	return MapCollect(&resp), nil
}

func (c *InSalesClient) UpdateCollect(ctx context.Context, collectID string, args map[string]any) (any, error) {
	payload := map[string]any{"collect": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/collects/%s.json", collectID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCollectResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode collect: %w", err)
	}
	return MapCollect(&resp), nil
}

func (c *InSalesClient) DeleteCollect(ctx context.Context, collectID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/collects/%s.json", collectID), nil)
	return err
}

// === PROPERTIES ===

func (c *InSalesClient) ListProperties(ctx context.Context) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/properties.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesPropertyResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode properties: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapProperty(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetProperty(ctx context.Context, propertyID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/properties/%s.json", propertyID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesPropertyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode property: %w", err)
	}
	return MapProperty(&resp), nil
}

func (c *InSalesClient) CreateProperty(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"property": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/properties.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesPropertyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode property: %w", err)
	}
	return MapProperty(&resp), nil
}

func (c *InSalesClient) UpdateProperty(ctx context.Context, propertyID string, args map[string]any) (any, error) {
	payload := map[string]any{"property": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/properties/%s.json", propertyID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesPropertyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode property: %w", err)
	}
	return MapProperty(&resp), nil
}

func (c *InSalesClient) DeleteProperty(ctx context.Context, propertyID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/properties/%s.json", propertyID), nil)
	return err
}

// === CHARACTERISTICS ===

func (c *InSalesClient) ListCharacteristics(ctx context.Context, propertyID string) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/properties/%s/characteristics.json", propertyID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCharacteristicResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode characteristics: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapCharacteristic(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetCharacteristic(ctx context.Context, propertyID, characteristicID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/properties/%s/characteristics/%s.json", propertyID, characteristicID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesCharacteristicResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode characteristic: %w", err)
	}
	return MapCharacteristic(&resp), nil
}

func (c *InSalesClient) CreateCharacteristic(ctx context.Context, propertyID string, args map[string]any) (any, error) {
	payload := map[string]any{"characteristic": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/properties/%s/characteristics.json", propertyID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCharacteristicResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode characteristic: %w", err)
	}
	return MapCharacteristic(&resp), nil
}

func (c *InSalesClient) UpdateCharacteristic(ctx context.Context, propertyID, characteristicID string, args map[string]any) (any, error) {
	payload := map[string]any{"characteristic": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/properties/%s/characteristics/%s.json", propertyID, characteristicID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCharacteristicResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode characteristic: %w", err)
	}
	return MapCharacteristic(&resp), nil
}

func (c *InSalesClient) DeleteCharacteristic(ctx context.Context, propertyID, characteristicID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/properties/%s/characteristics/%s.json", propertyID, characteristicID), nil)
	return err
}

// === REVIEWS ===

func (c *InSalesClient) ListReviews(ctx context.Context, page, perPage int) ([]any, error) {
	path := fmt.Sprintf("/reviews.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesReviewResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode reviews: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapReview(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetReview(ctx context.Context, reviewID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/reviews/%s.json", reviewID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesReviewResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode review: %w", err)
	}
	return MapReview(&resp), nil
}

func (c *InSalesClient) CreateReview(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"review": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/reviews.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesReviewResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode review: %w", err)
	}
	return MapReview(&resp), nil
}

func (c *InSalesClient) UpdateReview(ctx context.Context, reviewID string, args map[string]any) (any, error) {
	payload := map[string]any{"review": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/reviews/%s.json", reviewID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesReviewResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode review: %w", err)
	}
	return MapReview(&resp), nil
}

func (c *InSalesClient) DeleteReview(ctx context.Context, reviewID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/reviews/%s.json", reviewID), nil)
	return err
}

// === FILES ===

func (c *InSalesClient) ListFiles(ctx context.Context) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/files.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesFileResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode files: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapFile(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetFile(ctx context.Context, fileID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/files/%s.json", fileID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesFileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode file: %w", err)
	}
	return MapFile(&resp), nil
}

func (c *InSalesClient) CreateFile(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"file": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/files.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesFileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode file: %w", err)
	}
	return MapFile(&resp), nil
}

func (c *InSalesClient) DeleteFile(ctx context.Context, fileID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/files/%s.json", fileID), nil)
	return err
}

// === DOMAINS ===

func (c *InSalesClient) ListDomains(ctx context.Context) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/domains.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesDomainResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode domains: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapDomain(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetDomain(ctx context.Context, domainID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/domains/%s.json", domainID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesDomainResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode domain: %w", err)
	}
	return MapDomain(&resp), nil
}

func (c *InSalesClient) CreateDomain(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"domain": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/domains.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesDomainResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode domain: %w", err)
	}
	return MapDomain(&resp), nil
}

func (c *InSalesClient) UpdateDomain(ctx context.Context, domainID string, args map[string]any) (any, error) {
	payload := map[string]any{"domain": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/domains/%s.json", domainID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesDomainResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode domain: %w", err)
	}
	return MapDomain(&resp), nil
}

func (c *InSalesClient) DeleteDomain(ctx context.Context, domainID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/domains/%s.json", domainID), nil)
	return err
}

// === CLIENT GROUPS ===

func (c *InSalesClient) ListClientGroups(ctx context.Context) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/client_groups.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesClientGroupResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode client groups: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapClientGroup(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetClientGroup(ctx context.Context, groupID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/client_groups/%s.json", groupID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesClientGroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode client group: %w", err)
	}
	return MapClientGroup(&resp), nil
}

func (c *InSalesClient) CreateClientGroup(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"client_group": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/client_groups.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesClientGroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode client group: %w", err)
	}
	return MapClientGroup(&resp), nil
}

func (c *InSalesClient) UpdateClientGroup(ctx context.Context, groupID string, args map[string]any) (any, error) {
	payload := map[string]any{"client_group": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/client_groups/%s.json", groupID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesClientGroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode client group: %w", err)
	}
	return MapClientGroup(&resp), nil
}

func (c *InSalesClient) DeleteClientGroup(ctx context.Context, groupID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/client_groups/%s.json", groupID), nil)
	return err
}

// === CUSTOM STATUSES ===

func (c *InSalesClient) ListCustomStatuses(ctx context.Context) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/custom_statuses.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCustomStatusResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode custom statuses: %w", err)
	}
	result := make([]any, len(items))
	for i := range items {
		result[i] = MapCustomStatus(&items[i])
	}
	return result, nil
}

func (c *InSalesClient) GetCustomStatus(ctx context.Context, statusID string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/custom_statuses/%s.json", statusID), nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesCustomStatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode custom status: %w", err)
	}
	return MapCustomStatus(&resp), nil
}

func (c *InSalesClient) CreateCustomStatus(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"custom_status": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/custom_statuses.json", payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCustomStatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode custom status: %w", err)
	}
	return MapCustomStatus(&resp), nil
}

func (c *InSalesClient) UpdateCustomStatus(ctx context.Context, statusID string, args map[string]any) (any, error) {
	payload := map[string]any{"custom_status": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/custom_statuses/%s.json", statusID), payload)
	if err != nil {
		return nil, err
	}
	var resp InSalesCustomStatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode custom status: %w", err)
	}
	return MapCustomStatus(&resp), nil
}

func (c *InSalesClient) DeleteCustomStatus(ctx context.Context, statusID string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/custom_statuses/%s.json", statusID), nil)
	return err
}

// Generic CRUD helpers to reduce boilerplate

func (c *InSalesClient) crudList(ctx context.Context, path string) ([]any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []json.RawMessage
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode list: %w", err)
	}
	result := make([]any, len(items))
	for i, item := range items {
		var m map[string]any
		if err := json.Unmarshal(item, &m); err != nil {
			continue
		}
		result[i] = m
	}
	return result, nil
}

func (c *InSalesClient) crudGet(ctx context.Context, path string) (any, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) crudCreate(ctx context.Context, path string, key string, args map[string]any) (any, error) {
	payload := map[string]any{key: args}
	body, _, err := c.doRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) crudUpdate(ctx context.Context, path string, key string, args map[string]any) (any, error) {
	payload := map[string]any{key: args}
	body, _, err := c.doRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) crudDelete(ctx context.Context, path string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

// === PRODUCT FIELDS ===

func (c *InSalesClient) ListProductFields(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/product_fields.json")
}
func (c *InSalesClient) GetProductField(ctx context.Context, fieldID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/product_fields/%s.json", fieldID))
}
func (c *InSalesClient) CreateProductField(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/product_fields.json", "product_field", args)
}
func (c *InSalesClient) UpdateProductField(ctx context.Context, fieldID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/product_fields/%s.json", fieldID), "product_field", args)
}
func (c *InSalesClient) DeleteProductField(ctx context.Context, fieldID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/product_fields/%s.json", fieldID))
}

// === VARIANT FIELDS ===

func (c *InSalesClient) ListVariantFields(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/variant_fields.json")
}
func (c *InSalesClient) GetVariantField(ctx context.Context, fieldID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/variant_fields/%s.json", fieldID))
}
func (c *InSalesClient) CreateVariantField(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/variant_fields.json", "variant_field", args)
}
func (c *InSalesClient) UpdateVariantField(ctx context.Context, fieldID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/variant_fields/%s.json", fieldID), "variant_field", args)
}
func (c *InSalesClient) DeleteVariantField(ctx context.Context, fieldID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/variant_fields/%s.json", fieldID))
}

// === COLLECTION FIELDS ===

func (c *InSalesClient) ListCollectionFields(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/collection_fields.json")
}
func (c *InSalesClient) GetCollectionField(ctx context.Context, fieldID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/collection_fields/%s.json", fieldID))
}
func (c *InSalesClient) CreateCollectionField(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/collection_fields.json", "collection_field", args)
}
func (c *InSalesClient) UpdateCollectionField(ctx context.Context, fieldID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/collection_fields/%s.json", fieldID), "collection_field", args)
}
func (c *InSalesClient) DeleteCollectionField(ctx context.Context, fieldID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/collection_fields/%s.json", fieldID))
}

// === COLLECTION FILTERS ===

func (c *InSalesClient) ListCollectionFilters(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/collection_filters.json")
}
func (c *InSalesClient) GetCollectionFilter(ctx context.Context, filterID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/collection_filters/%s.json", filterID))
}
func (c *InSalesClient) CreateCollectionFilter(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/collection_filters.json", "collection_filter", args)
}
func (c *InSalesClient) UpdateCollectionFilter(ctx context.Context, filterID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/collection_filters/%s.json", filterID), "collection_filter", args)
}
func (c *InSalesClient) DeleteCollectionFilter(ctx context.Context, filterID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/collection_filters/%s.json", filterID))
}

// === JS TAGS ===

func (c *InSalesClient) ListJsTags(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/js_tags.json")
}
func (c *InSalesClient) GetJsTag(ctx context.Context, tagID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/js_tags/%s.json", tagID))
}
func (c *InSalesClient) CreateJsTag(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/js_tags.json", "js_tag", args)
}
func (c *InSalesClient) UpdateJsTag(ctx context.Context, tagID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/js_tags/%s.json", tagID), "js_tag", args)
}
func (c *InSalesClient) DeleteJsTag(ctx context.Context, tagID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/js_tags/%s.json", tagID))
}

// === EXTERNAL DISCOUNTS ===

func (c *InSalesClient) ListExternalDiscounts(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/external_discounts.json")
}
func (c *InSalesClient) GetExternalDiscount(ctx context.Context, discountID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/external_discounts/%s.json", discountID))
}
func (c *InSalesClient) CreateExternalDiscount(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/external_discounts.json", "external_discount", args)
}
func (c *InSalesClient) UpdateExternalDiscount(ctx context.Context, discountID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/external_discounts/%s.json", discountID), "external_discount", args)
}
func (c *InSalesClient) DeleteExternalDiscount(ctx context.Context, discountID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/external_discounts/%s.json", discountID))
}

// === DISCOUNT COLLECTIONS GROUPS ===

func (c *InSalesClient) ListDiscountCollectionsGroups(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/discount_collections_groups.json")
}
func (c *InSalesClient) GetDiscountCollectionsGroup(ctx context.Context, groupID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/discount_collections_groups/%s.json", groupID))
}
func (c *InSalesClient) CreateDiscountCollectionsGroup(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/discount_collections_groups.json", "discount_collections_group", args)
}
func (c *InSalesClient) UpdateDiscountCollectionsGroup(ctx context.Context, groupID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/discount_collections_groups/%s.json", groupID), "discount_collections_group", args)
}
func (c *InSalesClient) DeleteDiscountCollectionsGroup(ctx context.Context, groupID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/discount_collections_groups/%s.json", groupID))
}

// === WAREHOUSES ===

func (c *InSalesClient) ListWarehouses(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/warehouses.json")
}
func (c *InSalesClient) GetWarehouse(ctx context.Context, warehouseID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/warehouses/%s.json", warehouseID))
}
func (c *InSalesClient) CreateWarehouse(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/warehouses.json", "warehouse", args)
}
func (c *InSalesClient) UpdateWarehouse(ctx context.Context, warehouseID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/warehouses/%s.json", warehouseID), "warehouse", args)
}
func (c *InSalesClient) DeleteWarehouse(ctx context.Context, warehouseID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/warehouses/%s.json", warehouseID))
}

// === SALES CHANNELS ===

func (c *InSalesClient) ListSalesChannels(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/sales_channels.json")
}
func (c *InSalesClient) GetSalesChannel(ctx context.Context, channelID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/sales_channels/%s.json", channelID))
}
func (c *InSalesClient) CreateSalesChannel(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/sales_channels.json", "sales_channel", args)
}
func (c *InSalesClient) UpdateSalesChannel(ctx context.Context, channelID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/sales_channels/%s.json", channelID), "sales_channel", args)
}
func (c *InSalesClient) DeleteSalesChannel(ctx context.Context, channelID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/sales_channels/%s.json", channelID))
}

// === STOCK CURRENCIES ===

func (c *InSalesClient) ListStockCurrencies(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/stock_currencies.json")
}
func (c *InSalesClient) GetStockCurrency(ctx context.Context, currencyID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/stock_currencies/%s.json", currencyID))
}
func (c *InSalesClient) CreateStockCurrency(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/stock_currencies.json", "stock_currency", args)
}
func (c *InSalesClient) UpdateStockCurrency(ctx context.Context, currencyID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/stock_currencies/%s.json", currencyID), "stock_currency", args)
}
func (c *InSalesClient) DeleteStockCurrency(ctx context.Context, currencyID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/stock_currencies/%s.json", currencyID))
}

// === PRICE KINDS ===

func (c *InSalesClient) ListPriceKinds(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/price_kinds.json")
}
func (c *InSalesClient) GetPriceKind(ctx context.Context, kindID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/price_kinds/%s.json", kindID))
}
func (c *InSalesClient) CreatePriceKind(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/price_kinds.json", "price_kind", args)
}
func (c *InSalesClient) UpdatePriceKind(ctx context.Context, kindID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/price_kinds/%s.json", kindID), "price_kind", args)
}
func (c *InSalesClient) DeletePriceKind(ctx context.Context, kindID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/price_kinds/%s.json", kindID))
}

// === OPTION NAMES ===

func (c *InSalesClient) ListOptionNames(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/option_names.json")
}
func (c *InSalesClient) GetOptionName(ctx context.Context, optionNameID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/option_names/%s.json", optionNameID))
}
func (c *InSalesClient) CreateOptionName(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/option_names.json", "option_name", args)
}
func (c *InSalesClient) UpdateOptionName(ctx context.Context, optionNameID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/option_names/%s.json", optionNameID), "option_name", args)
}
func (c *InSalesClient) DeleteOptionName(ctx context.Context, optionNameID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/option_names/%s.json", optionNameID))
}

// === OPTION VALUES ===

func (c *InSalesClient) ListOptionValues(ctx context.Context, optionNameID string) ([]any, error) {
	return c.crudList(ctx, fmt.Sprintf("/option_names/%s/option_values.json", optionNameID))
}
func (c *InSalesClient) GetOptionValue(ctx context.Context, optionNameID, optionValueID string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/option_names/%s/option_values/%s.json", optionNameID, optionValueID))
}
func (c *InSalesClient) CreateOptionValue(ctx context.Context, optionNameID string, args map[string]any) (any, error) {
	return c.crudCreate(ctx, fmt.Sprintf("/option_names/%s/option_values.json", optionNameID), "option_value", args)
}
func (c *InSalesClient) UpdateOptionValue(ctx context.Context, optionNameID, optionValueID string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/option_names/%s/option_values/%s.json", optionNameID, optionValueID), "option_value", args)
}
func (c *InSalesClient) DeleteOptionValue(ctx context.Context, optionNameID, optionValueID string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/option_names/%s/option_values/%s.json", optionNameID, optionValueID))
}

// === APPLICATION ACTIONS ===

func (c *InSalesClient) ListApplicationActions(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/application_actions.json")
}
func (c *InSalesClient) GetApplicationAction(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/application_actions/%s.json", id))
}
func (c *InSalesClient) CreateApplicationAction(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/application_actions.json", "application_action", args)
}
func (c *InSalesClient) UpdateApplicationAction(ctx context.Context, id string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/application_actions/%s.json", id), "application_action", args)
}
func (c *InSalesClient) DeleteApplicationAction(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/application_actions/%s.json", id))
}

// === APPLICATION CHARGES ===

func (c *InSalesClient) ListApplicationCharges(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/application_charges.json")
}
func (c *InSalesClient) GetApplicationCharge(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/application_charges/%s.json", id))
}
func (c *InSalesClient) CreateApplicationCharge(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/application_charges.json", "application_charge", args)
}
func (c *InSalesClient) DeclineApplicationCharge(ctx context.Context, id string) error {
	_, _, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/application_charges/%s/decline.json", id), nil)
	return err
}

// === APPLICATION WIDGETS ===

func (c *InSalesClient) ListApplicationWidgets(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/application_widgets.json")
}
func (c *InSalesClient) GetApplicationWidget(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/application_widgets/%s.json", id))
}
func (c *InSalesClient) CreateApplicationWidget(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/application_widgets.json", "application_widget", args)
}
func (c *InSalesClient) UpdateApplicationWidget(ctx context.Context, id string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/application_widgets/%s.json", id), "application_widget", args)
}
func (c *InSalesClient) DeleteApplicationWidget(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/application_widgets/%s.json", id))
}

// === RECURRING APPLICATION CHARGES ===

func (c *InSalesClient) ListRecurringApplicationCharges(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/recurring_application_charges.json")
}
func (c *InSalesClient) GetRecurringApplicationCharge(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/recurring_application_charges/%s.json", id))
}
func (c *InSalesClient) CreateRecurringApplicationCharge(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/recurring_application_charges.json", "recurring_application_charge", args)
}
func (c *InSalesClient) DeleteRecurringApplicationCharge(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/recurring_application_charges/%s.json", id))
}

// === BONUS SYSTEM TRANSACTIONS ===

func (c *InSalesClient) ListBonusSystemTransactions(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/bonus_system_transactions.json")
}
func (c *InSalesClient) GetBonusSystemTransaction(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/bonus_system_transactions/%s.json", id))
}
func (c *InSalesClient) CreateBonusSystemTransaction(ctx context.Context, clientID int, args map[string]any) (any, error) {
	payload := map[string]any{"bonus_system_transaction": args}
	body, _, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/clients/%d/bonus_system_transactions.json", clientID), payload)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return result, nil
}
func (c *InSalesClient) ListClientBonusSystemTransactions(ctx context.Context, clientID string) ([]any, error) {
	return c.crudList(ctx, fmt.Sprintf("/clients/%s/bonus_system_transactions.json", clientID))
}

// === OAUTH APPS ===

func (c *InSalesClient) ListOAuthApps(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/o_auth_apps.json")
}
func (c *InSalesClient) GetOAuthApp(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/o_auth_apps/%s.json", id))
}
func (c *InSalesClient) CreateOAuthApp(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/o_auth_apps.json", "o_auth_app", args)
}
func (c *InSalesClient) UpdateOAuthApp(ctx context.Context, id string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/o_auth_apps/%s.json", id), "o_auth_app", args)
}
func (c *InSalesClient) DeleteOAuthApp(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/o_auth_apps/%s.json", id))
}

// === USER TASKS ===

func (c *InSalesClient) ListUserTasks(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/user_tasks.json")
}
func (c *InSalesClient) GetUserTask(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/user_tasks/%s.json", id))
}
func (c *InSalesClient) CreateUserTask(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/user_tasks.json", "user_task", args)
}
func (c *InSalesClient) UpdateUserTask(ctx context.Context, id string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/user_tasks/%s.json", id), "user_task", args)
}
func (c *InSalesClient) DeleteUserTask(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/user_tasks/%s.json", id))
}

// === SUPPLEMENTARY ===

func (c *InSalesClient) ListSupplementaries(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/supplementaries.json")
}
func (c *InSalesClient) GetSupplementary(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/supplementaries/%s.json", id))
}
func (c *InSalesClient) CreateSupplementary(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/supplementaries.json", "supplementary", args)
}
func (c *InSalesClient) DeleteSupplementary(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/supplementaries/%s.json", id))
}

// === SIMILARS ===

func (c *InSalesClient) ListSimilars(ctx context.Context, productID string) ([]any, error) {
	return c.crudList(ctx, fmt.Sprintf("/products/%s/similars.json", productID))
}
func (c *InSalesClient) CreateSimilar(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/similars.json", "similar", args)
}
func (c *InSalesClient) DeleteSimilar(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/similars/%s.json", id))
}

// === UNITS ===

func (c *InSalesClient) ListUnits(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/units.json")
}

// === OUTLETS ===

func (c *InSalesClient) ListOutlets(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/outlets.json")
}
func (c *InSalesClient) GetOutlet(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/outlets/%s.json", id))
}
func (c *InSalesClient) CreateOutlet(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/outlets.json", "outlet", args)
}
func (c *InSalesClient) UpdateOutlet(ctx context.Context, id string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/outlets/%s.json", id), "outlet", args)
}
func (c *InSalesClient) DeleteOutlet(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/outlets/%s.json", id))
}

// === MARKETPLACES ===

func (c *InSalesClient) ListMarketplaces(ctx context.Context) ([]any, error) {
	return c.crudList(ctx, "/marketplaces.json")
}
func (c *InSalesClient) GetMarketplace(ctx context.Context, id string) (any, error) {
	return c.crudGet(ctx, fmt.Sprintf("/marketplaces/%s.json", id))
}
func (c *InSalesClient) CreateMarketplace(ctx context.Context, args map[string]any) (any, error) {
	return c.crudCreate(ctx, "/marketplaces.json", "marketplace", args)
}
func (c *InSalesClient) UpdateMarketplace(ctx context.Context, id string, args map[string]any) (any, error) {
	return c.crudUpdate(ctx, fmt.Sprintf("/marketplaces/%s.json", id), "marketplace", args)
}
func (c *InSalesClient) DeleteMarketplace(ctx context.Context, id string) error {
	return c.crudDelete(ctx, fmt.Sprintf("/marketplaces/%s.json", id))
}

// === PROVIDER TOKEN ===

func (c *InSalesClient) CreateProviderToken(ctx context.Context, messengerType string) (any, error) {
	payload := map[string]any{"messenger_type": messengerType}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/dialogs/provider_tokens.json", payload)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode provider token: %w", err)
	}
	return result, nil
}

// === ACCOUNT ===

func (c *InSalesClient) UpdateAccount(ctx context.Context, args map[string]any) (any, error) {
	payload := map[string]any{"account": args}
	body, _, err := c.doRequest(ctx, http.MethodPut, "/account.json", payload)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode account: %w", err)
	}
	return result, nil
}

// === COUNTS ===

func (c *InSalesClient) GetOrdersCount(ctx context.Context, filters map[string]any) (int, error) {
	path := "/orders/count.json"
	if len(filters) > 0 {
		params := url.Values{}
		for k, v := range filters {
			switch val := v.(type) {
			case string:
				if val != "" {
					params.Set(k, val)
				}
			case []any:
				for _, item := range val {
					if s, ok := item.(string); ok {
						params.Add(k+"[]", s)
					}
				}
			}
		}
		if qs := params.Encode(); qs != "" {
			path += "?" + qs
		}
	}
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("decode orders count: %w", err)
	}
	return resp.Count, nil
}

func (c *InSalesClient) GetProductsCount(ctx context.Context) (int, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/products/count.json", nil)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("decode products count: %w", err)
	}
	return resp.Count, nil
}
