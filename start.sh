#!/bin/sh

set -e

echo "Run migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start app"
exec "$@"