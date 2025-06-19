package middleware

import (
	scope "github.com/ebobola-dev/socially-app-go-server/internal/di"
	"github.com/gofiber/fiber/v2"
)

func AppScope(appScope *scope.AppScope) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("scope", appScope)
		return c.Next()
	}
}

func GetAppScope(c *fiber.Ctx) *scope.AppScope {
	scope, ok := c.Locals("scope").(*scope.AppScope)
	if !ok {
		panic("App scope not found in context")
	}
	return scope
}
