# Load environment variables
include .env

# Get repository name from current directory
REPO_NAME ?= $(shell basename "$$(pwd)")

# Coverage directory
COVERAGE_DIR ?= .coverage

# Golang-migrate version
MIGRATE_VERSION := v4.16.2  # Change to the latest version if needed

# Golang-lint version
LINT_VERSION := v1.59.1  # Change to the latest version if needed

# Installation directory for binaries
INSTALL_DIR := $(HOME)/.local/bin

# Formatting for beautiful terminal output
BLUE=\033[1;34m
GREEN=\033[1;32m
YELLOW=\033[1;33m
NC=\033[0m  # No Color

# ───────────────────────────────────────────────────────────
# 🏃 RUN APPLICATION
# ───────────────────────────────────────────────────────────
run:
	@echo -e "$(BLUE)🚀 Running the application...$(NC)"
	@go run cmd/server/main.go

# ───────────────────────────────────────────────────────────
# 🔎 LINT CODE (golangci-lint)
# ───────────────────────────────────────────────────────────
lint:
	@echo -e "$(YELLOW)🔎 Running golangci-lint...$(NC)"
	@which golangci-lint >/dev/null 2>&1 || (echo -e "$(RED)❌ golangci-lint not installed! Installing now...$(NC)" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINT_VERSION))
	@golangci-lint run ./...

# ───────────────────────────────────────────────────────────
# 📢 STATIC CODE ANALYSIS (staticcheck)
# ───────────────────────────────────────────────────────────
staticcheck:
	@echo -e "$(YELLOW)📢 Running staticcheck...$(NC)"
	@which staticcheck >/dev/null 2>&1 || (echo -e "$(RED)❌ staticcheck not installed! Installing now...$(NC)" && go install honnef.co/go/tools/cmd/staticcheck@latest)
	@staticcheck ./...

# ───────────────────────────────────────────────────────────
# ✅ RUN TESTS
# ───────────────────────────────────────────────────────────
test:
	@echo -e "$(YELLOW)🔍 Running tests...$(NC)"
	@go test -v ./...

# ───────────────────────────────────────────────────────────
# 📊 GENERATE COVERAGE REPORT
# ───────────────────────────────────────────────────────────
html-coverage: $(COVERAGE_DIR)/.combined.html
	@echo -e "$(GREEN)📊 Generating HTML coverage report...$(NC)"
	@go tool cover -html=$(COVERAGE_DIR)/.combined.html

$(COVERAGE_DIR)/.combined.html: $(COVERAGE_DIR)/coverage.out | $(COVERAGE_DIR)
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/.combined.html

$(COVERAGE_DIR)/coverage.out: | $(COVERAGE_DIR)
	@echo -e "$(YELLOW)📈 Running coverage analysis...$(NC)"
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...

# Ensure .coverage directory exists
$(COVERAGE_DIR):
	@mkdir -p $(COVERAGE_DIR)

# ───────────────────────────────────────────────────────────
# 🏗️ BUILD PROJECT
# ───────────────────────────────────────────────────────────
build: clean
	@echo -e "$(BLUE)🏗️ Building the project...$(NC)"
	@go build -o build/ cmd/server/main.go

# ───────────────────────────────────────────────────────────
# 🧹 CLEAN BUILD & COVERAGE FILES
# ───────────────────────────────────────────────────────────
clean:
	@echo -e "$(YELLOW)🧹 Cleaning up build and coverage files...$(NC)"
	@rm -rf build/*
	@rm -rf $(COVERAGE_DIR)

# ───────────────────────────────────────────────────────────
# 🔍 CHECK MIGRATION VERSION
# ───────────────────────────────────────────────────────────
version:
	@echo -e "$(BLUE)🔍 Checking installed migrate version...$(NC)"
	@$(INSTALL_DIR)/migrate -version

# ───────────────────────────────────────────────────────────
# 📥 INSTALL GOLANG-MIGRATE
# ───────────────────────────────────────────────────────────
install_migration:
	@echo -e "$(GREEN)📥 Installing golang-migrate ($(MIGRATE_VERSION))...$(NC)"
	@mkdir -p $(INSTALL_DIR)
	@curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.linux-amd64.tar.gz -o migrate.tar.gz
	@tar -xvf migrate.tar.gz
	@mv migrate $(INSTALL_DIR)/migrate
	@chmod +x $(INSTALL_DIR)/migrate
	@rm -f migrate.tar.gz
	@echo -e "$(GREEN)✅ Installation complete. Ensure $(INSTALL_DIR) is in your PATH.$(NC)"

# ───────────────────────────────────────────────────────────
# 📦 CREATE A NEW DATABASE MIGRATION
# ───────────────────────────────────────────────────────────
create_migration:
	@echo -e "$(YELLOW)📦 Creating a new database migration...$(NC)"
	@$(INSTALL_DIR)/migrate create -ext=sql -dir=package/database/migrations -seq init

# ───────────────────────────────────────────────────────────
# ⬆️ APPLY DATABASE MIGRATIONS
# ───────────────────────────────────────────────────────────
migrate_up:
	@echo -e "$(GREEN)⬆️ Applying database migrations...$(NC)"
	@$(INSTALL_DIR)/migrate -path=package/database/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
		-verbose up

# ───────────────────────────────────────────────────────────
# ⬇️ ROLLBACK DATABASE MIGRATIONS
# ───────────────────────────────────────────────────────────
migrate_down:
	@echo -e "$(RED)⬇️ Rolling back database migrations...$(NC)"
	@$(INSTALL_DIR)/migrate -path=package/database/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
		-verbose down

# Mark these targets as non-file targets
.PHONY: clean build test install_migration create_migration migrate_up migrate_down version
