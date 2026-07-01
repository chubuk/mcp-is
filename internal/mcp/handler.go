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

func (h *Handler) HandleTool(ctx context.Context, toolName string, args map[string]any) (any, error) {
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
		updateArgs := make(map[string]any)
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
		updateArgs := make(map[string]any)
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

	// === ВАРИАНТЫ ТОВАРОВ ===

	case "list_variants":
		productID, _ := args["product_id"].(string)
		if productID == "" {
			return nil, fmt.Errorf("product_id is required")
		}
		return h.service.ListVariants(ctx, productID)

	case "get_variant":
		productID, _ := args["product_id"].(string)
		variantID, _ := args["variant_id"].(string)
		if productID == "" || variantID == "" {
			return nil, fmt.Errorf("product_id and variant_id are required")
		}
		return h.service.GetVariant(ctx, productID, variantID)

	case "update_variant":
		productID, _ := args["product_id"].(string)
		variantID, _ := args["variant_id"].(string)
		if productID == "" || variantID == "" {
			return nil, fmt.Errorf("product_id and variant_id are required")
		}
		return h.service.UpdateVariant(ctx, productID, variantID, args)

	case "delete_variant":
		productID, _ := args["product_id"].(string)
		variantID, _ := args["variant_id"].(string)
		if productID == "" || variantID == "" {
			return nil, fmt.Errorf("product_id and variant_id are required")
		}
		return nil, h.service.DeleteVariant(ctx, productID, variantID)

	// === ЗНАЧЕНИЯ СВОЙСТВ ===

	case "list_property_values":
		accessoryID, _ := args["accessory_id"].(string)
		if accessoryID == "" {
			return nil, fmt.Errorf("accessory_id is required")
		}
		return h.service.ListPropertyValues(ctx, accessoryID)

	case "create_property_value":
		accessoryID, _ := args["accessory_id"].(string)
		title, _ := args["title"].(string)
		if accessoryID == "" || title == "" {
			return nil, fmt.Errorf("accessory_id and title are required")
		}
		return h.service.CreatePropertyValue(ctx, accessoryID, title)

	case "update_property_value":
		propertyValueID, _ := args["property_value_id"].(string)
		if propertyValueID == "" {
			return nil, fmt.Errorf("property_value_id is required")
		}
		return h.service.UpdatePropertyValue(ctx, propertyValueID, args)

	case "delete_property_value":
		propertyValueID, _ := args["property_value_id"].(string)
		if propertyValueID == "" {
			return nil, fmt.Errorf("property_value_id is required")
		}
		return nil, h.service.DeletePropertyValue(ctx, propertyValueID)

	// === МЕНЮ ===

	case "list_menus":
		return h.service.ListMenus(ctx)

	case "get_menu":
		menuID, _ := args["menu_id"].(string)
		if menuID == "" {
			return nil, fmt.Errorf("menu_id is required")
		}
		return h.service.GetMenu(ctx, menuID)

	case "create_menu":
		title, _ := args["title"].(string)
		handle, _ := args["handle"].(string)
		if title == "" {
			return nil, fmt.Errorf("title is required")
		}
		return h.service.CreateMenu(ctx, title, handle)

	case "update_menu":
		menuID, _ := args["menu_id"].(string)
		if menuID == "" {
			return nil, fmt.Errorf("menu_id is required")
		}
		return h.service.UpdateMenu(ctx, menuID, args)

	case "delete_menu":
		menuID, _ := args["menu_id"].(string)
		if menuID == "" {
			return nil, fmt.Errorf("menu_id is required")
		}
		return nil, h.service.DeleteMenu(ctx, menuID)

	// === РЕДИРЕКТЫ ===

	case "list_redirects":
		return h.service.ListRedirects(ctx)

	case "create_redirect":
		url, _ := args["url"].(string)
		target, _ := args["target"].(string)
		if url == "" || target == "" {
			return nil, fmt.Errorf("url and target are required")
		}
		return h.service.CreateRedirect(ctx, url, target)

	case "delete_redirect":
		redirectID, _ := args["redirect_id"].(string)
		if redirectID == "" {
			return nil, fmt.Errorf("redirect_id is required")
		}
		return nil, h.service.DeleteRedirect(ctx, redirectID)

	// === БЛОГИ ===

	case "list_blogs":
		return h.service.ListBlogs(ctx)

	case "get_blog":
		blogID, _ := args["blog_id"].(string)
		if blogID == "" {
			return nil, fmt.Errorf("blog_id is required")
		}
		return h.service.GetBlog(ctx, blogID)

	case "create_blog":
		title, _ := args["title"].(string)
		handle, _ := args["handle"].(string)
		if title == "" {
			return nil, fmt.Errorf("title is required")
		}
		return h.service.CreateBlog(ctx, title, handle)

	case "delete_blog":
		blogID, _ := args["blog_id"].(string)
		if blogID == "" {
			return nil, fmt.Errorf("blog_id is required")
		}
		return nil, h.service.DeleteBlog(ctx, blogID)

	// === ЗАПИСИ БЛОГА ===

	case "list_blog_posts":
		blogID, _ := args["blog_id"].(string)
		if blogID == "" {
			return nil, fmt.Errorf("blog_id is required")
		}
		return h.service.ListBlogPosts(ctx, blogID)

	case "get_blog_post":
		blogID, _ := args["blog_id"].(string)
		postID, _ := args["post_id"].(string)
		if blogID == "" || postID == "" {
			return nil, fmt.Errorf("blog_id and post_id are required")
		}
		return h.service.GetBlogPost(ctx, blogID, postID)

	case "create_blog_post":
		blogID, _ := args["blog_id"].(string)
		title, _ := args["title"].(string)
		bodyHTML, _ := args["body_html"].(string)
		if blogID == "" || title == "" || bodyHTML == "" {
			return nil, fmt.Errorf("blog_id, title, and body_html are required")
		}
		return h.service.CreateBlogPost(ctx, blogID, title, bodyHTML)

	case "update_blog_post":
		blogID, _ := args["blog_id"].(string)
		postID, _ := args["post_id"].(string)
		if blogID == "" || postID == "" {
			return nil, fmt.Errorf("blog_id and post_id are required")
		}
		return h.service.UpdateBlogPost(ctx, blogID, postID, args)

	case "delete_blog_post":
		blogID, _ := args["blog_id"].(string)
		postID, _ := args["post_id"].(string)
		if blogID == "" || postID == "" {
			return nil, fmt.Errorf("blog_id and post_id are required")
		}
		return nil, h.service.DeleteBlogPost(ctx, blogID, postID)

	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}
