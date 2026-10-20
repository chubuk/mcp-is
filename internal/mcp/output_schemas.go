package mcp

import "github.com/chubuk/mcp-is/internal/integration/insales"

// Output shapes for core tools (WithOutputSchema). ids are strings; lists use the house envelope.

type ProductsPage struct {
	Items   []insales.Product `json:"items"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
}

type OrdersPage struct {
	Items   []insales.Order `json:"items"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
}

type ClientsPage struct {
	Items   []any `json:"items"`
	Total   int   `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}

type SearchProductDetailsOut struct {
	Product    any `json:"product,omitempty"`
	Products   any `json:"products,omitempty"`
	MetaFields any `json:"meta_fields,omitempty"`
	Message    any `json:"message,omitempty"`
}

type OrderFullDetailsOut struct {
	Order      any `json:"order,omitempty"`
	Client     any `json:"client,omitempty"`
	MetaFields any `json:"meta_fields,omitempty"`
}

type ClientFullProfileOut struct {
	Client            any `json:"client,omitempty"`
	RecentOrders      any `json:"recent_orders,omitempty"`
	MetaFields        any `json:"meta_fields,omitempty"`
	BonusTransactions any `json:"bonus_transactions,omitempty"`
}

type CatalogSearchOut struct {
	Items   []insales.Product `json:"items"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
}

type CollectionProductsOut struct {
	Collection any               `json:"collection,omitempty"`
	Products   []insales.Product `json:"products,omitempty"`
	Items      []insales.Product `json:"items,omitempty"`
	Total      int               `json:"total,omitempty"`
	Page       int               `json:"page,omitempty"`
	PerPage    int               `json:"per_page,omitempty"`
}
