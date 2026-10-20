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
	available := false
	if resp.Available != nil {
		available = *resp.Available
	}
	return &Product{
		ID:               fmt.Sprintf("%d", resp.ID),
		Title:            resp.Title,
		Description:      resp.Description,
		ShortDescription: resp.ShortDescription,
		Permalink:        resp.Permalink,
		IsHidden:         resp.IsHidden,
		Available:        available,
		SKU:              resp.SKU,
		Price:            resp.Price,
		Quantity:         resp.Quantity,
		Vendor:           resp.Vendor,
		CategoryID:       resp.CategoryID,
		CreatedAt:        resp.CreatedAt,
		UpdatedAt:        resp.UpdatedAt,
		Variants:         mapVariants(resp.Variants),
		Images:           mapImages(resp.Images),
		FieldValues:      mapFieldValues(resp.ProductFieldValues),
		Properties:       mapProductProperties(resp.Properties),
		Characteristics:  mapCharacteristics(resp.Characteristics),
		OptionNames:      mapOptionNames(resp.OptionNames),
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

func mapImages(images []InSalesImage) []any {
	result := make([]any, len(images))
	for i, img := range images {
		result[i] = map[string]any{
			"id":         img.ID,
			"url":        img.URL,
			"filename":   img.Filename,
			"position":   img.Position,
			"product_id": img.ProductID,
		}
	}
	return result
}

func mapFieldValues(fvs []InSalesFieldValue) []any {
	result := make([]any, len(fvs))
	for i, fv := range fvs {
		item := map[string]any{
			"id":               fmt.Sprintf("%d", fv.ID),
			"product_field_id": fmt.Sprintf("%d", fv.ProductFieldID),
			"value":            fv.Value,
		}
		if fv.Handle != "" {
			item["handle"] = fv.Handle
		}
		if fv.Title != "" {
			item["title"] = fv.Title
		} else if fv.Name != "" {
			item["title"] = fv.Name
		}
		result[i] = item
	}
	return result
}

// EnrichFieldValues joins handle/title from a product-field catalog into field_values.
// catalog items: {id, handle, title} (as returned by /product_fields.json).
func EnrichFieldValues(values []any, catalog []any) []any {
	if len(values) == 0 || len(catalog) == 0 {
		return values
	}
	byID := map[string]map[string]any{}
	for _, raw := range catalog {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id := fmt.Sprintf("%v", m["id"])
		byID[id] = m
	}
	for _, raw := range values {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id := fmt.Sprintf("%v", m["product_field_id"])
		if m["product_field_id"] == nil {
			id = fmt.Sprintf("%v", m["field_id"])
		}
		if cat, ok := byID[id]; ok {
			if _, has := m["handle"]; !has {
				if h, ok := cat["handle"].(string); ok && h != "" {
					m["handle"] = h
				}
			}
			if _, has := m["title"]; !has {
				if t, ok := cat["title"].(string); ok && t != "" {
					m["title"] = t
				}
			}
		}
	}
	return values
}

func mapProductProperties(props []InSalesProductProperty) []any {
	result := make([]any, len(props))
	for i, p := range props {
		result[i] = map[string]any{
			"id":       p.ID,
			"title":    p.Title,
			"position": p.Position,
		}
	}
	return result
}

func mapCharacteristics(chars []InSalesCharacteristic) []any {
	result := make([]any, len(chars))
	for i, c := range chars {
		result[i] = map[string]any{
			"id":          c.ID,
			"property_id": c.PropertyID,
			"title":       c.Title,
			"position":    c.Position,
		}
	}
	return result
}

func mapOptionNames(opts []InSalesOptionName) []any {
	result := make([]any, len(opts))
	for i, o := range opts {
		result[i] = map[string]any{
			"id":       o.ID,
			"title":    o.Title,
			"position": o.Position,
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

	var shippingAddress map[string]any
	if resp.ShippingAddress != nil {
		shippingAddress = map[string]any{
			"id":                    resp.ShippingAddress.ID,
			"name":                  resp.ShippingAddress.Name,
			"surname":               resp.ShippingAddress.Surname,
			"middlename":            resp.ShippingAddress.Middlename,
			"phone":                 resp.ShippingAddress.Phone,
			"formatted_phone":       resp.ShippingAddress.FormattedPhone,
			"full_name":             resp.ShippingAddress.FullName,
			"full_delivery_address": resp.ShippingAddress.FullDeliveryAddress,
		}
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
		FieldsValues:      mapFieldValues(resp.FieldsValues),
		ShippingAddress:   shippingAddress,
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
		"id":                       resp.ID,
		"email":                    resp.Email,
		"name":                     resp.Name,
		"surname":                  resp.Surname,
		"middlename":               resp.Middlename,
		"first_name":               resp.FirstName,
		"last_name":                resp.LastName,
		"phone":                    resp.Phone,
		"full_name":                resp.FullName,
		"comment":                  resp.Comment,
		"registered":               resp.Registered,
		"subscribe":                resp.Subscribe,
		"client_group_id":          resp.ClientGroupID,
		"bonus_points":             resp.BonusPoints,
		"type":                     resp.Type,
		"progressive_discount":     resp.ProgressiveDiscount,
		"group_discount":           resp.GroupDiscount,
		"consent_to_personal_data": resp.ConsentToPersonalData,
		"created_at":               resp.CreatedAt,
		"updated_at":               resp.UpdatedAt,
		"fields_values":            mapFieldValues(resp.FieldsValues),
		"client_tags":              resp.ClientTags,
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
			URL:      item.URL,
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
		"id":                  resp.ID,
		"code":                resp.Code,
		"type_id":             resp.TypeID,
		"discount":            resp.Discount,
		"description":         resp.Description,
		"disabled":            resp.Disabled,
		"act_once":            resp.ActOnce,
		"act_once_for_client": resp.ActOnceForClient,
		"min_price":           resp.MinPrice,
		"expired_at":          resp.ExpiredAt,
		"created_at":          resp.CreatedAt,
		"updated_at":          resp.UpdatedAt,
	}
}

func MapCollect(resp *InSalesCollectResponse) map[string]any {
	return map[string]any{
		"id":            resp.ID,
		"collection_id": resp.CollectionID,
		"product_id":    resp.ProductID,
		"position":      resp.Position,
		"created_at":    resp.CreatedAt,
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
	fieldType := resp.Type
	if fieldType == "" {
		fieldType = resp.FieldType
	}
	return map[string]any{
		"id":         resp.ID,
		"title":      resp.Title,
		"handle":     resp.Handle,
		"type":       fieldType,
		"field_type": fieldType,
		"position":   resp.Position,
		"is_hidden":  resp.IsHidden,
		"is_indexed": resp.IsIndexed,
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
		"id":            resp.ID,
		"title":         resp.Title,
		"filter_type":   resp.FilterType,
		"position":      resp.Position,
		"collection_id": resp.CollectionID,
		"created_at":    resp.CreatedAt,
		"updated_at":    resp.UpdatedAt,
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
		"id":           resp.ID,
		"title":        resp.Title,
		"position":     resp.Position,
		"navigational": resp.Navigational,
		"permalink":    resp.Permalink,
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
