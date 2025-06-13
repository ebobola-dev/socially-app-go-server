package middleware

import (
	"time"

	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware(log logger.ILogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()
		duration := time.Since(start).Milliseconds()

		if status < 400 {

			log.Info("%s %s -> %d (%d ms)\n", method, path, status, duration)
		} else {
			log.Warn("%s %s -> %d (%d ms)\n", method, path, status, duration)
		}

		return err
	}
}
