package main

import (
	"golang-rest-api-template/config"
	"golang-rest-api-template/package/logger"
)

func main() {

	// Initialize the configuration
	config := config.NewService()
	if err := config.Init(); err != nil {
		logger.NewLogger(logger.DefaultOptions()).Error(err.Error())
		panic(err)
	}

	// Run the application
	if err := config.Run(); err != nil {
		logger.NewLogger(logger.DefaultOptions()).Error(err.Error())
		panic(err)
	}

}
