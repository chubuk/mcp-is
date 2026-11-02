package insales

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SEOFields is the SEO subset InSales entities carry (collection also has
// seo_description). Pointers keep "not set" (null) and "cleared" ("") apart.
type SEOFields struct {
	HTMLTitle       *string `json:"html_title"`
	MetaKeywords    *string `json:"meta_keywords"`
	MetaDescription *string `json:"meta_description"`
	SEODescription  *string `json:"seo_description,omitempty"`
}

var seoCommonFields = map[string]bool{
	"html_title": true, "meta_keywords": true, "meta_description": true,
}

var seoCollectionFields = map[string]bool{
	"html_title": true, "meta_keywords": true, "meta_description": true,
	"seo_description": true,
}

// seoEndpoint resolves the entity URL, payload key and writable SEO fields.
func seoEndpoint(resourceType, id, blogID string) (path, key string, allowed map[string]bool, err error) {
	switch resourceType {
	case "product":
		return fmt.Sprintf("/products/%s.json", id), "product", seoCommonFields, nil
	case "page":
		return fmt.Sprintf("/pages/%s.json", id), "page", seoCommonFields, nil
	case "article":
		if blogID == "" {
			return "", "", nil, fmt.Errorf("blog_id is required for resource_type=article")
		}
		return fmt.Sprintf("/blogs/%s/articles/%s.json", blogID, id), "article", seoCommonFields, nil
	case "blog":
		return fmt.Sprintf("/blogs/%s.json", id), "blog", seoCommonFields, nil
	case "collection":
		return fmt.Sprintf("/collections/%s.json", id), "collection", seoCollectionFields, nil
	default:
		return "", "", nil, fmt.Errorf("unknown resource_type %q (product|page|article|blog|collection)", resourceType)
	}
}

func (c *InSalesClient) GetSEO(ctx context.Context, resourceType, id, blogID string) (*SEOFields, error) {
	path, _, _, err := seoEndpoint(resourceType, id, blogID)
	if err != nil {
		return nil, err
	}
	body, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var set SEOFields
	if err := json.Unmarshal(body, &set); err != nil {
		return nil, fmt.Errorf("decode seo: %w", err)
	}
	return &set, nil
}

// UpdateSEO writes only whitelisted SEO fields; empties are dropped (clearing is
// explicit via clear_fields).
func (c *InSalesClient) UpdateSEO(ctx context.Context, resourceType, id, blogID string, args map[string]any) (*SEOFields, error) {
	path, key, allowed, err := seoEndpoint(resourceType, id, blogID)
	if err != nil {
		return nil, err
	}
	payload := make(map[string]any)
	for k, v := range cleanArgs(args) {
		if allowed[k] {
			payload[k] = v
		}
	}
	if cf, ok := args["clear_fields"]; ok {
		if err := applyClearFields(payload, allowed, cf); err != nil {
			return nil, err
		}
	}
	body, _, err := c.doRequest(ctx, http.MethodPut, path, map[string]any{key: payload})
	if err != nil {
		return nil, err
	}
	var set SEOFields
	if err := json.Unmarshal(body, &set); err != nil {
		return nil, fmt.Errorf("decode seo: %w", err)
	}
	return &set, nil
}
