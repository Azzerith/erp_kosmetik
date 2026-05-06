#!/bin/bash

# Database migration script

set -e

echo "Running database migrations..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run migrations using Go binary
go run cmd/api/main.go migrate

echo "Migrations completed successfully!"