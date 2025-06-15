package middleware

import (
	"errors"
	"time"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/go-playground/validator/v10"

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

		if err != nil {
			var apiErr api_error.ApiError
			var valiationErr validator.ValidationErrors
			switch {
			case errors.As(err, &apiErr):
				log.Warn("%s %s -> %d (%d ms) %s\n", method, path, apiErr.StatusCode(), duration, err)
				return c.Status(apiErr.StatusCode()).JSON(apiErr.Response().ToJSON())
			case errors.As(err, &valiationErr):
				log.Warn("%s %s -> 400 (%d ms) %s\n", method, path, duration, err)
				errResp := response.ParseValidationErrors(err)
				return c.Status(fiber.StatusBadRequest).JSON(errResp.ToJSON())
			default:
				log.Error(err)
				log.Warn("%s %s -> 500 (%d ms) %s\n", method, path, duration, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"_message": "Unexcepted server error",
				})
			}
		} else {
			log.Info("%s %s -> %d (%d ms)\n", method, path, status, duration)
		}

		return nil
	}
}
