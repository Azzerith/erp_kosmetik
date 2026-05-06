#!/bin/bash

# Deployment script

set -e

echo "Starting deployment..."

# Build the application
echo "Building application..."
go build -o bin/api cmd/api/main.go

# Run database migrations
echo "Running migrations..."
./bin/api migrate

# Run seeder (optional, comment if not needed)
# echo "Running seeder..."
# ./bin/api seed

# Restart service (adjust based on your process manager)
echo "Restarting service..."
if command -v systemctl &> /dev/null; then
    sudo systemctl restart erp-backend
elif command -v supervisorctl &> /dev/null; then
    sudo supervisorctl restart erp-backend
else
    echo "No process manager found. Please restart manually."
fi

echo "Deployment completed successfully!"