#!/bin/sh
set -e

BIN_NAME=rustdesk-api-server-pro
BIN_PATH="/app/${BIN_NAME}"
BIN_RECORD="/app/.bin.${BIN_NAME}"
LOG='/var/log/rustdesk-api-server-pro.log'

if [ ! -x "$BIN_PATH" ]; then
  printf 'Binary not found at %s\n' "$BIN_PATH" >&2
  exit 1
fi

mkdir -p /app/data /var/log
chown "$(id -u):$(id -g)" /app/data /var/log 2>/dev/null || true

cd /app/data

if [ ! -f "$BIN_RECORD" ] && [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
  "$BIN_PATH" user add "$ADMIN_USER" "$ADMIN_PASS" --admin 2>"$LOG" && touch "$BIN_RECORD" || true
fi

"$BIN_PATH" sync 2>"$LOG" || true
exec "$BIN_PATH" start >"$LOG" 2>&1
