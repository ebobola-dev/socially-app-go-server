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

func Logging(log logger.ILogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()
		duration := time.Since(start).Milliseconds()

		if err != nil {
			var apiErr api_error.IApiError
			var valiationErr validator.ValidationErrors
			var fiberErr *fiber.Error
			switch {
			case errors.As(err, &apiErr):
				log.Warning("%s %s -> %d (%d ms) %s\n", method, path, apiErr.StatusCode(), duration, err)
				return c.Status(apiErr.StatusCode()).JSON(apiErr.Response().ToJSON())
			case errors.As(err, &valiationErr):
				log.Warning("%s %s -> 400 (%d ms) %s\n", method, path, duration, err)
				errResp := response.ParseValidationErrors(err)
				return c.Status(fiber.StatusBadRequest).JSON(errResp.ToJSON())
			case errors.As(err, &fiberErr):
				log.Error("%s %s -> %d (%d ms) %+v\n", method, path, fiberErr.Code, duration, err)
				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"_message": fiberErr.Message,
				})
			default:
				log.Error("%s %s -> 500 (%d ms) %+v\n", method, path, duration, err)
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
