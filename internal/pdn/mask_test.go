package pdn

import "testing"

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal email", "ivan@example.com", "i***@example.com"},
		{"short local", "ab@domain.com", "a***@domain.com"},
		{"no at sign", "invalid-email", "invalid-email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskEmail(tt.input)
			if got != tt.want {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Russian phone", "+79991234567", "+7 (***) ***-**-**"},
		{"Short phone", "123", "***"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskPhone(tt.input)
			if got != tt.want {
				t.Errorf("MaskPhone(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaskName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Russian name", "Иван", "И***"},
		{"English name", "John", "J***"},
		{"Empty name", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskName(tt.input)
			if got != tt.want {
				t.Errorf("MaskName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaskClientData(t *testing.T) {
	input := map[string]any{
		"email":      "ivan@example.com",
		"phone":      "+79991234567",
		"first_name": "Иван",
		"last_name":  "Петров",
		"order_id":   "12345",
	}

	result := MaskClientData(input)

	if result["email"] != "i***@example.com" {
		t.Errorf("email masked incorrectly: %v", result["email"])
	}
	if result["phone"] != "+7 (***) ***-**-**" {
		t.Errorf("phone masked incorrectly: %v", result["phone"])
	}
	if result["first_name"] != "И***" {
		t.Errorf("first_name masked incorrectly: %v", result["first_name"])
	}
	if result["last_name"] != "П***" {
		t.Errorf("last_name masked incorrectly: %v", result["last_name"])
	}
	if result["order_id"] != "12345" {
		t.Errorf("non-PII field should not be masked: %v", result["order_id"])
	}
}

func TestMaskClientData_Recursive(t *testing.T) {
	input := map[string]any{
		"id": "100",
		"client": map[string]any{
			"full_name": "Иванов Иван",
			"email":     "ivan@example.com",
		},
		"shipping_address": map[string]any{
			"phone": "+79991234567",
			"name":  "Иванов",
		},
		"items": []any{
			map[string]any{"client_name": "Петров Пётр"},
		},
	}
	out := MaskClientData(input)

	client := out["client"].(map[string]any)
	if client["full_name"] != "И***" || client["email"] != "i***@example.com" {
		t.Errorf("nested client not masked: %v", client)
	}
	addr := out["shipping_address"].(map[string]any)
	if addr["phone"] != "+7 (***) ***-**-**" || addr["name"] != "И***" {
		t.Errorf("nested address not masked: %v", addr)
	}
	items := out["items"].([]any)
	if items[0].(map[string]any)["client_name"] != "П***" {
		t.Errorf("item client_name not masked: %v", items[0])
	}
	if out["id"] != "100" {
		t.Errorf("id modified: %v", out["id"])
	}
}

func TestMaskClientData_NonStringValuesKept(t *testing.T) {
	input := map[string]any{
		"client_name": nil,
		"email":       12345, // не строка — значение не теряем
		"nested":      map[string]any{"full_name": "Петров"},
	}
	out := MaskClientData(input)
	if v, ok := out["client_name"]; !ok || v != nil {
		t.Errorf("nil client_name must stay present, got %v (ok=%v)", v, ok)
	}
	if out["email"] != 12345 {
		t.Errorf("non-string email must be kept, got %v", out["email"])
	}
	if out["nested"].(map[string]any)["full_name"] != "П***" {
		t.Errorf("nested under non-PII key not masked: %v", out["nested"])
	}
}

func TestMaskAny_SliceTree(t *testing.T) {
	in := []any{map[string]any{"name": "Иван"}, "plain"}
	out := MaskAny(in).([]any)
	if out[0].(map[string]any)["name"] != "И***" {
		t.Errorf("map in slice not masked: %v", out[0])
	}
	if out[1] != "plain" {
		t.Errorf("plain value modified: %v", out[1])
	}
}
