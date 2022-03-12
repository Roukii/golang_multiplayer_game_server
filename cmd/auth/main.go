package main

import (
	"log"

	"github.com/Roukii/pock_multiplayer/config"
	"github.com/Roukii/pock_multiplayer/internal/gateway"
)

func main() {
	// Configuration
	err := config.LoadConfig("./", "config")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	gateway.RunGateway()
}
