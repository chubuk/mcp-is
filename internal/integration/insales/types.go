package insales

type InSalesShopResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Domain   string `json:"domain"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
}

type InSalesProductResponse struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Price       string                 `json:"price"`
	CategoryID  *int                   `json:"category_id"`
	Vendor      string                 `json:"vendor"`
	SKU         string                 `json:"sku"`
	Quantity    int                    `json:"quantity"`
	Weight      string                 `json:"weight"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Variants    []InSalesVariant       `json:"variants"`
	Images      []InSalesImage         `json:"images"`
}

type InSalesVariant struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Price     string `json:"price"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	Weight    string `json:"weight"`
	ProductID int    `json:"product_id"`
}

type InSalesImage struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	Position  int    `json:"position"`
	ProductID int    `json:"product_id"`
}

type InSalesOrderResponse struct {
	ID                int                    `json:"id"`
	Number            int                    `json:"number"`
	FulfillmentStatus string                 `json:"fulfillment_status"`
	FinancialStatus   string                 `json:"financial_status"`
	TotalPrice        float64                `json:"total_price"`
	ItemsPrice        float64                `json:"items_price"`
	DeliveryPrice     float64                `json:"delivery_price"`
	Discount          map[string]interface{} `json:"discount"`
	Discounts         []interface{}          `json:"discounts"`
	Client            *InSalesClientRef      `json:"client"`
	DeliveryTitle     string                 `json:"delivery_title"`
	PaymentTitle      string                 `json:"payment_title"`
	Comment           string                 `json:"comment"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
	OrderLines        []InSalesOrderLine     `json:"order_lines"`
}

type InSalesClientRef struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type InSalesOrderLine struct {
	ID           int     `json:"id"`
	ProductID    int     `json:"product_id"`
	VariantID    int     `json:"variant_id"`
	Title        string  `json:"title"`
	Quantity     int     `json:"quantity"`
	SalePrice    float64 `json:"sale_price"`
	TotalPrice   float64 `json:"total_price"`
	DiscountsAmount float64 `json:"discounts_amount"`
	SKU          string  `json:"sku"`
}

type InSalesClientResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesCategoryResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	ParentID  *int   `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesCollectionResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesMetaFieldResponse struct {
	ID           int    `json:"id"`
	ResourceType string `json:"resource_type"`
	ResourceID   int    `json:"resource_id"`
	Key          string `json:"key"`
	Value        string `json:"value"`
}

type InSalesWebhookResponse struct {
	ID     int    `json:"id"`
	Event  string `json:"event"`
	URL    string `json:"url"`
	Format string `json:"format"`
}
