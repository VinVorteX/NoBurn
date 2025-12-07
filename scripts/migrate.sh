#!/bin/bash

set -e

DB_URL="${DB_URL:-postgres://noburn_user:noburn_pass@localhost:5432/noburn_db?sslmode=disable}"

case "$1" in
  up)
    echo "Running migrations..."
    migrate -path migrations -database "$DB_URL" up
    ;;
  down)
    echo "Rolling back migrations..."
    migrate -path migrations -database "$DB_URL" down
    ;;
  create)
    if [ -z "$2" ]; then
      echo "Usage: $0 create <migration_name>"
      exit 1
    fi
    echo "Creating migration: $2"
    migrate create -ext sql -dir migrations -seq "$2"
    ;;
  version)
    migrate -path migrations -database "$DB_URL" version
    ;;
  *)
    echo "Usage: $0 {up|down|create|version}"
    exit 1
    ;;
esac