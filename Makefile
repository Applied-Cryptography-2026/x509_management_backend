.PHONY: build run start migrate clean test lint fmt deps

# Binary names
APP_BINARY=x509-app
CLI_BINARY=x509-cli
MIGRATE_BINARY=x509-migrate

# Directories
BUILD_DIR=./bin

## Build all binaries
build: build-app build-cli

build-app:
	go build -o $(BUILD_DIR)/$(APP_BINARY) ./cmd/app

build-cli:
	go build -o $(BUILD_DIR)/$(CLI_BINARY) ./cmd/cli

## Run the application with hot-reload (requires air)
start:
	air

## Run the application directly
run:
	go run ./cmd/app

## Run database migrations
migrate:
	go run ./cmd/migration

## Run all tests
test:
	go test -v -race ./...

## Run tests with coverage
test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Lint the codebase
lint:
	golangci-lint run ./...

## Format the codebase
fmt:
	go fmt ./...
	goimports -w .

## Download/update dependencies
deps:
	go mod tidy
	go mod download

## Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## Set up the development database
setup-db:
	docker compose -f docker/docker-compose.yml up -d
	sleep 5
	docker compose -f docker/docker-compose.yml exec mysql \
		mysql -uroot -proot < db/migrations/001_setup.sql

## Start MySQL via docker compose
db-up:
	docker compose -f docker/docker-compose.yml up -d

## Stop MySQL
db-down:
	docker compose -f docker/docker-compose.yml down

## Build the CLI tool
cli:
	go run ./cmd/cli -cmd info -cert ./testdata/test-cert.pem
