# golang-rest-api-template
Go Rest API Templates

## Overview
This is a template for building production-ready REST APIs using Go. It follows best practices and includes a standardized project structure with all necessary components for building scalable microservices.

## Project Structure
```
golang-microservice-template/
│── cmd/
│   ├── server/                # Main entry point for the service
│   │   ├── main.go
│── config/
│   ├── config.go              # Application configuration
│── internal/
│   ├── handlers/              # HTTP handlers
│   ├── services/              # Business logic
│   ├── repository/            # Data access layer
│   ├── models/                # Structs & DTOs
│   │   ├── user.go
│── pkg/                       # Utility packages (logging, middleware, etc.)
│   ├── logger/
│   │   ├── logger.go
│   ├── middleware/
│   │   ├── auth.go
│── api/
│   ├── swagger/               # API documentation
│── test/
│   ├── e2e/                   # End-to-end tests
│   │   ├── user_test.go
│── migrations/                # Database migrations
│── scripts/                   # Automation scripts
│   ├── entrypoint.sh          # Docker entrypoint script
│── Dockerfile                 # Docker build configuration
│── docker-compose.yml          # Docker Compose setup
│── Makefile                    # Build automation
│── go.mod                      # Go module dependencies
│── go.sum                      # Dependencies lock file
│── Jenkinsfile                 # Jenkins CI/CD pipeline
│── README.md                   # Project documentation
```

## Features
- Structured logging
- Middleware support (authentication, etc.)
- Configuration management
- API documentation with Swagger
- Docker support
- CI/CD pipeline with Jenkins
- Database migrations
- End-to-end testing
- Makefile for common operations

## Getting Started

### Prerequisites
- Go 1.16 or higher
- Docker and Docker Compose
- Make

### Running the Application
1. Clone the repository
```bash
git clone https://github.com/MitulShah1/golang-rest-api-template.git
```

2. Start the application using Docker Compose
```bash
docker-compose up
```

### Development
Build the application:
```bash
make build
```

Run tests:
```bash
make test
```

Generate API documentation:
```bash
make swagger
```

## Configuration
Configuration is managed through `config/config.yaml`. Environment variables can override these settings.

## Testing
- Unit tests are alongside the code
- Integration tests are in the `test/` directory
- Run all tests with `make test`

## Deployment
The project includes:
- Dockerfile for containerization
- docker-compose.yml for local development
- Jenkinsfile for CI/CD pipeline

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details


