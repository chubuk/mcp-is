package mcp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/chubuk/mcp-is/internal/service"
)

type MCPServer struct {
	mcpServer  *server.MCPServer
	httpServer *server.StreamableHTTPServer
	readOnly   bool
	autoPDN       bool
	serverCountry string
}

type mcpContextKey string

const (
	mcpShopNameKey      mcpContextKey = "shop_name"
	mcpKeyNameKey       mcpContextKey = "key_name"
	mcpReadOnlyKey      mcpContextKey = "read_only"
	mcpShopPDNKey       mcpContextKey = "shop_pdn_enabled"
	mcpAIKeyPDNKey      mcpContextKey = "ai_key_pdn_enabled"
)

func NewMCPServer(svc *service.Service, readOnly, autoPDN bool, serverCountry string) *MCPServer {
	s := server.NewMCPServer(
		"mcp-is",
		"0.4.3",
		server.WithToolCapabilities(true),
		server.WithToolFilter(toolFilter(readOnly)),
		server.WithRecovery(),
		server.WithInstructions("MCP-сервер для управления магазином InSales. Доступны инструменты для работы с товарами, заказами, клиентами, категориями и другими ресурсами."),
	)

	RegisterTools(s, svc)

	httpServer := server.NewStreamableHTTPServer(s)

	return &MCPServer{
		mcpServer:  s,
		httpServer: httpServer,
		readOnly:   readOnly,
		autoPDN:       autoPDN,
		serverCountry: serverCountry,
	}
}

func toolFilter(globalReadOnly bool) server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		readOnly := globalReadOnly
		if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}

		if !readOnly {
			return tools
		}

		filtered := make([]mcp.Tool, 0, len(tools))
		for _, t := range tools {
			isWrite := t.Annotations.ReadOnlyHint != nil && !*t.Annotations.ReadOnlyHint
			if !isWrite {
				filtered = append(filtered, t)
			}
		}
		return filtered
	}
}

func (s *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Limit body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	ctx := r.Context()

	shopName, _ := ctx.Value(core.ShopNameKey).(string)
	keyName, _ := ctx.Value(auth.KeyNameKey).(string)
	aiKeyPDN, _ := ctx.Value(auth.KeyPDNKey).(bool)
	readOnly := s.readOnly
	if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
		readOnly = readOnly || v
	}

	// Auto-PDN logic:
	// 1. If auto-PDN disabled → no protection
	// 2. If server NOT in Russia AND shop IS Russian → protect
	// 3. If server IS in Russia AND shop IS Russian AND LLM is foreign → protect
	shopJurisdiction := "INTL"
	llmJurisdiction := "UNKNOWN"
	if shopName != "" {
		shopJurisdiction = pdn.DetectJurisdiction(shopName)
	}
	if keyName != "" {
		llmJurisdiction = pdn.DetectLLMProvider(keyName)
	}
	shopPDN := false
	if s.autoPDN {
		serverInRU := s.serverCountry == "RU"
		shopInRU := shopJurisdiction == "RU"
		llmForeign := llmJurisdiction != "RU"
		
		if (!serverInRU && shopInRU) || (serverInRU && shopInRU && llmForeign) {
			shopPDN = true
		}
	}

	// Manual PDN override takes precedence
	if shopPDN || aiKeyPDN {
		slog.Debug("PDN protection active",
			"shop_jurisdiction", shopJurisdiction,
			"llm_provider", llmJurisdiction,
			"auto_protected", shopPDN,
			"manual_ai_key", aiKeyPDN,
		)
	}

	if shopName != "" || keyName != "" {
		ctx = context.WithValue(ctx, mcpShopNameKey, shopName)
		ctx = context.WithValue(ctx, mcpKeyNameKey, keyName)
		ctx = context.WithValue(ctx, mcpReadOnlyKey, readOnly)
		ctx = context.WithValue(ctx, mcpShopPDNKey, shopPDN)
		ctx = context.WithValue(ctx, mcpAIKeyPDNKey, aiKeyPDN)
	}

	s.httpServer.ServeHTTP(w, r.WithContext(ctx))
}

func (s *MCPServer) NotifyToolsChanged() {
	s.mcpServer.SendNotificationToAllClients(
		"notifications/tools/list_changed",
		nil,
	)
	slog.Info("sent tools/list_changed notification")
}
