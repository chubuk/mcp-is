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

func MaskClientData(data map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range data {
		switch k {
		case "email":
			if s, ok := v.(string); ok {
				result[k] = MaskEmail(s)
			}
		case "phone":
			if s, ok := v.(string); ok {
				result[k] = MaskPhone(s)
			}
		case "first_name":
			if s, ok := v.(string); ok {
				result[k] = MaskName(s)
			}
		case "last_name":
			if s, ok := v.(string); ok {
				result[k] = MaskName(s)
			}
		default:
			result[k] = v
		}
	}
	return result
}
