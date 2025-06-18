package middleware

import (
	"errors"
	"strings"

	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//% Ну да костыль и что
		//% Мир не идеален, этот код тоже
		if strings.HasPrefix(c.Path(), "/api/v2/users/check_username") {
			return c.Next()
		}

		s := GetAppScope(c)
		tx := GetTX(c)
		authHeaders := c.Request().Header.Peek("Authorization")
		if len(authHeaders) == 0 {
			return auth_error.ErrMissingHeader
		}
		headerValue := string(authHeaders)
		parts := strings.SplitN(headerValue, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return auth_error.ErrWrongFormat
		}
		token := parts[1]
		if token == "" {
			return auth_error.ErrNoToken
		}
		userData, token_err := s.JwtService.ValidateUserAccess(token)
		if token_err != nil {
			if errors.Is(token_err, jwt.ErrTokenExpired) {
				return auth_error.ErrExpired
			}
			return auth_error.ErrInvalidToken
		}
		userId := userData.ID
		_, get_err := s.UserRepository.GetByID(tx, userId)
		if get_err != nil {
			if errors.Is(get_err, gorm.ErrRecordNotFound) {
				return auth_error.NewUserNotFoundError(userId.String())
			}
			return get_err
		}
		c.Locals("user_id", userId)
		return c.Next()
	}
}

func GetUserId(c *fiber.Ctx) uuid.UUID {
	userId, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		panic("user_id not found in context")
	}
	return userId
}
