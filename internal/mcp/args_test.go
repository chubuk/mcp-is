package mcp

import (
	"encoding/json"
	"testing"
)

func TestGetFloatArg_CoercesAndErrors(t *testing.T) {
	tests := []struct {
		name    string
		arg     any
		want    float64
		wantErr bool
	}{
		{"float64", float64(3), 3, false},
		{"int", 3, 3, false},
		{"json.Number", json.Number("3.5"), 3.5, false},
		{"string", "3", 3, false},
		{"string spaces", "  3 ", 3, false},
		{"garbage string", "abc", 0, true},
		{"bool", true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFloatArg(map[string]any{"x": tt.arg}, "x", 25)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %v", got)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Errorf("got (%v, %v), want (%v, nil)", got, err, tt.want)
			}
		})
	}
	if got, err := getFloatArg(map[string]any{}, "x", 25); err != nil || got != 25 {
		t.Errorf("absent key: got (%v, %v), want (25, nil)", got, err)
	}
	if got, err := getFloatArg(map[string]any{"x": nil}, "x", 25); err != nil || got != 25 {
		t.Errorf("nil key: got (%v, %v), want (25, nil)", got, err)
	}
}

func TestPageArgs_ClampsAndParsesStrings(t *testing.T) {
	page, perPage, err := pageArgs(map[string]any{"page": "3", "per_page": "5"})
	if err != nil || page != 3 || perPage != 5 {
		t.Errorf("string args: got (%d, %d, %v)", page, perPage, err)
	}
	page, perPage, err = pageArgs(map[string]any{"page": float64(0), "per_page": float64(300)})
	if err != nil || page != 1 || perPage != 250 {
		t.Errorf("clamps: got (%d, %d, %v)", page, perPage, err)
	}
	page, perPage, err = pageArgs(map[string]any{"per_page": float64(0)})
	if err != nil || page != 1 || perPage != 25 {
		t.Errorf("zero per_page: got (%d, %d, %v)", page, perPage, err)
	}
	if _, _, err = pageArgs(map[string]any{"per_page": "abc"}); err == nil {
		t.Error("garbage per_page must be a tool error, not a silent default")
	}
	page, perPage, err = pageArgs(map[string]any{})
	if err != nil || page != 1 || perPage != 25 {
		t.Errorf("defaults: got (%d, %d, %v)", page, perPage, err)
	}
}

func TestResolveKeywordsAlias(t *testing.T) {
	args := map[string]any{"keywords": "k"}
	if err := resolveKeywordsAlias(args); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if args["meta_keywords"] != "k" {
		t.Errorf("alias not applied: %v", args)
	}
	if _, ok := args["keywords"]; ok {
		t.Errorf("alias key must be removed: %v", args)
	}
	if err := resolveKeywordsAlias(map[string]any{"keywords": "k", "meta_keywords": "m"}); err == nil {
		t.Error("conflict must error")
	}
	if err := resolveKeywordsAlias(map[string]any{"meta_keywords": "m"}); err != nil {
		t.Errorf("no alias must be a no-op: %v", err)
	}
}

func TestFilterKnownAndIgnoredKeys(t *testing.T) {
	args := map[string]any{"title": "x", "zzz": 1, "id": 2}
	ignored := filterKnown(args, map[string]bool{"title": true, "product_id": true})
	if len(args) != 1 || args["title"] != "x" {
		t.Errorf("whitelist filtering wrong: %v", args)
	}
	if len(ignored) != 2 || ignored[0] != "id" || ignored[1] != "zzz" {
		t.Errorf("ignored not sorted: %v", ignored)
	}
	out := withIgnoredKeys(map[string]any{"title": "x"}, ignored).(map[string]any)
	if got, ok := out["ignored_keys"].([]string); !ok || len(got) != 2 {
		t.Errorf("ignored_keys missing: %v", out)
	}
	if withIgnoredKeys(map[string]any{"a": 1}, nil).(map[string]any)["ignored_keys"] != nil {
		t.Error("empty ignored must not add the key")
	}
}
