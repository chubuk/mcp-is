package mcp

import (
	"context"
	"fmt"

	"github.com/your-project/mcp-insales/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) HandleTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	switch toolName {
	case "get_shop_info":
		return h.service.GetStoreDetails(ctx)

	case "list_products":
		page := 1
		perPage := 25
		if v, ok := args["page"].(float64); ok {
			page = int(v)
		}
		if v, ok := args["per_page"].(float64); ok {
			perPage = int(v)
		}
		return h.service.GetProducts(ctx, page, perPage)

	case "get_product":
		productID, _ := args["product_id"].(string)
		if productID == "" {
			return nil, fmt.Errorf("missing required parameter: product_id")
		}
		return h.service.GetProduct(ctx, productID)

	case "create_product":
		title, _ := args["title"].(string)
		price, _ := args["price"].(string)
		if title == "" || price == "" {
			return nil, fmt.Errorf("missing required parameters: title, price")
		}
		return h.service.CreateProduct(ctx, title, price, args)

	case "update_product":
		productID, _ := args["product_id"].(string)
		if productID == "" {
			return nil, fmt.Errorf("missing required parameter: product_id")
		}
		return h.service.UpdateProduct(ctx, productID, args)

	case "delete_product":
		productID, _ := args["product_id"].(string)
		if productID == "" {
			return nil, fmt.Errorf("missing required parameter: product_id")
		}
		return nil, h.service.DeleteProduct(ctx, productID)

	case "list_orders":
		page := 1
		perPage := 25
		if v, ok := args["page"].(float64); ok {
			page = int(v)
		}
		if v, ok := args["per_page"].(float64); ok {
			perPage = int(v)
		}
		return h.service.GetOrders(ctx, page, perPage)

	case "get_order":
		orderID, _ := args["order_id"].(string)
		if orderID == "" {
			return nil, fmt.Errorf("missing required parameter: order_id")
		}
		return h.service.GetOrderDetails(ctx, orderID)

	case "update_order":
		orderID, _ := args["order_id"].(string)
		if orderID == "" {
			return nil, fmt.Errorf("missing required parameter: order_id")
		}
		return h.service.UpdateOrder(ctx, orderID, args)

	case "create_order":
		clientID, _ := args["client_id"].(float64)
		if clientID == 0 {
			return nil, fmt.Errorf("missing required parameter: client_id")
		}
		return h.service.CreateOrder(ctx, int(clientID), args)

	case "list_clients":
		page := 1
		perPage := 25
		if v, ok := args["page"].(float64); ok {
			page = int(v)
		}
		if v, ok := args["per_page"].(float64); ok {
			perPage = int(v)
		}
		return h.service.GetClients(ctx, page, perPage)

	case "get_client":
		clientID, _ := args["client_id"].(string)
		if clientID == "" {
			return nil, fmt.Errorf("missing required parameter: client_id")
		}
		return h.service.GetClient(ctx, clientID)

	case "create_client":
		email, _ := args["email"].(string)
		if email == "" {
			return nil, fmt.Errorf("missing required parameter: email")
		}
		return h.service.CreateClient(ctx, email, args)

	case "update_client":
		clientID, _ := args["client_id"].(string)
		if clientID == "" {
			return nil, fmt.Errorf("missing required parameter: client_id")
		}
		return h.service.UpdateClient(ctx, clientID, args)

	case "delete_client":
		clientID, _ := args["client_id"].(string)
		if clientID == "" {
			return nil, fmt.Errorf("missing required parameter: client_id")
		}
		return nil, h.service.DeleteClient(ctx, clientID)

	case "list_categories":
		page := 1
		perPage := 25
		if v, ok := args["page"].(float64); ok {
			page = int(v)
		}
		if v, ok := args["per_page"].(float64); ok {
			perPage = int(v)
		}
		return h.service.GetCategories(ctx, page, perPage)

	case "get_category":
		categoryID, _ := args["category_id"].(string)
		if categoryID == "" {
			return nil, fmt.Errorf("missing required parameter: category_id")
		}
		return h.service.GetCategory(ctx, categoryID)

	case "create_category":
		title, _ := args["title"].(string)
		if title == "" {
			return nil, fmt.Errorf("missing required parameter: title")
		}
		return h.service.CreateCategory(ctx, title, args)

	case "update_category":
		categoryID, _ := args["category_id"].(string)
		if categoryID == "" {
			return nil, fmt.Errorf("missing required parameter: category_id")
		}
		return h.service.UpdateCategory(ctx, categoryID, args)

	case "delete_category":
		categoryID, _ := args["category_id"].(string)
		if categoryID == "" {
			return nil, fmt.Errorf("missing required parameter: category_id")
		}
		return nil, h.service.DeleteCategory(ctx, categoryID)

	case "list_collections":
		page := 1
		perPage := 25
		if v, ok := args["page"].(float64); ok {
			page = int(v)
		}
		if v, ok := args["per_page"].(float64); ok {
			perPage = int(v)
		}
		return h.service.GetCollections(ctx, page, perPage)

	case "get_collection":
		collectionID, _ := args["collection_id"].(string)
		if collectionID == "" {
			return nil, fmt.Errorf("missing required parameter: collection_id")
		}
		return h.service.GetCollection(ctx, collectionID)

	case "create_collection":
		title, _ := args["title"].(string)
		if title == "" {
			return nil, fmt.Errorf("missing required parameter: title")
		}
		return h.service.CreateCollection(ctx, title)

	case "update_collection":
		collectionID, _ := args["collection_id"].(string)
		if collectionID == "" {
			return nil, fmt.Errorf("missing required parameter: collection_id")
		}
		title, _ := args["title"].(string)
		return h.service.UpdateCollection(ctx, collectionID, title)

	case "delete_collection":
		collectionID, _ := args["collection_id"].(string)
		if collectionID == "" {
			return nil, fmt.Errorf("missing required parameter: collection_id")
		}
		return nil, h.service.DeleteCollection(ctx, collectionID)

	case "list_images":
		productID, _ := args["product_id"].(string)
		if productID == "" {
			return nil, fmt.Errorf("missing required parameter: product_id")
		}
		return h.service.ListImages(ctx, productID)

	case "add_image":
		productID, _ := args["product_id"].(string)
		url, _ := args["url"].(string)
		if productID == "" || url == "" {
			return nil, fmt.Errorf("missing required parameters: product_id, url")
		}
		return h.service.AddImage(ctx, productID, url, args)

	case "delete_image":
		productID, _ := args["product_id"].(string)
		imageID, _ := args["image_id"].(string)
		if productID == "" || imageID == "" {
			return nil, fmt.Errorf("missing required parameters: product_id, image_id")
		}
		return nil, h.service.DeleteImage(ctx, productID, imageID)

	case "list_meta_fields":
		return h.service.ListMetaFields(ctx, args)

	case "create_meta_field":
		return h.service.CreateMetaField(ctx, args)

	case "update_meta_field":
		metaFieldID, _ := args["meta_field_id"].(string)
		if metaFieldID == "" {
			return nil, fmt.Errorf("missing required parameter: meta_field_id")
		}
		return h.service.UpdateMetaField(ctx, metaFieldID, args)

	case "delete_meta_field":
		metaFieldID, _ := args["meta_field_id"].(string)
		if metaFieldID == "" {
			return nil, fmt.Errorf("missing required parameter: meta_field_id")
		}
		return nil, h.service.DeleteMetaField(ctx, metaFieldID)

	case "list_webhooks":
		return h.service.ListWebhooks(ctx)

	case "create_webhook":
		event, _ := args["event"].(string)
		url, _ := args["url"].(string)
		if event == "" || url == "" {
			return nil, fmt.Errorf("missing required parameters: event, url")
		}
		return h.service.CreateWebhook(ctx, event, url, args)

	case "update_webhook":
		webhookID, _ := args["webhook_id"].(string)
		if webhookID == "" {
			return nil, fmt.Errorf("missing required parameter: webhook_id")
		}
		return h.service.UpdateWebhook(ctx, webhookID, args)

	case "delete_webhook":
		webhookID, _ := args["webhook_id"].(string)
		if webhookID == "" {
			return nil, fmt.Errorf("missing required parameter: webhook_id")
		}
		return nil, h.service.DeleteWebhook(ctx, webhookID)

	case "list_currencies":
		return h.service.ListCurrencies(ctx)

	case "list_delivery_variants":
		return h.service.ListDeliveryVariants(ctx)

	case "list_payment_gateways":
		return h.service.ListPaymentGateways(ctx)

	case "list_pages":
		page, _ := args["page"].(float64)
		perPage, _ := args["per_page"].(float64)
		return h.service.ListPages(ctx, int(page), int(perPage))

	case "get_page":
		pageID, ok := args["page_id"].(string)
		if !ok {
			return nil, fmt.Errorf("page_id is required")
		}
		return h.service.GetPage(ctx, pageID)

	case "create_page":
		title, _ := args["title"].(string)
		content, _ := args["content"].(string)
		handle, _ := args["handle"].(string)
		return h.service.CreatePage(ctx, title, content, handle)

	case "update_page":
		pageID, ok := args["page_id"].(string)
		if !ok {
			return nil, fmt.Errorf("page_id is required")
		}
		updateArgs := make(map[string]interface{})
		for k, v := range args {
			if k != "page_id" {
				updateArgs[k] = v
			}
		}
		return h.service.UpdatePage(ctx, pageID, updateArgs)

	case "delete_page":
		pageID, ok := args["page_id"].(string)
		if !ok {
			return nil, fmt.Errorf("page_id is required")
		}
		return nil, h.service.DeletePage(ctx, pageID)

	case "list_accessories":
		page, _ := args["page"].(float64)
		perPage, _ := args["per_page"].(float64)
		return h.service.ListAccessories(ctx, int(page), int(perPage))

	case "get_accessory":
		accessoryID, ok := args["accessory_id"].(string)
		if !ok {
			return nil, fmt.Errorf("accessory_id is required")
		}
		return h.service.GetAccessory(ctx, accessoryID)

	case "create_accessory":
		title, _ := args["title"].(string)
		return h.service.CreateAccessory(ctx, title)

	case "update_accessory":
		accessoryID, ok := args["accessory_id"].(string)
		if !ok {
			return nil, fmt.Errorf("accessory_id is required")
		}
		updateArgs := make(map[string]interface{})
		for k, v := range args {
			if k != "accessory_id" {
				updateArgs[k] = v
			}
		}
		return h.service.UpdateAccessory(ctx, accessoryID, updateArgs)

	case "delete_accessory":
		accessoryID, ok := args["accessory_id"].(string)
		if !ok {
			return nil, fmt.Errorf("accessory_id is required")
		}
		return nil, h.service.DeleteAccessory(ctx, accessoryID)

	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}
