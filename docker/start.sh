#!/bin/sh
set -e

BIN_NAME=rustdesk-api-server-pro
BIN_PATH="/app/${BIN_NAME}"
LOG='/var/log/rustdesk-api-server-pro.log'

if [ ! -x "$BIN_PATH" ]; then
  printf 'Binary not found at %s\n' "$BIN_PATH" >&2
  exit 1
fi

mkdir -p /app/data /var/log
chown -R appuser:appgroup /app/data /var/log

cd /app

su-exec appuser:appgroup "$BIN_PATH" sync 2>"$LOG" || true

if [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
  su-exec appuser:appgroup "$BIN_PATH" user add "$ADMIN_USER" "$ADMIN_PASS" --admin 2>"$LOG" || true
fi

exec su-exec appuser:appgroup "$BIN_PATH" start >"$LOG" 2>&1
