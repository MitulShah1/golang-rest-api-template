# Load environment variables
include .env

# Get repository name from current directory
REPO_NAME ?= $(shell basename "$$(pwd)")
BUILD_NAME=server
BUILD_DIR=build
CMD_DIR=cmd/server
GO_FLAGS=-ldflags "-s -w"  # Strip debug info for a smaller binary

# Coverage directory
COVERAGE_DIR ?= .coverage

# Golang-migrate version
MIGRATE_VERSION ?= v4.16.2  # Change to the latest version if needed

# Golang-lint version
LINT_VERSION ?= v1.59.1  # Change to the latest version if needed

# Swag version
SWAG_VERSION ?= v1.16.4  # Change to the latest version if needed

# Installation directory for binaries
INSTALL_DIR ?= $(HOME)/.local/bin

# Go Imports Vesrsion
IMPORTS_VERSION ?= v0.24.0

# Go Vulncheck Version
VULN_VERSION ?= v1.1.3

# Formatting for beautiful terminal output
BLUE=\033[1;34m
GREEN=\033[1;32m
YELLOW=\033[1;33m
NC=\033[0m  # No Color

# ───────────────────────────────────────────────────────────
# 📝 CHECK & COPY .env IF MISSING
# ───────────────────────────────────────────────────────────
env:
	@echo -e "$(YELLOW)🔍 Checking for .env file...$(NC)"
	@if [ ! -f .env ]; then \
		echo -e "$(RED)⚠️  .env file not found! Creating from .env.example...$(NC)"; \
		cp .env.example .env; \
		echo -e "$(GREEN)✅ .env file created successfully!$(NC)"; \
	else \
		echo -e "$(GREEN)✅ .env file exists!$(NC)"; \
	fi

# ───────────────────────────────────────────────────────────
# 🎨 FORMAT CODE (gofmt & goimports)
# ───────────────────────────────────────────────────────────
format:
	@echo -e "$(YELLOW)🎨 Formatting Go code...$(NC)"
	@gofmt -w .
	@go install golang.org/x/tools/cmd/goimports@$(IMPORTS_VERSION)
	@goimports -w .
	@echo -e "$(GREEN)✅ Code formatted successfully!$(NC)"

# ───────────────────────────────────────────────────────────
# 🔍 RUN GO VET (Code Inspection)
# ───────────────────────────────────────────────────────────
vet:
	@echo -e "$(YELLOW)🔍 Running go vet...$(NC)"
	@go vet ./...
	@echo -e "$(GREEN)✅ go vet completed!$(NC)"

# ───────────────────────────────────────────────────────────
# 🛡️ SECURITY SCAN (govulncheck)
# ───────────────────────────────────────────────────────────
security_scan:
	@echo -e "$(RED)🛡️ Running security vulnerability scan...$(NC)"
	@go install golang.org/x/vuln/cmd/govulncheck@$(VULN_VERSION)
	@govulncheck ./...
	@echo -e "$(GREEN)✅ Security scan completed!$(NC)"

# ───────────────────────────────────────────────────────────
# 🔄 Install DEPENDENCIES (go mod tidy & upgrade)
# ───────────────────────────────────────────────────────────
install_deps:
	@echo -e "$(YELLOW)🔄 Install Go dependencies....$(NC)"
	@go mod tidy	
	@echo -e "$(GREEN)✅ Dependencies updated!$(NC)"
	
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
# 🏃 RUN APPLICATION
# ───────────────────────────────────────────────────────────
run:
	@echo -e "$(BLUE)🚀 Running the application...$(NC)"
	@go run cmd/server/main.go

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
# 📥 INSTALL SWAG CLI TOOL & PACKAGES
# ───────────────────────────────────────────────────────────
install_swag:
	@echo -e "$(GREEN)📥 Installing Swag CLI and dependencies...$(NC)"
	@which swag >/dev/null 2>&1 || (echo -e "$(RED)❌ Swag CLI not found! Installing now...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest)
	@echo -e "$(YELLOW)🔄 Updating project dependencies for Swag...$(NC)"
	@go mod tidy
	@go mod download
	@echo -e "$(GREEN)✅ Swag installation complete!$(NC)"

# ───────────────────────────────────────────────────────────
# 📜 GENERATE API DOCUMENTATION
# ───────────────────────────────────────────────────────────
generate_docs: install_swag
	@echo -e "$(YELLOW)📜 Generating API documentation using Swag...$(NC)"
	@swag init -g ./cmd/server/main.go -o ./docs
	@echo -e "$(GREEN)✅ API documentation generated successfully!$(NC)"

# ───────────────────────────────────────────────────────────
# 🏗️ BUILD PROJECT
# ───────────────────────────────────────────────────────────
build:	
	@echo -e "$(BLUE)🏗️ Building the Go application...$(NC)"
	@mkdir -p $(BUILD_DIR)  # ✅ Ensure the build directory exists
	@CGO_ENABLED=0 GOOS=linux go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BUILD_NAME) $(CMD_DIR)/main.go
	@ls -lh $(BUILD_DIR)  # ✅ Debug: List contents of the build directory
	@echo -e "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BUILD_NAME)$(NC)"

# ───────────────────────────────────────────────────────────
# 🧹 CLEAN BUILD & COVERAGE FILES
# ───────────────────────────────────────────────────────────
clean:
	@echo -e "$(YELLOW)🧹 Cleaning up build and coverage files...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@echo -e "$(GREEN)✅ Cleanup complete!$(NC)"

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

# ───────────────────────────────────────────────────────────
# 🐳 BUILD DOCKER IMAGE
# ───────────────────────────────────────────────────────────
docker_build: env docker_down
	@echo -e "$(BLUE)🐳 Building Docker image...$(NC)"
	@sudo docker-compose build
	@echo -e "$(GREEN)✅ Docker image built successfully!$(NC)"

# ───────────────────────────────────────────────────────────
# 🚀 START DOCKER CONTAINERS
# ───────────────────────────────────────────────────────────
docker_up: docker_build
	@echo -e "$(BLUE)🚀 Starting Docker containers...$(NC)"
	@sudo docker-compose up -d
	@echo -e "$(GREEN)✅ Docker containers started successfully!$(NC)"

# ───────────────────────────────────────────────────────────
# 📦 STOP & REMOVE DOCKER CONTAINERS
# ───────────────────────────────────────────────────────────
docker_down:
	@echo -e "$(YELLOW)🛑 Stopping and removing Docker containers...$(NC)"
	@sudo docker-compose down
	@echo -e "$(GREEN)✅ Docker containers stopped and removed!$(NC)"

# ───────────────────────────────────────────────────────────
# 📜 VIEW DOCKER LOGS
# ───────────────────────────────────────────────────────────
docker_logs:
	@echo -e "$(YELLOW)📜 Viewing Docker logs...$(NC)"
	@sudo docker-compose logs -f

# ───────────────────────────────────────────────────────────
# ✅ CLEAN DOCKER IMAGES & CONTAINERS
# ───────────────────────────────────────────────────────────
docker_clean: docker_down
	@echo -e "$(RED)🗑️ Cleaning up Docker images and containers...$(NC)"
	@sudo docker system prune -af
	@echo -e "$(GREEN)✅ Docker cleanup complete!$(NC)"

# ───────────────────────────────────────────────────────────
# 🚀 CI/CD PRE-CHECK (Runs everything before deployment)
# ───────────────────────────────────────────────────────────
ci_check: env format vet lint staticcheck security_scan test
	@echo -e "$(GREEN)✅ CI/CD pre-check passed successfully!$(NC)"

# Mark these targets as non-file targets
.PHONY: env clean build test install_migration create_migration \
		migrate_up migrate_down version install_swag generate_docs \
		ci_check format vet lint staticcheck security_scan \
		install_deps html-coverage docker_build docker_up \
		docker_down docker_logs docker_clean
