# Template Setup Guide

This guide helps you set up your new repository after creating it from this template.

## 🚀 Initial Setup

### 1. Update Module Name
Edit `go.mod` and change the module name:
```go
// Change from:
module github.com/MitulShah1/golang-rest-api-template

// To your module name:
module github.com/YOUR_USERNAME/YOUR_REPO_NAME
```

### 2. Update Repository References
Search and replace these references in your codebase:
- `github.com/MitulShah1/golang-rest-api-template` → `github.com/YOUR_USERNAME/YOUR_REPO_NAME`
- `MitulShah1/golang-rest-api-template` → `YOUR_USERNAME/YOUR_REPO_NAME`

### 3. Update README.md
- Change the project title
- Update badges to point to your repository
- Modify the description to match your project

### 4. Update GitHub Actions
In `.github/workflows/go.yml`, update the repository name in badges if needed.

### 5. Update Docker Configuration
In `docker-compose.yml`, consider updating service names to match your project.

## 🔧 Customization Options

### Database Configuration
- Update database connection settings in `config/config.go`
- Modify migration files in `package/database/migrations/`
- Update database driver if switching from MySQL

### API Endpoints
- Modify existing handlers in `internal/handlers/`
- Add new endpoints following the established pattern
- Update Swagger documentation for new endpoints

### Middleware
- Customize middleware in `package/middleware/`
- Add authentication/authorization as needed
- Configure CORS settings for your domain

### Environment Variables
- Update `.env.example` with your specific configuration
- Add new environment variables as needed
- Document all required environment variables

## 📝 Best Practices for Template Usage

### 1. Keep the Structure
- Maintain the established project structure
- Follow the existing patterns for handlers, services, and repositories
- Use the provided middleware and utilities

### 2. Testing
- Write tests for all new functionality
- Follow the existing test patterns
- Maintain high test coverage

### 3. Documentation
- Update Swagger documentation for new endpoints
- Keep README.md up to date
- Document any new configuration options

### 4. CI/CD
- The GitHub Actions workflow is ready to use
- Update repository secrets as needed
- Configure deployment targets

## 🎯 Quick Commands

After setup, use these commands to get started:

```bash
# Install dependencies
go mod tidy

# Run tests
make test

# Start development environment
make docker_up

# Generate API documentation
make generate_docs

# Build the application
make build
```

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Swagger Documentation](https://swagger.io/docs/)
- [Docker Documentation](https://docs.docker.com/)

## 🤝 Support

If you encounter issues with the template:
1. Check the existing issues in this repository
2. Create a new issue with detailed information
3. Consider contributing back improvements

---

**Happy Coding! 🚀** 