#!/bin/bash
# test-smtp.sh — live SMTP test using env vars or .env file
# Usage:
#   export SMTP_HOST=smtp.mail.ru SMTP_PORT=587 SMTP_USER=... SMTP_PASS=... SMTP_FROM='MCP IS <user@mail.ru>'
#   ./scripts/test-smtp.sh recipient@example.com
set -euo pipefail

ENV_FILE="${ENV_FILE:-/opt/mcp-is/.env}"
if [ -f "${ENV_FILE}" ]; then
    # shellcheck disable=SC1090
    set -a; . "${ENV_FILE}"; set +a
fi

TO="${1:-${SMTP_USER}}"
HOST="${SMTP_HOST:-}"
PORT="${SMTP_PORT:-587}"
USER="${SMTP_USER:-}"
PASS="${SMTP_PASS:-}"
FROM="${SMTP_FROM:-$USER}"

if [ -z "${HOST}" ] || [ -z "${USER}" ]; then
    echo "ERROR: set SMTP_HOST and SMTP_USER (env or ${ENV_FILE})"
    exit 1
fi
if [ -z "${TO}" ]; then
    echo "Usage: $0 recipient@example.com"
    exit 1
fi

SUBJECT="SMTP Test MCP IS $(date -u +%Y-%m-%dT%H:%M:%SZ)"
BODY="Тестовое письмо SMTP от MCP IS.
Host: ${HOST}:${PORT}
From: ${FROM}
To: ${TO}
Time: $(date -R)

Если вы видите это письмо — SMTP работает.
"

echo "Sending via ${HOST}:${PORT} as ${USER} → ${TO} ..."

if command -v swaks >/dev/null 2>&1; then
    swaks --to "${TO}" --from "${FROM}" --server "${HOST}" --port "${PORT}" \
        --auth LOGIN --auth-user "${USER}" --auth-password "${PASS}" \
        --tls --header "Subject: ${SUBJECT}" --body "${BODY}"
else
    python3 - "$HOST" "$PORT" "$USER" "$PASS" "$FROM" "$TO" "$SUBJECT" "$BODY" <<'PY'
import smtplib, sys
from email.message import EmailMessage
host, port, user, pwd, fr, to, subj, body = sys.argv[1:9]
msg = EmailMessage()
msg["From"] = fr
msg["To"] = to
msg["Subject"] = subj
msg.set_content(body)
# Prefer SSL on 465, STARTTLS otherwise
if int(port) == 465:
    s = smtplib.SMTP_SSL(host, int(port), timeout=20)
else:
    s = smtplib.SMTP(host, int(port), timeout=20)
    s.starttls()
s.login(user, pwd)
s.send_message(msg)
s.quit()
print("OK: message accepted by", host)
PY
fi

echo ""
echo "Done. Check mailbox (and spam) at ${TO}."
