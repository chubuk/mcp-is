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

# Run linter (if installed)
golangci-lint run
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

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated and passing
- [ ] Documentation updated (if applicable)
- [ ] No breaking changes (or clearly documented)
- [ ] Commit messages follow conventions

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
```

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (BSL 1.1).
