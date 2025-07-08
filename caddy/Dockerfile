# Caddy 2 for WebEnable CMS
FROM caddy:2-alpine

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD caddy validate --config /etc/caddy/Caddyfile || exit 1

EXPOSE 80 443
