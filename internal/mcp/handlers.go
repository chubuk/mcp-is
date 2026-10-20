package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
)

func formatResult(v any) string {
	if v == nil {
		return "null"
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

// toolResult returns the same payload as pretty JSON text and as structuredContent
// (MCP dual form; clients may use either).
func toolResult(v any) *mcp.CallToolResult {
	return mcp.NewToolResultStructured(v, formatResult(v))
}

func isHTMLResponse(s string) bool {
	trimmed := strings.TrimSpace(s)
	return strings.HasPrefix(trimmed, "<!DOCTYPE") ||
		strings.HasPrefix(trimmed, "<html") ||
		strings.HasPrefix(trimmed, "<HTML") ||
		strings.Contains(s, "<!DOCTYPE html") ||
		strings.Contains(s, "<!DOCTYPE HTML")
}

func formatResultOrError(v any) (string, bool) {
	s := formatResult(v)
	if isHTMLResponse(s) {
		return "InSales API returned an error page. This endpoint may not be available for your API key.", true
	}
	return s, false
}

func cleanError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if isHTMLResponse(msg) {
		// Extract status code if present (e.g. "API error 404: <html...")
		if idx := strings.Index(msg, "<!"); idx > 0 {
			return fmt.Errorf("InSales API error: %s", strings.TrimSpace(msg[:idx]))
		}
		return fmt.Errorf("InSales API returned an error page (endpoint may not be available for this API key)")
	}
	return err
}

func toolError(err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(cleanError(err).Error())
}

func getFloat(args map[string]any, key string, def float64) float64 {
	if v, ok := args[key].(float64); ok {
		return v
	}
	return def
}

func applyPDNToResult(ctx context.Context, result any, toolName string) any {
	shopPDN, _ := ctx.Value(mcpShopPDNKey).(bool)
	aiKeyPDN, _ := ctx.Value(mcpAIKeyPDNKey).(bool)

	if !shopPDN && !aiKeyPDN {
		return result
	}

	if dataMap, ok := result.(map[string]any); ok {
		return pdn.MaskClientData(dataMap)
	}
	if dataArray, ok := result.([]any); ok {
		masked := make([]any, len(dataArray))
		for i, item := range dataArray {
			if dataMap, ok := item.(map[string]any); ok {
				masked[i] = pdn.MaskClientData(dataMap)
			} else {
				masked[i] = item
			}
		}
		return masked
	}
	return result
}

// === ТОВАРЫ ===

func listProductsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.GetProducts(ctx, page, perPage, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetProduct(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		price, err := req.RequireString("price")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateProduct(ctx, title, price, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateProduct(ctx, productID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteProduct(ctx, productID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Товар удалён"), nil
	}
}

// === ЗАКАЗЫ ===

func listOrdersHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		filters := make(map[string]any)
		if v, ok := args["fulfillment_status"].(string); ok && v != "" {
			filters["fulfillment_status"] = v
		}
		if v, ok := args["updated_since"].(string); ok && v != "" {
			filters["updated_since"] = v
		}
		if v, ok := args["from_id"].(float64); ok && v > 0 {
			filters["from_id"] = int(v)
		}
		result, err := svc.GetOrders(ctx, page, perPage, filters)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "list_orders")
		return toolResult(masked), nil
	}
}

func getOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetOrderDetails(ctx, orderID)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "get_order")
		return toolResult(masked), nil
	}
}

func updateOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "order_id")
		result, err := svc.UpdateOrder(ctx, orderID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		clientID := int(getFloat(args, "client_id", 0))
		// client_id optional when client / client_attributes present
		if clientID == 0 {
			if _, hasClient := args["client"]; !hasClient {
				if _, hasAttrs := args["client_attributes"]; !hasAttrs {
					if _, hasEmail := args["client_email"]; !hasEmail {
						return mcp.NewToolResultError("client_id or client (name/email/phone) is required"), nil
					}
				}
			}
		}
		result, err := svc.CreateOrder(ctx, clientID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === КЛИЕНТЫ ===

func listClientsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.GetClients(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "list_clients")
		return toolResult(masked), nil
	}
}

func getClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetClient(ctx, clientID)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "get_client")
		return toolResult(masked), nil
	}
}

func createClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateClient(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "create_client")
		return toolResult(masked), nil
	}
}

func updateClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateClient(ctx, clientID, args)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "update_client")
		return toolResult(masked), nil
	}
}

func deleteClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteClient(ctx, clientID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Клиент удалён"), nil
	}
}

// === КАТЕГОРИИ ===

func listCategoriesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.GetCategories(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCategory(ctx, categoryID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateCategory(ctx, title, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateCategory(ctx, categoryID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteCategory(ctx, categoryID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Категория удалена"), nil
	}
}

// === КОЛЛЕКЦИИ ===

func listCollectionsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.GetCollections(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCollection(ctx, collectionID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateCollection(ctx, title)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		title, _ := req.RequireString("title")
		result, err := svc.UpdateCollection(ctx, collectionID, title)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteCollection(ctx, collectionID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Коллекция удалена"), nil
	}
}

// === ИЗОБРАЖЕНИЯ ===

func listImagesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListImages(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func addImageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.AddImage(ctx, productID, url, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteImageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		imageID, err := req.RequireString("image_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteImage(ctx, productID, imageID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Изображение удалено"), nil
	}
}

// === МЕТАПОЛЯ ===

func listMetaFieldsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.ListMetaFields(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateMetaField(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metaFieldID, err := req.RequireString("meta_field_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateMetaField(ctx, metaFieldID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metaFieldID, err := req.RequireString("meta_field_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteMetaField(ctx, metaFieldID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Метаполе удалено"), nil
	}
}

// === ВЕБХУКИ ===

func listWebhooksHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListWebhooks(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		event, err := req.RequireString("event")
		if err != nil {
			return toolError(err), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateWebhook(ctx, event, url, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhookID, err := req.RequireString("webhook_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateWebhook(ctx, webhookID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhookID, err := req.RequireString("webhook_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteWebhook(ctx, webhookID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Вебхук удалён"), nil
	}
}

// === СИСТЕМНЫЕ ===

func getShopInfoHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.GetStoreDetails(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listCurrenciesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListCurrencies(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listDeliveryVariantsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListDeliveryVariants(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listPaymentGatewaysHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListPaymentGateways(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === СТРАНИЦЫ ===

func listPagesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.ListPages(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getPageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetPage(ctx, pageID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createPageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		content, err := req.RequireString("content")
		if err != nil {
			return toolError(err), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreatePage(ctx, title, content, handle)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updatePageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdatePage(ctx, pageID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deletePageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeletePage(ctx, pageID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Страница удалена"), nil
	}
}

// === АКСЕССУАРЫ ===

func listAccessoriesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.ListAccessories(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetAccessory(ctx, accessoryID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateAccessory(ctx, title, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateAccessory(ctx, accessoryID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteAccessory(ctx, accessoryID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Аксессуар удалён"), nil
	}
}

// === ВАРИАНТЫ ===

func listVariantsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListVariants(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetVariant(ctx, productID, variantID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateVariant(ctx, productID, variantID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteVariant(ctx, productID, variantID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Вариант удалён"), nil
	}
}

// === ЗНАЧЕНИЯ СВОЙСТВ ===

func listPropertyValuesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListPropertyValues(ctx, accessoryID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createPropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreatePropertyValue(ctx, accessoryID, title)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updatePropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		propertyValueID, err := req.RequireString("property_value_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdatePropertyValue(ctx, accessoryID, propertyValueID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deletePropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		propertyValueID, err := req.RequireString("property_value_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeletePropertyValue(ctx, accessoryID, propertyValueID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Значение удалено"), nil
	}
}

// === МЕНЮ ===

func listMenusHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListMenus(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetMenu(ctx, menuID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreateMenu(ctx, title, handle)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateMenu(ctx, menuID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteMenu(ctx, menuID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Меню удалено"), nil
	}
}

// === РЕДИРЕКТЫ ===

func listRedirectsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListRedirects(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createRedirectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := req.RequireString("url")
		if err != nil {
			return toolError(err), nil
		}
		target, err := req.RequireString("target")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateRedirect(ctx, url, target)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteRedirectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		redirectID, err := req.RequireString("redirect_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteRedirect(ctx, redirectID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Редирект удалён"), nil
	}
}

// === БЛОГИ ===

func listBlogsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListBlogs(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetBlog(ctx, blogID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreateBlog(ctx, title, handle)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteBlog(ctx, blogID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Блог удалён"), nil
	}
}

func updateBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateBlog(ctx, blogID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === ЗАПИСИ БЛОГА ===

func listBlogPostsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListBlogPosts(ctx, blogID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetBlogPost(ctx, blogID, postID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		bodyHTML, err := req.RequireString("body_html")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateBlogPost(ctx, blogID, title, bodyHTML)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateBlogPost(ctx, blogID, postID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteBlogPost(ctx, blogID, postID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Запись удалена"), nil
	}
}

// === DISCOUNT CODES ===

func listDiscountCodesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.ListDiscountCodes(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codeID, err := req.RequireString("discount_code_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetDiscountCode(ctx, codeID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateDiscountCode(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codeID, err := req.RequireString("discount_code_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateDiscountCode(ctx, codeID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codeID, err := req.RequireString("discount_code_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteDiscountCode(ctx, codeID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Промокод удалён"), nil
	}
}

// === COLLECTS ===

func listCollectsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.ListCollects(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCollectsByProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCollectsByProduct(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCollectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		collectionID := int(getFloat(args, "collection_id", 0))
		productID := int(getFloat(args, "product_id", 0))
		if collectionID == 0 || productID == 0 {
			return toolError(fmt.Errorf("collection_id and product_id are required")), nil
		}
		result, err := svc.CreateCollect(ctx, collectionID, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteCollectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectID, err := req.RequireString("collect_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteCollect(ctx, collectID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Связь удалена"), nil
	}
}

// === PROPERTIES ===

func listPropertiesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListProperties(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getPropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetProperty(ctx, propertyID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createPropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateProperty(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updatePropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateProperty(ctx, propertyID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deletePropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteProperty(ctx, propertyID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Свойство удалено"), nil
	}
}

// === CHARACTERISTICS ===

func listCharacteristicsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListCharacteristics(ctx, propertyID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		charID, err := req.RequireString("characteristic_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCharacteristic(ctx, propertyID, charID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateCharacteristic(ctx, propertyID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		charID, err := req.RequireString("characteristic_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateCharacteristic(ctx, propertyID, charID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		charID, err := req.RequireString("characteristic_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteCharacteristic(ctx, propertyID, charID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Характеристика удалена"), nil
	}
}

// === REVIEWS ===

func listReviewsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page := int(getFloat(args, "page", 1))
		perPage := int(getFloat(args, "per_page", 25))
		result, err := svc.ListReviews(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reviewID, err := req.RequireString("review_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetReview(ctx, reviewID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateReview(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reviewID, err := req.RequireString("review_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateReview(ctx, reviewID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reviewID, err := req.RequireString("review_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteReview(ctx, reviewID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Отзыв удалён"), nil
	}
}

// === FILES ===

func listFilesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListFiles(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getFileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID, err := req.RequireString("file_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetFile(ctx, fileID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createFileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateFile(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteFileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID, err := req.RequireString("file_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteFile(ctx, fileID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Файл удалён"), nil
	}
}

// Generic CRUD handlers for round 3+ resources

func genericListHandler(svc *service.Service, listFn func(ctx context.Context) ([]any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := listFn(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericGetHandler(svc *service.Service, param string, getFn func(ctx context.Context, id string) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString(param)
		if err != nil {
			return toolError(err), nil
		}
		result, err := getFn(ctx, id)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericCreateHandler(svc *service.Service, createFn func(ctx context.Context, args map[string]any) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := createFn(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericUpdateHandler(svc *service.Service, param string, updateFn func(ctx context.Context, id string, args map[string]any) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString(param)
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := updateFn(ctx, id, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericDeleteHandler(deleteFn func(ctx context.Context, id string) error, paramName, deletedMsg string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString(paramName)
		if err != nil {
			return toolError(err), nil
		}
		if err := deleteFn(ctx, id); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(deletedMsg), nil
	}
}

// nested helpers for parent-scoped resources (product_id + value_id, etc.)

func nestedListHandler(listFn func(ctx context.Context, parentID string) ([]any, error), parentParam string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		result, err := listFn(ctx, parentID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(map[string]any{
			"items": result,
			"total": len(result),
		}), nil
	}
}

func nestedGetHandler(getFn func(ctx context.Context, parentID, key string) (any, error), parentParam, keyParam string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		key, err := req.RequireString(keyParam)
		if err != nil {
			return toolError(err), nil
		}
		result, err := getFn(ctx, parentID, key)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func nestedUpdateHandler(updateFn func(ctx context.Context, parentID, key string, args map[string]any) (any, error), parentParam, keyParam string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		key, err := req.RequireString(keyParam)
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, parentParam)
		delete(args, keyParam)
		// body is only the value payload
		payload := map[string]any{}
		if v, ok := args["value"]; ok {
			payload["value"] = v
		}
		for k, v := range args {
			if k != "value" && k != "name" && k != "args" {
				payload[k] = v
			}
		}
		result, err := updateFn(ctx, parentID, key, payload)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func nestedDeleteHandler(deleteFn func(ctx context.Context, parentID, key string) error, parentParam, keyParam, deletedMsg string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		key, err := req.RequireString(keyParam)
		if err != nil {
			return toolError(err), nil
		}
		if err := deleteFn(ctx, parentID, key); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(deletedMsg), nil
	}
}

func createProductFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		if _, ok := args["type"]; !ok {
			if ft, ok := args["field_type"].(string); ok && ft != "" {
				args["type"] = ft
			}
		}
		delete(args, "field_type")
		if _, ok := args["handle"].(string); !ok {
			return mcp.NewToolResultError("handle is required (language-invariant slug)"), nil
		}
		if _, ok := args["type"].(string); !ok {
			return mcp.NewToolResultError("type is required: ProductField::TextField | ProductField::TextArea"), nil
		}
		result, err := svc.CreateProductField(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateProductFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fieldID, err := req.RequireString("field_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "field_id")
		delete(args, "field_type")
		if ft, ok := req.GetArguments()["field_type"].(string); ok && ft != "" {
			if _, has := args["type"]; !has {
				args["type"] = ft
			}
		}
		result, err := svc.UpdateProductField(ctx, fieldID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createProductFieldValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "product_id")
		if _, ok := args["product_field_id"]; !ok {
			if _, ok := args["product_field_handle"]; !ok {
				return mcp.NewToolResultError("product_field_id or product_field_handle is required"), nil
			}
		}
		result, err := svc.CreateProductFieldValue(ctx, productID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createVariantFieldValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "variant_id")
		if _, ok := args["variant_field_id"]; !ok {
			if _, ok := args["variant_field_handle"]; !ok {
				return mcp.NewToolResultError("variant_field_id or variant_field_handle is required"), nil
			}
		}
		result, err := svc.CreateVariantFieldValue(ctx, variantID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === ORDER SUB-OPERATIONS (map into order[...] via UpdateOrder) ===

func orderSubHandler(svc *service.Service, build func(req mcp.CallToolRequest) (orderID string, body map[string]any, err error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, body, err := build(req)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.UpdateOrder(ctx, orderID, body)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func strArg(req mcp.CallToolRequest, key string) string {
	return req.GetString(key, "")
}

func floatArg(req mcp.CallToolRequest, key string) (float64, bool) {
	args := req.GetArguments()
	v, ok := args[key].(float64)
	return v, ok
}

func addOrderLineHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		productID := strArg(req, "product_id")
		variantID := strArg(req, "variant_id")
		if productID == "" && variantID == "" {
			return "", nil, errString("product_id or variant_id is required")
		}
		if productID != "" && variantID != "" {
			return "", nil, errString("pass only one of product_id or variant_id")
		}
		line := map[string]any{}
		if productID != "" {
			line["product_id"] = productID
		}
		if variantID != "" {
			line["variant_id"] = variantID
		}
		if q, ok := floatArg(req, "quantity"); ok {
			line["quantity"] = q
		} else {
			return "", nil, errString("quantity is required")
		}
		if v, ok := floatArg(req, "vat"); ok {
			line["vat"] = v
		}
		if p, ok := floatArg(req, "sale_price"); ok {
			line["sale_price"] = p
		}
		return orderID, map[string]any{"order_lines_attributes": []any{line}}, nil
	})
}

func updateOrderLineHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		lineID, err := req.RequireString("line_id")
		if err != nil {
			return "", nil, err
		}
		line := map[string]any{"id": lineID}
		if q, ok := floatArg(req, "quantity"); ok {
			line["quantity"] = q
		}
		if v, ok := floatArg(req, "vat"); ok {
			line["vat"] = v
		}
		if p, ok := floatArg(req, "sale_price"); ok {
			line["sale_price"] = p
		}
		if c := strArg(req, "comment"); c != "" {
			line["comment"] = c
		}
		return orderID, map[string]any{"order_lines_attributes": []any{line}}, nil
	})
}

func removeOrderLineHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		lineID, err := req.RequireString("line_id")
		if err != nil {
			return "", nil, err
		}
		return orderID, map[string]any{
			"order_lines_attributes": []any{
				map[string]any{"id": lineID, "_destroy": true},
			},
		}, nil
	})
}

func addOrderDiscountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		d, ok := floatArg(req, "discount")
		if !ok {
			return "", nil, errString("discount is required")
		}
		item := map[string]any{"discount": d}
		if t, ok := floatArg(req, "type_id"); ok {
			item["type_id"] = t
		}
		if desc := strArg(req, "description"); desc != "" {
			item["description"] = desc
		}
		if ids, ok := req.GetArguments()["line_ids"].([]any); ok && len(ids) > 0 {
			item["discount_order_lines_ids"] = ids
		}
		return orderID, map[string]any{"discounts_attributes": []any{item}}, nil
	})
}

func updateOrderDiscountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		discountID, err := req.RequireString("discount_id")
		if err != nil {
			return "", nil, err
		}
		item := map[string]any{"id": discountID}
		if d, ok := floatArg(req, "discount"); ok {
			item["discount"] = d
		}
		if t, ok := floatArg(req, "type_id"); ok {
			item["type_id"] = t
		}
		if desc := strArg(req, "description"); desc != "" {
			item["description"] = desc
		}
		return orderID, map[string]any{"discounts_attributes": []any{item}}, nil
	})
}

func removeOrderDiscountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		discountID, err := req.RequireString("discount_id")
		if err != nil {
			return "", nil, err
		}
		return orderID, map[string]any{
			"discounts_attributes": []any{
				map[string]any{"id": discountID, "_destroy": true},
			},
		}, nil
	})
}

func updateOrderShippingAddressHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		addr := map[string]any{}
		args := req.GetArguments()
		for _, k := range []string{"address", "full_locality_name", "country", "region", "city", "zip", "phone", "full_name", "comment"} {
			if v, ok := args[k]; ok && v != nil && v != "" {
				addr[k] = v
			}
		}
		if len(addr) == 0 {
			return "", nil, errString("at least one address field is required")
		}
		return orderID, map[string]any{"shipping_address_attributes": addr}, nil
	})
}

func updateOrderCustomStatusHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		permalink, err := req.RequireString("custom_status_permalink")
		if err != nil {
			return "", nil, err
		}
		body := map[string]any{"custom_status_permalink": permalink}
		if fs := strArg(req, "fulfillment_status"); fs != "" {
			body["fulfillment_status"] = fs
		}
		return orderID, body, nil
	})
}

type errString string

func (e errString) Error() string { return string(e) }

func createMenuItemHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "menu_id")
		rel, _ := args["related_object_id"].(string)
		if rel == "" {
			return mcp.NewToolResultError("related_object_id is required (Blog_1, Page_2, Collection_3, or URL)"), nil
		}
		if !looksLikeRelatedObjectID(rel) {
			return mcp.NewToolResultError(`related_object_id must be composite like "Blog_1", "Page_2", "Collection_3", or an absolute URL`), nil
		}
		result, err := svc.CreateMenuItem(ctx, menuID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func looksLikeRelatedObjectID(rel string) bool {
	if strings.HasPrefix(rel, "http://") || strings.HasPrefix(rel, "https://") {
		return true
	}
	// Type_id composite: Blog_1, Page_2, Collection_3, ...
	i := strings.IndexByte(rel, '_')
	if i <= 0 || i == len(rel)-1 {
		return false
	}
	for _, r := range rel[i+1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// === OPTION VALUES handlers (nested resource) ===

func genericListOptionValuesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListOptionValues(ctx, optionNameID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericGetOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		optionValueID, err := req.RequireString("option_value_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetOptionValue(ctx, optionNameID, optionValueID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericCreateOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateOptionValue(ctx, optionNameID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericUpdateOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		optionValueID, err := req.RequireString("option_value_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateOptionValue(ctx, optionNameID, optionValueID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericDeleteOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		optionValueID, err := req.RequireString("option_value_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeleteOptionValue(ctx, optionNameID, optionValueID); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Значение удалено"), nil
	}
}

// === ROUND 4: Special handlers ===

func declineApplicationChargeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("charge_id")
		if err != nil {
			return toolError(err), nil
		}
		if err := svc.DeclineApplicationCharge(ctx, id); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("Счёт отклонён"), nil
	}
}

func createBonusSystemTransactionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		clientID := int(getFloat(args, "client_id", 0))
		if clientID == 0 {
			return toolError(fmt.Errorf("client_id is required")), nil
		}
		result, err := svc.CreateBonusSystemTransaction(ctx, clientID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listClientBonusTransactionsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListClientBonusSystemTransactions(ctx, clientID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listSimilarsByProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListSimilars(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === PROVIDER TOKEN handler ===

func createProviderTokenHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		messengerType, err := req.RequireString("messenger_type")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateProviderToken(ctx, messengerType)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === ACCOUNT handler ===

func updateAccountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.UpdateAccount(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === COUNT handlers ===

func getOrdersCountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		count, err := svc.GetOrdersCount(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(map[string]any{"count": count}), nil
	}
}

func getProductsCountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		count, err := svc.GetProductsCount(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(map[string]any{"count": count}), nil
	}
}

// === PIPELINE HANDLERS ===

func searchProductDetailsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		query, _ := args["query"].(string)
		productID, _ := args["product_id"].(string)
		if query == "" && productID == "" {
			return toolError(fmt.Errorf("query или product_id обязательны")), nil
		}
		result, err := svc.SearchProductDetails(ctx, query, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func orderFullDetailsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.OrderFullDetails(ctx, orderID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func clientFullProfileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ClientFullProfile(ctx, clientID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func productCatalogSearchHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		query, _ := args["query"].(string)
		result, err := svc.ProductCatalogSearch(ctx, query, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func collectionProductsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CollectionProducts(ctx, collectionID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === EXPORT handlers ===

func exportProductsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		csv, err := svc.ExportProducts(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(csv), nil
	}
}

func exportOrdersHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		csv, err := svc.ExportOrders(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(csv), nil
	}
}

func exportClientsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		csv, err := svc.ExportClients(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(csv), nil
	}
}
