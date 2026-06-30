package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadOnlyBlocksWriteTools(t *testing.T) {
	writeToolNames := []string{
		"create_product", "update_product", "delete_product",
		"create_order", "update_order",
		"create_client", "update_client", "delete_client",
		"create_category", "update_category", "delete_category",
		"create_collection", "update_collection", "delete_collection",
		"add_image", "delete_image",
		"create_meta_field", "update_meta_field", "delete_meta_field",
		"create_webhook", "update_webhook", "delete_webhook",
		"create_page", "update_page", "delete_page",
		"create_accessory", "update_accessory", "delete_accessory",
	}

	for _, toolName := range writeToolNames {
		t.Run(toolName, func(t *testing.T) {
			server := NewServer(nil, true)

			body, _ := json.Marshal(map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "tools/call",
				"id":      1,
				"params": map[string]interface{}{
					"name":      toolName,
					"arguments": map[string]interface{}{},
				},
			})

			req := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.HandleMCP(w, req)

			var resp map[string]interface{}
			json.NewDecoder(w.Body).Decode(&resp)

			if _, hasError := resp["error"]; !hasError {
				t.Errorf("tool '%s' should be blocked in read-only mode, but no error returned", toolName)
			}
		})
	}
}

func TestReadOnlyAllowsReadTools(t *testing.T) {
	readToolNames := []string{
		"list_products", "get_product",
		"list_orders", "get_order",
		"list_clients", "get_client",
		"list_categories", "get_category",
		"list_collections", "get_collection",
		"list_images",
		"list_meta_fields",
		"list_webhooks",
		"list_pages", "get_page",
		"list_accessories", "get_accessory",
		"get_shop_info",
		"list_currencies", "list_delivery_variants", "list_payment_gateways",
	}

	for _, toolName := range readToolNames {
		t.Run(toolName, func(t *testing.T) {
			if writeTools[toolName] {
				t.Errorf("read tool '%s' is incorrectly marked as write tool", toolName)
			}
		})
	}
}

func TestWriteToolsListComplete(t *testing.T) {
	expectedWriteTools := map[string]bool{
		"create_product":      true,
		"update_product":      true,
		"delete_product":      true,
		"create_order":        true,
		"update_order":        true,
		"create_client":       true,
		"update_client":       true,
		"delete_client":       true,
		"create_category":     true,
		"update_category":     true,
		"delete_category":     true,
		"create_collection":   true,
		"update_collection":   true,
		"delete_collection":   true,
		"add_image":           true,
		"delete_image":        true,
		"create_meta_field":   true,
		"update_meta_field":   true,
		"delete_meta_field":   true,
		"create_webhook":      true,
		"update_webhook":      true,
		"delete_webhook":      true,
		"create_page":         true,
		"update_page":         true,
		"delete_page":         true,
		"create_accessory":    true,
		"update_accessory":    true,
		"delete_accessory":    true,
	}

	for tool := range expectedWriteTools {
		if !writeTools[tool] {
			t.Errorf("tool '%s' missing from writeTools map", tool)
		}
	}

	for tool := range writeTools {
		if !expectedWriteTools[tool] {
			t.Errorf("tool '%s' in writeTools but not in expected list", tool)
		}
	}
}
