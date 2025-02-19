# golang-rest-api-template
Go Rest API Templates


golang-microservice-template/
│── cmd/
│   ├── server/                # Main entry point for the service
│   │   ├── main.go
│── config/
│   ├── config.yaml            # Configuration file
│   ├── config.go              # Load configuration
│── internal/
│   ├── handlers/              # HTTP handlers
│   │   ├── user_handler.go
│   │   ├── health_handler.go
│   ├── services/              # Business logic
│   │   ├── user_service.go
│   ├── repository/            # Data access layer
│   │   ├── user_repository.go
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
│── Jenkinsfile                  # Jenkins CI/CD pipeline
│── README.md                   # Project documentation
