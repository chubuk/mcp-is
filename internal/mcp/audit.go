package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// hashArgs is a stable hash of the tool arguments (raw args are never stored — PDN).
func hashArgs(args map[string]any) string {
	if len(args) == 0 {
		return ""
	}
	b, err := json.Marshal(args)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// auditFromContext fills identity columns from auth context (empty when unknown).
func auditFromContext(ctx context.Context, toolName, action string, args map[string]any, ok bool) storage.AuditLogEntry {
	userID, _ := ctx.Value(auth.UserIDKey).(string)
	keyID, _ := ctx.Value(auth.KeyIDKey).(string)
	shopID, _ := ctx.Value(auth.ShopIDKey).(string)
	reqID, _ := ctx.Value(auth.RequestIDKey).(string)
	return storage.AuditLogEntry{
		RequestID: reqID,
		UserID:    userID,
		APIKeyID:  keyID,
		ShopID:    shopID,
		ToolName:  toolName,
		Action:    action,
		ArgsHash:  hashArgs(args),
		ResultOK:  ok,
	}
}

// auditMiddleware logs write-tool calls (ReadOnlyHint false or nil) to audit_log.
// Audit failures never fail the tool call (PDN-safe, best-effort).
func auditMiddleware(repo *storage.AuditRepository, next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := req.Params.Name
		res, err := next(ctx, req)
		// re-check write at call time via registered tool annotations when available
		action := "call"
		ok := true
		if err != nil {
			action = "error"
			ok = false
		} else if res != nil && res.IsError {
			action = "error"
			ok = false
		}
		entry := auditFromContext(ctx, name, action, req.GetArguments(), ok)
		if insertErr := repo.Insert(ctx, entry); insertErr != nil {
			slog.Warn("audit insert failed", "tool", name, "error", insertErr)
		}
		return res, err
	}
}

// isWriteTool reports whether the tool mutates state (explicit annotation is required house-wide).
func isWriteTool(t *mcp.Tool) bool {
	return t.Annotations.ReadOnlyHint == nil || !*t.Annotations.ReadOnlyHint
}

// auditServerMiddleware wraps handlers and skips read-only tools (lookup at call time).
func auditServerMiddleware(s *server.MCPServer, repo *storage.AuditRepository) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			st := s.GetTool(req.Params.Name)
			if st != nil && !isWriteTool(&st.Tool) {
				return next(ctx, req)
			}
			return auditMiddleware(repo, next)(ctx, req)
		}
	}
}

// recordDenied logs a blocked write attempt (read-only session).
func recordDenied(repo *storage.AuditRepository, ctx context.Context, toolName string, args map[string]any) {
	entry := auditFromContext(ctx, toolName, "denied", args, false)
	if err := repo.Insert(ctx, entry); err != nil {
		slog.Warn("audit insert failed", "tool", toolName, "error", err)
	}
}
