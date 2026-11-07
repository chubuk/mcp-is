package mcp

import (
	"encoding/json"

	"github.com/chubuk/mcp-is/internal/integration/insales"
)

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

// ── Raw outputSchema для мутационных инструментов (компромисс, план 3c #4) ──

// mutationResultSchema — форма ответа update/delete-инструментов (см. MutationResult).
var mutationResultSchema = json.RawMessage(`{
  "type": "object",
  "properties": {
    "ok": {"type": "boolean", "description": "Успех операции"},
    "action": {"type": "string", "enum": ["create", "update", "delete"], "description": "Тип мутации"},
    "resource": {"type": "string", "description": "Тип ресурса (product, order, ...)"},
    "id": {"type": "string", "description": "ID затронутого объекта"},
    "title": {"type": "string", "description": "Заголовок объекта (если применимо)"},
    "changed_fields": {
      "type": "object",
      "description": "Изменённые поля: {поле: {old, new}} (update)",
      "additionalProperties": {
        "type": "object",
        "properties": {"old": {}, "new": {}},
        "required": ["old", "new"]
      }
    },
    "ignored_keys": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Неизвестные ключи запроса, проигнорированные сервером"
    },
    "data": {
      "type": "object",
      "description": "Полное итоговое состояние объекта после update",
      "additionalProperties": true
    },
    "deleted": {
      "type": "object",
      "description": "Снимок объекта до удаления (delete); null если не найден",
      "additionalProperties": true
    },
    "message": {"type": "string", "description": "Человекочитаемое пояснение"}
  },
  "required": ["ok", "action", "resource", "id"]
}`)

// createdObjectSchema — форма ответа create-инструментов (созданный объект InSales).
var createdObjectSchema = json.RawMessage(`{
  "type": "object",
  "description": "Созданный объект InSales (id и поля ресурса)",
  "additionalProperties": true
}`)

// outputSchemaFor возвращает raw outputSchema по имени инструмента:
// update_*/delete_*/add_*/decline_* → MutationResult; create_* → generic object; иначе nil.
// Вызывается из addTool — дрейф невозможен: новые мутационные инструменты
// автоматически получают схему по классу.
func outputSchemaFor(name string) json.RawMessage {
	switch ToolClass(name) {
	case ClassDelete:
		return mutationResultSchema
	case ClassWrite:
		if len(name) >= 7 && name[:7] == "create_" {
			return createdObjectSchema
		}
		return mutationResultSchema
	default:
		return nil
	}
}
