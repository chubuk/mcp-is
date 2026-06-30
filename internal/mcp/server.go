package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/your-project/mcp-insales/internal/auth"
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

	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	method, _ := req["method"].(string)
	id := req["id"]
	params, _ := req["params"].(map[string]interface{})

	slog.Info("mcp request", "method", method, "id", id)

	var result interface{}
	var rpcError interface{}

	switch method {
	case "initialize":
		result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "mcp-insales",
				"version": "0.1.0",
			},
		}

	case "tools/list":
		result = map[string]interface{}{
			"tools": GetAllTools(),
		}

	case "tools/call":
		toolName, _ := params["name"].(string)
		arguments, _ := params["arguments"].(map[string]interface{})

		if s.readOnly && isWriteTool(toolName) {
			rpcError = map[string]interface{}{
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
			rpcError = map[string]interface{}{
				"code":    -32000,
				"message": err.Error(),
			}
		} else {
			result = map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": formatResult(output),
					},
				},
			}
		}

	default:
		rpcError = map[string]interface{}{
			"code":    -32601,
			"message": fmt.Sprintf("method not found: %s", method),
		}
	}

	resp := map[string]interface{}{
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

func formatResult(v interface{}) string {
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
	"create_product":      true,
	"update_product":      true,
	"delete_product":      true,
	"create_order":        true,
	"update_order":        true,
	"create_client":       true,
	"update_client":       true,
	"delete_client":       true,
	"create_category":     true,
	"update_category":     true,
	"delete_category":     true,
	"create_collection":   true,
	"update_collection":   true,
	"delete_collection":   true,
	"add_image":           true,
	"delete_image":        true,
	"create_meta_field":   true,
	"update_meta_field":   true,
	"delete_meta_field":   true,
	"create_webhook":      true,
	"update_webhook":      true,
	"delete_webhook":      true,
	"create_page":         true,
	"update_page":         true,
	"delete_page":         true,
	"create_accessory":    true,
	"update_accessory":    true,
	"delete_accessory":    true,
}

func isWriteTool(name string) bool {
	return writeTools[name]
}
