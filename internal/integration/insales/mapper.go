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

func MapClient(resp *InSalesClientResponse) map[string]interface{} {
	return map[string]interface{}{
		"id":         resp.ID,
		"email":      resp.Email,
		"first_name": resp.FirstName,
		"last_name":  resp.LastName,
		"phone":      resp.Phone,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCategory(resp *InSalesCategoryResponse) map[string]interface{} {
	return map[string]interface{}{
		"id":         resp.ID,
		"title":      resp.Title,
		"parent_id":  resp.ParentID,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCollection(resp *InSalesCollectionResponse) map[string]interface{} {
	return map[string]interface{}{
		"id":         resp.ID,
		"title":      resp.Title,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapMetaField(resp *InSalesMetaFieldResponse) map[string]interface{} {
	return map[string]interface{}{
		"id":            resp.ID,
		"resource_type": resp.ResourceType,
		"resource_id":   resp.ResourceID,
		"key":           resp.Key,
		"value":         resp.Value,
	}
}

func MapWebhook(resp *InSalesWebhookResponse) map[string]interface{} {
	return map[string]interface{}{
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
