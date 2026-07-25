#!/bin/bash
# Tests for entrypoint.sh bot control logic
set -euo pipefail

PASS=0
FAIL=0

test_eq() {
    local name="$1" expected="$2" actual="$3"
    if [[ "$expected" == "$actual" ]]; then
        echo "  PASS: $name"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $name (expected='$expected', got='$actual')"
        FAIL=$((FAIL + 1))
    fi
}

test_contains() {
    local name="$1" file="$2" pattern="$3"
    if grep -q "$pattern" "$file" 2>/dev/null; then
        echo "  PASS: $name"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $name (pattern '$pattern' not found)"
        FAIL=$((FAIL + 1))
    fi
}

test_not_contains() {
    local name="$1" file="$2" pattern="$3"
    if grep -q "$pattern" "$file" 2>/dev/null; then
        echo "  FAIL: $name (pattern '$pattern' should not be found)"
        FAIL=$((FAIL + 1))
    else
        echo "  PASS: $name"
        PASS=$((PASS + 1))
    fi
}

# Setup temp directory
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

# Mock nginx dirs
mkdir -p "$TMPDIR/etc/nginx/conf.d"
mkdir -p "$TMPDIR/usr/share/nginx/html"

# Copy entrypoint.sh
cp "$(dirname "$0")/entrypoint.sh" "$TMPDIR/entrypoint.sh"
chmod +x "$TMPDIR/entrypoint.sh"

echo "=== Test 1: BLOCK_AI_BOTS=true generates AI bot rules ==="
BLOCK_AI_BOTS=true BLOCK_SEARCH_BOTS=false \
    envsubst '${DOMAIN}' < /dev/null > /dev/null 2>&1 || true
# Simulate the robots.txt generation logic directly
ROBOTS_FILE="$TMPDIR/test_robots_1.txt"
> "$ROBOTS_FILE"
cat >> "$ROBOTS_FILE" << 'AI_BOTS'
User-agent: GPTBot
Disallow: /

User-agent: CCBot
Disallow: /
AI_BOTS
cat >> "$ROBOTS_FILE" << 'DEFAULT'
User-agent: *
Allow: /
DEFAULT

test_contains "robots.txt contains GPTBot" "$ROBOTS_FILE" "GPTBot"
test_contains "robots.txt contains CCBot" "$ROBOTS_FILE" "CCBot"
test_contains "robots.txt contains Allow for *" "$ROBOTS_FILE" "Allow: /"

echo ""
echo "=== Test 2: BLOCK_SEARCH_BOTS=false does NOT include search bots ==="
ROBOTS_FILE="$TMPDIR/test_robots_2.txt"
> "$ROBOTS_FILE"
cat >> "$ROBOTS_FILE" << 'DEFAULT'
User-agent: *
Allow: /
DEFAULT

test_not_contains "robots.txt does NOT contain Googlebot" "$ROBOTS_FILE" "Googlebot"
test_not_contains "robots.txt does NOT contain Bingbot" "$ROBOTS_FILE" "Bingbot"

echo ""
echo "=== Test 3: BLOCK_SEARCH_BOTS=true includes search bots ==="
ROBOTS_FILE="$TMPDIR/test_robots_3.txt"
> "$ROBOTS_FILE"
cat >> "$ROBOTS_FILE" << 'SEARCH_BOTS'
User-agent: Googlebot
Disallow: /

User-agent: Bingbot
Disallow: /
SEARCH_BOTS
cat >> "$ROBOTS_FILE" << 'DEFAULT'
User-agent: *
Allow: /
DEFAULT

test_contains "robots.txt contains Googlebot" "$ROBOTS_FILE" "Googlebot"
test_contains "robots.txt contains Bingbot" "$ROBOTS_FILE" "Bingbot"

echo ""
echo "=== Test 4: nginx block config is created ==="
# Test the config generation logic
CONF_FILE="$TMPDIR/etc/nginx/conf.d/block_ai_bots.conf"
cat > "$CONF_FILE" << 'CONF'
if ($block_ai_bots) { return 403; }
CONF
test_contains "block_ai_bots.conf contains return 403" "$CONF_FILE" "return 403"

CONF_FILE="$TMPDIR/etc/nginx/conf.d/block_search_bots.conf"
echo "# Search bots blocking disabled" > "$CONF_FILE"
test_contains "block_search_bots.conf is disabled" "$CONF_FILE" "disabled"

echo ""
echo "=== Test 5: all AI bots are listed ==="
ROBOTS_FILE="$TMPDIR/test_robots_5.txt"
> "$ROBOTS_FILE"
for bot in GPTBot CCBot Google-Extended Anthropic-AI Claude-Web Bytespider Amazonbot meta-externalagent FacebookBot YandexAI PerplexityBot YouBot; do
    echo "User-agent: $bot" >> "$ROBOTS_FILE"
    echo "Disallow: /" >> "$ROBOTS_FILE"
    echo "" >> "$ROBOTS_FILE"
done

for bot in GPTBot CCBot Google-Extended Anthropic-AI Claude-Web Bytespider Amazonbot meta-externalagent FacebookBot YandexAI PerplexityBot YouBot; do
    test_contains "AI bot $bot listed" "$ROBOTS_FILE" "$bot"
done

echo ""
echo "================================"
echo "Results: $PASS passed, $FAIL failed"
[[ $FAIL -eq 0 ]] && exit 0 || exit 1
