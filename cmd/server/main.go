package main

import (
	"golang-rest-api-template/config"
	_ "golang-rest-api-template/internal/handlers/category/model"
	"golang-rest-api-template/package/logger"
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

	// Initialize the logger
	log := logger.NewLogger(logger.DefaultOptions())

	// Initialize the configuration
	config := config.NewService()
	if err := config.Init(); err != nil {
		log.Fatal("error while initize app", "error", err.Error())
	}

	// Run the application
	if err := config.Run(); err != nil {
		log.Fatal("error while run app", "error", err.Error())
	}

}
