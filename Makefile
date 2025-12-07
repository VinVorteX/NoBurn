.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version

DB_URL ?= postgres://noburn_user:noburn_pass@localhost:5433/noburn_db?sslmode=disable

# Migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path migrations -database "$(DB_URL)" force $$version

migrate-version:
	@migrate -path migrations -database "$(DB_URL)" version

# Development (local)
dev-api:
	go run cmd/api/main.go

dev-worker:
	go run cmd/worker/main.go

# Docker
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-restart:
	docker compose restart

docker-clean:
	docker compose down -v

# Docker with rebuild
docker-rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d

# Build binaries
build-api:
	go build -o bin/api cmd/api/main.go

build-worker:
	go build -o bin/worker cmd/worker/main.go

build-all: build-api build-worker

# Test
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
