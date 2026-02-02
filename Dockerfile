# ───────────────────────────────────────────────────────────
# 🌟 STAGE 1: BUILD GO APPLICATION
# ───────────────────────────────────────────────────────────
FROM golang:1.25 AS builder

# Build args (no secrets - app gets config at runtime via env_file)
ARG SWAG_VERSION
ARG MIGRATE_VERSION

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod tidy && go mod download

# Install Swagger (for docs) and golang-migrate (for DB migrations)
RUN go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}
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

EXPOSE 8080

# Run the application
CMD ["/app/server"]
