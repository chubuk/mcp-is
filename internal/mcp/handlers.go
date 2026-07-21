package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/chubuk/mcp-is/internal/service"
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
		result, err := svc.GetProducts(ctx, page, perPage)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetProduct(ctx, productID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		price, err := req.RequireString("price")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateProduct(ctx, title, price, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateProduct(ctx, productID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteProduct(ctx, productID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
		result, err := svc.GetOrders(ctx, page, perPage)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		masked := applyPDNToResult(ctx, result, "list_orders")
		return mcp.NewToolResultText(formatResult(masked)), nil
	}
}

func getOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetOrderDetails(ctx, orderID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		masked := applyPDNToResult(ctx, result, "get_order")
		return mcp.NewToolResultText(formatResult(masked)), nil
	}
}

func updateOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateOrder(ctx, orderID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		clientID := int(getFloat(args, "client_id", 0))
		if clientID == 0 {
			return mcp.NewToolResultError("client_id is required"), nil
		}
		result, err := svc.CreateOrder(ctx, clientID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}
		masked := applyPDNToResult(ctx, result, "list_clients")
		return mcp.NewToolResultText(formatResult(masked)), nil
	}
}

func getClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetClient(ctx, clientID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		masked := applyPDNToResult(ctx, result, "get_client")
		return mcp.NewToolResultText(formatResult(masked)), nil
	}
}

func createClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		email, err := req.RequireString("email")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateClient(ctx, email, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		masked := applyPDNToResult(ctx, result, "create_client")
		return mcp.NewToolResultText(formatResult(masked)), nil
	}
}

func updateClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateClient(ctx, clientID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		masked := applyPDNToResult(ctx, result, "update_client")
		return mcp.NewToolResultText(formatResult(masked)), nil
	}
}

func deleteClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteClient(ctx, clientID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetCategory(ctx, categoryID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateCategory(ctx, title, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateCategory(ctx, categoryID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteCategory(ctx, categoryID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetCollection(ctx, collectionID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.CreateCollection(ctx, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		title, _ := req.RequireString("title")
		result, err := svc.UpdateCollection(ctx, collectionID, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteCollection(ctx, collectionID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Коллекция удалена"), nil
	}
}

// === ИЗОБРАЖЕНИЯ ===

func listImagesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.ListImages(ctx, productID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func addImageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.AddImage(ctx, productID, url, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteImageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		imageID, err := req.RequireString("image_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteImage(ctx, productID, imageID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateMetaField(ctx, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metaFieldID, err := req.RequireString("meta_field_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateMetaField(ctx, metaFieldID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metaFieldID, err := req.RequireString("meta_field_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteMetaField(ctx, metaFieldID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Метаполе удалено"), nil
	}
}

// === ВЕБХУКИ ===

func listWebhooksHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListWebhooks(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		event, err := req.RequireString("event")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateWebhook(ctx, event, url, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhookID, err := req.RequireString("webhook_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateWebhook(ctx, webhookID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhookID, err := req.RequireString("webhook_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteWebhook(ctx, webhookID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Вебхук удалён"), nil
	}
}

// === СИСТЕМНЫЕ ===

func getShopInfoHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.GetStoreDetails(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func listCurrenciesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListCurrencies(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func listDeliveryVariantsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListDeliveryVariants(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func listPaymentGatewaysHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListPaymentGateways(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getPageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetPage(ctx, pageID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createPageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := req.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreatePage(ctx, title, content, handle)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updatePageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdatePage(ctx, pageID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deletePageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeletePage(ctx, pageID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetAccessory(ctx, accessoryID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.CreateAccessory(ctx, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateAccessory(ctx, accessoryID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteAccessory(ctx, accessoryID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Аксессуар удалён"), nil
	}
}

// === ВАРИАНТЫ ===

func listVariantsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.ListVariants(ctx, productID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetVariant(ctx, productID, variantID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateVariant(ctx, productID, variantID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteVariant(ctx, productID, variantID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Вариант удалён"), nil
	}
}

// === ЗНАЧЕНИЯ СВОЙСТВ ===

func listPropertyValuesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.ListPropertyValues(ctx, accessoryID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createPropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.CreatePropertyValue(ctx, accessoryID, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updatePropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyValueID, err := req.RequireString("property_value_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdatePropertyValue(ctx, propertyValueID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deletePropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyValueID, err := req.RequireString("property_value_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeletePropertyValue(ctx, propertyValueID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Значение удалено"), nil
	}
}

// === МЕНЮ ===

func listMenusHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListMenus(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetMenu(ctx, menuID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreateMenu(ctx, title, handle)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateMenu(ctx, menuID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteMenu(ctx, menuID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Меню удалено"), nil
	}
}

// === РЕДИРЕКТЫ ===

func listRedirectsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListRedirects(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createRedirectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := req.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		target, err := req.RequireString("target")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.CreateRedirect(ctx, url, target)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteRedirectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		redirectID, err := req.RequireString("redirect_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteRedirect(ctx, redirectID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Редирект удалён"), nil
	}
}

// === БЛОГИ ===

func listBlogsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListBlogs(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetBlog(ctx, blogID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreateBlog(ctx, title, handle)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteBlog(ctx, blogID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Блог удалён"), nil
	}
}

// === ЗАПИСИ БЛОГА ===

func listBlogPostsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.ListBlogPosts(ctx, blogID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func getBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.GetBlogPost(ctx, blogID, postID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func createBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		bodyHTML, err := req.RequireString("body_html")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := svc.CreateBlogPost(ctx, blogID, title, bodyHTML)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func updateBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		result, err := svc.UpdateBlogPost(ctx, blogID, postID, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatResult(result)), nil
	}
}

func deleteBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := svc.DeleteBlogPost(ctx, blogID, postID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Запись удалена"), nil
	}
}
