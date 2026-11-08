package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chubuk/mcp-is/internal/integration/insales"
)

// SettingsItem — строка снимка настроек (элемент списка в админке InSales).
type SettingsItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Handle string `json:"handle,omitempty"`
}

// SettingsGroup — снимок одной группы «Настройки товаров».
type SettingsGroup struct {
	Title string         `json:"title"`
	Items []SettingsItem `json:"items,omitempty"`
	Error string         `json:"error,omitempty"`
}

// SettingsSnapshot — снимок групп «Настройки товаров» (ключ группы → данные);
// хранится в shops.settings_snapshot и служит базой для диффа при проверке.
type SettingsSnapshot map[string]SettingsGroup

// SettingsGroupDiff — изменения группы с прошлой проверки.
type SettingsGroupDiff struct {
	Title   string         `json:"title"`
	Total   int            `json:"total"`
	Added   []SettingsItem `json:"added,omitempty"`
	Removed []SettingsItem `json:"removed,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// settingsClient — только чтение списков настроек; *insales.InSalesClient удовлетворяет.
type settingsClient interface {
	ListProductFields(ctx context.Context) ([]any, error)
	ListVariantFields(ctx context.Context) ([]any, error)
	ListCollectionFields(ctx context.Context) ([]any, error)
	ListProperties(ctx context.Context) ([]any, error)
	ListOptionNames(ctx context.Context) ([]any, error)
	ListDiscountCollectionsGroups(ctx context.Context) ([]any, error)
	ListPriceKinds(ctx context.Context) ([]any, error)
	ListAccessories(ctx context.Context, page, perPage int) ([]insales.Accessory, error)
}

// settingsGroupDefs — блок «Настройки товаров» админки InSales, в порядке её страницы.
var settingsGroupDefs = []struct{ key, title string }{
	{"product_fields", "Дополнительные поля: Товары"},
	{"variant_fields", "Дополнительные поля: Варианты"},
	{"collection_fields", "Дополнительные поля: Категории"},
	{"properties", "Параметры товаров"},
	{"option_names", "Свойства вариантов"},
	{"discount_collections_groups", "Группы категорий"},
	{"price_kinds", "Типы цен"},
	{"accessories", "Опции"},
}

// collectShopSettings собирает снимок групп; ошибка одной группы не валит снимок.
func collectShopSettings(ctx context.Context, cli settingsClient) SettingsSnapshot {
	snap := make(SettingsSnapshot, len(settingsGroupDefs))
	for _, g := range settingsGroupDefs {
		items, err := fetchSettingsGroup(ctx, cli, g.key)
		grp := SettingsGroup{Title: g.title}
		if err != nil {
			grp.Error = err.Error()
		} else {
			grp.Items = items
		}
		snap[g.key] = grp
	}
	return snap
}

func fetchSettingsGroup(ctx context.Context, cli settingsClient, key string) ([]SettingsItem, error) {
	var list any
	var err error
	switch key {
	case "product_fields":
		list, err = cli.ListProductFields(ctx)
	case "variant_fields":
		list, err = cli.ListVariantFields(ctx)
	case "collection_fields":
		list, err = cli.ListCollectionFields(ctx)
	case "properties":
		list, err = cli.ListProperties(ctx)
	case "option_names":
		list, err = cli.ListOptionNames(ctx)
	case "discount_collections_groups":
		list, err = cli.ListDiscountCollectionsGroups(ctx)
	case "price_kinds":
		list, err = cli.ListPriceKinds(ctx)
	case "accessories":
		var all []insales.Accessory
		for page := 1; page <= 20; page++ {
			batch, perr := cli.ListAccessories(ctx, page, 100)
			if perr != nil {
				return nil, perr
			}
			all = append(all, batch...)
			if len(batch) < 100 {
				break
			}
		}
		list = all
	default:
		return nil, fmt.Errorf("unknown settings group %q", key)
	}
	if err != nil {
		return nil, err
	}
	return extractSettingsItems(list), nil
}

// extractSettingsItems нормализует list-ответ клиента (map или struct) в строки снимка.
func extractSettingsItems(list any) []SettingsItem {
	b, err := json.Marshal(list)
	if err != nil {
		return nil
	}
	var rows []map[string]any
	if err := json.Unmarshal(b, &rows); err != nil {
		return nil
	}
	items := make([]SettingsItem, 0, len(rows))
	for _, m := range rows {
		items = append(items, SettingsItem{
			ID:     anyToString(m["id"]),
			Title:  anyToString(m["title"]),
			Handle: anyToString(m["handle"]),
		})
	}
	return items
}

func anyToString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// diffSettings сравнивает текущий снимок с сохранённым: добавленные/удалённые строки по группам.
// Пустой prev (первый прогон) означает, что всё текущее — added.
func diffSettings(prev, curr SettingsSnapshot) []SettingsGroupDiff {
	out := make([]SettingsGroupDiff, 0, len(settingsGroupDefs))
	for _, g := range settingsGroupDefs {
		p, hasP := prev[g.key]
		c, hasC := curr[g.key]
		d := SettingsGroupDiff{Title: g.title}
		if hasC {
			d.Total = len(c.Items)
			d.Error = c.Error
		}
		comparable := hasC && c.Error == "" && (!hasP || p.Error == "")
		if comparable {
			prevIDs := make(map[string]bool, len(p.Items))
			for _, it := range p.Items {
				prevIDs[it.ID] = true
			}
			currIDs := make(map[string]bool, len(c.Items))
			for _, it := range c.Items {
				currIDs[it.ID] = true
				if !prevIDs[it.ID] {
					d.Added = append(d.Added, it)
				}
			}
			for _, it := range p.Items {
				if !currIDs[it.ID] {
					d.Removed = append(d.Removed, it)
				}
			}
		}
		out = append(out, d)
	}
	return out
}
