.PHONY: help run build clean test migrate seed docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make run         - Run the application"
	@echo "  make build       - Build the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make migrate     - Run database migrations"
	@echo "  make seed        - Run database seeder"
	@echo "  make docker-up   - Start Docker containers"
	@echo "  make docker-down - Stop Docker containers"
	@echo "  make dev         - Run in development mode with hot reload"

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

clean:
	rm -rf bin/
	rm -rf tmp/

test:
	go test -v ./...

migrate:
	go run cmd/api/main.go migrate

seed:
	go run cmd/api/main.go seed

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

dev:
	air -c .air.toml