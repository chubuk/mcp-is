package insales

type InSalesShopResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Domain   string `json:"domain"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
}

type InSalesProductResponse struct {
	ID                    int                      `json:"id"`
	Title                 string                   `json:"title"`
	Description           string                   `json:"description"`
	ShortDescription      string                   `json:"short_description"`
	Price                 string                   `json:"price"`
	CategoryID            *int                     `json:"category_id"`
	Vendor                string                   `json:"vendor"`
	SKU                   string                   `json:"sku"`
	Quantity              int                      `json:"quantity"`
	Weight                string                   `json:"weight"`
	Permalink             string                   `json:"permalink"`
	IsHidden              bool                     `json:"is_hidden"`
	Available             *bool                    `json:"available"`
	Archived              bool                     `json:"archived"`
	CreatedAt             string                   `json:"created_at"`
	UpdatedAt             string                   `json:"updated_at"`
	Variants              []InSalesVariant         `json:"variants"`
	Images                []InSalesImage           `json:"images"`
	ProductFieldValues    []InSalesFieldValue      `json:"product_field_values"`
	Properties            []InSalesProductProperty `json:"properties"`
	Characteristics       []InSalesCharacteristic  `json:"characteristics"`
	OptionNames           []InSalesOptionName      `json:"option_names"`
	VideoLinks            []any                    `json:"video_links"`
	ProductAccessoryLinks []any                    `json:"product_accessory_links"`
	CollectionsIds        []int                    `json:"collections_ids"`
}

type InSalesFieldValue struct {
	ID             int    `json:"id"`
	ProductFieldID int    `json:"product_field_id"`
	Value          string `json:"value"`
	Handle         string `json:"handle,omitempty"`
	Title          string `json:"title,omitempty"`
	Name           string `json:"name,omitempty"`
}

type InSalesProductProperty struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Position int    `json:"position"`
}

type InSalesCharacteristic struct {
	ID         int    `json:"id"`
	PropertyID int    `json:"property_id"`
	Title      string `json:"title"`
	Position   int    `json:"position"`
}

type InSalesOptionName struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Position int    `json:"position"`
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
	ID                int                     `json:"id"`
	Number            int                     `json:"number"`
	FulfillmentStatus string                  `json:"fulfillment_status"`
	FinancialStatus   string                  `json:"financial_status"`
	TotalPrice        float64                 `json:"total_price"`
	ItemsPrice        float64                 `json:"items_price"`
	DeliveryPrice     float64                 `json:"delivery_price"`
	Discount          map[string]any          `json:"discount"`
	Discounts         []any                   `json:"discounts"`
	Client            *InSalesClientRef       `json:"client"`
	DeliveryTitle     string                  `json:"delivery_title"`
	PaymentTitle      string                  `json:"payment_title"`
	Comment           string                  `json:"comment"`
	CreatedAt         string                  `json:"created_at"`
	UpdatedAt         string                  `json:"updated_at"`
	OrderLines        []InSalesOrderLine      `json:"order_lines"`
	FieldsValues      []InSalesFieldValue     `json:"fields_values"`
	ShippingAddress   *InSalesShippingAddress `json:"shipping_address"`
}

type InSalesShippingAddress struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	Middlename          string `json:"middlename"`
	Phone               string `json:"phone"`
	FormattedPhone      string `json:"formatted_phone"`
	FullName            string `json:"full_name"`
	FullDeliveryAddress string `json:"full_delivery_address"`
}

type InSalesClientRef struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type InSalesOrderLine struct {
	ID              int     `json:"id"`
	ProductID       int     `json:"product_id"`
	VariantID       int     `json:"variant_id"`
	Title           string  `json:"title"`
	Quantity        int     `json:"quantity"`
	SalePrice       float64 `json:"sale_price"`
	TotalPrice      float64 `json:"total_price"`
	DiscountsAmount float64 `json:"discounts_amount"`
	SKU             string  `json:"sku"`
}

type InSalesClientResponse struct {
	ID                    int                   `json:"id"`
	Email                 string                `json:"email"`
	Name                  string                `json:"name"`
	Surname               string                `json:"surname"`
	Middlename            string                `json:"middlename"`
	FirstName             string                `json:"first_name"`
	LastName              string                `json:"last_name"`
	Phone                 string                `json:"phone"`
	FullName              string                `json:"full_name"`
	Comment               string                `json:"comment"`
	Registered            bool                  `json:"registered"`
	Subscribe             bool                  `json:"subscribe"`
	ClientGroupID         *int                  `json:"client_group_id"`
	BonusPoints           int                   `json:"bonus_points"`
	Type                  string                `json:"type"`
	ProgressiveDiscount   string                `json:"progressive_discount"`
	GroupDiscount         string                `json:"group_discount"`
	ConsentToPersonalData bool                  `json:"consent_to_personal_data"`
	CreatedAt             string                `json:"created_at"`
	UpdatedAt             string                `json:"updated_at"`
	FieldsValues          []InSalesFieldValue   `json:"fields_values"`
	ClientTags            []string              `json:"client_tags"`
	DefaultAddress        *InSalesClientAddress `json:"default_address"`
}

type InSalesClientAddress struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	Phone               string `json:"phone"`
	FullDeliveryAddress string `json:"full_delivery_address"`
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

type InSalesPageResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Handle    string `json:"handle"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesAccessoryResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesPropertyValueResponse struct {
	ID          int    `json:"id"`
	AccessoryID int    `json:"accessory_id"`
	Title       string `json:"title"`
	Position    int    `json:"position"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type InSalesMenuResponse struct {
	ID        int               `json:"id"`
	Title     string            `json:"title"`
	Handle    string            `json:"handle"`
	Items     []InSalesMenuItem `json:"items"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type InSalesMenuItem struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Position int    `json:"position"`
}

type InSalesRedirectResponse struct {
	ID        int    `json:"id"`
	Url       string `json:"url"`
	Target    string `json:"target"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesBlogResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Handle    string `json:"handle"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesBlogPostResponse struct {
	ID          int    `json:"id"`
	BlogID      int    `json:"blog_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Permalink   string `json:"permalink"`
	Author      string `json:"author"`
	Handle      string `json:"handle"`
	BodyHTML    string `json:"body_html"`
	PublishedAt string `json:"published_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DiscountCode
type InSalesDiscountCodeResponse struct {
	ID               int    `json:"id"`
	Code             string `json:"code"`
	TypeID           int    `json:"type_id"`
	Discount         string `json:"discount"`
	Description      string `json:"description"`
	Disabled         bool   `json:"disabled"`
	ActOnce          bool   `json:"act_once"`
	ActOnceForClient bool   `json:"act_once_for_client"`
	MinPrice         string `json:"min_price"`
	ExpiredAt        string `json:"expired_at"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Collect
type InSalesCollectResponse struct {
	ID           int    `json:"id"`
	CollectionID int    `json:"collection_id"`
	ProductID    int    `json:"product_id"`
	Position     int    `json:"position"`
	CreatedAt    string `json:"created_at"`
}

// Property
type InSalesPropertyResponse struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Position       int    `json:"position"`
	Permalink      string `json:"permalink"`
	Backoffice     bool   `json:"backoffice"`
	IsHidden       bool   `json:"is_hidden"`
	IsNavigational bool   `json:"is_navigational"`
}

// Characteristic
type InSalesCharacteristicResponse struct {
	ID         int    `json:"id"`
	PropertyID int    `json:"property_id"`
	Title      string `json:"title"`
	Position   int    `json:"position"`
	Permalink  string `json:"permalink"`
}

// Review
type InSalesReviewResponse struct {
	ID           int    `json:"id"`
	ProductID    int    `json:"product_id"`
	Author       string `json:"author"`
	Content      string `json:"content"`
	Email        string `json:"email"`
	Rating       int    `json:"rating"`
	Published    bool   `json:"published"`
	ManagerName  string `json:"manager_name"`
	ManagerReply string `json:"manager_reply"`
	Spam         bool   `json:"spam"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// File
type InSalesFileResponse struct {
	ID          int    `json:"id"`
	AbsoluteURL string `json:"absolute_url"`
	CreatedAt   string `json:"created_at"`
}

// Domain
type InSalesDomainResponse struct {
	ID        int    `json:"id"`
	Domain    string `json:"domain"`
	Main      bool   `json:"main"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ClientGroup
type InSalesClientGroupResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Discount  string `json:"discount"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CustomStatus
type InSalesCustomStatusResponse struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	StatusType string `json:"status_type"`
	Position   int    `json:"position"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ProductField
type InSalesProductFieldResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Handle    string `json:"handle"`
	Type      string `json:"type"`
	FieldType string `json:"field_type"`
	Position  int    `json:"position"`
	IsHidden  *bool  `json:"is_hidden"`
	IsIndexed *bool  `json:"is_indexed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// VariantField
type InSalesVariantFieldResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Handle    string `json:"handle"`
	FieldType string `json:"field_type"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CollectionField
type InSalesCollectionFieldResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Handle    string `json:"handle"`
	FieldType string `json:"field_type"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CollectionFilter
type InSalesCollectionFilterResponse struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	FilterType   string `json:"filter_type"`
	Position     int    `json:"position"`
	CollectionID int    `json:"collection_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// JsTag
type InSalesJsTagResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Position  int    `json:"position"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ExternalDiscount
type InSalesExternalDiscountResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Percent   string `json:"percent"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// DiscountCollectionsGroup
type InSalesDiscountCollectionsGroupResponse struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	AccountID     int    `json:"account_id"`
	CollectionIDs []int  `json:"collection_ids"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Warehouse
type InSalesWarehouseResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SalesChannel
type InSalesSalesChannelResponse struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	SystemName string `json:"system_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// StockCurrency
type InSalesStockCurrencyResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Code      string `json:"code"`
	Rate      string `json:"rate"`
	Main      bool   `json:"main"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PriceKind
type InSalesPriceKindResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// OptionName
type InSalesOptionNameResponse struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Position     int    `json:"position"`
	Navigational bool   `json:"navigational"`
	Permalink    string `json:"permalink"`
}

// OptionValue
type InSalesOptionValueResponse struct {
	ID           int    `json:"id"`
	OptionNameID int    `json:"option_name_id"`
	Title        string `json:"title"`
	Position     int    `json:"position"`
	Permalink    string `json:"permalink"`
	ImageURL     string `json:"image_url"`
}

// Round 4 types

type InSalesApplicationActionResponse struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Handle string `json:"handle"`
	Entity string `json:"entity"`
	URL    string `json:"url"`
}

type InSalesApplicationChargeResponse struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ConfirmationURL string `json:"confirmation_url"`
	Test            bool   `json:"test"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type InSalesApplicationWidgetResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesRecurringApplicationChargeResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesBonusSystemTransactionResponse struct {
	ID        int    `json:"id"`
	ClientID  int    `json:"client_id"`
	Amount    string `json:"amount"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesOAuthAppResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	AppKey    string `json:"app_key"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesUserTaskResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesSupplementaryResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Price     string `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InSalesSimilarResponse struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	SimilarID int    `json:"similar_product_id"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
}

type InSalesUnitResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Unit  string `json:"unit"`
}
