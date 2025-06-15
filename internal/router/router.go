package router

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/handler"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/ebobola-dev/socially-app-go-server/internal/validation"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

func New(cfg *config.Config, log logger.ILogger) *fiber.App {
	app := fiber.New()

	validate := validator.New()
	validate.RegisterValidation("otp_value", validation.OtpValueValidator)
	app.Use(middleware.LoggingMiddleware(log))

	registrationHandler := handler.NewRegistrationHandler(log, validate)
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
