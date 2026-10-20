package insales

// BuildOrderBody maps MCP tool args to InSales `order` object keys.
// Only known API keys are copied; path/query IDs (order_id) never enter the body.
// `status` is a compatibility alias for `financial_status`.
func BuildOrderBody(args map[string]any) map[string]any {
	if args == nil {
		return map[string]any{}
	}
	order := map[string]any{}

	for _, k := range []string{
		"comment", "manager_comment",
		"financial_status", "custom_status_permalink", "fulfillment_status",
		"coupon",
		"delivery_variant_id", "payment_gateway_id", "warehouse_id",
		"client_id",
	} {
		if v, ok := args[k]; ok && v != nil {
			order[k] = v
		}
	}

	// alias: status → financial_status (only if financial_status not set)
	if v, ok := args["status"]; ok && v != nil {
		if _, exists := order["financial_status"]; !exists {
			order["financial_status"] = v
		}
	}

	for _, k := range []string{
		"order_lines_attributes",
		"discounts_attributes",
		"shipping_address_attributes",
		"client_attributes",
		"client",
		"fields_values_attributes",
	} {
		if v, ok := args[k]; ok && v != nil {
			order[k] = normalizeOrderAttr(v)
		}
	}

	return order
}

// normalizeOrderAttr copies nested attribute structures (array of objects or single object)
// into a plain map/slice so JSON encoding matches InSales examples.
func normalizeOrderAttr(v any) any {
	switch t := v.(type) {
	case []any:
		out := make([]any, 0, len(t))
		for _, item := range t {
			out = append(out, normalizeOrderAttr(item))
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[k] = normalizeOrderAttr(val)
		}
		return out
	default:
		return v
	}
}
