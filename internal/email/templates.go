package email

import (
	"embed"
	"fmt"
	"sort"
)

// Template is a mail template stored in memory or DB.
type Template struct {
	Key          string
	Name         string
	Subject      string
	HTMLBody     string
	TextBody     string
	Description  string
	RequiredVars []string
	IsModified   bool
}

// TemplateStore provides templates by key.
type TemplateStore interface {
	Get(key string) (Template, error)
	List() ([]Template, error)
}

// FuncStore adapts function lookups (e.g. DB repository) to TemplateStore.
type FuncStore struct {
	GetFunc  func(key string) (Template, error)
	ListFunc func() ([]Template, error)
}

func (f FuncStore) Get(key string) (Template, error) {
	if f.GetFunc == nil {
		return Template{}, fmt.Errorf("template %q not found", key)
	}
	return f.GetFunc(key)
}

func (f FuncStore) List() ([]Template, error) {
	if f.ListFunc == nil {
		return nil, nil
	}
	return f.ListFunc()
}

//go:embed templates/default/*.html templates/default/*.txt
var defaultTemplates embed.FS

// MemoryStore serves built-in defaults (fallback when DB store is not set).
type MemoryStore struct {
	items map[string]Template
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{items: make(map[string]Template)}
	for _, t := range builtinTemplates() {
		s.items[t.Key] = t
	}
	return s
}

func (s *MemoryStore) Get(key string) (Template, error) {
	t, ok := s.items[key]
	if !ok {
		return Template{}, fmt.Errorf("template %q not found", key)
	}
	return t, nil
}

func (s *MemoryStore) List() ([]Template, error) {
	out := make([]Template, 0, len(s.items))
	for _, t := range s.items {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}

// Default returns the embedded seed for a key (used by "reset to default").
func Default(key string) (Template, error) {
	for _, t := range builtinTemplates() {
		if t.Key == key {
			return t, nil
		}
	}
	return Template{}, fmt.Errorf("default template %q not found", key)
}

// DefaultList returns all embedded seeds.
func DefaultList() []Template {
	return builtinTemplates()
}

func builtinTemplates() []Template {
	html := func(name string) string {
		b, err := defaultTemplates.ReadFile("templates/default/" + name)
		if err != nil {
			return ""
		}
		return string(b)
	}
	txt := func(name string) string {
		b, err := defaultTemplates.ReadFile("templates/default/" + name)
		if err != nil {
			return ""
		}
		return string(b)
	}

	return []Template{
		{
			Key:          "reset_password",
			Name:         "Сброс пароля",
			Subject:      "Сброс пароля — MCP IS",
			HTMLBody:     html("reset_password.html"),
			TextBody:     txt("reset_password.txt"),
			Description:  "Ссылка для смены пароля (1 час)",
			RequiredVars: []string{"Name", "ResetURL"},
		},
		{
			Key:          "verify_email",
			Name:         "Подтверждение email",
			Subject:      "Подтвердите email — MCP IS",
			HTMLBody:     html("verify_email.html"),
			TextBody:     txt("verify_email.txt"),
			Description:  "Подтверждение адреса почты",
			RequiredVars: []string{"Name", "VerifyURL"},
		},
		{
			Key:          "login_code",
			Name:         "Код подтверждения",
			Subject:      "Код подтверждения — MCP IS",
			HTMLBody:     html("login_code.html"),
			TextBody:     txt("login_code.txt"),
			Description:  "Одноразовый код: подтверждение email при регистрации и привязка магазина",
			RequiredVars: []string{"Name", "Code"},
		},
		{
			Key:          "welcome",
			Name:         "Приветствие",
			Subject:      "Добро пожаловать в MCP IS",
			HTMLBody:     html("welcome.html"),
			TextBody:     txt("welcome.txt"),
			Description:  "После установки приложения / регистрации",
			RequiredVars: []string{"Name", "MCPURL"},
		},
		{
			Key:          "shop_linked",
			Name:         "Магазин привязан",
			Subject:      "Магазин привязан — MCP IS",
			HTMLBody:     html("shop_linked.html"),
			TextBody:     txt("shop_linked.txt"),
			Description:  "Уведомление о привязке магазина InSales к аккаунту",
			RequiredVars: []string{"Name", "Shop", "Date"},
		},
		{
			Key:          "shop_unlinked",
			Name:         "Магазин отвязан",
			Subject:      "Магазин отвязан — MCP IS",
			HTMLBody:     html("shop_unlinked.html"),
			TextBody:     txt("shop_unlinked.txt"),
			Description:  "Уведомление об отвязке магазина InSales (удаление приложения)",
			RequiredVars: []string{"Name", "Shop", "Date"},
		},
		{
			Key:          "test",
			Name:         "Тестовое письмо",
			Subject:      "SMTP Test — MCP IS",
			HTMLBody:     html("test.html"),
			TextBody:     txt("test.txt"),
			Description:  "Проверка SMTP из админки",
			RequiredVars: []string{"Name"},
		},
	}
}
