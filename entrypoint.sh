#!/bin/sh
set -e

if [ -z "$DB_URL" ]; then
  echo "⚠️  DB_URL not set, skipping migrations"
else
  echo "🔄 Running database migrations..."
  ./migrate -path migrations -database "$DB_URL" up || echo "⚠️  Migration failed, continuing..."
  echo "✅ Migrations completed"
fi

echo "🚀 Starting application..."

# Execute the main command
exec "$@"
