package router

import (
	scope "github.com/ebobola-dev/socially-app-go-server/internal/di"
	"github.com/ebobola-dev/socially-app-go-server/internal/handler"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/validation"

	"github.com/gofiber/fiber/v2"
)

func New(appScope scope.IAppScope, repositoriesScope scope.IRepositoriesScope, servicesScope scope.IServicesScope) *fiber.App {
	app := fiber.New()
	log := appScope.GetLogger()
	db := appScope.GetDB()
	otpRepository := repositoriesScope.GetOtpRepository()
	userRepository := repositoriesScope.GetUserRepository()
	emailService := servicesScope.GetEmailService()
	jwtService := servicesScope.GetJwtService()
	hashService := servicesScope.GetHashService()

	validate := validation.NewValidator()

	app.Use(middleware.LoggingMiddleware(log))
	app.Use(middleware.DatabaseSessionMiddleware(db))

	registrationHandler := handler.NewRegistrationHandler(log, validate, otpRepository, userRepository, emailService, jwtService, hashService)
	userHandler := handler.NewUserHandler(log, validate, userRepository)

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
