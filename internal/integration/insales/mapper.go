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

func MapDiscountCode(resp *InSalesDiscountCodeResponse) map[string]any {
	return map[string]any{
		"id":                resp.ID,
		"code":              resp.Code,
		"type_id":           resp.TypeID,
		"discount":          resp.Discount,
		"description":       resp.Description,
		"disabled":          resp.Disabled,
		"act_once":          resp.ActOnce,
		"act_once_for_client": resp.ActOnceForClient,
		"min_price":         resp.MinPrice,
		"expired_at":        resp.ExpiredAt,
		"created_at":        resp.CreatedAt,
		"updated_at":        resp.UpdatedAt,
	}
}

func MapCollect(resp *InSalesCollectResponse) map[string]any {
	return map[string]any{
		"id":             resp.ID,
		"collection_id":  resp.CollectionID,
		"product_id":     resp.ProductID,
		"position":       resp.Position,
		"created_at":     resp.CreatedAt,
	}
}

func MapProperty(resp *InSalesPropertyResponse) map[string]any {
	return map[string]any{
		"id":              resp.ID,
		"title":           resp.Title,
		"position":        resp.Position,
		"permalink":       resp.Permalink,
		"backoffice":      resp.Backoffice,
		"is_hidden":       resp.IsHidden,
		"is_navigational": resp.IsNavigational,
	}
}

func MapCharacteristic(resp *InSalesCharacteristicResponse) map[string]any {
	return map[string]any{
		"id":          resp.ID,
		"property_id": resp.PropertyID,
		"title":       resp.Title,
		"position":    resp.Position,
		"permalink":   resp.Permalink,
	}
}

func MapReview(resp *InSalesReviewResponse) map[string]any {
	return map[string]any{
		"id":            resp.ID,
		"product_id":    resp.ProductID,
		"author":        resp.Author,
		"content":       resp.Content,
		"email":         resp.Email,
		"rating":        resp.Rating,
		"published":     resp.Published,
		"manager_name":  resp.ManagerName,
		"manager_reply": resp.ManagerReply,
		"spam":          resp.Spam,
		"created_at":    resp.CreatedAt,
		"updated_at":    resp.UpdatedAt,
	}
}

func MapFile(resp *InSalesFileResponse) map[string]any {
	return map[string]any{
		"id":           resp.ID,
		"absolute_url": resp.AbsoluteURL,
		"created_at":   resp.CreatedAt,
	}
}

func MapDomain(resp *InSalesDomainResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"domain":     resp.Domain,
		"main":       resp.Main,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapClientGroup(resp *InSalesClientGroupResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"discount":   resp.Discount,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCustomStatus(resp *InSalesCustomStatusResponse) map[string]any {
	return map[string]any{
		"id":          resp.ID,
		"title":       resp.Title,
		"status_type": resp.StatusType,
		"position":    resp.Position,
		"created_at":  resp.CreatedAt,
		"updated_at":  resp.UpdatedAt,
	}
}

func MapProductField(resp *InSalesProductFieldResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"handle":     resp.Handle,
		"field_type": resp.FieldType,
		"position":   resp.Position,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapVariantField(resp *InSalesVariantFieldResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"handle":     resp.Handle,
		"field_type": resp.FieldType,
		"position":   resp.Position,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCollectionField(resp *InSalesCollectionFieldResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"handle":     resp.Handle,
		"field_type": resp.FieldType,
		"position":   resp.Position,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapCollectionFilter(resp *InSalesCollectionFilterResponse) map[string]any {
	return map[string]any{
		"id":             resp.ID,
		"title":          resp.Title,
		"filter_type":    resp.FilterType,
		"position":       resp.Position,
		"collection_id":  resp.CollectionID,
		"created_at":     resp.CreatedAt,
		"updated_at":     resp.UpdatedAt,
	}
}

func MapJsTag(resp *InSalesJsTagResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"content":    resp.Content,
		"position":   resp.Position,
		"enabled":    resp.Enabled,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapExternalDiscount(resp *InSalesExternalDiscountResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"percent":    resp.Percent,
		"amount":     resp.Amount,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapDiscountCollectionsGroup(resp *InSalesDiscountCollectionsGroupResponse) map[string]any {
	return map[string]any{
		"id":             resp.ID,
		"title":          resp.Title,
		"account_id":     resp.AccountID,
		"collection_ids": resp.CollectionIDs,
		"created_at":     resp.CreatedAt,
		"updated_at":     resp.UpdatedAt,
	}
}

func MapWarehouse(resp *InSalesWarehouseResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"address":    resp.Address,
		"active":     resp.Active,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapSalesChannel(resp *InSalesSalesChannelResponse) map[string]any {
	return map[string]any{
		"id":          resp.ID,
		"title":       resp.Title,
		"system_name": resp.SystemName,
		"created_at":  resp.CreatedAt,
		"updated_at":  resp.UpdatedAt,
	}
}

func MapStockCurrency(resp *InSalesStockCurrencyResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"code":       resp.Code,
		"rate":       resp.Rate,
		"main":       resp.Main,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapPriceKind(resp *InSalesPriceKindResponse) map[string]any {
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"position":   resp.Position,
		"created_at": resp.CreatedAt,
		"updated_at": resp.UpdatedAt,
	}
}

func MapOptionName(resp *InSalesOptionNameResponse) map[string]any {
	return map[string]any{
		"id":            resp.ID,
		"title":         resp.Title,
		"position":      resp.Position,
		"navigational":  resp.Navigational,
		"permalink":     resp.Permalink,
	}
}

func MapOptionValue(resp *InSalesOptionValueResponse) map[string]any {
	return map[string]any{
		"id":             resp.ID,
		"option_name_id": resp.OptionNameID,
		"title":          resp.Title,
		"position":       resp.Position,
		"permalink":      resp.Permalink,
		"image_url":      resp.ImageURL,
	}
}

func MapApplicationAction(resp *InSalesApplicationActionResponse) map[string]any {
	return map[string]any{"id": resp.ID, "title": resp.Title, "handle": resp.Handle, "entity": resp.Entity, "url": resp.URL}
}

func MapApplicationCharge(resp *InSalesApplicationChargeResponse) map[string]any {
	return map[string]any{"id": resp.ID, "name": resp.Name, "price": resp.Price, "status": resp.Status, "confirmation_url": resp.ConfirmationURL, "test": resp.Test, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapApplicationWidget(resp *InSalesApplicationWidgetResponse) map[string]any {
	return map[string]any{"id": resp.ID, "title": resp.Title, "enabled": resp.Enabled, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapRecurringApplicationCharge(resp *InSalesRecurringApplicationChargeResponse) map[string]any {
	return map[string]any{"id": resp.ID, "name": resp.Name, "price": resp.Price, "status": resp.Status, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapBonusSystemTransaction(resp *InSalesBonusSystemTransactionResponse) map[string]any {
	return map[string]any{"id": resp.ID, "client_id": resp.ClientID, "amount": resp.Amount, "comment": resp.Comment, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapOAuthApp(resp *InSalesOAuthAppResponse) map[string]any {
	return map[string]any{"id": resp.ID, "title": resp.Title, "app_key": resp.AppKey, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapUserTask(resp *InSalesUserTaskResponse) map[string]any {
	return map[string]any{"id": resp.ID, "title": resp.Title, "status": resp.Status, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapSupplementary(resp *InSalesSupplementaryResponse) map[string]any {
	return map[string]any{"id": resp.ID, "title": resp.Title, "price": resp.Price, "created_at": resp.CreatedAt, "updated_at": resp.UpdatedAt}
}

func MapSimilar(resp *InSalesSimilarResponse) map[string]any {
	return map[string]any{"id": resp.ID, "product_id": resp.ProductID, "similar_product_id": resp.SimilarID, "position": resp.Position, "created_at": resp.CreatedAt}
}

func MapUnit(resp *InSalesUnitResponse) map[string]any {
	return map[string]any{"id": resp.ID, "title": resp.Title, "unit": resp.Unit}
}
