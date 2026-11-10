package insales

import (
	"reflect"
	"testing"
)

func TestBuildBlogArticle_RenameAndStrip(t *testing.T) {
	got := BuildBlogArticle(map[string]any{
		"blog_id":         "1",
		"post_id":         "7",
		"idempotency_key": "abc",
		"title":           "t",
		"body_html":       "<p>x</p>",
		"published_at":    "2026-12-01",
		"image_url":       "https://shop.example/i.jpg",
		"notice":          "",
	})
	want := map[string]any{
		"title":            "t",
		"content":          "<p>x</p>",
		"published_at":     "2026-12-01",
		"image_attributes": map[string]any{"src": "https://shop.example/i.jpg"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildBlogArticle = %#v, want %#v", got, want)
	}
}

func TestBuildBlogArticle_EmptyImageIgnored(t *testing.T) {
	got := BuildBlogArticle(map[string]any{"image_url": "", "title": "t"})
	if _, ok := got["image_attributes"]; ok {
		t.Errorf("empty image_url must not produce image_attributes: %#v", got)
	}
	if _, ok := got["image_url"]; ok {
		t.Errorf("image_url must not leak into article: %#v", got)
	}
}

func TestBuildBlogArticle_SeOAndTagsPassThrough(t *testing.T) {
	got := BuildBlogArticle(map[string]any{
		"html_title":       "SEO",
		"meta_keywords":    "k",
		"meta_description": "d",
		"all_tags":         []any{"a", "b"},
		"pinned":           true,
	})
	want := map[string]any{
		"html_title":       "SEO",
		"meta_keywords":    "k",
		"meta_description": "d",
		"all_tags":         []any{"a", "b"},
		"pinned":           true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildBlogArticle = %#v, want %#v", got, want)
	}
}
