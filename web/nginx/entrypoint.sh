#!/bin/sh
set -e

DOMAIN=${DOMAIN:-localhost}
EMAIL=${LETSENCRYPT_EMAIL:-admin@${DOMAIN}}
BLOCK_AI_BOTS=${BLOCK_AI_BOTS:-true}
BLOCK_SEARCH_BOTS=${BLOCK_SEARCH_BOTS:-false}

# Create certbot webroot
mkdir -p /var/www/certbot

# If no SSL certs exist, generate self-signed for development
if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo "No SSL certificates found. Generating self-signed cert for development..."
    mkdir -p /etc/letsencrypt/live/${DOMAIN}
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /etc/letsencrypt/live/${DOMAIN}/privkey.pem \
        -out /etc/letsencrypt/live/${DOMAIN}/fullchain.pem \
        -subj "/CN=${DOMAIN}"
    echo "Self-signed certificate generated for ${DOMAIN}"
fi

# Generate AI bots blocking config
if [ "$BLOCK_AI_BOTS" = "true" ]; then
    cat > /etc/nginx/conf.d/block_ai_bots.conf << 'CONF'
if ($block_ai_bots) { return 403; }
CONF
    echo "AI bots blocking: ENABLED"
else
    echo "# AI bots blocking disabled" > /etc/nginx/conf.d/block_ai_bots.conf
    echo "AI bots blocking: DISABLED"
fi

# Generate search bots blocking config
if [ "$BLOCK_SEARCH_BOTS" = "true" ]; then
    cat > /etc/nginx/conf.d/block_search_bots.conf << 'CONF'
if ($block_search_bots) { return 403; }
CONF
    echo "Search bots blocking: ENABLED"
else
    echo "# Search bots blocking disabled" > /etc/nginx/conf.d/block_search_bots.conf
    echo "Search bots blocking: DISABLED"
fi

# Generate robots.txt (writable location, not in :ro volume)
ROBOTS_FILE="/etc/nginx/robots.txt"
> "$ROBOTS_FILE"

if [ "$BLOCK_AI_BOTS" = "true" ]; then
    cat >> "$ROBOTS_FILE" << 'AI_BOTS'
# AI crawlers - blocked
User-agent: GPTBot
Disallow: /

User-agent: CCBot
Disallow: /

User-agent: Google-Extended
Disallow: /

User-agent: Anthropic-AI
Disallow: /

User-agent: Claude-Web
Disallow: /

User-agent: Bytespider
Disallow: /

User-agent: Amazonbot
Disallow: /

User-agent: meta-externalagent
Disallow: /

User-agent: FacebookBot
Disallow: /

User-agent: YandexAI
Disallow: /

User-agent: PerplexityBot
Disallow: /

User-agent: YouBot
Disallow: /

AI_BOTS
fi

if [ "$BLOCK_SEARCH_BOTS" = "true" ]; then
    cat >> "$ROBOTS_FILE" << 'SEARCH_BOTS'
# Search engine crawlers - blocked
User-agent: Googlebot
Disallow: /

User-agent: Bingbot
Disallow: /

User-agent: YandexBot
Disallow: /

User-agent: Baiduspider
Disallow: /

User-agent: DuckDuckBot
Disallow: /

User-agent: Applebot
Disallow: /

SEARCH_BOTS
fi

cat >> "$ROBOTS_FILE" << 'DEFAULT'
# All other bots
User-agent: *
Allow: /
DEFAULT

echo "robots.txt generated"

# Copy nginx.conf and substitute env vars
cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.tmpl
envsubst '${DOMAIN}' < /etc/nginx/nginx.conf.tmpl > /etc/nginx/nginx.conf.final
mv /etc/nginx/nginx.conf.final /etc/nginx/nginx.conf

# Watch for domain changes and reload nginx
(
    while true; do
        sleep 60
        CURRENT_DOMAIN=$(cat /etc/nginx/domain.conf 2>/dev/null || echo "${DOMAIN}")
        if [ "$CURRENT_DOMAIN" != "${DOMAIN}" ]; then
            DOMAIN="$CURRENT_DOMAIN"
            envsubst '${DOMAIN}' < /etc/nginx/nginx.conf.tmpl > /etc/nginx/nginx.conf.final
            mv /etc/nginx/nginx.conf.final /etc/nginx/nginx.conf
            nginx -s reload 2>/dev/null || true
            echo "Domain updated to ${DOMAIN}, nginx reloaded"
        fi
    done
) &

exec "$@"
