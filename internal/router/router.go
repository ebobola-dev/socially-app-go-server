package router

import (
	scope "github.com/ebobola-dev/socially-app-go-server/internal/di"
	"github.com/ebobola-dev/socially-app-go-server/internal/handler"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func New(appScope *scope.AppScope) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, //? 100 MB
	})

	app.Use(middleware.AppScope(appScope))
	app.Use(middleware.Logging(appScope.Log))
	app.Use(middleware.DatabaseSession(appScope.Db))
	app.Get("/metrics", monitor.New())

	registrationHandler := handler.NewRegistrationHandler()
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	privilegesHandler := handler.NewPrivilegeHandler()
	mediaHandler := handler.NewMediaHandler()

	apiV2 := app.Group("/api/v2")

	registration := apiV2.Group("/registration")
	{
		registration.Post("/", registrationHandler.Registration)
		registration.Post("/verify_otp", registrationHandler.VerifyOtp)
		registration.Post("/complete", registrationHandler.CompleteRegistration)
	}

	auth := apiV2.Group("/auth", middleware.DeviceId())
	{
		auth.Post("/login", authHandler.Login)
		auth.Patch("/refresh", authHandler.Refresh)
		auth.Post("/logout", middleware.Authentication(), authHandler.Logout)
	}

	users := apiV2.Group("/users", middleware.Authentication())
	{
		users.Get("/check_username", userHandler.CheckUsername)
		users.Delete("/delete_my_account", userHandler.DeleteMyAccount)
		users.Get("/search", middleware.Pagination(), userHandler.Search)
		users.Patch("/", userHandler.UpdateProfile)
		users.Patch("/password", userHandler.UpdatePassword)
		users.Patch("/avatar", middleware.ContentType("multipart/form-data"), userHandler.UpdateAvatar)
		users.Delete("/avatar", userHandler.DeleteAvatar)
		users.Get("/privileges", middleware.Pagination(), userHandler.GetPrivileges)
		users.Get("/:user_id", userHandler.GetById) //% must be last
	}

	privileges := apiV2.Group("/privileges", middleware.Authentication())
	{
		privileges.Get("/", middleware.Pagination(), privilegesHandler.GetAll)
		privileges.Get("/users", middleware.Pagination(), privilegesHandler.GetUsers)
		privileges.Post("/", middleware.AllPrivileges("owner"), privilegesHandler.Create)
		privileges.Delete("/", middleware.AllPrivileges("owner"), privilegesHandler.Delete)
	}

	media := apiV2.Group("/media")
	media.Get("/:bucket/*", mediaHandler.Get)

	return app
}
