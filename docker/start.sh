#!/bin/sh
set -e

if [ ! -f /usr/local/bin/rustdesk-api-server-pro ]; then
    ln -s /app/rustdesk-api-server-pro /usr/local/bin/rustdesk-api-server-pro
fi

mkdir -p /app/data

cd /app/data

if [ ! -f /app/data/.init.lock ] && [ -n "$ADMIN_USER" ] && [ -n "$ADMIN_PASS" ]; then
    /app/rustdesk-api-server-pro user add $ADMIN_USER $ADMIN_PASS --admin
    touch /app/data/.init.lock
fi

/app/rustdesk-api-server-pro sync
/app/rustdesk-api-server-pro start
