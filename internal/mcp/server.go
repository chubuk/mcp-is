package mcp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/core"
	"github.com/your-project/mcp-insales/internal/service"
)

type MCPServer struct {
	mcpServer  *server.MCPServer
	httpServer *server.StreamableHTTPServer
	readOnly   bool
}

func NewMCPServer(svc *service.Service, readOnly bool) *MCPServer {
	s := server.NewMCPServer(
		"mcp-insales",
		"0.2.0",
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
			isReadOnly := t.Annotations.ReadOnlyHint != nil && *t.Annotations.ReadOnlyHint
			if isReadOnly {
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
	readOnly := s.readOnly
	if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
		readOnly = readOnly || v
	}

	if shopName != "" || keyName != "" {
		ctx = context.WithValue(ctx, "shop_name", shopName)
		ctx = context.WithValue(ctx, "key_name", keyName)
		ctx = context.WithValue(ctx, "read_only", readOnly)
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
