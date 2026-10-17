package insales

import (
	"context"
)

type Client interface {
	GetStoreInfo(ctx context.Context) (*StoreInfo, error)
	ListProducts(ctx context.Context, page, perPage int, filters map[string]any) ([]Product, error)
	GetProduct(ctx context.Context, productID string) (*Product, error)
	CreateProduct(ctx context.Context, title, price string, args map[string]any) (*Product, error)
	UpdateProduct(ctx context.Context, productID string, args map[string]any) (*Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListOrders(ctx context.Context, page, perPage int, filters map[string]any) ([]Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	UpdateOrder(ctx context.Context, orderID string, args map[string]any) (*Order, error)
	CreateOrder(ctx context.Context, clientID int, args map[string]any) (*Order, error)
	ListClients(ctx context.Context, page, perPage int) ([]any, error)
	GetClient(ctx context.Context, clientID string) (any, error)
	CreateClient(ctx context.Context, args map[string]any) (any, error)
	UpdateClient(ctx context.Context, clientID string, args map[string]any) (any, error)
	DeleteClient(ctx context.Context, clientID string) error
	ListCategories(ctx context.Context, page, perPage int) ([]any, error)
	GetCategory(ctx context.Context, categoryID string) (any, error)
	CreateCategory(ctx context.Context, title string, args map[string]any) (any, error)
	UpdateCategory(ctx context.Context, categoryID string, args map[string]any) (any, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	ListCollections(ctx context.Context, page, perPage int) ([]any, error)
	GetCollection(ctx context.Context, collectionID string) (any, error)
	CreateCollection(ctx context.Context, title string) (any, error)
	UpdateCollection(ctx context.Context, collectionID, title string) (any, error)
	DeleteCollection(ctx context.Context, collectionID string) error
	ListImages(ctx context.Context, productID string) ([]any, error)
	AddImage(ctx context.Context, productID, url string, args map[string]any) (any, error)
	DeleteImage(ctx context.Context, productID, imageID string) error
	ListMetaFields(ctx context.Context, args map[string]any) ([]any, error)
	CreateMetaField(ctx context.Context, args map[string]any) (any, error)
	UpdateMetaField(ctx context.Context, metaFieldID string, args map[string]any) (any, error)
	DeleteMetaField(ctx context.Context, metaFieldID string) error
	ListWebhooks(ctx context.Context) ([]any, error)
	CreateWebhook(ctx context.Context, event, url string, args map[string]any) (any, error)
	UpdateWebhook(ctx context.Context, webhookID string, args map[string]any) (any, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
	ListCurrencies(ctx context.Context) (any, error)
	ListDeliveryVariants(ctx context.Context) (any, error)
	ListPaymentGateways(ctx context.Context) (any, error)
	ListPages(ctx context.Context, page, perPage int) ([]Page, error)
	GetPage(ctx context.Context, pageID string) (*Page, error)
	CreatePage(ctx context.Context, title, content, handle string) (*Page, error)
	UpdatePage(ctx context.Context, pageID string, args map[string]any) (*Page, error)
	DeletePage(ctx context.Context, pageID string) error
	ListAccessories(ctx context.Context, page, perPage int) ([]Accessory, error)
	GetAccessory(ctx context.Context, accessoryID string) (*Accessory, error)
	CreateAccessory(ctx context.Context, title string, args map[string]any) (*Accessory, error)
	UpdateAccessory(ctx context.Context, accessoryID string, args map[string]any) (*Accessory, error)
	DeleteAccessory(ctx context.Context, accessoryID string) error

	// Variants
	ListVariants(ctx context.Context, productID string) ([]Variant, error)
	GetVariant(ctx context.Context, productID, variantID string) (*Variant, error)
	UpdateVariant(ctx context.Context, productID, variantID string, args map[string]any) (*Variant, error)
	DeleteVariant(ctx context.Context, productID, variantID string) error

	// Property Values
	ListPropertyValues(ctx context.Context, accessoryID string) ([]PropertyValue, error)
	CreatePropertyValue(ctx context.Context, accessoryID, title string) (*PropertyValue, error)
	UpdatePropertyValue(ctx context.Context, accessoryID, propertyValueID string, args map[string]any) (*PropertyValue, error)
	DeletePropertyValue(ctx context.Context, accessoryID, propertyValueID string) error

	// Menus
	ListMenus(ctx context.Context) ([]Menu, error)
	GetMenu(ctx context.Context, menuID string) (*Menu, error)
	CreateMenu(ctx context.Context, title, handle string) (*Menu, error)
	UpdateMenu(ctx context.Context, menuID string, args map[string]any) (*Menu, error)
	DeleteMenu(ctx context.Context, menuID string) error

	// Redirects
	ListRedirects(ctx context.Context) ([]Redirect, error)
	CreateRedirect(ctx context.Context, url, target string) (*Redirect, error)
	DeleteRedirect(ctx context.Context, redirectID string) error

	// Blogs
	ListBlogs(ctx context.Context) ([]Blog, error)
	GetBlog(ctx context.Context, blogID string) (*Blog, error)
	CreateBlog(ctx context.Context, title, handle string) (*Blog, error)
	DeleteBlog(ctx context.Context, blogID string) error
	ListBlogPosts(ctx context.Context, blogID string) ([]BlogPost, error)
	GetBlogPost(ctx context.Context, blogID, postID string) (*BlogPost, error)
	CreateBlogPost(ctx context.Context, blogID, title, bodyHTML string) (*BlogPost, error)
	UpdateBlogPost(ctx context.Context, blogID, postID string, args map[string]any) (*BlogPost, error)
	DeleteBlogPost(ctx context.Context, blogID, postID string) error

	// DiscountCodes
	ListDiscountCodes(ctx context.Context, page, perPage int) ([]any, error)
	GetDiscountCode(ctx context.Context, codeID string) (any, error)
	CreateDiscountCode(ctx context.Context, args map[string]any) (any, error)
	UpdateDiscountCode(ctx context.Context, codeID string, args map[string]any) (any, error)
	DeleteDiscountCode(ctx context.Context, codeID string) error

	// Collects
	ListCollects(ctx context.Context, page, perPage int) ([]any, error)
	GetCollectsByProduct(ctx context.Context, productID string) ([]any, error)
	CreateCollect(ctx context.Context, collectionID, productID int) (any, error)
	UpdateCollect(ctx context.Context, collectID string, args map[string]any) (any, error)
	DeleteCollect(ctx context.Context, collectID string) error

	// Properties
	ListProperties(ctx context.Context) ([]any, error)
	GetProperty(ctx context.Context, propertyID string) (any, error)
	CreateProperty(ctx context.Context, args map[string]any) (any, error)
	UpdateProperty(ctx context.Context, propertyID string, args map[string]any) (any, error)
	DeleteProperty(ctx context.Context, propertyID string) error

	// Characteristics
	ListCharacteristics(ctx context.Context, propertyID string) ([]any, error)
	GetCharacteristic(ctx context.Context, propertyID, characteristicID string) (any, error)
	CreateCharacteristic(ctx context.Context, propertyID string, args map[string]any) (any, error)
	UpdateCharacteristic(ctx context.Context, propertyID, characteristicID string, args map[string]any) (any, error)
	DeleteCharacteristic(ctx context.Context, propertyID, characteristicID string) error

	// Reviews
	ListReviews(ctx context.Context, page, perPage int) ([]any, error)
	GetReview(ctx context.Context, reviewID string) (any, error)
	CreateReview(ctx context.Context, args map[string]any) (any, error)
	UpdateReview(ctx context.Context, reviewID string, args map[string]any) (any, error)
	DeleteReview(ctx context.Context, reviewID string) error

	// Files
	ListFiles(ctx context.Context) ([]any, error)
	GetFile(ctx context.Context, fileID string) (any, error)
	CreateFile(ctx context.Context, args map[string]any) (any, error)
	DeleteFile(ctx context.Context, fileID string) error

	// Domains
	ListDomains(ctx context.Context) ([]any, error)
	GetDomain(ctx context.Context, domainID string) (any, error)
	CreateDomain(ctx context.Context, args map[string]any) (any, error)
	UpdateDomain(ctx context.Context, domainID string, args map[string]any) (any, error)
	DeleteDomain(ctx context.Context, domainID string) error

	// ClientGroups
	ListClientGroups(ctx context.Context) ([]any, error)
	GetClientGroup(ctx context.Context, groupID string) (any, error)
	CreateClientGroup(ctx context.Context, args map[string]any) (any, error)
	UpdateClientGroup(ctx context.Context, groupID string, args map[string]any) (any, error)
	DeleteClientGroup(ctx context.Context, groupID string) error

	// CustomStatuses
	ListCustomStatuses(ctx context.Context) ([]any, error)
	GetCustomStatus(ctx context.Context, statusID string) (any, error)
	CreateCustomStatus(ctx context.Context, args map[string]any) (any, error)
	UpdateCustomStatus(ctx context.Context, statusID string, args map[string]any) (any, error)
	DeleteCustomStatus(ctx context.Context, statusID string) error

	// ProductFields
	ListProductFields(ctx context.Context) ([]any, error)
	GetProductField(ctx context.Context, fieldID string) (any, error)
	CreateProductField(ctx context.Context, args map[string]any) (any, error)
	UpdateProductField(ctx context.Context, fieldID string, args map[string]any) (any, error)
	DeleteProductField(ctx context.Context, fieldID string) error

	// VariantFields
	ListVariantFields(ctx context.Context) ([]any, error)
	GetVariantField(ctx context.Context, fieldID string) (any, error)
	CreateVariantField(ctx context.Context, args map[string]any) (any, error)
	UpdateVariantField(ctx context.Context, fieldID string, args map[string]any) (any, error)
	DeleteVariantField(ctx context.Context, fieldID string) error

	// CollectionFields
	ListCollectionFields(ctx context.Context) ([]any, error)
	GetCollectionField(ctx context.Context, fieldID string) (any, error)
	CreateCollectionField(ctx context.Context, args map[string]any) (any, error)
	UpdateCollectionField(ctx context.Context, fieldID string, args map[string]any) (any, error)
	DeleteCollectionField(ctx context.Context, fieldID string) error

	// CollectionFilters
	ListCollectionFilters(ctx context.Context) ([]any, error)
	GetCollectionFilter(ctx context.Context, filterID string) (any, error)
	CreateCollectionFilter(ctx context.Context, args map[string]any) (any, error)
	UpdateCollectionFilter(ctx context.Context, filterID string, args map[string]any) (any, error)
	DeleteCollectionFilter(ctx context.Context, filterID string) error

	// JsTags
	ListJsTags(ctx context.Context) ([]any, error)
	GetJsTag(ctx context.Context, tagID string) (any, error)
	CreateJsTag(ctx context.Context, args map[string]any) (any, error)
	UpdateJsTag(ctx context.Context, tagID string, args map[string]any) (any, error)
	DeleteJsTag(ctx context.Context, tagID string) error

	// ExternalDiscounts
	ListExternalDiscounts(ctx context.Context) ([]any, error)
	GetExternalDiscount(ctx context.Context, discountID string) (any, error)
	CreateExternalDiscount(ctx context.Context, args map[string]any) (any, error)
	UpdateExternalDiscount(ctx context.Context, discountID string, args map[string]any) (any, error)
	DeleteExternalDiscount(ctx context.Context, discountID string) error

	// DiscountCollectionsGroups
	ListDiscountCollectionsGroups(ctx context.Context) ([]any, error)
	GetDiscountCollectionsGroup(ctx context.Context, groupID string) (any, error)
	CreateDiscountCollectionsGroup(ctx context.Context, args map[string]any) (any, error)
	UpdateDiscountCollectionsGroup(ctx context.Context, groupID string, args map[string]any) (any, error)
	DeleteDiscountCollectionsGroup(ctx context.Context, groupID string) error

	// Warehouses
	ListWarehouses(ctx context.Context) ([]any, error)
	GetWarehouse(ctx context.Context, warehouseID string) (any, error)
	CreateWarehouse(ctx context.Context, args map[string]any) (any, error)
	UpdateWarehouse(ctx context.Context, warehouseID string, args map[string]any) (any, error)
	DeleteWarehouse(ctx context.Context, warehouseID string) error

	// SalesChannels
	ListSalesChannels(ctx context.Context) ([]any, error)
	GetSalesChannel(ctx context.Context, channelID string) (any, error)
	CreateSalesChannel(ctx context.Context, args map[string]any) (any, error)
	UpdateSalesChannel(ctx context.Context, channelID string, args map[string]any) (any, error)
	DeleteSalesChannel(ctx context.Context, channelID string) error

	// StockCurrencies
	ListStockCurrencies(ctx context.Context) ([]any, error)
	GetStockCurrency(ctx context.Context, currencyID string) (any, error)
	CreateStockCurrency(ctx context.Context, args map[string]any) (any, error)
	UpdateStockCurrency(ctx context.Context, currencyID string, args map[string]any) (any, error)
	DeleteStockCurrency(ctx context.Context, currencyID string) error

	// PriceKinds
	ListPriceKinds(ctx context.Context) ([]any, error)
	GetPriceKind(ctx context.Context, kindID string) (any, error)
	CreatePriceKind(ctx context.Context, args map[string]any) (any, error)
	UpdatePriceKind(ctx context.Context, kindID string, args map[string]any) (any, error)
	DeletePriceKind(ctx context.Context, kindID string) error

	// OptionNames
	ListOptionNames(ctx context.Context) ([]any, error)
	GetOptionName(ctx context.Context, optionNameID string) (any, error)
	CreateOptionName(ctx context.Context, args map[string]any) (any, error)
	UpdateOptionName(ctx context.Context, optionNameID string, args map[string]any) (any, error)
	DeleteOptionName(ctx context.Context, optionNameID string) error

	// OptionValues
	ListOptionValues(ctx context.Context, optionNameID string) ([]any, error)
	GetOptionValue(ctx context.Context, optionNameID, optionValueID string) (any, error)
	CreateOptionValue(ctx context.Context, optionNameID string, args map[string]any) (any, error)
	UpdateOptionValue(ctx context.Context, optionNameID, optionValueID string, args map[string]any) (any, error)
	DeleteOptionValue(ctx context.Context, optionNameID, optionValueID string) error

	// ApplicationActions
	ListApplicationActions(ctx context.Context) ([]any, error)
	GetApplicationAction(ctx context.Context, id string) (any, error)
	CreateApplicationAction(ctx context.Context, args map[string]any) (any, error)
	UpdateApplicationAction(ctx context.Context, id string, args map[string]any) (any, error)
	DeleteApplicationAction(ctx context.Context, id string) error

	// ApplicationCharges
	ListApplicationCharges(ctx context.Context) ([]any, error)
	GetApplicationCharge(ctx context.Context, id string) (any, error)
	CreateApplicationCharge(ctx context.Context, args map[string]any) (any, error)
	DeclineApplicationCharge(ctx context.Context, id string) error

	// ApplicationWidgets
	ListApplicationWidgets(ctx context.Context) ([]any, error)
	GetApplicationWidget(ctx context.Context, id string) (any, error)
	CreateApplicationWidget(ctx context.Context, args map[string]any) (any, error)
	UpdateApplicationWidget(ctx context.Context, id string, args map[string]any) (any, error)
	DeleteApplicationWidget(ctx context.Context, id string) error

	// RecurringApplicationCharges
	ListRecurringApplicationCharges(ctx context.Context) ([]any, error)
	GetRecurringApplicationCharge(ctx context.Context, id string) (any, error)
	CreateRecurringApplicationCharge(ctx context.Context, args map[string]any) (any, error)
	DeleteRecurringApplicationCharge(ctx context.Context, id string) error

	// BonusSystemTransactions
	ListBonusSystemTransactions(ctx context.Context) ([]any, error)
	GetBonusSystemTransaction(ctx context.Context, id string) (any, error)
	CreateBonusSystemTransaction(ctx context.Context, clientID int, args map[string]any) (any, error)
	ListClientBonusSystemTransactions(ctx context.Context, clientID string) ([]any, error)

	// OAuthApps
	ListOAuthApps(ctx context.Context) ([]any, error)
	GetOAuthApp(ctx context.Context, id string) (any, error)
	CreateOAuthApp(ctx context.Context, args map[string]any) (any, error)
	UpdateOAuthApp(ctx context.Context, id string, args map[string]any) (any, error)
	DeleteOAuthApp(ctx context.Context, id string) error

	// UserTasks
	ListUserTasks(ctx context.Context) ([]any, error)
	GetUserTask(ctx context.Context, id string) (any, error)
	CreateUserTask(ctx context.Context, args map[string]any) (any, error)
	UpdateUserTask(ctx context.Context, id string, args map[string]any) (any, error)
	DeleteUserTask(ctx context.Context, id string) error

	// Supplementary
	ListSupplementaries(ctx context.Context) ([]any, error)
	GetSupplementary(ctx context.Context, id string) (any, error)
	CreateSupplementary(ctx context.Context, args map[string]any) (any, error)
	DeleteSupplementary(ctx context.Context, id string) error

	// Similars
	ListSimilars(ctx context.Context, productID string) ([]any, error)
	CreateSimilar(ctx context.Context, args map[string]any) (any, error)
	DeleteSimilar(ctx context.Context, id string) error

	// Units
	ListUnits(ctx context.Context) ([]any, error)

	// Outlets
	ListOutlets(ctx context.Context) ([]any, error)
	GetOutlet(ctx context.Context, outletID string) (any, error)
	CreateOutlet(ctx context.Context, args map[string]any) (any, error)
	UpdateOutlet(ctx context.Context, outletID string, args map[string]any) (any, error)
	DeleteOutlet(ctx context.Context, outletID string) error

	// Marketplaces
	ListMarketplaces(ctx context.Context) ([]any, error)
	GetMarketplace(ctx context.Context, marketplaceID string) (any, error)
	CreateMarketplace(ctx context.Context, args map[string]any) (any, error)
	UpdateMarketplace(ctx context.Context, marketplaceID string, args map[string]any) (any, error)
	DeleteMarketplace(ctx context.Context, marketplaceID string) error

	// ProviderToken
	CreateProviderToken(ctx context.Context, messengerType string) (any, error)

	// Account
	UpdateAccount(ctx context.Context, args map[string]any) (any, error)

	// Counts
	GetOrdersCount(ctx context.Context, filters map[string]any) (int, error)
	GetProductsCount(ctx context.Context) (int, error)
}

type StoreInfo struct {
	ID       string `json:"id"`
	Name     string `json:"title"`
	Domain   string `json:"domain"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
}

type Product struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description,omitempty"`
	ShortDescription string    `json:"short_description,omitempty"`
	Permalink        string    `json:"permalink,omitempty"`
	IsHidden         bool      `json:"is_hidden,omitempty"`
	Available        bool      `json:"available,omitempty"`
	SKU              string    `json:"sku"`
	Price            string    `json:"price"`
	Quantity         int       `json:"quantity"`
	Vendor           string    `json:"vendor"`
	CategoryID       *int      `json:"category_id"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
	Variants         []Variant `json:"variants"`
	Images           []any     `json:"images,omitempty"`
	FieldValues      []any     `json:"field_values,omitempty"`
	Properties       []any     `json:"properties,omitempty"`
	Characteristics  []any     `json:"characteristics,omitempty"`
	OptionNames      []any     `json:"option_names,omitempty"`
}

type Variant struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Price     string `json:"price"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	Weight    string `json:"weight"`
	ProductID int    `json:"product_id"`
}

type Order struct {
	ID                string         `json:"id"`
	Number            string         `json:"number"`
	FulfillmentStatus string         `json:"fulfillment_status"`
	FinancialStatus   string         `json:"financial_status"`
	TotalPrice        string         `json:"total_price"`
	ItemsPrice        string         `json:"items_price"`
	DeliveryPrice     string         `json:"delivery_price"`
	DiscountAmount    string         `json:"discount_amount"`
	ClientID          int            `json:"client_id"`
	ClientName        string         `json:"client_name"`
	DeliveryTitle     string         `json:"delivery_title"`
	PaymentTitle      string         `json:"payment_title"`
	Comment           string         `json:"comment"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
	Lines             []OrderLine    `json:"lines"`
	FieldsValues      []any          `json:"fields_values,omitempty"`
	ShippingAddress   map[string]any `json:"shipping_address,omitempty"`
}

type OrderLine struct {
	ID        string `json:"id"`
	ProductID int    `json:"product_id"`
	VariantID int    `json:"variant_id"`
	Title     string `json:"title"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
	Total     string `json:"total"`
	SKU       string `json:"sku"`
}

type Page struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Handle    string `json:"handle"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Accessory struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PropertyValue struct {
	ID          string `json:"id"`
	AccessoryID int    `json:"accessory_id"`
	Title       string `json:"title"`
	Position    int    `json:"position"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Menu struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Handle    string     `json:"handle"`
	Items     []MenuItem `json:"items"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

type MenuItem struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Position int    `json:"position"`
}

type Redirect struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Target    string `json:"target"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Blog struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Handle    string `json:"handle"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BlogPost struct {
	ID          string `json:"id"`
	BlogID      int    `json:"blog_id"`
	Title       string `json:"title"`
	Handle      string `json:"handle"`
	BodyHTML    string `json:"body_html"`
	Author      string `json:"author"`
	PublishedAt string `json:"published_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
