include .env

MIGRATION_NAME ?= unnamed_migration

.PHONY: create_migration migrate_up migrate_down local local-down local-logs dev run build test clean help

create_migration:
	migrate create -ext=sql -dir=internal/adapter/repository/postgres/migrations -seq $(MIGRATION_NAME)

migrate_up:
	migrate -path=internal/adapter/repository/postgres/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose up

migrate_down:
	migrate -path=internal/adapter/repository/postgres/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose down 1


local:
	@echo "Starting local environment (Redis and PostgreSQL)"
	@echo "Services will be available at:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	docker compose -f docker-compose.local.yml up --build -d

local-down:
	docker compose -f docker-compose.local.yml down

local-logs:
	docker compose -f docker-compose.local.yml logs -f

dev:
	@echo "Starting development server with hot reload (Air)..."
	@AIR_BIN=$$(go env GOPATH)/bin/air; \
	if [ ! -f "$$AIR_BIN" ]; then \
		echo "Error: Air is not installed. Run: go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi; \
	$$AIR_BIN

run:
	@echo "Running application (no hot reload)..."
	go run ./cmd/http/main.go

build:
	@echo "Building binary..."
	go build -o bin/http ./cmd/http/main.go

test:
	go test -cover ./...

clean:
	docker system prune -f

help:
	@echo "Migrations usage:"
	@echo "  make create_migration MIGRATION_NAME=migration_name"
	@echo "  make create_migration  # uses 'unnamed_migration' as default"
	@echo "  make migrate_up        # executes all pending migrations"
	@echo "  make migrate_down      # reverts the last migration"