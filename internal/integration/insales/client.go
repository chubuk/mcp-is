package insales

import (
	"context"
)

type Client interface {
	GetStoreInfo(ctx context.Context) (*StoreInfo, error)
	ListProducts(ctx context.Context, page, perPage int) ([]Product, error)
	GetProduct(ctx context.Context, productID string) (*Product, error)
	CreateProduct(ctx context.Context, title, price string, args map[string]any) (*Product, error)
	UpdateProduct(ctx context.Context, productID string, args map[string]any) (*Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListOrders(ctx context.Context, page, perPage int) ([]Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	UpdateOrder(ctx context.Context, orderID string, args map[string]any) (*Order, error)
	CreateOrder(ctx context.Context, clientID int, args map[string]any) (*Order, error)
	ListClients(ctx context.Context, page, perPage int) ([]any, error)
	GetClient(ctx context.Context, clientID string) (any, error)
	CreateClient(ctx context.Context, email string, args map[string]any) (any, error)
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
	CreateAccessory(ctx context.Context, title string) (*Accessory, error)
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
	UpdatePropertyValue(ctx context.Context, propertyValueID string, args map[string]any) (*PropertyValue, error)
	DeletePropertyValue(ctx context.Context, propertyValueID string) error

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
}

type StoreInfo struct {
	ID       string `json:"id"`
	Name     string `json:"title"`
	Domain   string `json:"domain"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
}

type Product struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	SKU        string    `json:"sku"`
	Price      string    `json:"price"`
	Quantity   int       `json:"quantity"`
	Vendor     string    `json:"vendor"`
	CategoryID *int      `json:"category_id"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	Variants   []Variant `json:"variants"`
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
	ID                string      `json:"id"`
	Number            string      `json:"number"`
	FulfillmentStatus string      `json:"fulfillment_status"`
	FinancialStatus   string      `json:"financial_status"`
	TotalPrice        string      `json:"total_price"`
	ItemsPrice        string      `json:"items_price"`
	DeliveryPrice     string      `json:"delivery_price"`
	DiscountAmount    string      `json:"discount_amount"`
	ClientID          int         `json:"client_id"`
	ClientName        string      `json:"client_name"`
	DeliveryTitle     string      `json:"delivery_title"`
	PaymentTitle      string      `json:"payment_title"`
	Comment           string      `json:"comment"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
	Lines             []OrderLine `json:"lines"`
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
	URL     string `json:"url"`
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
