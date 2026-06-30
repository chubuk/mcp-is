package insales

import (
	"context"
)

type Client interface {
	GetStoreInfo(ctx context.Context) (*StoreInfo, error)
	ListProducts(ctx context.Context, page, perPage int) ([]Product, error)
	GetProduct(ctx context.Context, productID string) (*Product, error)
	CreateProduct(ctx context.Context, title, price string, args map[string]interface{}) (*Product, error)
	UpdateProduct(ctx context.Context, productID string, args map[string]interface{}) (*Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListOrders(ctx context.Context, page, perPage int) ([]Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	UpdateOrder(ctx context.Context, orderID string, args map[string]interface{}) (*Order, error)
	CreateOrder(ctx context.Context, clientID int, args map[string]interface{}) (*Order, error)
	ListClients(ctx context.Context, page, perPage int) ([]interface{}, error)
	GetClient(ctx context.Context, clientID string) (interface{}, error)
	CreateClient(ctx context.Context, email string, args map[string]interface{}) (interface{}, error)
	UpdateClient(ctx context.Context, clientID string, args map[string]interface{}) (interface{}, error)
	DeleteClient(ctx context.Context, clientID string) error
	ListCategories(ctx context.Context, page, perPage int) ([]interface{}, error)
	GetCategory(ctx context.Context, categoryID string) (interface{}, error)
	CreateCategory(ctx context.Context, title string, args map[string]interface{}) (interface{}, error)
	UpdateCategory(ctx context.Context, categoryID string, args map[string]interface{}) (interface{}, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	ListCollections(ctx context.Context, page, perPage int) ([]interface{}, error)
	GetCollection(ctx context.Context, collectionID string) (interface{}, error)
	CreateCollection(ctx context.Context, title string) (interface{}, error)
	UpdateCollection(ctx context.Context, collectionID, title string) (interface{}, error)
	DeleteCollection(ctx context.Context, collectionID string) error
	ListImages(ctx context.Context, productID string) ([]interface{}, error)
	AddImage(ctx context.Context, productID, url string, args map[string]interface{}) (interface{}, error)
	DeleteImage(ctx context.Context, productID, imageID string) error
	ListMetaFields(ctx context.Context, args map[string]interface{}) ([]interface{}, error)
	CreateMetaField(ctx context.Context, args map[string]interface{}) (interface{}, error)
	UpdateMetaField(ctx context.Context, metaFieldID string, args map[string]interface{}) (interface{}, error)
	DeleteMetaField(ctx context.Context, metaFieldID string) error
	ListWebhooks(ctx context.Context) ([]interface{}, error)
	CreateWebhook(ctx context.Context, event, url string, args map[string]interface{}) (interface{}, error)
	UpdateWebhook(ctx context.Context, webhookID string, args map[string]interface{}) (interface{}, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
	ListCurrencies(ctx context.Context) (interface{}, error)
	ListDeliveryVariants(ctx context.Context) (interface{}, error)
	ListPaymentGateways(ctx context.Context) (interface{}, error)
	ListPages(ctx context.Context, page, perPage int) ([]Page, error)
	GetPage(ctx context.Context, pageID string) (*Page, error)
	CreatePage(ctx context.Context, title, content, handle string) (*Page, error)
	UpdatePage(ctx context.Context, pageID string, args map[string]interface{}) (*Page, error)
	DeletePage(ctx context.Context, pageID string) error
	ListAccessories(ctx context.Context, page, perPage int) ([]Accessory, error)
	GetAccessory(ctx context.Context, accessoryID string) (*Accessory, error)
	CreateAccessory(ctx context.Context, title string) (*Accessory, error)
	UpdateAccessory(ctx context.Context, accessoryID string, args map[string]interface{}) (*Accessory, error)
	DeleteAccessory(ctx context.Context, accessoryID string) error
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
