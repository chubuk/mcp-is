package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/mark3labs/mcp-go/mcp"
)

// fieldValuesToolSpec — к какой группе settings_snapshot относится
// параметр field_values тула; destinies фильтрует поля Field (заказ/клиент).
type fieldValuesToolSpec struct {
	group     string
	destinies []int
}

var fieldValuesTools = map[string]fieldValuesToolSpec{
	"create_product":    {group: "product_fields"},
	"update_product":    {group: "product_fields"},
	"update_variant":    {group: "variant_fields"},
	"create_category":   {group: "collection_fields"},
	"update_category":   {group: "collection_fields"},
	"create_collection": {group: "collection_fields"},
	"update_collection": {group: "collection_fields"},
	"create_order":      {group: "fields", destinies: []int{3}},
	"update_order":      {group: "fields", destinies: []int{3}},
	"create_client":     {group: "fields", destinies: []int{2, 5, 6}},
	"update_client":     {group: "fields", destinies: []int{2, 5, 6}},
}

type settingsGroupJSON struct {
	Title string `json:"title"`
	Items []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Handle  string `json:"handle"`
		Destiny int    `json:"destiny"`
	} `json:"items"`
}

// augmentFieldValuesSchemas дописывает в описание параметра field_values
// актуальные handle доп. полей магазина из settings_snapshot (ctx): схема
// меняется вслед за полями магазина (снимок обновляется при «Проверить
// соединение/доступы»). Меняет копии тулов — Properties разделяются между
// запросами, поэтому карты копируются перед записью.
func augmentFieldValuesSchemas(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
	raw, _ := ctx.Value(auth.SettingsSnapshotJSONKey).(string)
	if raw == "" {
		return tools
	}
	var snap map[string]settingsGroupJSON
	if err := json.Unmarshal([]byte(raw), &snap); err != nil {
		return tools
	}
	out := make([]mcp.Tool, len(tools))
	copy(out, tools)
	for i := range out {
		spec, ok := fieldValuesTools[out[i].Name]
		if !ok {
			continue
		}
		props, ok := out[i].InputSchema.Properties["field_values"].(map[string]any)
		if !ok {
			continue
		}
		newProps := make(map[string]any, len(props))
		for k, v := range props {
			newProps[k] = v
		}
		newProps["description"] = fieldValuesDescription(snap[spec.group], spec.destinies)
		newSchemaProps := make(map[string]any, len(out[i].InputSchema.Properties))
		for k, v := range out[i].InputSchema.Properties {
			newSchemaProps[k] = v
		}
		newSchemaProps["field_values"] = newProps
		out[i].InputSchema.Properties = newSchemaProps
	}
	return out
}

func fieldValuesDescription(group settingsGroupJSON, destinies []int) string {
	base := "Значения дополнительных полей магазина: {handle: значение}."
	var parts []string
	for _, it := range group.Items {
		if it.Handle == "" {
			continue
		}
		if len(destinies) > 0 {
			match := false
			for _, d := range destinies {
				if d == it.Destiny {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		title := it.Title
		if title == "" {
			title = it.Handle
		}
		parts = append(parts, fmt.Sprintf("%s (%s)", it.Handle, title))
	}
	if len(parts) == 0 {
		return base + " Поля этого магазина не найдены — создайте их в админке InSales и нажмите «Проверить соединение»."
	}
	return base + " Поля магазина: " + strings.Join(parts, ", ") + "."
}
