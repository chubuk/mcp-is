# Contributing to MCP IS

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project follows our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Bugs

1. Check [existing issues](../../issues) to avoid duplicates
2. Open a new issue using the **Bug Report** template
3. Include:
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment (OS, Docker version, Go version)

### Suggesting Features

1. Open a new issue using the **Feature Request** template
2. Describe the problem you're trying to solve
3. Suggest a solution if you have one

### Contributing Code

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Write or update tests
5. Ensure all tests pass
6. Submit a Pull Request

## Development Setup

### Prerequisites

- Go 1.21+
- Docker and Docker Compose v2
- PostgreSQL 15+ (or use Docker)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/your-username/mcp-is.git
cd mcp-is

# Start development environment
docker compose up -d

# Run tests
go test ./...

# Run linter (required; CI runs the same golangci-lint version and config)
make lint

# Auto-fix formatting issues (gofmt)
make lint-fix

# UI regression smoke (local docker on :8443; sticky checkboxes + Quill remount)
./scripts/ui-smoke.sh
```

### Project Structure

```
├── cmd/server/       # Entry point
├── internal/         # Application code
│   ├── auth/         # JWT, middleware, rate limiter
│   ├── handlers/     # HTTP handlers
│   ├── mcp/          # MCP server and tools
│   ├── storage/      # Database layer
│   └── ...
├── docs/             # Documentation
├── packaging/        # Debian packaging
└── web/              # Frontend (nginx, static files)
```

## Code Style

### Go

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` and `goimports`
- Run `make lint` before committing — CI enforces the same golangci-lint config
- Names: camelCase for variables, PascalCase for exported
- Comments: exported functions must have doc comments
- Errors: use `fmt.Errorf("context: %w", err)` for wrapping

### Commits

- Use [Conventional Commits](https://www.conventionalcommits.org/)
- Format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Examples:
  - `feat(auth): add JWT token refresh`
  - `fix(api): handle nil response from upstream`
  - `docs(readme): update installation instructions`

## Pull Request Process

1. **Title**: Use conventional commit format (`feat: add new feature`)
2. **Description**: Explain what and why (not how)
3. **Linked issues**: Reference related issues (`Closes #123`)
4. **Tests**: Include tests for new functionality
5. **Documentation**: Update docs if behavior changes
6. **Breaking changes**: Clearly mark in description

### PR Checklist (Definition of Done)

A change is **done** only when all applicable checks are green:

- [ ] `make lint` (golangci-lint: govet, staticcheck, errcheck, ineffassign, gofmt)
- [ ] `go vet ./...`
- [ ] `go test ./...`
- [ ] `scripts/check-schema.sh` — **required** if Go structs or SQL/migrations changed
- [ ] `./scripts/ui-smoke.sh` — **required** if UI changed (local docker on `https://127.0.0.1:8443`)
- [ ] Code follows project style guidelines
- [ ] Tests added/updated and passing
- [ ] Documentation updated (if behavior changes)
- [ ] No breaking changes (or clearly documented)
- [ ] Commit messages follow Conventional Commits

### Review Process

1. Maintainers will review within 48 hours
2. Address review comments
3. Once approved, maintainer will merge

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./internal/handlers/ -v -run TestLogin

# Run with coverage
go test -cover ./...

# Integration tests (needs PostgreSQL; never point at production data)
TEST_DB_URL='postgres://admin:password@127.0.0.1:55432/mcp_is?sslmode=disable' \
  go test -tags integration ./internal/storage/... -count=1
```

### UI

Frontend is a vanilla JS SPA (no bundler). Prefer shared helpers — do not invent a new form/pager per page:

| Helper | Where | Use for |
|--------|--------|---------|
| `App.readFields` / `App.check` | `web/static/js/app.js` | reading form fields (`bool` / `text` / `trim` / `int`) |
| `renderPagination` / `showPage` / `getPerPage` | `web/static/js/pagination.js` | lists and page navigation |
| `.form-check-hint` | `web/static/css/style.css` | checkbox/radio caption hint under the control |

UI regression smoke (login, sticky checkboxes, Quill remount after SPA nav):

```bash
./scripts/ui-smoke.sh   # default https://127.0.0.1:8443
```

## Go standards (summary)

Full house rules live in `.project_rules.md` / `.github` workflow gates. In short:

- Errors: wrap with `fmt.Errorf("…: %w", err)`; never discard non-trivial errors
- Logging: `log/slog`; user-facing API/UI errors in Russian
- Tests: table-driven where it pays off; `got` / `want` in failures
- No `panic` for control flow; `context.Context` as first parameter; no magic reflection kits
- Deps: stdlib first; add a module only if it clearly pays for itself
- SQL: parameterized queries only; migrations versioned and reversible

## What does **not** belong in public docs

Agent/workspace tooling (internal skills, SSH hosts, deploy tokens, private paths) stays in local workspace config and project rules — not in README/CONTRIBUTING. Publish only what every contributor needs (lint, tests, smoke, style).

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (BSL 1.1). Commercial use of the project requires a separate agreement with the Licensor (see `LICENSE` and `README.md`).
