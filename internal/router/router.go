package router

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/handler"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"

	"github.com/gofiber/fiber/v2"
)

func New(cfg *config.Config, log logger.ILogger) *fiber.App {
	app := fiber.New()

	app.Use(middleware.LoggingMiddleware(log))

	registrationHandler := handler.NewRegistrationHandler(log)
	userHandler := handler.NewUserHandler(log) //

	app.Post("/registration", registrationHandler.Registration)
	app.Post("/registration/verify_otp", registrationHandler.VerifyOtp)
	app.Post("/registration/complete", registrationHandler.CompleteRegistration)

	app.Get("/users/check_username", userHandler.CheckUsername)

	return app
}
