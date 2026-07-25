package insales

import "fmt"

func MapShop(resp *InSalesShopResponse) *StoreInfo {
	return &StoreInfo{
		ID:       fmt.Sprintf("%d", resp.ID),
		Name:     resp.Title,
		Domain:   resp.Domain,
		Email:    resp.Email,
		Currency: resp.Currency,
	}
}

func MapProduct(resp *InSalesProductResponse) *Product {
	return &Product{
		ID:        fmt.Sprintf("%d", resp.ID),
		Title:     resp.Title,
		SKU:       resp.SKU,
		Price:     resp.Price,
		Quantity:  resp.Quantity,
		Vendor:    resp.Vendor,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
		Variants:  mapVariants(resp.Variants),
	}
}

func mapVariants(v []InSalesVariant) []Variant {
	result := make([]Variant, len(v))
	for i, rv := range v {
		result[i] = Variant{
			ID:        fmt.Sprintf("%d", rv.ID),
			Title:     rv.Title,
			Price:     rv.Price,
			SKU:       rv.SKU,
			Quantity:  rv.Quantity,
			Weight:    rv.Weight,
			ProductID: rv.ProductID,
		}
	}
	return result
}

func MapOrder(resp *InSalesOrderResponse) *Order {
	var clientID int
	var clientName string
	if resp.Client != nil {
		clientID = resp.Client.ID
		clientName = resp.Client.Name
	}

	discountAmount := 0.0
	if d, ok := resp.Discount["amount"].(float64); ok {
		discountAmount = d
	}

	return &Order{
		ID:                fmt.Sprintf("%d", resp.ID),
		Number:            fmt.Sprintf("%d", resp.Number),
		FulfillmentStatus: resp.FulfillmentStatus,
		FinancialStatus:   resp.FinancialStatus,
		TotalPrice:        fmt.Sprintf("%.2f", resp.TotalPrice),
		ItemsPrice:        fmt.Sprintf("%.2f", resp.ItemsPrice),
		DeliveryPrice:     fmt.Sprintf("%.2f", resp.DeliveryPrice),
		DiscountAmount:    fmt.Sprintf("%.2f", discountAmount),
		ClientID:          clientID,
		ClientName:        clientName,
		DeliveryTitle:     resp.DeliveryTitle,
		PaymentTitle:      resp.PaymentTitle,
		Comment:           resp.Comment,
		CreatedAt:         resp.CreatedAt,
		UpdatedAt:         resp.UpdatedAt,
		Lines:             mapOrderLines(resp.OrderLines),
	}
}

func mapOrderLines(l []InSalesOrderLine) []OrderLine {
	result := make([]OrderLine, len(l))
	for i, rl := range l {
		result[i] = OrderLine{
			ID:        fmt.Sprintf("%d", rl.ID),
			ProductID: rl.ProductID,
			VariantID: rl.VariantID,
			Title:     rl.Title,
			Quantity:  rl.Quantity,
			Price:     fmt.Sprintf("%.2f", rl.SalePrice),
			Total:     fmt.Sprintf("%.2f", rl.TotalPrice),
			SKU:       rl.SKU,
		}
	}
	return result
}

func MapClient(resp *InSalesClientResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"email":      resp.Email,
		"first_name": resp.FirstName,
		"last_name":  resp.LastName,
		"phone":      resp.Phone,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCategory(resp *InSalesCategoryResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"parent_id":  resp.ParentID,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCollection(resp *InSalesCollectionResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapMetaField(resp *InSalesMetaFieldResponse) map[string]any {
	return map[string]any{
		"id":            resp.ID,
		"resource_type": resp.ResourceType,
		"resource_id":   resp.ResourceID,
		"key":           resp.Key,
		"value":         resp.Value,
	}
}

func MapWebhook(resp *InSalesWebhookResponse) map[string]any {
	return map[string]any{
		"id":     resp.ID,
		"event":  resp.Event,
		"url":    resp.URL,
		"format": resp.Format,
	}
}

func MapPage(resp *InSalesPageResponse) *Page {
	return &Page{
		ID:        fmt.Sprintf("%d", resp.ID),
		Title:     resp.Title,
		Content:   resp.Content,
		Handle:    resp.Handle,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func MapAccessory(resp *InSalesAccessoryResponse) *Accessory {
	return &Accessory{
		ID:        fmt.Sprintf("%d", resp.ID),
		Title:     resp.Title,
		Position:  resp.Position,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func MapPropertyValue(resp *InSalesPropertyValueResponse) *PropertyValue {
	return &PropertyValue{
		ID:          fmt.Sprintf("%d", resp.ID),
		AccessoryID: resp.AccessoryID,
		Title:       resp.Title,
		Position:    resp.Position,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}
}

func MapMenu(resp *InSalesMenuResponse) *Menu {
	items := make([]MenuItem, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = MenuItem{
			ID:       item.ID,
			Title:    item.Title,
			URL:     item.URL,
			Position: item.Position,
		}
	}
	return &Menu{
		ID:        fmt.Sprintf("%d", resp.ID),
		Title:     resp.Title,
		Handle:    resp.Handle,
		Items:     items,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func MapRedirect(resp *InSalesRedirectResponse) *Redirect {
	return &Redirect{
		ID:        fmt.Sprintf("%d", resp.ID),
		URL:       resp.Url,
		Target:    resp.Target,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func MapBlog(resp *InSalesBlogResponse) *Blog {
	return &Blog{
		ID:        fmt.Sprintf("%d", resp.ID),
		Title:     resp.Title,
		Handle:    resp.Handle,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func MapBlogPost(resp *InSalesBlogPostResponse) *BlogPost {
	// InSales API uses "content" but we expose as "body_html"
	bodyHTML := resp.Content
	if bodyHTML == "" {
		bodyHTML = resp.BodyHTML
	}
	// InSales API uses "permalink" but we expose as "handle"
	handle := resp.Permalink
	if handle == "" {
		handle = resp.Handle
	}
	return &BlogPost{
		ID:          fmt.Sprintf("%d", resp.ID),
		BlogID:      resp.BlogID,
		Title:       resp.Title,
		Handle:      handle,
		BodyHTML:    bodyHTML,
		Author:      resp.Author,
		PublishedAt: resp.PublishedAt,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}
}
