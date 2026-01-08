# Load environment variables (optional, won't fail if .env doesn't exist)
-include .env

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

include help.mk  # place after ALL target and before all other targets

# ───────────────────────────────────────────────────────────
# 📝 CHECK & COPY .env IF MISSING
# ───────────────────────────────────────────────────────────
env: ## 📝 CHECK & COPY .env IF MISSING
	@echo -e "$(YELLOW)🔍 Checking for .env file...$(NC)"
	@if [ ! -f .env ]; then \
		echo -e "$(RED)⚠️  .env file not found! Creating from .env.example...$(NC)"; \
		cp .env.example .env; \
		echo -e "$(GREEN)✅ .env file created successfully!$(NC)"; \
		echo -e "$(YELLOW)📝 Please update Redis configuration in .env file:$(NC)"; \
		echo -e "$(BLUE)   REDIS_HOST=localhost$(NC)"; \
		echo -e "$(BLUE)   REDIS_PORT=6379$(NC)"; \
		echo -e "$(BLUE)   REDIS_PASSWORD=your_redis_password$(NC)"; \
		echo -e "$(BLUE)   REDIS_DB=0$(NC)"; \
	else \
		echo -e "$(GREEN)✅ .env file exists!$(NC)"; \
	fi

# ───────────────────────────────────────────────────────────
# 🎨 FORMAT CODE (gofmt & goimports)
# ───────────────────────────────────────────────────────────
format: ## 🎨 FORMAT CODE (gofmt & goimports)
	@echo -e "$(YELLOW)🎨 Formatting Go code...$(NC)"
	@gofmt -w .
	@go install golang.org/x/tools/cmd/goimports@$(IMPORTS_VERSION)
	@goimports -w .
	@echo -e "$(GREEN)✅ Code formatted successfully!$(NC)"

vet: ## 🔍 RUN GO VET (Code Inspection)
	@echo -e "$(YELLOW)🔍 Running go vet...$(NC)"
	@go vet ./...
	@echo -e "$(GREEN)✅ go vet completed!$(NC)"

security_scan: ## 🛡️ SECURITY SCAN (govulncheck)
	@echo -e "$(RED)🛡️ Running security vulnerability scan...$(NC)"
	@go install golang.org/x/vuln/cmd/govulncheck@$(VULN_VERSION)
	@govulncheck ./...
	@echo -e "$(GREEN)✅ Security scan completed!$(NC)"

install_deps: ## 🔄 Install DEPENDENCIES (go mod tidy & upgrade)
	@echo -e "$(YELLOW)🔄 Install Go dependencies....$(NC)"
	@go mod tidy	
	@echo -e "$(GREEN)✅ Dependencies updated!$(NC)"
	
lint: ## 🔎 LINT CODE (golangci-lint)
	@echo -e "$(YELLOW)🔎 Running golangci-lint...$(NC)"
	@which golangci-lint >/dev/null 2>&1 || (echo -e "$(RED)❌ golangci-lint not installed! Installing now...$(NC)" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINT_VERSION))
	@golangci-lint run ./...

staticcheck: ## 📢 STATIC CODE ANALYSIS (staticcheck)
	@echo -e "$(YELLOW)📢 Running staticcheck...$(NC)"
	@which staticcheck >/dev/null 2>&1 || (echo -e "$(RED)❌ staticcheck not installed! Installing now...$(NC)" && go install honnef.co/go/tools/cmd/staticcheck@latest)
	@staticcheck ./...

run: ## 🏃 RUN APPLICATION
	@echo -e "$(BLUE)🚀 Running the application...$(NC)"
	@go run cmd/server/main.go

test: ## ✅ RUN TESTS
	@echo -e "$(YELLOW)🔍 Running tests...$(NC)"
	@go test -v ./...

html-coverage: $(COVERAGE_DIR)/coverage.out ## 📊 GENERATE COVERAGE REPORT
	@echo -e "$(GREEN)📊 Generating HTML coverage report...$(NC)"
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo -e "$(GREEN)✅ HTML coverage report generated at $(COVERAGE_DIR)/coverage.html$(NC)"

	# Open the file based on OS
	@uname | grep -qi "darwin" && open $(COVERAGE_DIR)/coverage.html || xdg-open $(COVERAGE_DIR)/coverage.html

$(COVERAGE_DIR)/coverage.out: | $(COVERAGE_DIR)
	@echo -e "$(YELLOW)📈 Running coverage analysis...$(NC)"
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...

$(COVERAGE_DIR): ## Ensure .coverage directory exists
	@mkdir -p $(COVERAGE_DIR)

install_swag: ## 📥 INSTALL SWAG CLI TOOL & PACKAGES
	@echo -e "$(GREEN)📥 Installing Swag CLI and dependencies...$(NC)"
	@which swag >/dev/null 2>&1 || (echo -e "$(RED)❌ Swag CLI not found! Installing now...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest)
	@echo -e "$(YELLOW)🔄 Updating project dependencies for Swag...$(NC)"
	@go mod tidy
	@go mod download
	@echo -e "$(GREEN)✅ Swag installation complete!$(NC)"

generate_docs: install_swag ## 📜 GENERATE API DOCUMENTATION
	@echo -e "$(YELLOW)📜 Generating API documentation using Swag...$(NC)"
	@swag init --parseDependency  --parseInternal --parseDepth 1 -g ./cmd/server/main.go -o ./docs
	@echo -e "$(GREEN)✅ API documentation generated successfully!$(NC)"

build:	 ## 🏗️ BUILD PROJECT
	@echo -e "$(BLUE)🏗️ Building the Go application...$(NC)"
	@mkdir -p $(BUILD_DIR)  # ✅ Ensure the build directory exists
	@CGO_ENABLED=0 GOOS=linux go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BUILD_NAME) $(CMD_DIR)/main.go
	@ls -lh $(BUILD_DIR)  # ✅ Debug: List contents of the build directory
	@echo -e "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BUILD_NAME)$(NC)"

clean: ## 🧹 CLEAN BUILD & COVERAGE FILES
	@echo -e "$(YELLOW)🧹 Cleaning up build and coverage files...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@echo -e "$(GREEN)✅ Cleanup complete!$(NC)"

version: ## 🔍 CHECK MIGRATION VERSION
	@echo -e "$(BLUE)🔍 Checking installed migrate version...$(NC)"
	@$(INSTALL_DIR)/migrate -version

install_migration: ## 📥 INSTALL GOLANG-MIGRATE
	@echo -e "$(GREEN)📥 Installing golang-migrate ($(MIGRATE_VERSION))...$(NC)"
	@mkdir -p $(INSTALL_DIR)
	@curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.linux-amd64.tar.gz -o migrate.tar.gz
	@tar -xvf migrate.tar.gz
	@mv migrate $(INSTALL_DIR)/migrate
	@chmod +x $(INSTALL_DIR)/migrate
	@rm -f migrate.tar.gz
	@echo -e "$(GREEN)✅ Installation complete. Ensure $(INSTALL_DIR) is in your PATH.$(NC)"

create_migration: ## 📦 CREATE A NEW DATABASE MIGRATION
	@echo -e "$(YELLOW)📦 Creating a new database migration...$(NC)"
	@$(INSTALL_DIR)/migrate create -ext=sql -dir=package/database/migrations -seq init

migrate_up: ## ⬆️ APPLY DATABASE MIGRATIONS
	@echo -e "$(GREEN)⬆️ Applying database migrations...$(NC)"
	@$(INSTALL_DIR)/migrate -path=package/database/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
		-verbose up
migrate_down: ## ⬇️ ROLLBACK DATABASE MIGRATIONS
	@echo -e "$(RED)⬇️ Rolling back database migrations...$(NC)"
	@$(INSTALL_DIR)/migrate -path=package/database/migrations \
		-database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" \
		-verbose down
docker_build: env docker_down ## 🐳 BUILD DOCKER IMAGE
	@echo -e "$(BLUE)🐳 Building Docker image...$(NC)"
	@sudo docker-compose build
	@echo -e "$(GREEN)✅ Docker image built successfully!$(NC)"

docker_up: docker_build ## 🚀 START DOCKER CONTAINERS
	@echo -e "$(BLUE)🚀 Starting Docker containers...$(NC)"
	@sudo docker-compose up -d
	@echo -e "$(GREEN)✅ Docker containers started successfully!$(NC)"

docker_down: ## 🛑 STOP & REMOVE DOCKER CONTAINERS
	@echo -e "$(YELLOW)🛑 Stopping and removing Docker containers...$(NC)"
	@sudo docker-compose down
	@echo -e "$(GREEN)✅ Docker containers stopped and removed!$(NC)"

docker_logs: ## 📜 VIEW DOCKER LOGS
	@echo -e "$(YELLOW)📜 Viewing Docker logs...$(NC)"
	@sudo docker-compose logs -f

docker_clean: docker_down ## 🗑️ CLEAN DOCKER IMAGES & CONTAINERS
	@echo -e "$(RED)🗑️ Cleaning up Docker images and containers...$(NC)"
	@sudo docker system prune -af
	@echo -e "$(GREEN)✅ Docker cleanup complete!$(NC)"

ci_check: env format vet lint staticcheck security_scan test ## 🚀 CI/CD PRE-CHECK
	@echo -e "$(GREEN)✅ CI/CD pre-check passed successfully!$(NC)"

# Mark these targets as non-file targets
.PHONY: env clean build test install_migration create_migration \
		migrate_up migrate_down version install_swag generate_docs \
		ci_check format vet lint staticcheck security_scan \
		install_deps html-coverage docker_build docker_up \
		docker_down docker_logs docker_clean
