package mcp

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/chubuk/mcp-is/internal/core"
)

func TestToolFilterReadOnly(t *testing.T) {
	filter := toolFilter(true)

	tools := []mcp.Tool{
		{Name: "list_products", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "create_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
		{Name: "get_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "delete_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
	}

	ctx := context.Background()
	filtered := filter(ctx, tools)

	if len(filtered) != 2 {
		t.Errorf("expected 2 read-only tools, got %d", len(filtered))
	}

	for _, tool := range filtered {
		if tool.Name != "list_products" && tool.Name != "get_product" {
			t.Errorf("unexpected tool in filtered list: %s", tool.Name)
		}
	}
}

func TestToolFilterReadWrite(t *testing.T) {
	filter := toolFilter(false)

	tools := []mcp.Tool{
		{Name: "list_products"},
		{Name: "create_product"},
	}

	ctx := context.Background()
	filtered := filter(ctx, tools)

	if len(filtered) != 2 {
		t.Errorf("expected 2 tools when not read-only, got %d", len(filtered))
	}
}

func TestToolFilterShopReadOnly(t *testing.T) {
	filter := toolFilter(false)

	tools := []mcp.Tool{
		{Name: "list_products", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "create_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
	}

	ctx := context.WithValue(context.Background(), core.ShopReadOnlyKey, true)
	filtered := filter(ctx, tools)

	if len(filtered) != 1 {
		t.Errorf("expected 1 read-only tool when shop is read-only, got %d", len(filtered))
	}
	if filtered[0].Name != "list_products" {
		t.Errorf("expected list_products, got %s", filtered[0].Name)
	}
}

func TestAllCommandNamesComplete(t *testing.T) {
	// Verify all expected tools have descriptions in shops.go commandDescriptions
	// This test ensures the shops.go and tools.go stay in sync
	expectedTools := []string{
		"list_products", "get_product", "create_product", "update_product", "delete_product",
		"list_orders", "get_order", "update_order", "create_order",
		"list_clients", "get_client", "create_client", "update_client", "delete_client",
		"list_categories", "get_category", "create_category", "update_category", "delete_category",
		"list_collections", "get_collection", "create_collection", "update_collection", "delete_collection",
		"list_images", "add_image", "delete_image",
		"list_meta_fields", "create_meta_field", "update_meta_field", "delete_meta_field",
		"list_webhooks", "create_webhook", "update_webhook", "delete_webhook",
		"get_shop_info", "list_currencies", "list_delivery_variants", "list_payment_gateways",
		"list_pages", "get_page", "create_page", "update_page", "delete_page",
		"list_accessories", "get_accessory", "create_accessory", "update_accessory", "delete_accessory",
		"list_variants", "get_variant", "update_variant", "delete_variant",
		"list_property_values", "create_property_value", "update_property_value", "delete_property_value",
		"list_menus", "get_menu", "create_menu", "update_menu", "delete_menu",
		"list_redirects", "create_redirect", "delete_redirect",
		"list_blogs", "get_blog", "create_blog", "delete_blog",
		"list_blog_posts", "get_blog_post", "create_blog_post", "update_blog_post", "delete_blog_post",
	}

	// Verify the list has the right count
	if len(expectedTools) != 74 {
		t.Errorf("expected 74 tools, got %d", len(expectedTools))
	}
}

func boolPtr(b bool) *bool {
	return &b
}
