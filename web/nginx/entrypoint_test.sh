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

# Run entrypoint.sh in the image, dump the three files it generates
run_entrypoint() {
    docker run --rm --entrypoint sh -v "$ENTRYPOINT_SRC:/entrypoint.sh:ro" "$@" "$IMG" -c '
        sh /entrypoint.sh true >/dev/null 2>&1
        echo "===AI_CONF==="
        cat /etc/nginx/conf.d/block_ai_bots.conf
        echo "===SEARCH_CONF==="
        cat /etc/nginx/conf.d/block_search_bots.conf
        echo "===ROBOTS==="
        cat /etc/nginx/robots.txt'
}

echo "=== Scenario 1: BLOCK_AI_BOTS=true BLOCK_SEARCH_BOTS=false ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=true -e BLOCK_SEARCH_BOTS=false)
test_contains "AI conf enabled" "$OUT" '$block_ai_bots'
test_contains "search conf disabled" "$OUT" 'Search bots blocking disabled'
test_contains "robots.txt contains GPTBot" "$OUT" "GPTBot"
test_contains "robots.txt contains CCBot" "$OUT" "CCBot"
test_contains "robots.txt contains FacebookBot" "$OUT" "FacebookBot"
test_contains "robots.txt contains YandexAI" "$OUT" "YandexAI"
test_contains "robots.txt contains PerplexityBot" "$OUT" "PerplexityBot"
test_contains "robots.txt contains YouBot" "$OUT" "YouBot"
test_not_contains "robots.txt does NOT contain Googlebot" "$OUT" "Googlebot"
test_not_contains "robots.txt does NOT contain Bingbot" "$OUT" "Bingbot"
test_contains "robots.txt contains Allow for *" "$OUT" "Allow: /"
test_contains "robots.txt contains User-agent: *" "$OUT" "User-agent: *"

echo ""
echo "=== Scenario 2: BLOCK_AI_BOTS=false BLOCK_SEARCH_BOTS=true ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=false -e BLOCK_SEARCH_BOTS=true)
test_contains "AI conf disabled" "$OUT" 'AI bots blocking disabled'
test_contains "search conf enabled" "$OUT" '$block_search_bots'
test_not_contains "robots.txt does NOT contain GPTBot" "$OUT" "GPTBot"
test_contains "robots.txt contains Googlebot" "$OUT" "Googlebot"
test_contains "robots.txt contains Bingbot" "$OUT" "Bingbot"
test_contains "robots.txt contains YandexBot" "$OUT" "YandexBot"
test_contains "robots.txt contains Baiduspider" "$OUT" "Baiduspider"
test_contains "robots.txt contains DuckDuckBot" "$OUT" "DuckDuckBot"
test_contains "robots.txt contains Applebot" "$OUT" "Applebot"

echo ""
echo "=== Scenario 3: both enabled ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=true -e BLOCK_SEARCH_BOTS=true)
test_contains "AI conf enabled" "$OUT" '$block_ai_bots'
test_contains "search conf enabled" "$OUT" '$block_search_bots'
test_contains "robots.txt contains GPTBot" "$OUT" "GPTBot"
test_contains "robots.txt contains Googlebot" "$OUT" "Googlebot"

echo ""
echo "=== Scenario 4: both disabled ==="
OUT=$(run_entrypoint -e BLOCK_AI_BOTS=false -e BLOCK_SEARCH_BOTS=false)
test_contains "AI conf disabled" "$OUT" 'AI bots blocking disabled'
test_contains "search conf disabled" "$OUT" 'Search bots blocking disabled'
test_not_contains "robots.txt does NOT contain GPTBot" "$OUT" "GPTBot"
test_not_contains "robots.txt does NOT contain Googlebot" "$OUT" "Googlebot"
test_contains "robots.txt contains User-agent: *" "$OUT" "User-agent: *"

echo ""
echo "==============================="
echo "PASS: $PASS  FAIL: $FAIL"
if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
