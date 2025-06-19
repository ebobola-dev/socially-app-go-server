package middleware

import (
	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	"github.com/gofiber/fiber/v2"
)

func AnyPrivileges(privNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := GetUserId(c)
		tx := GetTX(c)
		s := GetAppScope(c)

		has, getErr := s.UserRepository.HasAnyPrivileges(tx, userId, privNames...)
		if getErr != nil {
			return getErr
		}
		if !has {
			return auth_error.NewNoAnyPrivilegeError(privNames...)
		}

		return c.Next()
	}
}

func AllPrivileges(privNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := GetUserId(c)
		tx := GetTX(c)
		s := GetAppScope(c)

		has, getErr := s.UserRepository.HasAllPrivileges(tx, userId, privNames...)
		if getErr != nil {
			return getErr
		}
		if !has {
			return auth_error.NewNoAllPrivilegeError(privNames...)
		}

		return c.Next()
	}
}
