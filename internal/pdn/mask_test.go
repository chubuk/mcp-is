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
