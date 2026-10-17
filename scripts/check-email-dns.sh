#!/bin/bash
# check-email-dns.sh — verify SPF/DKIM/DMARC/MX for a domain (default: mcp-is.ru)
set -euo pipefail

DOMAIN="${1:-mcp-is.ru}"
RESOLVER="${RESOLVER:-8.8.8.8}"
FAIL=0

echo "=== DNS email check: ${DOMAIN} (resolver ${RESOLVER}) ==="

txt() {
    dig +short TXT "$1" @"${RESOLVER}" 2>/dev/null | tr -d '"'
}

echo ""
echo "--- MX ---"
MX=$(dig +short MX "${DOMAIN}" @"${RESOLVER}" 2>/dev/null)
if [ -z "${MX}" ]; then
    echo "MISSING MX (ok for variant A: send-only via external SMTP)"
else
    echo "${MX}"
fi

echo ""
echo "--- SPF (TXT @) ---"
SPF=$(txt "${DOMAIN}" | grep -i 'v=spf1' || true)
if [ -z "${SPF}" ]; then
    echo "MISSING SPF — add TXT @  v=spf1 include:_spf.mail.ru ~all"
    FAIL=1
else
    echo "OK ${SPF}"
fi

echo ""
echo "--- DMARC (TXT _dmarc.${DOMAIN}) ---"
DMARC=$(txt "_dmarc.${DOMAIN}" | grep -i 'v=DMARC1' || true)
if [ -z "${DMARC}" ]; then
    echo "MISSING DMARC — add TXT _dmarc  v=DMARC1; p=none; rua=mailto:dmarc@${DOMAIN}"
    FAIL=1
else
    echo "OK ${DMARC}"
fi

echo ""
echo "--- DKIM (common selectors) ---"
DKIM_FOUND=0
for s in mail default selector1 selector2 google dkim m1 m2; do
    VAL=$(txt "${s}._domainkey.${DOMAIN}" | grep -i 'v=DKIM1' || true)
    if [ -n "${VAL}" ]; then
        echo "OK ${s}._domainkey: ${VAL:0:80}..."
        DKIM_FOUND=1
    fi
done
if [ "${DKIM_FOUND}" -eq 0 ]; then
    echo "MISSING DKIM — add <selector>._domainkey from mail provider panel"
    FAIL=1
fi

echo ""
echo "--- A ---"
dig +short A "${DOMAIN}" @"${RESOLVER}"

echo ""
if [ "${FAIL}" -eq 0 ]; then
    echo "RESULT: OK — DNS email records present"
else
    echo "RESULT: INCOMPLETE — see docs/setup/EMAIL-DNS.md"
fi
exit "${FAIL}"
