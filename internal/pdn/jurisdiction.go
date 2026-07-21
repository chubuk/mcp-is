package pdn

import "strings"

func DetectJurisdiction(shopURL string) string {
	switch {
	case strings.HasSuffix(shopURL, ".myinsales.ru"):
		return "RU"
	case strings.HasSuffix(shopURL, ".myinsales.kz"):
		return "KZ"
	case strings.HasSuffix(shopURL, ".myinsales.by"):
		return "BY"
	default:
		return "INTL"
	}
}

func DetectLLMProvider(clientName string) string {
	lower := strings.ToLower(clientName)
	switch {
	case strings.Contains(lower, "yandex"):
		return "RU"
	case strings.Contains(lower, "sber"):
		return "RU"
	case strings.Contains(lower, "claude"):
		return "US"
	case strings.Contains(lower, "openai") || strings.Contains(lower, "gpt"):
		return "US"
	default:
		return "UNKNOWN"
	}
}

func ShouldAutoProtectPDN(shopJurisdiction, llmJurisdiction string) bool {
	return shopJurisdiction == "RU" && llmJurisdiction != "RU"
}
