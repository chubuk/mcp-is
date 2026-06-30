package insales

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type InSalesClient struct {
	BaseURL    string
	ShopDomain string
	APIKey     string
	Password   string
	HTTPClient *http.Client
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
	}, nil
}

func (c *InSalesClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
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

func (c *InSalesClient) GetStoreInfo(ctx context.Context) (*StoreInfo, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/shops.json", nil)
	if err != nil {
		return nil, err
	}
	var resp InSalesShopResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode shop info: %w", err)
	}
	return MapShop(&resp), nil
}

func (c *InSalesClient) ListProducts(ctx context.Context, page, perPage int) ([]Product, error) {
	path := fmt.Sprintf("/products.json?page=%d&per_page=%d", page, perPage)
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

func (c *InSalesClient) CreateProduct(ctx context.Context, title, price string, args map[string]interface{}) (*Product, error) {
	payload := map[string]interface{}{
		"product": map[string]interface{}{
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

func (c *InSalesClient) UpdateProduct(ctx context.Context, productID string, args map[string]interface{}) (*Product, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) ListOrders(ctx context.Context, page, perPage int) ([]Order, error) {
	path := fmt.Sprintf("/orders.json?page=%d&per_page=%d", page, perPage)
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

func (c *InSalesClient) UpdateOrder(ctx context.Context, orderID string, args map[string]interface{}) (*Order, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) CreateOrder(ctx context.Context, clientID int, args map[string]interface{}) (*Order, error) {
	payload := map[string]interface{}{
		"order": map[string]interface{}{
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

func (c *InSalesClient) ListClients(ctx context.Context, page, perPage int) ([]interface{}, error) {
	path := fmt.Sprintf("/clients.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesClientResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode clients: %w", err)
	}
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = MapClient(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetClient(ctx context.Context, clientID string) (interface{}, error) {
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

func (c *InSalesClient) CreateClient(ctx context.Context, email string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
		"client": map[string]interface{}{
			"email":      email,
			"first_name": args["first_name"],
			"last_name":  args["last_name"],
			"phone":      args["phone"],
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

func (c *InSalesClient) UpdateClient(ctx context.Context, clientID string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) ListCategories(ctx context.Context, page, perPage int) ([]interface{}, error) {
	path := fmt.Sprintf("/categories.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCategoryResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode categories: %w", err)
	}
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = MapCategory(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetCategory(ctx context.Context, categoryID string) (interface{}, error) {
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

func (c *InSalesClient) CreateCategory(ctx context.Context, title string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
		"category": map[string]interface{}{
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

func (c *InSalesClient) UpdateCategory(ctx context.Context, categoryID string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) ListCollections(ctx context.Context, page, perPage int) ([]interface{}, error) {
	path := fmt.Sprintf("/collections.json?page=%d&per_page=%d", page, perPage)
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesCollectionResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode collections: %w", err)
	}
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = MapCollection(&item)
	}
	return result, nil
}

func (c *InSalesClient) GetCollection(ctx context.Context, collectionID string) (interface{}, error) {
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

func (c *InSalesClient) CreateCollection(ctx context.Context, title string) (interface{}, error) {
	payload := map[string]interface{}{
		"collection": map[string]interface{}{
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

func (c *InSalesClient) UpdateCollection(ctx context.Context, collectionID, title string) (interface{}, error) {
	payload := map[string]interface{}{
		"collection": map[string]interface{}{
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

func (c *InSalesClient) ListImages(ctx context.Context, productID string) ([]interface{}, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/products/%s/images.json", productID), nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesImage
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode images: %w", err)
	}
	result := make([]interface{}, len(items))
	for i := range items {
		result[i] = items[i]
	}
	return result, nil
}

func (c *InSalesClient) AddImage(ctx context.Context, productID, url string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
		"image": map[string]interface{}{
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

func (c *InSalesClient) ListMetaFields(ctx context.Context, args map[string]interface{}) ([]interface{}, error) {
	path := "/meta_fields.json"
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesMetaFieldResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode meta fields: %w", err)
	}
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = MapMetaField(&item)
	}
	return result, nil
}

func (c *InSalesClient) CreateMetaField(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) UpdateMetaField(ctx context.Context, metaFieldID string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) ListWebhooks(ctx context.Context) ([]interface{}, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/webhooks.json", nil)
	if err != nil {
		return nil, err
	}
	var items []InSalesWebhookResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode webhooks: %w", err)
	}
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = MapWebhook(&item)
	}
	return result, nil
}

func (c *InSalesClient) CreateWebhook(ctx context.Context, event, url string, args map[string]interface{}) (interface{}, error) {
	format, _ := args["format"].(string)
	if format == "" {
		format = "json"
	}
	payload := map[string]interface{}{
		"webhook": map[string]interface{}{
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

func (c *InSalesClient) UpdateWebhook(ctx context.Context, webhookID string, args map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{
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

func (c *InSalesClient) ListCurrencies(ctx context.Context) (interface{}, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/currencies.json", nil)
	if err != nil {
		return nil, err
	}
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode currencies: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) ListDeliveryVariants(ctx context.Context) (interface{}, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/delivery_variants.json", nil)
	if err != nil {
		return nil, err
	}
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode delivery variants: %w", err)
	}
	return result, nil
}

func (c *InSalesClient) ListPaymentGateways(ctx context.Context) (interface{}, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, "/payment_gateways.json", nil)
	if err != nil {
		return nil, err
	}
	var result interface{}
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
	payload := map[string]interface{}{
		"page": map[string]interface{}{
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

func (c *InSalesClient) UpdatePage(ctx context.Context, pageID string, args map[string]interface{}) (*Page, error) {
	payload := map[string]interface{}{"page": args}
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

func (c *InSalesClient) CreateAccessory(ctx context.Context, title string) (*Accessory, error) {
	payload := map[string]interface{}{
		"accessory": map[string]interface{}{
			"title": title,
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

func (c *InSalesClient) UpdateAccessory(ctx context.Context, accessoryID string, args map[string]interface{}) (*Accessory, error) {
	payload := map[string]interface{}{"accessory": args}
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
