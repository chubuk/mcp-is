package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// popFieldValues извлекает и удаляет "field_values" ({handle: значение}) из args.
func popFieldValues(args map[string]any) map[string]any {
	raw, ok := args["field_values"].(map[string]any)
	delete(args, "field_values")
	if !ok || len(raw) == 0 {
		return nil
	}
	return raw
}

// Destiny полей Field (docs InSales): 1=ShippingAddress, 2=Client::Individual,
// 3=Order, 4=OrderLine, 5=Client::Juridical, 6=Client.
var (
	orderFieldDestinies  = map[int]bool{3: true}
	clientFieldDestinies = map[int]bool{2: true, 5: true, 6: true}
)

// inlineFieldAttrs — форма для коллекций/категорий: [{handle, value}] (docs InSales).
func inlineFieldAttrs(fv map[string]any) []map[string]any {
	attrs := make([]map[string]any, 0, len(fv))
	for handle, value := range fv {
		attrs = append(attrs, map[string]any{"handle": handle, "value": value})
	}
	return attrs
}

// fieldsValuesAttrs превращает {handle: value} в fields_values_attributes
// [{field_id, value}] для заказов/клиентов: API принимает только field_id
// (handle молча игнорируется — живая проба 2026-11-10). Оставляет поля
// с указанными destiny; неизвестный handle → ошибка со списком доступных.
func (s *Service) fieldsValuesAttrs(ctx context.Context, destinies map[int]bool, fv map[string]any) ([]map[string]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	fields, err := client.ListFields(ctx)
	if err != nil {
		return nil, fmt.Errorf("list fields: %w", err)
	}
	idByHandle := make(map[string]any, len(fields))
	for _, f := range fields {
		m, ok := f.(map[string]any)
		if !ok {
			continue
		}
		destiny, _ := m["destiny"].(float64)
		if len(destinies) > 0 && !destinies[int(destiny)] {
			continue
		}
		if h, ok := m["handle"].(string); ok && h != "" {
			idByHandle[h] = m["id"]
		}
	}
	attrs := make([]map[string]any, 0, len(fv))
	var missing []string
	for handle, value := range fv {
		id, ok := idByHandle[handle]
		if !ok {
			missing = append(missing, handle)
			continue
		}
		attrs = append(attrs, map[string]any{"field_id": id, "value": value})
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		known := make([]string, 0, len(idByHandle))
		for h := range idByHandle {
			known = append(known, h)
		}
		sort.Strings(known)
		return nil, fmt.Errorf("неизвестные handle полей: %s; доступные: %s",
			strings.Join(missing, ", "), strings.Join(known, ", "))
	}
	return attrs, nil
}

// anyToID — JSON-числа приходят float64: %v даёт экспоненту на больших ID.
func anyToID(v any) string {
	switch n := v.(type) {
	case float64:
		return strconv.FormatFloat(n, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// applyNestedFieldValues пишет значения доп. полей товара/варианта через
// вложенные ресурсы (upsert: есть значение поля — обновляем, нет — создаём).
// Список значений приходит без handle, поэтому матчим по product_field_id:
// handle → id берём из каталога полей магазина.
func (s *Service) applyNestedFieldValues(ctx context.Context, kind, parentID string, fv map[string]any) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	var existing, catalog []any
	switch kind {
	case "product":
		existing, err = client.ListProductFieldValues(ctx, parentID)
		if err == nil {
			catalog, err = client.ListProductFields(ctx)
		}
	case "variant":
		existing, err = client.ListVariantFieldValues(ctx, parentID)
		if err == nil {
			catalog, err = client.ListVariantFields(ctx)
		}
	default:
		return fmt.Errorf("field values: неизвестный тип %q", kind)
	}
	if err != nil {
		return fmt.Errorf("list %s field values: %w", kind, err)
	}
	fieldIDByHandle := make(map[string]string, len(catalog))
	for _, item := range catalog {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if h, ok := m["handle"].(string); ok && h != "" {
			fieldIDByHandle[h] = anyToID(m["id"])
		}
	}
	valueIDByFieldID := make(map[string]string, len(existing))
	for _, item := range existing {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		fid := m["product_field_id"]
		if fid == nil {
			fid = m["variant_field_id"]
		}
		valueIDByFieldID[anyToID(fid)] = anyToID(m["id"])
	}
	var missing []string
	for handle, value := range fv {
		fieldID, ok := fieldIDByHandle[handle]
		if !ok {
			missing = append(missing, handle)
			continue
		}
		valueKey := valueIDByFieldID[fieldID]
		switch {
		case valueKey == "" && kind == "product":
			_, err = client.CreateProductFieldValue(ctx, parentID, map[string]any{"product_field_id": fieldID, "value": value})
		case valueKey == "":
			_, err = client.CreateVariantFieldValue(ctx, parentID, map[string]any{"variant_field_id": fieldID, "value": value})
		case kind == "product":
			_, err = client.UpdateProductFieldValue(ctx, parentID, valueKey, map[string]any{"value": value})
		default:
			_, err = client.UpdateVariantFieldValue(ctx, parentID, valueKey, map[string]any{"value": value})
		}
		if err != nil {
			return fmt.Errorf("field %q (%s %s): %w", handle, kind, parentID, err)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		known := make([]string, 0, len(fieldIDByHandle))
		for h := range fieldIDByHandle {
			known = append(known, h)
		}
		sort.Strings(known)
		return fmt.Errorf("неизвестные handle полей: %s; доступные: %s",
			strings.Join(missing, ", "), strings.Join(known, ", "))
	}
	return nil
}
