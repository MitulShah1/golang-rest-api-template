package main

import (
	"golang-rest-api-template/config"
)

func main() {

	// Initialize the configuration
	config := config.Service{}
	if err := config.Init(); err != nil {
		panic(err)
	}

	// Run the application
	if err := config.Run(); err != nil {
		panic(err)
	}

}
