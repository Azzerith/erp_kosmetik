#!/bin/bash

# Database seeder script

set -e

echo "Running database seeder..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run seeder using Go binary
go run cmd/api/main.go seed

echo "Seeder completed successfully!"