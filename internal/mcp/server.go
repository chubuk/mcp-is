package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/core"
	"github.com/your-project/mcp-insales/internal/service"
)

type Server struct {
	handler  *Handler
	readOnly bool
}

func NewServer(svc *service.Service, readOnly bool) *Server {
	return &Server{
		handler:  NewHandler(svc),
		readOnly: readOnly,
	}
}

func (s *Server) HandleMCP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}

	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	method, _ := req["method"].(string)
	id := req["id"]
	params, _ := req["params"].(map[string]any)

	slog.Info("mcp request", "method", method, "id", id)

	var result any
	var rpcError any

	switch method {
	case "initialize":
		ctx := r.Context()
		shopName, _ := ctx.Value(core.ShopNameKey).(string)
		keyName, _ := ctx.Value(auth.KeyNameKey).(string)
		readOnly := s.readOnly
		if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}

		serverInfo := map[string]any{
			"name":    "mcp-insales",
			"version": "0.1.0",
		}
		if shopName != "" {
			serverInfo["shop"] = shopName
		}
		if keyName != "" {
			serverInfo["key"] = keyName
		}
		serverInfo["read_only"] = readOnly

		result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": serverInfo,
		}

	case "tools/list":
		allTools := GetAllTools()
		ctx := r.Context()
		readOnly := s.readOnly
		if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}
		if readOnly {
			filtered := make([]Tool, 0)
			for _, t := range allTools {
				if !isWriteTool(t.Name) {
					filtered = append(filtered, t)
				}
			}
			allTools = filtered
		}
		result = map[string]any{
			"tools": allTools,
		}

	case "tools/call":
		toolName, _ := params["name"].(string)
		arguments, _ := params["arguments"].(map[string]any)

		if s.readOnly && isWriteTool(toolName) {
			rpcError = map[string]any{
				"code":    -32000,
				"message": fmt.Sprintf("tool '%s' is not available in read-only mode", toolName),
			}
			break
		}

		ctx := r.Context()
		userID, _ := ctx.Value(auth.UserIDKey).(string)
		if userID != "" {
			ctx = context.WithValue(ctx, "account_id", userID)
		} else {
			ctx = context.WithValue(ctx, "account_id", "default")
		}

		output, err := s.handler.HandleTool(ctx, toolName, arguments)
		if err != nil {
			rpcError = map[string]any{
				"code":    -32000,
				"message": err.Error(),
			}
		} else {
			result = map[string]any{
				"content": []map[string]any{
					{
						"type": "text",
						"text": formatResult(output),
					},
				},
			}
		}

	default:
		rpcError = map[string]any{
			"code":    -32601,
			"message": fmt.Sprintf("method not found: %s", method),
		}
	}

	resp := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
	}

	if rpcError != nil {
		resp["error"] = rpcError
	} else {
		resp["result"] = result
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "event: endpoint\ndata: /mcp\n\n")
	w.(http.Flusher).Flush()

	<-r.Context().Done()
}

func (s *Server) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	event := r.PathValue("event")
	slog.Info("webhook received", "event", event)
	w.WriteHeader(http.StatusOK)
}

func formatResult(v any) string {
	if v == nil {
		return "null"
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

var writeTools = map[string]bool{
	"create_product":        true,
	"update_product":        true,
	"delete_product":        true,
	"create_order":          true,
	"update_order":          true,
	"create_client":         true,
	"update_client":         true,
	"delete_client":         true,
	"create_category":       true,
	"update_category":       true,
	"delete_category":       true,
	"create_collection":     true,
	"update_collection":     true,
	"delete_collection":     true,
	"add_image":             true,
	"delete_image":          true,
	"create_meta_field":     true,
	"update_meta_field":     true,
	"delete_meta_field":     true,
	"create_webhook":        true,
	"update_webhook":        true,
	"delete_webhook":        true,
	"create_page":           true,
	"update_page":           true,
	"delete_page":           true,
	"create_accessory":      true,
	"update_accessory":      true,
	"delete_accessory":      true,
	"update_variant":        true,
	"delete_variant":        true,
	"create_property_value": true,
	"update_property_value": true,
	"delete_property_value": true,
	"create_menu":           true,
	"update_menu":           true,
	"delete_menu":           true,
	"create_redirect":       true,
	"delete_redirect":       true,
	"create_blog":           true,
	"delete_blog":           true,
	"create_blog_post":      true,
	"update_blog_post":      true,
	"delete_blog_post":      true,
}

func isWriteTool(name string) bool {
	return writeTools[name]
}
