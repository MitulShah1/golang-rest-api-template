# ───────────────────────────────────────────────────────────
# 🌟 STAGE 1: BUILD GO APPLICATION
# ───────────────────────────────────────────────────────────
FROM golang:1.21 AS builder

# Set environment variables for versions (same as Makefile)
ARG SWAG_VERSION=v1.16.4
ARG MIGRATE_VERSION=v4.16.2
ARG LINT_VERSION=v1.59.1
ARG IMPORTS_VERSION=v0.24.0
ARG VULN_VERSION=v1.1.3

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod tidy && go mod download

# Install Swagger, GolangCI-Lint, GoImports, and Govulncheck
RUN go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION} && \
    go install golang.org/x/tools/cmd/goimports@${IMPORTS_VERSION} && \
    go install golang.org/x/vuln/cmd/govulncheck@${VULN_VERSION} && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINT_VERSION}

# ───────────────────────────────────────────────────────────
# ✅ INSTALL `golang-migrate` USING OFFICIAL METHOD
# ───────────────────────────────────────────────────────────
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz -o migrate.tar.gz && \
    tar -xvf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    rm -f migrate.tar.gz

# Copy the entire project into the container
COPY . .

# Generate Swagger documentation
RUN make generate_docs

# Build the Go application 
RUN CGO_ENABLED=0 GOOS=linux make build

# ───────────────────────────────────────────────────────────
# 🌟 STAGE 2: CREATE A SMALLER FINAL IMAGE
# ───────────────────────────────────────────────────────────
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Install dependencies required to run the app
RUN apt-get update && apt-get install -y ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/build/server /app/server

# Copy the Swagger docs to serve them later
COPY --from=builder /app/docs /app/docs

# Copy the migration files
COPY --from=builder /app/package/database/migrations /app/migrations

# Copy the migrate tool from the builder stage
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/server"]
