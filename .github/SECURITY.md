# Security Policy

## Supported Versions

| Version | Supported |
| ------- | --------- |
| 0.7.x   | ✅ Yes |
| < 0.7   | ❌ No |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly.

**DO NOT** open a public GitHub issue for security vulnerabilities.

### How to Report

Send an email to **chubuk@mail.ru** with:

1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if any)

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 1 week
- **Fix release**: Depends on severity

### What to Expect

- We will acknowledge receipt of your report
- We will investigate and confirm the vulnerability
- We will develop and test a fix
- We will release the fix and credit you (unless you prefer anonymity)

## Security Features

- **Authentication**: JWT tokens with configurable expiry
- **Authorization**: Role-based access (user, admin, super_admin)
- **Password hashing**: bcrypt with default cost
- **Rate limiting**: 10 req/min on auth endpoints
- **CORS**: Configurable allowed origins
- **CSP**: Content Security Policy headers
- **Cache protection**: `Cache-Control: no-store` on HTML/API (client also checks JWT `exp` and re-validates on bfcache restore)
- **Input validation**: All endpoints validated
- **SQL injection**: Parameterized queries throughout
- **PDN protection**: Auto-protection for Russian shops + foreign LLMs

## Best Practices for Deployment

- Change default JWT_SECRET and ENCRYPTION_KEY
- Use HTTPS (Let's Encrypt or self-signed)
- Set `MCP_READ_ONLY=true` if write operations not needed
- Configure `ALLOWED_ORIGINS` instead of using `*`
- Use strong passwords (8+ chars, uppercase, digit)
- Keep dependencies updated (`govulncheck ./...`)
- Do not store server SSH credentials in GitHub Secrets (deploy is manual)
