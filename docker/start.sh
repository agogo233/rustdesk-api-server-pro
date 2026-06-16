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

# 以 root 运行，修复挂载卷权限，确保 appuser 可写
mkdir -p /app/data /var/log
chown -R appuser:appgroup /app/data /var/log

cd /app

if [ ! -f "$BIN_RECORD" ] && [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
  su-exec appuser:appgroup "$BIN_PATH" user add "$ADMIN_USER" "$ADMIN_PASS" --admin 2>"$LOG" && touch "$BIN_RECORD" || true
fi

su-exec appuser:appgroup "$BIN_PATH" sync 2>"$LOG" || true
exec su-exec appuser:appgroup "$BIN_PATH" start >"$LOG" 2>&1
