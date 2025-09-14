# ───────────────────────────────────────────────────────────
# 🌟 STAGE 1: BUILD GO APPLICATION
# ───────────────────────────────────────────────────────────
FROM golang:1.25 AS builder

# Accept build arguments
ARG SERVER_PORT
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG DEBUG
ARG DISABLE_LOGS
ARG LOG_FORMAT
ARG LOG_CALLER
ARG LOG_STACKTRACE
ARG SWAG_VERSION
ARG MIGRATE_VERSION
ARG LINT_VERSION
ARG IMPORTS_VERSION
ARG VULN_VERSION

# Set environment variables for versions (same as Makefile)
ENV SERVER_PORT=$SERVER_PORT
ENV DB_HOST=$DB_HOST
ENV DB_PORT=$DB_PORT
ENV DB_USER=$DB_USER
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_NAME=$DB_NAME
ENV DEBUG=$DEBUG
ENV DISABLE_LOGS=$DISABLE_LOGS
ENV LOG_FORMAT=$LOG_FORMAT
ENV LOG_CALLER=$LOG_CALLER
ENV LOG_STACKTRACE=$LOG_STACKTRACE
ARG SWAG_VERSION=$SWAG_VERSION
ARG MIGRATE_VERSION=$MIGRATE_VERSION
ARG LINT_VERSION=$LINT_VERSION
ARG IMPORTS_VERSION=$IMPORTS_VERSION
ARG VULN_VERSION=$VULN_VERSION

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

# 🔥 Ensure `/app/build/` exists
RUN mkdir -p /app/build

# Generate Swagger documentation
RUN make generate_docs

# 🔥 Build the Go application using Makefile
RUN make build

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
EXPOSE ${SERVER_PORT}

# Run the application
CMD ["/app/server"]
