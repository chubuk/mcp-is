# Brief: Стандарты MCP-серверов vs mcp-is

**Date:** 2026-10-20
**Depth:** standard (5 angles, 1 follow-up, 15+ sources)
**Workspace:** /home/chubuk/projects/mcp-is/research/mcp-standards-2026/

## Refined question
Какие отраслевые стандарты и best practices сложились для MCP-серверов (протокол, дизайн tools, результаты для LLM, безопасность, транспорт, экосистема) к 2026 году — и насколько им соответствует проект mcp-is (Go + mark3labs/mcp-go v0.56, Streamable HTTP, ~245 tools) и текущий план доработок v2?

## Scope
In: MCP spec revisions; official SDKs; tool design (naming, descriptions, schemas, annotations, StructuredContent); results/LLM ergonomics; auth/security/transport; production MCP servers at scale; tool-spawl critiques; critical mapping to mcp-is + plan v2.
Out: InSales API internals; implementation (plan mode); unrelated agent frameworks.

## Assumptions
- Stack: Go, mark3labs/mcp-go v0.56.0, Streamable HTTP, public repo chubuk/mcp-is, v0.5.3, ~245 tools (108 RO / 137 write)
- User cares: LLM discoverability/understanding; compliance with world standards; whether plan v2 priorities are right
- Primary sources preferred; report in Russian

## Angles
F1 — MCP protocol & official spec norms (2025–2026)
F2 — Tool authoring best practices & LLM ergonomics
F3 — Production MCP servers at scale (vendors)
F4 — Security, auth, multi-tenancy, transport
F5 — Ecosystem maturity & anti-patterns (tool sprawl, registry, context cost)

## Deliverables
REPORT.md + plan deltas + Russian summary for user
