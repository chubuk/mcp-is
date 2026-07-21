package pdn

import "testing"

func TestDetectJurisdiction(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"Russian shop", "myshop.myinsales.ru", "RU"},
		{"Kazakh shop", "myshop.myinsales.kz", "KZ"},
		{"Belarus shop", "myshop.myinsales.by", "BY"},
		{"International shop", "myshop.shopify.com", "INTL"},
		{"Empty URL", "", "INTL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectJurisdiction(tt.url)
			if got != tt.want {
				t.Errorf("DetectJurisdiction(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestDetectLLMProvider(t *testing.T) {
	tests := []struct {
		name string
		client string
		want string
	}{
		{"Claude Desktop", "Claude Desktop", "US"},
		{"OpenAI", "openai-gpt4", "US"},
		{"GPT", "gpt-chat", "US"},
		{"YandexGPT", "YandexGPT", "RU"},
		{"SberGPT", "SberGPT", "RU"},
		{"Unknown client", "my-custom-client", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectLLMProvider(tt.client)
			if got != tt.want {
				t.Errorf("DetectLLMProvider(%q) = %q, want %q", tt.client, got, tt.want)
			}
		})
	}
}

func TestShouldAutoProtectPDN(t *testing.T) {
	tests := []struct {
		name       string
		shop       string
		llm        string
		want       bool
	}{
		{"RU shop + US LLM", "RU", "US", true},
		{"RU shop + UNKNOWN LLM", "RU", "UNKNOWN", true},
		{"RU shop + RU LLM", "RU", "RU", false},
		{"KZ shop + US LLM", "KZ", "US", false},
		{"INTL shop + US LLM", "INTL", "US", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldAutoProtectPDN(tt.shop, tt.llm)
			if got != tt.want {
				t.Errorf("ShouldAutoProtectPDN(%q, %q) = %v, want %v", tt.shop, tt.llm, got, tt.want)
			}
		})
	}
}
