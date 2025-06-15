package main

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/database"
	scope "github.com/ebobola-dev/socially-app-go-server/internal/di"
	"github.com/ebobola-dev/socially-app-go-server/internal/router"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
)

func main() {
	cfg := config.Initialize()
	log := logger.Create(cfg)

	db := database.Connect(*cfg.Database)

	appScope := scope.NewAppScope(cfg, log, db)
	repositoriesScope := scope.NewRepositoriesScope()
	servicesScope := scope.NewServicesScope(cfg.SMTP, cfg.JWT)

	app := router.New(appScope, repositoriesScope, servicesScope)

	log.Info("BUILD_TYPE: %s", cfg.BuildType.String())
	log.Info("Server running on port: %s\n", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
