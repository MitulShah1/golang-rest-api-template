package main

import (
	"github.com/MitulShah1/golang-rest-api-template/internal/application"
	_ "github.com/MitulShah1/golang-rest-api-template/internal/handlers/category/model"
)

// @title           REST API Template Example
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Create and initialize the application
	app := application.NewApplication()

	// Initialize all application components
	if err := app.Initialize(); err != nil {
		app.GetLogger().Fatal("failed to initialize application", "error", err.Error())
	}

	// Run the application (this will handle graceful shutdown)
	if err := app.Run(); err != nil {
		app.GetLogger().Fatal("application failed", "error", err.Error())
	}
}
