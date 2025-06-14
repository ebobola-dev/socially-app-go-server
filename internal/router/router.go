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
	userHandler := handler.NewUserHandler(log)

	apiV2 := app.Group("/api/v2")

	registration := apiV2.Group("/registration")
	registration.Post("/", registrationHandler.Registration)
	registration.Post("/verify_otp", registrationHandler.VerifyOtp)
	registration.Post("/complete", registrationHandler.CompleteRegistration)

	users := apiV2.Group("/users")
	users.Get("/check_username", userHandler.CheckUsername)

	return app
}
