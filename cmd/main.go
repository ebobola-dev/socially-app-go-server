package main

import (
	"socially-app/internal/config"
	"socially-app/internal/router"
	logger "socially-app/internal/util"
)

func main() {
	cfg := config.Initialize()
	log := logger.Create(cfg)

	app := router.New(cfg, log)
	log.PrintConfig(cfg)
	log.Info("Server running on port: %s\n", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
