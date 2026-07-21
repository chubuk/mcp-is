# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.2.x   | ✅ Yes |
| < 0.2   | ❌ No |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly.

**DO NOT** open a public GitHub issue for security vulnerabilities.

### How to Report

Send an email to **chubuk@mail.ru** with:

1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if any)

### What to Include

- Affected component (auth, API, MCP, etc.)
- Attack vector
- Proof of concept (if possible)
- Your assessment of severity

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 1 week
- **Fix release**: Depends on severity

### What to Expect

- We will acknowledge receipt of your report
- We will investigate and confirm the vulnerability
- We will develop and test a fix
- We will release the fix and credit you (unless you prefer anonymity)

## Security Best Practices

### For Deployment

- Change default JWT_SECRET and ENCRYPTION_KEY
- Use HTTPS in production (Let's Encrypt or self-signed)
- Set `MCP_READ_ONLY=true` if write operations are not needed
- Configure `ALLOWED_ORIGINS` instead of using `*`
- Use strong passwords (8+ chars, uppercase, digit)

### For Development

- Never commit `.env` files
- Use environment variables for secrets
- Keep dependencies updated
- Run `govulncheck ./...` periodically

## Known Security Features

- **Rate limiting**: 10 req/min on login/register/forgot-password
- **Password hashing**: bcrypt with default cost
- **JWT tokens**: HS256 with configurable expiry
- **CORS**: Configurable allowed origins
- **Body size limit**: 1MB on MCP endpoint
- **SQL injection prevention**: Parameterized queries throughout
- **Password policy**: 8+ chars, 1 uppercase, 1 digit
