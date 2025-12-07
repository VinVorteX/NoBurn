#!/bin/sh
set -e

echo "🔄 Running database migrations..."
./migrate -path migrations -database "$DB_URL" up

echo "✅ Migrations completed"
echo "🚀 Starting application..."

# Execute the main command
exec "$@"
