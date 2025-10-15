package main

import (
	"log"

	"github.com/minhhoccode111/realworld-fiber-clean/config"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
