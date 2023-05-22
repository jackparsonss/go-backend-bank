#!/bin/sh

# exit if any command fails
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"
