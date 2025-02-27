package main

import (
	"golang-rest-api-template/config"
	"golang-rest-api-template/package/logger"
)

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
