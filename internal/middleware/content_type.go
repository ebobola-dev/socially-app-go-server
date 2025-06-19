package middleware

import (
	"strings"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/gofiber/fiber/v2"
)

func ContentType(requiredContentType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentTypeHeaders := c.Request().Header.Peek("Content-Type")
		if len(contentTypeHeaders) == 0 {
			return common_error.NewBadContentTypeErr("", requiredContentType)
		}
		contentType := string(contentTypeHeaders)
		if !strings.HasPrefix(contentType, requiredContentType) {
			return common_error.NewBadContentTypeErr(contentType, requiredContentType)
		}
		return c.Next()
	}
}
