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
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()

	apiV2 := app.Group("/api/v2")

	registration := apiV2.Group("/registration")
	{
		registration.Post("/", registrationHandler.Registration)
		registration.Post("/verify_otp", registrationHandler.VerifyOtp)
		registration.Post("/complete", registrationHandler.CompleteRegistration)
	}

	auth := apiV2.Group("/auth", middleware.DeviceIdMiddleware())
	{
		auth.Post("/login", authHandler.Login)
		auth.Patch("/refresh", authHandler.Refresh)
		auth.Post("/logout", middleware.AuthenticationMiddleware(), authHandler.Logout)
	}

	users := apiV2.Group("/users", middleware.AuthenticationMiddleware())
	{
		users.Get("/check_username", userHandler.CheckUsername)

		//% must be last
		users.Get("/:user_id", userHandler.GetById)
	}

	return app
}
