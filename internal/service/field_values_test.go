package service

import (
	"reflect"
	"testing"
)

func TestPopFieldValues(t *testing.T) {
	args := map[string]any{"title": "t", "field_values": map[string]any{"sostav": "x"}}
	fv := popFieldValues(args)
	if !reflect.DeepEqual(fv, map[string]any{"sostav": "x"}) {
		t.Errorf("fv = %#v", fv)
	}
	if _, ok := args["field_values"]; ok {
		t.Errorf("field_values must be removed from args: %#v", args)
	}
	if popFieldValues(args) != nil {
		t.Error("second pop must be nil")
	}
}

func TestInlineFieldAttrs(t *testing.T) {
	got := inlineFieldAttrs(map[string]any{"a": "1", "b": 2})
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	seen := map[string]any{}
	for _, m := range got {
		seen[m["handle"].(string)] = m["value"]
	}
	if seen["a"] != "1" || seen["b"] != 2 {
		t.Errorf("attrs = %#v", got)
	}
}
