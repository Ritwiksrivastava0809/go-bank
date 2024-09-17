#!/bin/sh

set -e

# Start the first process
echo "run db migration"
/app/pkg/db/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Start the second process
echo "start app"
exec "$@"
