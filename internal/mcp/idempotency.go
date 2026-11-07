package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// idempotencyStore — нужный срез методов репозитория идемпотентности;
// *storage.IdempotencyRepository удовлетворяет, тесты подменяют fake'ом.
type idempotencyStore interface {
	Begin(ctx context.Context, aiKeyID, key, tool, requestHash string) (*storage.IdempotencyRecord, bool, error)
	Finish(ctx context.Context, aiKeyID, key string, response []byte) error
	Cancel(ctx context.Context, aiKeyID, key string) error
}

// idemRepo — репозиторий идемпотентности; задаётся при сборке сервера (NewMCPServer).
// nil = обёртка прозрачна (тесты без БД). Пакетная глобаль: тесты не запускать
// параллельно (t.Parallel) с подменой репозитория.
var idemRepo idempotencyStore

// SetIdempotencyRepo подключает репозиторий идемпотентности к точке регистрации тулов.
// Параметр concrete, чтобы nil-проверка была честной (typed-nil в interface не проскакивает).
func SetIdempotencyRepo(r *storage.IdempotencyRepository) {
	idemRepo = nil
	if r != nil {
		idemRepo = r
	}
}

const maxIdempotencyKeyLen = 128

// withIdempotency — обёртка create-хендлеров: опциональный аргумент idempotency_key
// дедуплицирует повторные вызовы (retry MCP-клиента после таймаута).
// Вставка до выполнения (insert-before-execute) закрывает race параллельных retry:
// дубль получает сохранённый ответ или ошибку «in progress», а не второй заказ.
// Ключ вырезается из аргументов до вызова handler'а (не попадает в ignored_keys/InSales).
func withIdempotency(toolName string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		key := ""
		if args != nil {
			if raw, ok := args["idempotency_key"].(string); ok {
				key = strings.TrimSpace(raw)
			}
			delete(args, "idempotency_key")
		}
		repo := idemRepo
		if key == "" {
			return handler(ctx, req)
		}
		if len(key) > maxIdempotencyKeyLen {
			return mcp.NewToolResultError(fmt.Sprintf("idempotency_key длиннее %d символов", maxIdempotencyKeyLen)), nil
		}
		aiKeyID, _ := ctx.Value(auth.KeyIDKey).(string)
		if repo == nil || aiKeyID == "" {
			return handler(ctx, req) // без репозитория/ключа скоупа дедупликации нет
		}

		hash := hashArgs(args)
		rec, created, err := repo.Begin(ctx, aiKeyID, key, toolName, hash)
		if err != nil {
			// fail-open: без дедупликации поведение как раньше, но дубли возможны
			slog.Warn("idempotency begin failed, fail-open", "tool", toolName, "error", err)
			return handler(ctx, req)
		}
		if !created {
			switch {
			case rec.Status == "done" && rec.RequestHash == hash:
				return replayResult(rec.Response), nil
			case rec.RequestHash != hash:
				return mcp.NewToolResultError("idempotency_key уже использован с другими данными — смените ключ"), nil
			default:
				return mcp.NewToolResultError("запрос с таким idempotency_key уже выполняется — повторите позже"), nil
			}
		}

		res, herr := handler(ctx, req)
		if herr != nil || res == nil || res.IsError {
			if cerr := repo.Cancel(ctx, aiKeyID, key); cerr != nil {
				slog.Warn("idempotency cancel failed", "tool", toolName, "error", cerr)
			}
			return res, herr
		}
		payload, jerr := marshalResult(res)
		if jerr != nil {
			slog.Warn("idempotency marshal failed", "tool", toolName, "error", jerr)
			return res, nil
		}
		if ferr := repo.Finish(ctx, aiKeyID, key, payload); ferr != nil {
			slog.Warn("idempotency finish failed", "tool", toolName, "error", ferr)
		}
		return res, nil
	}
}

// savedResult — сериализованный успешный ответ тула.
type savedResult struct {
	Text       string `json:"text,omitempty"`
	Structured any    `json:"structured,omitempty"`
}

func marshalResult(res *mcp.CallToolResult) ([]byte, error) {
	sr := savedResult{Structured: res.StructuredContent}
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			sr.Text = tc.Text
			break
		}
	}
	if sr.Structured == nil && sr.Text == "" {
		return nil, errors.New("empty result")
	}
	return json.Marshal(sr)
}

// replayResult восстанавливает сохранённый ответ с флагом idempotent_replay.
func replayResult(payload []byte) *mcp.CallToolResult {
	var sr savedResult
	if err := json.Unmarshal(payload, &sr); err != nil {
		slog.Warn("idempotency replay unmarshal failed", "error", err)
		return mcp.NewToolResultError("не удалось восстановить сохранённый ответ idempotency_key")
	}
	if m, ok := sr.Structured.(map[string]any); ok {
		m["idempotent_replay"] = true
		return mcp.NewToolResultStructured(m, sr.Text)
	}
	if sr.Structured != nil {
		return mcp.NewToolResultStructured(sr.Structured, sr.Text)
	}
	return mcp.NewToolResultText(sr.Text + "\n[idempotent_replay: true]")
}
