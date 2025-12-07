#!/bin/bash

echo "Starting NoBurn background worker..."

# Ensure Redis is running
if ! redis-cli ping > /dev/null 2>&1; then
    echo "Redis is not running. Please start Redis first."
    exit 1
fi

# Start the worker
go run cmd/worker/main.go