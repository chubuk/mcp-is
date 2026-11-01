package pdn

import "strings"

func MaskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 0 {
		return email
	}
	return email[:1] + "***" + email[at:]
}

func MaskPhone(phone string) string {
	if len(phone) < 4 {
		return "***"
	}
	return phone[:2] + " (***) ***-**-**"
}

func MaskName(name string) string {
	runes := []rune(name)
	if len(runes) == 0 {
		return ""
	}
	return string(runes[:1]) + "***"
}

// piiKeys — ключи с персональными данными (152-ФЗ) на любом уровне JSON-дерева.
var piiKeys = map[string]bool{
	"email": true, "phone": true, "formatted_phone": true,
	"first_name": true, "last_name": true, "name": true,
	"surname": true, "middlename": true, "full_name": true,
	"client_name": true,
}

// MaskAny рекурсивно маскирует PII-ключи в произвольном JSON-дереве
// (map[string]any / []any), не меняя структуру и не теряя значения.
func MaskAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return MaskClientData(t)
	case []any:
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = MaskAny(item)
		}
		return out
	default:
		return v
	}
}

// MaskClientData маскирует PII-ключи одного уровня + рекурсит вглубь.
// Не-строковые значения по PII-ключам не теряются — маскируются внутри.
func MaskClientData(data map[string]any) map[string]any {
	result := make(map[string]any, len(data))
	for k, v := range data {
		if piiKeys[k] {
			result[k] = maskPII(k, v)
		} else {
			result[k] = MaskAny(v)
		}
	}
	return result
}

func maskPII(key string, v any) any {
	if s, ok := v.(string); ok {
		switch key {
		case "email":
			return MaskEmail(s)
		case "phone", "formatted_phone":
			return MaskPhone(s)
		default:
			return MaskName(s)
		}
	}
	return MaskAny(v)
}
