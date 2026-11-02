package insales

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSeoEndpoint_Validation(t *testing.T) {
	if _, _, _, err := seoEndpoint("product", "1", ""); err != nil {
		t.Errorf("product: %v", err)
	}
	if _, _, _, err := seoEndpoint("collection", "1", ""); err != nil {
		t.Errorf("collection: %v", err)
	}
	if _, _, _, err := seoEndpoint("article", "1", ""); err == nil {
		t.Error("article without blog_id must error")
	}
	path, key, allowed, err := seoEndpoint("article", "2", "1")
	if err != nil || path != "/blogs/1/articles/2.json" || key != "article" || !allowed["meta_keywords"] {
		t.Errorf("article: %v %q %v", err, path, allowed)
	}
	if _, _, _, err := seoEndpoint("widget", "1", ""); err == nil {
		t.Error("unknown resource_type must error")
	}
	if _, _, allowed, _ := seoEndpoint("collection", "1", ""); !allowed["seo_description"] {
		t.Error("collection must allow seo_description")
	}
	if _, _, allowed, _ := seoEndpoint("product", "1", ""); allowed["seo_description"] {
		t.Error("product must not allow seo_description")
	}
}

func TestUpdateSEO_WhitelistAndClearFields(t *testing.T) {
	var got map[string]map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/products/5.json" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		got = body
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"html_title": "h", "meta_keywords": null}`))
	}))
	defer srv.Close()
	client := newTestClient(srv)

	set, err := client.UpdateSEO(context.Background(), "product", "5", "", map[string]any{
		"resource_type":   "product",
		"resource_id":     "5",
		"html_title":      "h",
		"meta_keywords":   "",
		"seo_description": "nope",
		"zzz":             "junk",
		"clear_fields":    []any{"meta_description"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := got["product"]
	if p["html_title"] != "h" {
		t.Errorf("html_title = %v", p["html_title"])
	}
	if _, ok := p["meta_keywords"]; ok {
		t.Errorf("empty meta_keywords must be dropped, got %v", p)
	}
	if _, ok := p["seo_description"]; ok {
		t.Errorf("seo_description must not pass product whitelist, got %v", p)
	}
	if _, ok := p["zzz"]; ok {
		t.Errorf("junk key must not pass, got %v", p)
	}
	if v, ok := p["meta_description"]; !ok || v != "" {
		t.Errorf("clear_fields must write empty string, got %v", p)
	}
	if set.HTMLTitle == nil || *set.HTMLTitle != "h" || set.MetaKeywords != nil {
		t.Errorf("decode: %+v", set)
	}
}
