#!/bin/sh
set -e

BIN_NAME=rustdesk-api-server-pro
BIN_PATH="/app/${BIN_NAME}"
BIN_RECORD="/app/.bin.${BIN_NAME}"

if [ ! -x "$BIN_PATH" ]; then
  printf 'Binary not found at %s\n' "$BIN_PATH" >&2
  exit 1
fi

mkdir -p /app/data
cd /app/data

if [ ! -f "$BIN_RECORD" ] && [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
  "$BIN_PATH" user add "$ADMIN_USER" "$ADMIN_PASS" --admin || true
  touch "$BIN_RECORD"
fi

"$BIN_PATH" sync
exec "$BIN_PATH" start
