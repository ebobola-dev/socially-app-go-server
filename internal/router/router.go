package router

import (
	scope "github.com/ebobola-dev/socially-app-go-server/internal/di"
	"github.com/ebobola-dev/socially-app-go-server/internal/handler"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func New(appScope *scope.AppScope) *fiber.App {
	app := fiber.New()

	app.Use(middleware.ScopeMiddleware(appScope))
	app.Use(middleware.LoggingMiddleware(appScope.Log))
	app.Use(middleware.DatabaseSessionMiddleware(appScope.Db))

	registrationHandler := handler.NewRegistrationHandler()
	userHandler := handler.NewUserHandler()

	apiV2 := app.Group("/api/v2")

	registration := apiV2.Group("/registration")
	registration.Post("/", registrationHandler.Registration)
	registration.Post("/verify_otp", registrationHandler.VerifyOtp)
	registration.Post("/complete", registrationHandler.CompleteRegistration)

	users := apiV2.Group("/users")
	users.Get("/check_username", userHandler.CheckUsername)
	users.Get("/:user_id", userHandler.GetById)

	return app
}
