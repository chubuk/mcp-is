#!/bin/bash
# Real-execution tests for entrypoint.sh: runs it inside the nginx image and
# asserts on the files it actually generates. (The old version wrote heredocs
# mirroring the expected output and grepped those — testing nothing.)
set -euo pipefail

IMG=${NGINX_TEST_IMAGE:-ghcr.io/chubuk/mcp-is-nginx:latest}
ENTRYPOINT_SRC="$(cd "$(dirname "$0")" && pwd)/entrypoint.sh"

PASS=0
FAIL=0

test_contains() {
    local name="$1" haystack="$2" pattern="$3"
    if grep -qF -- "$pattern" <<< "$haystack" 2>/dev/null; then
        echo "  PASS: $name"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $name (pattern '$pattern' not found)"
        FAIL=$((FAIL + 1))
    fi
}

test_not_contains() {
    local name="$1" haystack="$2" pattern="$3"
    if grep -qF -- "$pattern" <<< "$haystack" 2>/dev/null; then
        echo "  FAIL: $name (pattern '$pattern' should not be found)"
        FAIL=$((FAIL + 1))
    else
        echo "  PASS: $name"
        PASS=$((PASS + 1))
    fi
}

# Run entrypoint.sh in the image, dump the confs it generates
# (robots.txt теперь генерирует приложение, nginx проксирует — entrypoint его не пишет)
run_entrypoint() {
    docker run --rm --entrypoint sh -v "$ENTRYPOINT_SRC:/entrypoint.sh:ro" "$@" "$IMG" -c '
        sh /entrypoint.sh true >/dev/null 2>&1
        echo "===AI_CONF==="
        cat /etc/nginx/conf.d/block_ai_bots.conf
        echo "===SEARCH_CONF==="
        cat /etc/nginx/conf.d/block_search_bots.conf'
}

echo "=== Scenario 1: BLOCK_AI_BOTS=true BLOCK_SEARCH_BOTS=false ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=true -e BLOCK_SEARCH_BOTS=false)
test_contains "AI conf enabled" "$OUT" '$block_ai_bots'
test_contains "search conf disabled" "$OUT" 'Search bots blocking disabled'

echo ""
echo "=== Scenario 2: BLOCK_AI_BOTS=false BLOCK_SEARCH_BOTS=true ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=false -e BLOCK_SEARCH_BOTS=true)
test_contains "AI conf disabled" "$OUT" 'AI bots blocking disabled'
test_contains "search conf enabled" "$OUT" '$block_search_bots'

echo ""
echo "=== Scenario 3: both enabled ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=true -e BLOCK_SEARCH_BOTS=true)
test_contains "AI conf enabled" "$OUT" '$block_ai_bots'
test_contains "search conf enabled" "$OUT" '$block_search_bots'

echo ""
echo "=== Scenario 4: both disabled ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=false -e BLOCK_SEARCH_BOTS=false)
test_contains "AI conf disabled" "$OUT" 'AI bots blocking disabled'
test_contains "search conf disabled" "$OUT" 'Search bots blocking disabled'

echo ""
echo "==============================="
echo "PASS: $PASS  FAIL: $FAIL"
if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
