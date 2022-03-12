package main

import (
	"log"

	"github.com/Roukii/pock_multiplayer/config"
	"github.com/Roukii/pock_multiplayer/internal/app"
)

func main() {
	// Configuration
	err := config.LoadConfig("./", "config")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.RunGateway()
}
