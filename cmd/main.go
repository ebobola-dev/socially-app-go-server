package main

import (
	"context"

	scope "github.com/ebobola-dev/socially-app-go-server/internal/di"
	"github.com/ebobola-dev/socially-app-go-server/internal/router"
)

func main() {
	ctx := context.Background()
	appScope := scope.ConfigureAppScope(ctx)
	log := appScope.Log
	cfg := appScope.Cfg

	app := router.New(appScope)

	log.Info("BUILD_TYPE: %s", cfg.BuildType.String())
	log.Info("Server running on port: %s\n", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
