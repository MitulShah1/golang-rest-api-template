include .env

REPO_NAME ?= $(shell basename "$$(pwd)")
COVERAGE_DIR ?= .coverage
# Set the golang-migrate version
MIGRATE_VERSION := v4.16.2# Change to the latest version if needed
# Installation directory
INSTALL_DIR := $(HOME)/.local/bin

run:
	go run cmd/server/main.go

test:
	go test -v ./...

html-coverage:
	go tool cover -html=$(COVERAGE_DIR)/combined.txt

build: clean
	go build -o build/ cmd/server/main.go

clean:
	@-rm -r build/*

# Check installed version
version:
	@$(INSTALL_DIR)/migrate -version

# Download and install golang-migrate
install_migration:
	@echo "Installing golang-migrate..."
	@mkdir -p $(INSTALL_DIR)
	echo "https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.linux-amd64.tar.gz"
	@curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.linux-amd64.tar.gz -o migrate.tar.gz
	@tar -xvf migrate.tar.gz
	@mv migrate $(INSTALL_DIR)/migrate
	@chmod +x $(INSTALL_DIR)/migrate
	@rm -f migrate.tar.gz
	@echo "Installation complete. Ensure $(INSTALL_DIR) is in your PATH."
create_migration:
	@$(INSTALL_DIR)/migrate create -ext=sql -dir=package/database/migrations -seq init

migrate_up:
	@$(INSTALL_DIR)/migrate -path=package/database/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
		-verbose up

migrate_down:
	@$(INSTALL_DIR)/migrate -path=package/database/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
		-verbose down

.PHONY: clean build test install_migration create_migration migrate_up migrate_down