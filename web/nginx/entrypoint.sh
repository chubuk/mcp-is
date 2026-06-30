#!/bin/sh
set -e

DOMAIN=${DOMAIN:-localhost}
EMAIL=${LETSENCRYPT_EMAIL:-admin@${DOMAIN}}

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

# Copy nginx.conf and substitute env vars
cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.tmpl
envsubst '${DOMAIN}' < /etc/nginx/nginx.conf.tmpl > /etc/nginx/nginx.conf.final
mv /etc/nginx/nginx.conf.final /etc/nginx/nginx.conf

exec "$@"
