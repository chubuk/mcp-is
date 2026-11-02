package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	mcpServer     *server.MCPServer
	httpServer    *server.StreamableHTTPServer
	readOnly      bool
	autoPDN       bool
	serverCountry string
}

type mcpContextKey string

const (
	mcpShopNameKey mcpContextKey = "shop_name"
	mcpKeyNameKey  mcpContextKey = "key_name"
	mcpReadOnlyKey mcpContextKey = "read_only"
	mcpShopPDNKey  mcpContextKey = "shop_pdn_enabled"
	mcpAIKeyPDNKey mcpContextKey = "ai_key_pdn_enabled"
)

const baseInstructions = `MCP-сервер для магазина InSales (mcp-is).

## Discovery
Default tools/list is a core surface (~30 tools): products, orders, clients, categories, collections, product fields, pipelines, export.
For anything else use search_insales_tools(query) then call_insales_tool(name, args). All ~280 InSales tools stay callable.
Set MCP_TOOLS_MODE=all to list every tool at once (more tokens, lower selection accuracy).

## Data conventions
- List results: {items, total, page, per_page}. total = items in this response, NOT the catalog size (catalog size: get_products_count / export_products).
- IDs are strings. Field values include handle/title when available.
- Argument values must be language-invariant: pass API ids/handles/enum codes as-is (do not translate).

## Read-only
In read-only sessions write tools are hidden and call_insales_tool rejects them.

## Rate limits
300 requests/min per API key. Near the limit responses slow down instead of failing; on HTTP 429 wait for Retry-After seconds and retry, then pace requests.

## Pipelines
search_product_details, order_full_details, client_full_profile, product_catalog_search, collection_products — multi-step reads in one call.`

// effectiveAlias implements the read-side default chain (S2.3):
// stored connection_name → slug(shop name) → "mcp-is".
func effectiveAlias(connName, shopName string) string {
	if connName != "" {
		return connName
	}
	if shopName != "" {
		return core.Slugify(shopName)
	}
	return "mcp-is"
}

// buildInstructions prepends the per-connection identity block (first ~512 chars
// are self-contained) to the static server guide.
func buildInstructions(ctx context.Context, base string) string {
	connName, _ := ctx.Value(auth.ConnectionNameKey).(string)
	shopName, _ := ctx.Value(auth.ShopNameKey).(string)
	shopURL, _ := ctx.Value(auth.ShopURLKey).(string)
	shopID, _ := ctx.Value(auth.ShopIDKey).(string)

	var header string
	if shopID == "" {
		header = "ВНИМАНИЕ: AI-ключ не привязан к магазину — запросы к данным отклоняются. " +
			"Попросите владельца привязать ключ к магазину в личном кабинете (раздел «AI-ключи»).\n\n"
	} else {
		connName = effectiveAlias(connName, shopName)
		header = fmt.Sprintf("Connection: %s\nStore: %q (%s)\nScope: read + write for this store only.\n"+
			"All tools operate exclusively on this store; no store parameter exists.\n"+
			"When the user names this store or %q, use this connection.\n"+
			"If the user names a different store, say so — do not guess and do not write.\n"+
			"Call store_info to confirm the target store before any write.\n\n",
			connName, shopName, shopURL, connName)
	}
	return header + base
}

func NewMCPServer(svc *service.Service, readOnly, autoPDN bool, serverCountry string, auditRepo *storage.AuditRepository) *MCPServer {
	hooks := &server.Hooks{}
	hooks.AddAfterInitialize(func(ctx context.Context, _ any, _ *mcp.InitializeRequest, result *mcp.InitializeResult) {
		result.Instructions = buildInstructions(ctx, baseInstructions)
	})

	s := server.NewMCPServer(
		"mcp-is",
		"0.9.5",
		server.WithToolCapabilities(true),
		server.WithToolFilter(toolFilter(readOnly)),
		server.WithRecovery(),
		server.WithHooks(hooks),
		server.WithOutputSchemaValidation(),
		server.WithInstructions(baseInstructions),
	)

	RegisterTools(s, svc)
	RegisterDiscoveryTools(s, svc, auditRepo)
	if auditRepo != nil {
		s.Use(auditServerMiddleware(s, auditRepo))
	}

	httpServer := server.NewStreamableHTTPServer(s)

	return &MCPServer{
		mcpServer:     s,
		httpServer:    httpServer,
		readOnly:      readOnly,
		autoPDN:       autoPDN,
		serverCountry: serverCountry,
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
