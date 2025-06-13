package main

import (
	"socially-app/internal/config"
	"socially-app/internal/router"
	logger "socially-app/internal/util"
)

func main() {
	cfg := config.Initialize()
	log := logger.Create()

	r := router.New(cfg, log)
	log.Info("Server running on port: " + cfg.Port)
	log.Fatal(r.ListenAndServe())
}
