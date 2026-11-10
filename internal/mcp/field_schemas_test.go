package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestAugmentFieldValuesSchemas_AddsHandles(t *testing.T) {
	snap := `{"product_fields":{"title":"Дополнительные поля: Товары","items":[
		{"id":"1","title":"Состав","handle":"sostav"},
		{"id":"2","title":"Без handle","handle":""}]},
	 "fields":{"title":"Поля","items":[
		{"id":"3","title":"Комментарий","handle":"comment","destiny":3},
		{"id":"4","title":"Клиентское","handle":"client_note","destiny":6}]}}`
	ctx := context.WithValue(context.Background(), auth.SettingsSnapshotJSONKey, snap)

	tools := []mcp.Tool{
		{Name: "update_product", InputSchema: mcp.ToolInputSchema{Properties: map[string]any{
			"field_values": map[string]any{"type": "object", "description": "base"},
		}}},
		{Name: "update_order", InputSchema: mcp.ToolInputSchema{Properties: map[string]any{
			"field_values": map[string]any{"type": "object", "description": "base"},
		}}},
		{Name: "unrelated", InputSchema: mcp.ToolInputSchema{Properties: map[string]any{
			"title": map[string]any{"type": "string"},
		}}},
	}
	out := augmentFieldValuesSchemas(ctx, tools)

	desc := out[0].InputSchema.Properties["field_values"].(map[string]any)["description"].(string)
	if !strings.Contains(desc, "sostav (Состав)") {
		t.Errorf("product desc missing handle: %q", desc)
	}
	if strings.Contains(desc, "Без handle") {
		t.Errorf("item without handle must be skipped: %q", desc)
	}

	orderDesc := out[1].InputSchema.Properties["field_values"].(map[string]any)["description"].(string)
	if !strings.Contains(orderDesc, "comment (Комментарий)") {
		t.Errorf("order desc missing destiny=3 field: %q", orderDesc)
	}
	if strings.Contains(orderDesc, "client_note") {
		t.Errorf("destiny=6 field must not appear on order tool: %q", orderDesc)
	}

	// исходные тулы не изменены (копии)
	if tools[0].InputSchema.Properties["field_values"].(map[string]any)["description"] != "base" {
		t.Error("shared tool map mutated — race risk")
	}
}

func TestFieldValuesDescription_Empty(t *testing.T) {
	desc := fieldValuesDescription(settingsGroupJSON{}, nil)
	if !strings.Contains(desc, "не найдены") {
		t.Errorf("empty group should hint at creating fields: %q", desc)
	}
}
