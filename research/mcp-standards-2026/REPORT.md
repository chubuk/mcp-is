# Стандарты MCP-серверов (2025–2026) vs mcp-is

> Generated 2026-10-20 · depth: standard (5 angles) · 69 findings · workspace: research/mcp-standards-2026/

## Executive summary

1. **Актуальная ревизия MCP — 2026-07-28** (не 2025-11-25). Это breaking-редизайн: stateless, без `initialize`/`Mcp-Session-Id`, обязателен `server/discover`, `resultType`, `ttlMs`/`cacheScope` [1][2]. mcp-is на mcp-go v0.56 — «legacy era» (handshake + сессии) [3][7].
2. **245 tools — это за пределами cliff.** Anthropic: точность выбора инструмента падает после 30–50 tools; ~55k tokens на определения в типовом multiserver-сетапе [20][21]. Sentry держит публичный бюджет **~20–25 tools**, хвост прячет за `search_sentry_tools`+`execute_sentry_tool` [12]. GitHub режет surface через **toolsets** в URL path + `/readonly` [11]. MCP roadmap (2026-08) прямо называет 100-tool servers проблемой протокола [22].
3. **Industry-ответ — progressive discovery**, а не «ещё CRUD-обёртки». Claude Code: `ENABLE_TOOL_SEARCH`, `_meta["anthropic/alwaysLoad"]`, порог ~10% контекста [23]. Spec client-best-practices: `search_tools` meta-tool, порог 1–5% контекста [24].
4. **Описания tools — главный рычаг качества.** Anthropic: ≥3–4 предложения (что делает / когда НЕ использовать / параметры / caveats); мелкие правки описаний дали SOTA на SWE-bench [25][26]. У mcp-is среднее описание ~27 символов — критический gap.
5. **`structuredContent` + `outputSchema`** — стандарт с 2025-06-18; JSON-зеркало в TextContent обязательно для back-compat [5][27]. GitHub включает typed schemas только на protocol ≥2026-07-28 [28]. У mcp-is — только `NewToolResultText` (pretty JSON).
6. **Annotations.** Sentry требует **явных** `readOnlyHint`+`destructiveHint`+`openWorldHint` у каждого tool (compile-time) [13]. Defaults контринтуитивны: `destructiveHint` default **true** [6]. У mcp-is 245×readOnlyHint, 47×destructiveHint, 0×idempotent/openWorld; nil ReadOnlyHint = утечка write в RO filter.
7. **1:1 API-wrapping — anti-pattern.** Anthropic: консолидировать multi-step workflows; GitHub: verb-family (`issue_read`/`projects_write`) + silent deprecation aliases [29][30][31]. План «+36 tools к 245» усиливает sprawl, если не закрыть discovery.
8. **Auth.** Для HTTP-транспорта spec 2026-07-28 делает OAuth 2.1 + RFC 8707 audience binding обязательным профилем; token passthrough запрещён [32][33]. API-key (как у mcp-is) допустим для private/self-hosted, но не «spec-preferred». Per-shop InSales credentials из БД — правильный паттерн (не passthrough).
9. **i18n.** English descriptions — de facto norm; MCP i18n отложен (Server Card WG out of scope) [34]. RU-описания не запрещены, но EN в schema + RU в UI/title безопаснее для tool-selection.
10. **План v2 нужно усилить**: (A) progressive discovery / tool surface architecture **до** добавления новых tools; (B) descriptions 3–4 предложения; (C) полный annotation set; (D) structuredContent; (E) roadmap на 2026-07-28 / OAuth — отдельно, не в этом спринте.

## Background & scope

Проект **mcp-is** — публичный Go MCP-сервер для InSales (mark3labs/mcp-go v0.56.0, Streamable HTTP, JWT/API-key, ~245 tools). Вопрос: какие стандарты MCP-серверов сложились к 2026 году, насколько mcp-is им соответствует, и что изменить в плане доработок v2 (полное покрытие InSales API + синхронизация + тесты). Out of scope: внутренности InSales API, реализация (plan mode).

**Источники:** 69 findings в 5 файлах (F1–F5), в основном primary (modelcontextprotocol spec repo, Anthropic engineering, github/github-mcp-server, getsentry/sentry-mcp, Notion, Claude docs, arXiv 2026). WebSearch был недоступен — прямые fetches primary URL.

---

## 1. Что такое «стандарт MCP-сервера» в 2026

### 1.1 Линия протокола

| Ревизия | Статус | Ключевое |
|---------|--------|----------|
| 2024-11-05 | legacy | initialize handshake, HTTP+SSE |
| 2025-03-26 | legacy | tool annotations |
| 2025-06-18 | legacy | `structuredContent`, `outputSchema` |
| 2025-11-25 | legacy (текущий у большинства SDK) | initialize MUST, `Mcp-Session-Id`, SEP-1303 validation→isError |
| **2026-07-28** | **latest stable** | stateless, `server/discover`, `resultType`, `ttlMs`/`cacheScope`, no sessions, OAuth 2.1+RFC8707, header/body mirror |

«Legacy» = handshake era; «modern» = 2026-07-28+ [1][2][3].

### 1.2 Contract tools (сводка MUST/SHOULD)

- `tools/list` → `name` + `inputSchema` (valid JSON Schema object) [4]
- Имя: 1–128 chars, `[A-Za-z0-9_.-]`, unique per server [8]
- Annotations: hints only, defaults `readOnlyHint=false`, `destructiveHint=true`, `idempotentHint=false`, `openWorldHint=true` [6]
- Результат: `structuredContent` (+ JSON в TextContent back-compat); `outputSchema` → server MUST conform [5][27]
- Ошибки: tool-level → `isError:true` (model self-correct); protocol-level → JSON-RPC [7][10]
- Pagination lists: opaque cursor (не page numbers) [9]
- `InitializeResult.instructions` — free-form guidance для system prompt [16]

### 1.3 Production-паттерны (vendors)

| Паттерн | Кто | Суть |
|---------|-----|------|
| Verb-family tools | GitHub | `issue_read`/`issue_write`, старые имена — silent aliases [30][31] |
| Toolsets in URL | GitHub | `/mcp/issues`, `/mcp/repos/readonly` [11] |
| Tool budget ~20–25 | Sentry | «Never exceed 25»; Cursor caps 45 slots/provider [12] |
| Catalog meta-tools | Sentry | `search_*` + `execute_*`, top-level only foundational [13] |
| Scoped multi-tenant URL | Sentry | `/mcp/{org}/{project}`, hide discovery when scoped [14] |
| Explicit annotations | Sentry | all hints, compile-time enforced [13] |
| structuredContent = truth | Sentry | content = compatibility view, no raw upstream JSON [15] |
| Compact DTOs | GitHub | strip hypermedia URLs [28] |
| Admin-curated subset | Notion | ~24 tools, `notion-` prefix, soft-degrade + `notices` [35][36] |
| Deprecation aliases | GitHub, Notion | old names callable after removal from tools/list [31][36] |

---

## 2. Сверка с mcp-is (что есть сейчас)

| Элемент | Стандарт 2026 | mcp-is сегодня | Оценка |
|---------|---------------|----------------|--------|
| Транспорт | Streamable HTTP (2026 = stateless) | StreamableHTTP + sessions (mcp-go v0.56) | OK short-term; migration roadmap needed |
| tools/list | name+schema+annotations | есть, RU descriptions | schema OK; descriptions too thin |
| Tool count | budget 20–50 visible | **245** (108 RO / 137 write) | **критично** |
| Progressive discovery | search_tools / toolsets / defer_loading | нет | **критично** |
| Descriptions | 3–4 sentences | ~27 chars | **высокий** |
| structuredContent | dual text+structured | only TextContent pretty JSON | **средний** (plan E) |
| outputSchema | optional, key tools | нет | **средний** |
| readOnlyHint | explicit | 245/245 | OK |
| destructiveHint | explicit | 47/245 | средний |
| idempotent / openWorld | explicit | 0 | низкий–средний |
| nil-RO leak | must not | toolFilter treats nil as RO | **высокий** (plan R) |
| isError | model-fixable text | NewToolResultError + cleanError | OK |
| Instructions | rich guidance | 1 generic RU sentence | высокий (plan A4) |
| ID types | semantic / string | mixed int/string | средний (plan E) |
| Pagination | cursor (lists) / page (data) | page/per_page + house envelope | OK for data tools |
| Auth | OAuth 2.1+audience (HTTP) | API-key Bearer / JWT | acceptable private; not spec-preferred |
| Token passthrough | FORBIDDEN | нет (creds from DB) | OK |
| Registry / server.json | reverse-DNS, versioned | нет | низкий (optional) |
| EN descriptions | de facto | RU only | решение пользователя |
| Deprecation policy | ≥12 months (SEP-2596) | нет aliases | низкий пока нет rename |

---

## 3. Критический ревью плана v2

### Что в плане правильно (оставить)

- **Шаг 0** — фикс CreateOrder/UpdateOrder payload: correctness важнее любого standards work.
- **B3 product custom fields** — реальный user need, данные без handle/title непригодны LLM.
- **R ReadOnlyHint mandatory + toolFilter** — Sentry-level annotation hygiene [13].
- **A1/A2 sync-тесты** — drift lists 96 vs 245 must die.
- **E structuredContent / envelope / string IDs / API.md protocolVersion** — aligned with [5][27][28].
- **Docs pass before push** — process rule, не спорить.
- **Не трогать OAuth/transport** в этом спринте — правильно по YAGNI.

### Что план упускает (добавить)

1. **Tool surface architecture (новый приоритет, ДО B*).**  
   +36 tools к 245 = хуже для LLM [20][22][24]. Без discovery-слоя «полное покрытие» вредит usability.  
   Минимальный viable подход (без ломки API):
   - **Core visible set** (~25–40 tools): пайплайны, list/get ключевых ресурсов, product fields (B3), search.
   - **Catalog meta-tools**: `search_insales_tools(query)` + `call_insales_tool(name, args)` ИЛИ mcp-go `WithToolFilter` + query-параметры URL toolsets (GitHub-style `/mcp?tools=products,orders&readonly`).
   - `instructions` объявляет discovery-механику.
   - Старые 245 остаются callable; в `tools/list` по умолчанию — core (или полный список, если клиент уже знает про tool search — fallback).
2. **Descriptions quality pass** — не только 16 пропущенных arg-description, а 3–4 предложения на ключевых tools (what / when-not / params / response shape) [25][26].
3. **Полный annotation set** — `openWorldHint` (false для InSales CRUD), `idempotentHint` (true у get/update без side effects), `destructiveHint` на всех delete/mutate.
4. **EN tool descriptions (или bilingual)** — decision point; EN в schema улучшает selection, RU остаётся в UI/docs. Альтернатива: EN `description` + RU `title`/`_meta`.
5. **Roadmap-строка (out of scope sprint):** MCP 2026-07-28 (`server/discover`, `resultType`), OAuth 2.1 + RFC 9728 metadata, deprecation aliases policy, Registry `server.json`.
6. **Не консолидировать** существующие имена в verb-family сейчас (breaking); при rename в будущем — GitHub-style aliases [31].

### Конфликт «всё покрыть» vs «tool sprawl»

Пользователь: «реализовывать надо все». Research: 245+ flat tools = anti-pattern.  
**Резолюция:** полное покрытие InSales — да, но **за discovery-слоем**; default `tools/list` = core + meta-tools; полный каталог — через search/execute или toolset URL. Так LLM получает и coverage, и точность выбора.

---

## 4. Дельты к плану (кратко)

| # | Дельта | Куда в план |
|---|--------|-------------|
| D1 | **P0 Tool surface / progressive discovery** (core set + catalog meta-tools или URL toolsets) | новый шаг перед B* |
| D2 | Descriptions 3–4 sentences на list/get/pipeline/field tools | расширить E |
| D3 | Explicit `openWorldHint`/`idempotentHint`/`destructiveHint` | расширить R/A2 |
| D4 | `structuredContent` helper + outputSchema на ключевых | уже в E — усилить priority |
| D5 | EN description / RU title — decision | в Open questions плана |
| D6 | Roadmap: 2026-07-28, OAuth 2.1, Registry | Out of scope + Future |
| D7 | Порядок: Step0 → **P0 discovery** → A sync → B3 → B1/B2/B4 → E → docs | D |

---

## Open questions

1. **mcp-go v0.56**: поддержка tool search / `defer_loading` / cache hints / `resultType`? (findings F1 follow-up) — влияет на реализацию P0.
2. **EN vs RU descriptions** — требует решения пользователя (правило RU vs tool-selection quality).
3. **Catalog vs toolsets URL** — что удобнее целевым клиентам (MiMoCode/Cursor/Claude Desktop)? Можно сделать catalog meta-tools + опциональный `?tools=` filter.
4. Empirics: A/B prefix vs suffix naming, JSON vs Markdown — spec говорит «evaluate yourself» [25].
5. Live-валидация order payload / field values против InSales sandbox (plan risk).

## Sources

[1] MCP changelog 2026-07-28 — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/specification/2026-07-28/changelog.mdx (accessed 2026-10-20)
[2] MCP docs.json version list — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/docs.json (accessed 2026-10-20)
[3] MCP lifecycle 2025-11-25 — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/specification/2025-11-25/basic/lifecycle.mdx (accessed 2026-10-20)
[4] MCP tools.mdx 2025-11-25 — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/specification/2025-11-25/server/tools.mdx (accessed 2026-10-20)
[5] MCP tools.mdx structuredContent — same as [4] (2025-06-18 feature)
[6] MCP schema.json annotations — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/schema/2025-11-25/schema.json (accessed 2026-10-20)
[7] MCP tools error contract — [4] / tools.mdx 2026-07-28
[8] MCP tool name charset SEP-986 — tools.mdx 2025-11-25
[9] MCP pagination — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/specification/2025-11-25/server/utilities/pagination.mdx (accessed 2026-10-20)
[10] MCP error channels 2026-07-28 — tools.mdx 2026-07-28
[11] GitHub remote-server toolsets — https://raw.githubusercontent.com/github/github-mcp-server/main/docs/remote-server.md (accessed 2026-10-20)
[12] Sentry adding-tools budget — https://raw.githubusercontent.com/getsentry/sentry-mcp/main/docs/contributing/adding-tools.md (accessed 2026-10-20)
[13] Sentry adding-tools catalog+annotations — same as [12]
[14] Sentry product MCP scoping — https://docs.sentry.io/product/sentry-mcp/ (accessed 2026-10-20)
[15] Sentry tool-responses — https://raw.githubusercontent.com/getsentry/sentry-mcp/main/docs/contributing/tool-responses.md (accessed 2026-10-20)
[16] MCP schema.mdx 2026-07-28 instructions/order — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/specification/2026-07-28/schema.mdx (accessed 2026-10-20)
[17] MCP resources/prompts capabilities — resources.mdx 2025-11-25
[18] MCP progress — basic/utilities/progress.mdx 2025-11-25
[19] MCP 2026-07-28 resultType/cache — changelog.mdx 2026-07-28
[20] Anthropic tool-search 55k tokens / 30–50 cliff — https://docs.claude.com/en/docs/agents-and-tools/tool-use/tool-search-tool (accessed 2026-10-20)
[21] Claude Platform define-tools 3–4 sentences — https://platform.claude.com/docs/en/agents-and-tools/tool-use/define-tools (accessed 2026-10-20)
[22] MCP roadmap progressive discovery — https://blog.modelcontextprotocol.io/posts/mcp-roadmap/ (2026-08-22, accessed 2026-10-20)
[23] Claude Code MCP ENABLE_TOOL_SEARCH — https://docs.claude.com/en/docs/claude-code/mcp (accessed 2026-10-20)
[24] MCP client-best-practices search_tools — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/docs/2026-07-28/develop/clients/client-best-practices.mdx (accessed 2026-10-20)
[25] Anthropic writing-tools-for-agents — https://www.anthropic.com/engineering/writing-tools-for-agents (2025-09-11, accessed 2026-10-20)
[26] Anthropic SWE-bench description refinements — same as [25]
[27] MCP structuredContent MUST — tools.mdx 2025-11-25
[28] GitHub typed-tool-schemas — https://raw.githubusercontent.com/github/github-mcp-server/main/docs/typed-tool-schemas.md (accessed 2026-10-20)
[29] Anthropic no 1:1 API wrap — same as [25]
[30] GitHub tool-renaming verb-family — https://raw.githubusercontent.com/github/github-mcp-server/main/docs/tool-renaming.md (accessed 2026-10-20)
[31] GitHub DeprecatedToolAliases — same as [30]
[32] MCP authorization 2026-07-28 — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/specification/2026-07-28/basic/authorization/index.mdx (accessed 2026-10-20)
[33] MCP token passthrough forbidden — .../basic/authorization/security-considerations.mdx (accessed 2026-10-20)
[34] Server Card WG i18n out of scope — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/community/working-groups/server-card.mdx (accessed 2026-10-20)
[35] Notion MCP supported tools — https://developers.notion.com/guides/mcp/mcp-supported-tools.md (accessed 2026-10-20)
[36] Notion deprecation callable-by-name — same as [35]
[37] MCP Registry preview — https://blog.modelcontextprotocol.io/posts/2025-09-08-mcp-registry-preview/ (2025-09-08)
[38] MCP Registry schema reverse-DNS — https://static.modelcontextprotocol.io/schemas/2025-12-11/server.schema.json
[39] SEP-2596 deprecation ≥12 months — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/seps/2596-spec-feature-lifecycle-and-deprecation.md
[40] MCP deprecated.mdx 2026-07-28 — .../docs/specification/2026-07-28/deprecated.mdx
[41] Streamable HTTP 2026 stateless / Origin MUST — .../basic/transports/streamable-http.mdx
[42] MCP security best practices state handles / scopes — .../docs/2026-07-28/tutorials/security/security_best_practices.mdx
[43] arXiv MCP-FSD filesystem laundering — https://arxiv.org/abs/2610.22301 (2026-10-15)
[44] arXiv Approve What Exactly — https://arxiv.org/abs/2610.21063 (2026-10-15)
[45] mcp-go README WithToolFilter — https://raw.githubusercontent.com/mark3labs/mcp-go/main/README.md (accessed 2026-10-20)
[46] GitHub error-handling user-actionable — https://raw.githubusercontent.com/github/github-mcp-server/main/docs/error-handling.md
[47] GitHub scope-filtering — https://raw.githubusercontent.com/github/github-mcp-server/main/docs/scope-filtering.md
[48] GitHub README pagination page/perPage — https://raw.githubusercontent.com/github/github-mcp-server/main/README.md
[49] Gemini function-calling enums — https://ai.google.dev/gemini-api/docs/function-calling
[50] bex.co progressive discovery 100-tool (secondary) — https://bex.co/blog/2026/09/23/mcp-progressive-discovery-100-tool-servers (2026-09-23)
[51] SEP-2127 MCP Server Cards — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/seps/2127-mcp-server-cards.md
[52] MCP Registry FAQ immutability — https://raw.githubusercontent.com/modelcontextprotocol/modelcontextprotocol/main/docs/registry/faq.mdx

